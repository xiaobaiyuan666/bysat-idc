package service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	providerDomain "idc-finance/internal/modules/provider/domain"
)

func (service *Service) ListAccounts(providerType string) []providerDomain.Account {
	if service.accounts == nil {
		return nil
	}
	return service.accounts.ListAccounts(strings.ToUpper(strings.TrimSpace(providerType)))
}

func (service *Service) GetAccount(id int64) (providerDomain.Account, bool) {
	if service.accounts == nil {
		return providerDomain.Account{}, false
	}
	return service.accounts.GetAccountByID(id)
}

func (service *Service) CreateAccount(request AccountUpsertRequest) (providerDomain.Account, error) {
	if service.accounts == nil {
		return providerDomain.Account{}, fmt.Errorf("未初始化自动化接口仓储")
	}
	account := normalizeAccountPayload(request)
	if account.Name == "" {
		return providerDomain.Account{}, fmt.Errorf("接口名称不能为空")
	}
	if account.ProviderType == "" {
		return providerDomain.Account{}, fmt.Errorf("接口类型不能为空")
	}
	if account.ProviderType != string(providerDomain.ProviderTypeManual) && account.BaseURL == "" {
		return providerDomain.Account{}, fmt.Errorf("接口地址不能为空")
	}
	created := service.accounts.CreateAccount(account)
	if created.ID == 0 {
		return providerDomain.Account{}, fmt.Errorf("创建自动化接口失败")
	}
	return created, nil
}

func (service *Service) UpdateAccount(id int64, request AccountUpsertRequest) (providerDomain.Account, error) {
	if service.accounts == nil {
		return providerDomain.Account{}, fmt.Errorf("未初始化自动化接口仓储")
	}
	_, exists := service.accounts.GetAccountByID(id)
	if !exists {
		return providerDomain.Account{}, fmt.Errorf("自动化接口不存在")
	}
	account := normalizeAccountPayload(request)
	account.ID = id
	if account.Name == "" {
		return providerDomain.Account{}, fmt.Errorf("接口名称不能为空")
	}
	if account.ProviderType == "" {
		return providerDomain.Account{}, fmt.Errorf("接口类型不能为空")
	}
	updated, ok := service.accounts.UpdateAccount(id, account)
	if !ok {
		return providerDomain.Account{}, fmt.Errorf("更新自动化接口失败")
	}
	service.clearSessionCacheForAccount(updated.ProviderType, updated.ID)
	return updated, nil
}

func (service *Service) DeleteAccount(id int64) error {
	if service.accounts == nil {
		return fmt.Errorf("未初始化自动化接口仓储")
	}
	account, exists := service.accounts.GetAccountByID(id)
	if !exists {
		return fmt.Errorf("自动化接口不存在")
	}
	if service.db != nil {
		count, err := service.countAccountBindings(id)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("仍有商品或服务正在使用该接口，不能删除")
		}
	}
	if !service.accounts.DeleteAccount(id) {
		return fmt.Errorf("删除自动化接口失败")
	}
	service.clearSessionCacheForAccount(account.ProviderType, account.ID)
	return nil
}

func (service *Service) ensureDefaultAccounts() {
	if service.accounts == nil {
		return
	}

	if strings.TrimSpace(service.config.BaseURL) != "" {
		service.ensureDefaultAccount(providerDomain.Account{
			ProviderType:       string(providerDomain.ProviderTypeMofangCloud),
			Name:               "默认魔方云",
			BaseURL:            service.config.BaseURL,
			Username:           service.config.Username,
			Password:           service.config.Password,
			SourceName:         "魔方云",
			Lang:               firstNonEmpty(strings.TrimSpace(service.config.Lang), "zh-cn"),
			ListPath:           service.config.ListPath,
			DetailPath:         service.config.DetailPath,
			InsecureSkipVerify: service.config.InsecureSkipVerify,
			AutoUpdate:         true,
			Status:             providerDomain.AccountStatusActive,
		})
	}

	if strings.TrimSpace(service.financeConfig.BaseURL) != "" {
		service.ensureDefaultAccount(providerDomain.Account{
			ProviderType:       string(providerDomain.ProviderTypeZjmfAPI),
			Name:               firstNonEmpty(strings.TrimSpace(service.financeConfig.SourceName), "默认上下游财务"),
			BaseURL:            service.financeConfig.BaseURL,
			Username:           service.financeConfig.Username,
			Password:           service.financeConfig.Password,
			SourceName:         firstNonEmpty(strings.TrimSpace(service.financeConfig.SourceName), "上下游财务"),
			InsecureSkipVerify: service.financeConfig.InsecureSkipVerify,
			AutoUpdate:         true,
			Status:             providerDomain.AccountStatusActive,
		})
	}
}

func (service *Service) ensureDefaultAccount(account providerDomain.Account) {
	if service.accounts == nil || strings.TrimSpace(account.BaseURL) == "" {
		return
	}
	if _, exists := service.accounts.FindByProviderAndBaseURL(account.ProviderType, account.BaseURL); exists {
		return
	}
	_ = service.accounts.CreateAccount(account)
}

func normalizeAccountPayload(request AccountUpsertRequest) providerDomain.Account {
	return providerDomain.Account{
		ProviderType:       strings.ToUpper(strings.TrimSpace(request.ProviderType)),
		Name:               strings.TrimSpace(request.Name),
		BaseURL:            strings.TrimSpace(request.BaseURL),
		Username:           strings.TrimSpace(request.Username),
		Password:           strings.TrimSpace(request.Password),
		SourceName:         strings.TrimSpace(request.SourceName),
		ContactWay:         strings.TrimSpace(request.ContactWay),
		Description:        strings.TrimSpace(request.Description),
		AccountMode:        strings.TrimSpace(request.AccountMode),
		Lang:               firstNonEmpty(strings.TrimSpace(request.Lang), "zh-cn"),
		ListPath:           strings.TrimSpace(request.ListPath),
		DetailPath:         strings.TrimSpace(request.DetailPath),
		InsecureSkipVerify: request.InsecureSkipVerify,
		AutoUpdate:         request.AutoUpdate,
		Status:             providerDomain.AccountStatus(firstNonEmpty(strings.TrimSpace(request.Status), string(providerDomain.AccountStatusActive))),
		ExtraConfig:        strings.TrimSpace(request.ExtraConfig),
	}
}

func (service *Service) resolveMofangAccount(accountID int64) (resolvedProviderAccount, error) {
	if accountID > 0 && service.accounts != nil {
		account, exists := service.accounts.GetAccountByID(accountID)
		if !exists {
			return resolvedProviderAccount{}, fmt.Errorf("未找到魔方云接口账号")
		}
		if !strings.EqualFold(account.ProviderType, string(providerDomain.ProviderTypeMofangCloud)) {
			return resolvedProviderAccount{}, fmt.Errorf("所选接口不是魔方云类型")
		}
		return resolvedProviderAccount{
			AccountID:          account.ID,
			ProviderType:       account.ProviderType,
			Name:               account.Name,
			BaseURL:            account.BaseURL,
			Username:           account.Username,
			Password:           account.Password,
			SourceName:         firstNonEmpty(account.SourceName, account.Name),
			Lang:               firstNonEmpty(account.Lang, "zh-cn"),
			ListPath:           firstNonEmpty(account.ListPath, service.config.ListPath),
			DetailPath:         firstNonEmpty(account.DetailPath, service.config.DetailPath),
			InsecureSkipVerify: account.InsecureSkipVerify,
			AuthMode:           "mofang-jwt",
			Status:             account.Status,
		}, nil
	}

	if strings.TrimSpace(service.config.BaseURL) == "" {
		return resolvedProviderAccount{}, fmt.Errorf("未配置魔方云接口")
	}

	return resolvedProviderAccount{
		AccountID:          0,
		ProviderType:       string(providerDomain.ProviderTypeMofangCloud),
		Name:               "环境变量魔方云",
		BaseURL:            service.config.BaseURL,
		Username:           service.config.Username,
		Password:           service.config.Password,
		SourceName:         "环境变量魔方云",
		Lang:               firstNonEmpty(service.config.Lang, "zh-cn"),
		ListPath:           service.config.ListPath,
		DetailPath:         service.config.DetailPath,
		InsecureSkipVerify: service.config.InsecureSkipVerify,
		AuthMode:           "mofang-jwt",
		Status:             providerDomain.AccountStatusActive,
	}, nil
}

func (service *Service) resolveFinanceAccount(accountID int64) (resolvedProviderAccount, error) {
	if accountID > 0 && service.accounts != nil {
		account, exists := service.accounts.GetAccountByID(accountID)
		if !exists {
			return resolvedProviderAccount{}, fmt.Errorf("未找到上下游接口账号")
		}
		if !strings.EqualFold(account.ProviderType, string(providerDomain.ProviderTypeZjmfAPI)) &&
			!strings.EqualFold(account.ProviderType, string(providerDomain.ProviderTypeWHMCS)) &&
			!strings.EqualFold(account.ProviderType, string(providerDomain.ProviderTypeManual)) {
			return resolvedProviderAccount{}, fmt.Errorf("所选接口不是财务上下游类型")
		}
		return resolvedProviderAccount{
			AccountID:          account.ID,
			ProviderType:       account.ProviderType,
			Name:               account.Name,
			BaseURL:            account.BaseURL,
			Username:           account.Username,
			Password:           account.Password,
			SourceName:         firstNonEmpty(account.SourceName, account.Name),
			Lang:               firstNonEmpty(account.Lang, "zh-cn"),
			InsecureSkipVerify: account.InsecureSkipVerify,
			AuthMode:           "finance-bearer",
			Status:             account.Status,
		}, nil
	}

	if strings.TrimSpace(service.financeConfig.BaseURL) == "" {
		return resolvedProviderAccount{}, fmt.Errorf("未配置上下游财务接口")
	}

	return resolvedProviderAccount{
		AccountID:          0,
		ProviderType:       string(providerDomain.ProviderTypeZjmfAPI),
		Name:               firstNonEmpty(service.financeConfig.SourceName, "环境变量上下游"),
		BaseURL:            service.financeConfig.BaseURL,
		Username:           service.financeConfig.Username,
		Password:           service.financeConfig.Password,
		SourceName:         firstNonEmpty(service.financeConfig.SourceName, "环境变量上下游"),
		Lang:               "zh-cn",
		InsecureSkipVerify: service.financeConfig.InsecureSkipVerify,
		AuthMode:           "finance-bearer",
		Status:             providerDomain.AccountStatusActive,
	}, nil
}

func (service *Service) requestWithAccount(account resolvedProviderAccount, method, path string, payload map[string]any) (rawResponse, error) {
	token, activeURL, err := service.getAccountToken(account, false)
	if err != nil {
		return rawResponse{}, err
	}

	requestURL := joinURL(activeURL, path)
	bodyString := ""
	if strings.EqualFold(method, http.MethodGet) && len(payload) > 0 {
		values := url.Values{}
		for key, value := range payload {
			if strings.TrimSpace(fmt.Sprint(value)) == "" {
				continue
			}
			values.Set(key, fmt.Sprint(value))
		}
		if encoded := values.Encode(); encoded != "" {
			separator := "?"
			if strings.Contains(requestURL, "?") {
				separator = "&"
			}
			requestURL += separator + encoded
		}
	} else if len(payload) > 0 {
		values := url.Values{}
		for key, value := range payload {
			if strings.TrimSpace(fmt.Sprint(value)) == "" {
				continue
			}
			values.Set(key, fmt.Sprint(value))
		}
		bodyString = values.Encode()
	}

	request, err := http.NewRequest(method, requestURL, strings.NewReader(bodyString))
	if err != nil {
		return rawResponse{}, err
	}
	if strings.HasPrefix(account.AuthMode, "mofang") {
		request.Header.Set("think-lang", firstNonEmpty(account.Lang, "zh-cn"))
		request.Header.Set("Authorization", "JWT "+token)
		request.Header.Set("access-token", token)
	} else {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	if bodyString != "" {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}

	response, err := service.httpClient(account.InsecureSkipVerify).Do(request)
	if err != nil {
		return rawResponse{}, err
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	var parsed any
	_ = json.Unmarshal(body, &parsed)

	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusMethodNotAllowed {
		service.clearSessionCacheForAccount(account.ProviderType, account.AccountID)
		return service.requestWithAccount(account, method, path, payload)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		detail := strings.TrimSpace(pickString(extractRecord(parsed), "error", "msg", "message"))
		if detail == "" {
			detail = strings.TrimSpace(string(body))
		}
		if detail != "" {
			return rawResponse{}, fmt.Errorf("%s 请求失败: HTTP %d - %s", account.SourceName, response.StatusCode, detail)
		}
		return rawResponse{}, fmt.Errorf("%s 请求失败: HTTP %d", account.SourceName, response.StatusCode)
	}

	if !strings.HasPrefix(account.AuthMode, "mofang") {
		root, _ := parsed.(map[string]any)
		if status := pickInt(root, "status"); status != 0 && status != 200 {
			return rawResponse{}, fmt.Errorf(firstNonEmpty(pickString(root, "msg", "message"), "上游接口返回业务错误"))
		}
	}

	return rawResponse{
		StatusCode: response.StatusCode,
		Text:       string(body),
		Parsed:     parsed,
	}, nil
}

func (service *Service) getAccountToken(account resolvedProviderAccount, forceRefresh bool) (string, string, error) {
	key := service.accountSessionKey(account.ProviderType, account.AccountID)
	service.sessionStoreMu.Lock()
	cached, ok := service.sessionStore[key]
	if ok && !forceRefresh && strings.TrimSpace(cached.Token) != "" {
		service.sessionStoreMu.Unlock()
		return cached.Token, cached.ActiveURL, nil
	}
	service.sessionStoreMu.Unlock()

	var (
		token     string
		activeURL string
		err       error
	)
	if strings.HasPrefix(account.AuthMode, "mofang") {
		token, activeURL, err = service.loginMofangAccount(account)
	} else {
		token, activeURL, err = service.loginFinanceAccount(account)
	}
	if err != nil {
		return "", "", err
	}

	service.sessionStoreMu.Lock()
	service.sessionStore[key] = providerSessionState{
		Token:      token,
		ActiveURL:  activeURL,
		LastAuthAt: time.Now(),
	}
	service.sessionStoreMu.Unlock()
	return token, activeURL, nil
}

func (service *Service) loginMofangAccount(account resolvedProviderAccount) (string, string, error) {
	if strings.TrimSpace(account.Username) == "" || strings.TrimSpace(account.Password) == "" {
		return "", "", fmt.Errorf("未配置魔方云登录账号")
	}
	values := url.Values{}
	values.Set("username", account.Username)
	values.Set("password", account.Password)
	values.Set("customfield[google_code]", "")

	var lastErr error
	for _, baseURL := range baseURLCandidates(account.BaseURL) {
		request, err := http.NewRequest(http.MethodPost, joinURL(baseURL, "/v1/login"), strings.NewReader(values.Encode()))
		if err != nil {
			lastErr = err
			continue
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		request.Header.Set("think-lang", firstNonEmpty(account.Lang, "zh-cn"))

		response, err := service.httpClient(account.InsecureSkipVerify).Do(request)
		if err != nil {
			lastErr = err
			continue
		}
		body, _ := io.ReadAll(response.Body)
		_ = response.Body.Close()

		var parsed any
		_ = json.Unmarshal(body, &parsed)
		token := extractToken(parsed)
		if token == "" {
			token = strings.Trim(strings.TrimSpace(string(body)), "\"")
		}
		if response.StatusCode >= 200 && response.StatusCode < 300 && token != "" {
			return token, baseURL, nil
		}
		lastErr = fmt.Errorf("魔方云登录失败: HTTP %d", response.StatusCode)
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("魔方云登录失败")
	}
	return "", "", lastErr
}

func (service *Service) loginFinanceAccount(account resolvedProviderAccount) (string, string, error) {
	if strings.TrimSpace(account.Username) == "" || strings.TrimSpace(account.Password) == "" {
		return "", "", fmt.Errorf("未配置上游财务登录账号")
	}
	values := url.Values{}
	values.Set("username", account.Username)
	values.Set("password", account.Password)

	var lastErr error
	for _, baseURL := range baseURLCandidates(account.BaseURL) {
		request, err := http.NewRequest(http.MethodPost, joinURL(baseURL, "/zjmf_api_login"), strings.NewReader(values.Encode()))
		if err != nil {
			lastErr = err
			continue
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

		response, err := service.httpClient(account.InsecureSkipVerify).Do(request)
		if err != nil {
			lastErr = err
			continue
		}
		body, _ := io.ReadAll(response.Body)
		_ = response.Body.Close()

		var parsed any
		_ = json.Unmarshal(body, &parsed)
		token := extractFinanceToken(parsed)
		if response.StatusCode >= 200 && response.StatusCode < 300 && token != "" {
			return token, baseURL, nil
		}
		lastErr = fmt.Errorf("上游财务登录失败: HTTP %d", response.StatusCode)
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("上游财务登录失败")
	}
	return "", "", lastErr
}

func (service *Service) clearSessionCacheForAccount(providerType string, accountID int64) {
	service.sessionStoreMu.Lock()
	defer service.sessionStoreMu.Unlock()
	delete(service.sessionStore, service.accountSessionKey(providerType, accountID))
}

func (service *Service) accountSessionKey(providerType string, accountID int64) string {
	return strings.ToUpper(strings.TrimSpace(providerType)) + ":" + fmt.Sprint(accountID)
}

func (service *Service) httpClient(insecureSkipVerify bool) *http.Client {
	return &http.Client{
		Timeout: 20 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify}, //nolint:gosec
		},
	}
}

func baseURLCandidates(baseURL string) []string {
	trimmed := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if trimmed == "" {
		return nil
	}
	if strings.HasPrefix(trimmed, "https://") {
		return []string{trimmed, "http://" + strings.TrimPrefix(trimmed, "https://")}
	}
	return []string{trimmed}
}

func (service *Service) countAccountBindings(accountID int64) (int, error) {
	if service.db == nil || accountID == 0 {
		return 0, nil
	}

	var serviceCount int
	if err := service.db.QueryRow(`
SELECT COUNT(*)
FROM services
WHERE provider_account_id = ?`, accountID).Scan(&serviceCount); err != nil {
		return 0, err
	}

	var productCount int
	rows, err := service.db.Query(`SELECT automation_config, upstream_mapping FROM products`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var automationJSON []byte
		var upstreamJSON []byte
		if err := rows.Scan(&automationJSON, &upstreamJSON); err != nil {
			return 0, err
		}
		if matchAccountIDInJSON(automationJSON, accountID) || matchAccountIDInJSON(upstreamJSON, accountID) {
			productCount++
		}
	}
	return serviceCount + productCount, nil
}

func matchAccountIDInJSON(payload []byte, accountID int64) bool {
	if len(payload) == 0 || accountID == 0 {
		return false
	}
	var record map[string]any
	if err := json.Unmarshal(payload, &record); err != nil {
		return false
	}
	return int64(pickInt(record, "providerAccountId")) == accountID
}

func (service *Service) CheckHealth(accountIDs ...int64) HealthResponse {
	accountID := firstID(accountIDs...)
	if accountID == 0 {
		return service.legacyCheckHealth()
	}

	account, err := service.resolveMofangAccount(accountID)
	if err != nil {
		return HealthResponse{
			Enabled:   false,
			Connected: false,
			AuthMode:  "disabled",
			Message:   err.Error(),
		}
	}

	_, activeURL, err := service.getAccountToken(account, false)
	if err != nil {
		return HealthResponse{
			Enabled:    true,
			Connected:  false,
			BaseURL:    account.BaseURL,
			ActiveURL:  activeURL,
			AuthMode:   "jwt",
			Message:    err.Error(),
			LastAuthAt: service.lookupSessionLastAuth(account.ProviderType, account.AccountID),
		}
	}

	return HealthResponse{
		Enabled:    true,
		Connected:  true,
		BaseURL:    account.BaseURL,
		ActiveURL:  activeURL,
		AuthMode:   "jwt",
		Message:    "魔方云连接正常",
		LastAuthAt: service.lookupSessionLastAuth(account.ProviderType, account.AccountID),
	}
}

func (service *Service) ListInstances(accountIDs ...int64) ([]InstanceSummary, error) {
	accountID := firstID(accountIDs...)
	if accountID == 0 {
		return service.legacyListInstances()
	}

	account, err := service.resolveMofangAccount(accountID)
	if err != nil {
		return nil, err
	}

	page := 1
	totalPages := 1
	result := make([]InstanceSummary, 0)
	for page <= totalPages {
		raw, err := service.requestWithAccount(account, http.MethodGet, buildListPathForAccount(account, page, 100), nil)
		if err != nil {
			return nil, err
		}

		records := extractRecords(raw.Parsed)
		for _, record := range records {
			result = append(result, mapInstance(record))
		}

		meta := extractMeta(raw.Parsed)
		if meta.totalPages <= 1 {
			break
		}
		totalPages = meta.totalPages
		page++
	}

	return result, nil
}

func (service *Service) GetInstanceDetail(remoteID string, accountIDs ...int64) (InstanceSummary, error) {
	accountID := firstID(accountIDs...)
	if accountID == 0 {
		return service.legacyGetInstanceDetail(remoteID)
	}

	account, err := service.resolveMofangAccount(accountID)
	if err != nil {
		return InstanceSummary{}, err
	}

	raw, err := service.requestWithAccount(account, http.MethodGet, buildDetailPathForAccount(account, remoteID), nil)
	if err != nil {
		return InstanceSummary{}, err
	}

	record := extractRecord(raw.Parsed)
	if record == nil {
		return InstanceSummary{}, fmt.Errorf("远程实例详情为空")
	}

	item := mapInstance(record)
	if item.RemoteID == "" {
		item.RemoteID = remoteID
	}
	item.Raw = record
	return item, nil
}

func (service *Service) ExecuteAction(remoteID, action string, request ActionRequest, accountIDs ...int64) (ActionResponse, error) {
	accountID := firstID(accountIDs...)
	if accountID == 0 {
		return service.legacyExecuteAction(remoteID, action, request)
	}

	account, err := service.resolveMofangAccount(accountID)
	if err != nil {
		return ActionResponse{}, err
	}

	specs := map[string]struct {
		Method string
		Path   string
		Body   map[string]any
	}{
		"activate":       {Method: http.MethodPost, Path: "/v1/clouds/:id/on"},
		"power-on":       {Method: http.MethodPost, Path: "/v1/clouds/:id/on"},
		"power-off":      {Method: http.MethodPost, Path: "/v1/clouds/:id/off"},
		"reboot":         {Method: http.MethodPost, Path: "/v1/clouds/:id/reboot"},
		"hard-power-off": {Method: http.MethodPost, Path: "/admin/v1/mf_cloud/:id/hard_off"},
		"hard-reboot":    {Method: http.MethodPost, Path: "/admin/v1/mf_cloud/:id/hard_reboot"},
		"suspend":        {Method: http.MethodPost, Path: "/v1/clouds/:id/suspend"},
		"unsuspend":      {Method: http.MethodPost, Path: "/v1/clouds/:id/unsuspend"},
		"terminate":      {Method: http.MethodDelete, Path: "/v1/clouds/:id"},
		"reset-password": {Method: http.MethodPut, Path: "/v1/clouds/:id/password", Body: map[string]any{"password": request.Password}},
		"reinstall":      {Method: http.MethodPut, Path: "/v1/clouds/:id/reinstall", Body: map[string]any{"image_name": request.ImageName}},
		"vnc":            {Method: http.MethodPost, Path: "/v1/clouds/:id/vnc"},
		"get-vnc":        {Method: http.MethodPost, Path: "/v1/clouds/:id/vnc"},
	}

	spec, ok := specs[action]
	if !ok {
		return ActionResponse{}, fmt.Errorf("不支持的魔方云动作: %s", action)
	}

	raw, err := service.requestWithAccount(account, spec.Method, strings.ReplaceAll(spec.Path, ":id", remoteID), spec.Body)
	if err != nil {
		return ActionResponse{}, err
	}

	record := extractRecord(raw.Parsed)
	message := pickString(record, "msg", "message", "status_desc")
	if message == "" {
		message = "操作已提交"
	}

	return ActionResponse{
		OK:       raw.StatusCode >= 200 && raw.StatusCode < 300,
		Action:   action,
		RemoteID: remoteID,
		Status:   normalizeStatus(pickString(record, "status", "instance_status", "cloud_status", "power_status")),
		Message:  message,
		Response: record,
	}, nil
}

func buildListPathForAccount(account resolvedProviderAccount, page, limit int) string {
	basePath := strings.TrimSpace(account.ListPath)
	if basePath == "" {
		basePath = "/v1/clouds"
	}
	separator := "?"
	if strings.Contains(basePath, "?") {
		separator = "&"
	}
	return fmt.Sprintf("%s%spage=%d&limit=%d", basePath, separator, page, limit)
}

func buildDetailPathForAccount(account resolvedProviderAccount, remoteID string) string {
	basePath := strings.TrimSpace(account.DetailPath)
	if basePath == "" {
		basePath = "/v1/clouds/:id"
	}
	return strings.ReplaceAll(basePath, ":id", remoteID)
}

func firstID(ids ...int64) int64 {
	if len(ids) == 0 {
		return 0
	}
	return ids[0]
}

func (service *Service) lookupSessionLastAuth(providerType string, accountID int64) time.Time {
	service.sessionStoreMu.Lock()
	defer service.sessionStoreMu.Unlock()
	return service.sessionStore[service.accountSessionKey(providerType, accountID)].LastAuthAt
}

func (service *Service) CheckFinanceHealth(accountIDs ...int64) HealthResponse {
	accountID := firstID(accountIDs...)
	if accountID == 0 {
		return service.legacyCheckFinanceHealth()
	}

	account, err := service.resolveFinanceAccount(accountID)
	if err != nil {
		return HealthResponse{
			Enabled:   false,
			Connected: false,
			AuthMode:  "disabled",
			Message:   err.Error(),
		}
	}

	_, activeURL, err := service.getAccountToken(account, false)
	if err != nil {
		return HealthResponse{
			Enabled:    true,
			Connected:  false,
			BaseURL:    account.BaseURL,
			ActiveURL:  activeURL,
			AuthMode:   "bearer",
			Message:    err.Error(),
			LastAuthAt: service.lookupSessionLastAuth(account.ProviderType, account.AccountID),
		}
	}

	return HealthResponse{
		Enabled:    true,
		Connected:  true,
		BaseURL:    account.BaseURL,
		ActiveURL:  activeURL,
		AuthMode:   "bearer",
		Message:    "上下游接口连接正常",
		LastAuthAt: service.lookupSessionLastAuth(account.ProviderType, account.AccountID),
	}
}

func (service *Service) ListUpstreamProducts(accountIDs ...int64) ([]UpstreamCatalogGroup, error) {
	accountID := firstID(accountIDs...)
	if accountID == 0 {
		return service.legacyListUpstreamProducts()
	}

	account, err := service.resolveFinanceAccount(accountID)
	if err != nil {
		return nil, err
	}

	raw, err := service.requestWithAccount(account, http.MethodGet, "/cart/all", nil)
	if err != nil {
		return nil, err
	}

	root, ok := raw.Parsed.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("?????????????")
	}

	return parseUpstreamCatalogGroups(root)
}

func parseUpstreamCatalogGroups(root map[string]any) ([]UpstreamCatalogGroup, error) {
	if root == nil {
		return nil, fmt.Errorf("??????????")
	}

	groups := extractUpstreamCatalogGroupRecords(root["data"])
	if len(groups) == 0 {
		return []UpstreamCatalogGroup{}, nil
	}

	result := make([]UpstreamCatalogGroup, 0, len(groups))
	for _, groupRecord := range groups {
		group := UpstreamCatalogGroup{
			GroupID:   pickString(groupRecord, "id", "group_id"),
			GroupName: pickString(groupRecord, "name", "group_name"),
			Products:  make([]UpstreamCatalogProduct, 0),
		}

		for _, productRecord := range extractUpstreamCatalogProductRecords(groupRecord["products"]) {
			group.Products = append(group.Products, UpstreamCatalogProduct{
				RemoteProductCode: pickString(productRecord, "id"),
				GroupName:         group.GroupName,
				Name:              pickString(productRecord, "name"),
				ProductType:       strings.ToUpper(pickString(productRecord, "type")),
				Description:       htmlText(pickString(productRecord, "description")),
			})
		}

		if len(group.Products) > 0 {
			result = append(result, group)
		}
	}
	return result, nil
}

func extractUpstreamCatalogGroupRecords(input any) []map[string]any {
	switch typed := input.(type) {
	case []any:
		result := make([]map[string]any, 0, len(typed))
		for _, item := range typed {
			if record, ok := item.(map[string]any); ok {
				result = append(result, record)
			}
		}
		return result
	case map[string]any:
		if groups, ok := typed["groups"].([]any); ok {
			result := make([]map[string]any, 0, len(groups))
			for _, item := range groups {
				if record, ok := item.(map[string]any); ok {
					result = append(result, record)
				}
			}
			return result
		}
		if products, ok := typed["products"].([]any); ok {
			looksLikeGroups := false
			for _, item := range products {
				if record, ok := item.(map[string]any); ok {
					if _, hasChildren := record["products"]; hasChildren {
						looksLikeGroups = true
						break
					}
					if _, hasGroupName := record["group_name"]; hasGroupName {
						looksLikeGroups = true
						break
					}
					if _, hasGroupID := record["group_id"]; hasGroupID {
						looksLikeGroups = true
						break
					}
				}
			}
			if looksLikeGroups {
				result := make([]map[string]any, 0, len(products))
				for _, item := range products {
					if record, ok := item.(map[string]any); ok {
						result = append(result, record)
					}
				}
				return result
			}
			group := map[string]any{
				"id":       pickString(typed, "id", "group_id"),
				"name":     pickString(typed, "name", "group_name"),
				"products": products,
			}
			return []map[string]any{group}
		}
	}
	return nil
}

func extractUpstreamCatalogProductRecords(input any) []map[string]any {
	items, ok := input.([]any)
	if !ok {
		return nil
	}
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if record, ok := item.(map[string]any); ok {
			result = append(result, record)
		}
	}
	return result
}

func (service *Service) GetUpstreamProductTemplate(remoteProductCode string, accountIDs ...int64) (UpstreamProductTemplate, error) {
	accountID := firstID(accountIDs...)
	if accountID == 0 {
		return service.legacyGetUpstreamProductTemplate(remoteProductCode)
	}

	account, err := service.resolveFinanceAccount(accountID)
	if err != nil {
		return UpstreamProductTemplate{}, err
	}
	if strings.TrimSpace(remoteProductCode) == "" {
		return UpstreamProductTemplate{}, fmt.Errorf("????????")
	}

	raw, err := service.requestWithAccount(account, http.MethodGet, "/cart/get_product_config", map[string]any{
		"pid": remoteProductCode,
	})
	if err != nil {
		return UpstreamProductTemplate{}, err
	}

	root, ok := raw.Parsed.(map[string]any)
	if !ok {
		return UpstreamProductTemplate{}, fmt.Errorf("?????????????")
	}

	data, _ := root["data"].(map[string]any)
	productRecord, _ := data["products"].(map[string]any)

	template := UpstreamProductTemplate{
		RemoteProductCode: remoteProductCode,
		Name:              pickString(productRecord, "name"),
		Description:       htmlText(pickString(productRecord, "description")),
		ProductType:       strings.ToUpper(pickString(productRecord, "type")),
		Currency:          firstNonEmpty(firstNonEmpty(pickString(data, "currency"), pickString(root, "currency")), "CNY"),
		Pricing:           buildUpstreamPricing(data),
		ConfigOptions:     buildUpstreamConfigOptions(data),
	}

	if groups, err := service.ListUpstreamProducts(accountID); err == nil {
		for _, group := range groups {
			for _, item := range group.Products {
				if item.RemoteProductCode == remoteProductCode {
					template.GroupName = group.GroupName
					if template.Name == "" {
						template.Name = item.Name
					}
					if template.Description == "" {
						template.Description = item.Description
					}
					if template.ProductType == "" {
						template.ProductType = item.ProductType
					}
					return template, nil
				}
			}
		}
	}

	return template, nil
}
