package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type UpstreamCatalogProduct struct {
	RemoteProductCode string `json:"remoteProductCode"`
	GroupName         string `json:"groupName"`
	Name              string `json:"name"`
	ProductType       string `json:"productType"`
	Description       string `json:"description"`
}

type UpstreamCatalogGroup struct {
	GroupID   string                  `json:"groupId"`
	GroupName string                  `json:"groupName"`
	Products  []UpstreamCatalogProduct `json:"products"`
}

type UpstreamCyclePrice struct {
	CycleCode string  `json:"cycleCode"`
	CycleName string  `json:"cycleName"`
	Price     float64 `json:"price"`
	SetupFee  float64 `json:"setupFee"`
}

type UpstreamConfigChoice struct {
	Value      string  `json:"value"`
	Label      string  `json:"label"`
	PriceDelta float64 `json:"priceDelta"`
}

type UpstreamConfigOption struct {
	Code         string                `json:"code"`
	Name         string                `json:"name"`
	InputType    string                `json:"inputType"`
	Required     bool                  `json:"required"`
	DefaultValue string                `json:"defaultValue"`
	Description  string                `json:"description"`
	Choices      []UpstreamConfigChoice `json:"choices"`
}

type UpstreamProductTemplate struct {
	RemoteProductCode string                `json:"remoteProductCode"`
	GroupName         string                `json:"groupName"`
	Name              string                `json:"name"`
	Description       string                `json:"description"`
	ProductType       string                `json:"productType"`
	Currency          string                `json:"currency"`
	Pricing           []UpstreamCyclePrice  `json:"pricing"`
	ConfigOptions     []UpstreamConfigOption `json:"configOptions"`
}

func (service *Service) legacyCheckFinanceHealth() HealthResponse {
	if strings.TrimSpace(service.financeConfig.BaseURL) == "" {
		return HealthResponse{
			Enabled:   false,
			Connected: false,
			AuthMode:  "disabled",
			Message:   "未配置上游财务地址",
		}
	}

	_, err := service.getFinanceSessionToken(false)
	if err != nil {
		return HealthResponse{
			Enabled:    true,
			Connected:  false,
			BaseURL:    service.financeConfig.BaseURL,
			ActiveURL:  service.financeActiveURL,
			AuthMode:   "bearer",
			Message:    err.Error(),
			LastAuthAt: service.financeLastAuthAt,
		}
	}

	return HealthResponse{
		Enabled:    true,
		Connected:  true,
		BaseURL:    service.financeConfig.BaseURL,
		ActiveURL:  service.financeActiveURL,
		AuthMode:   "bearer",
		Message:    "上游财务连接正常",
		LastAuthAt: service.financeLastAuthAt,
	}
}

func (service *Service) legacyListUpstreamProducts() ([]UpstreamCatalogGroup, error) {
	if strings.TrimSpace(service.financeConfig.BaseURL) == "" {
		return nil, fmt.Errorf("未配置上游财务地址")
	}

	raw, err := service.financeRequest(http.MethodGet, "/cart/all", nil)
	if err != nil {
		return nil, err
	}

	root, ok := raw.Parsed.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("上游商品目录返回格式不正确")
	}

	data, _ := root["data"].([]any)
	result := make([]UpstreamCatalogGroup, 0, len(data))
	for _, groupAny := range data {
		groupRecord, ok := groupAny.(map[string]any)
		if !ok {
			continue
		}

		group := UpstreamCatalogGroup{
			GroupID:   pickString(groupRecord, "id"),
			GroupName: pickString(groupRecord, "name"),
			Products:  make([]UpstreamCatalogProduct, 0),
		}

		productItems, _ := groupRecord["products"].([]any)
		for _, productAny := range productItems {
			productRecord, ok := productAny.(map[string]any)
			if !ok {
				continue
			}
			group.Products = append(group.Products, UpstreamCatalogProduct{
				RemoteProductCode: pickString(productRecord, "id"),
				GroupName:         group.GroupName,
				Name:              pickString(productRecord, "name"),
				ProductType:       strings.ToUpper(pickString(productRecord, "type")),
				Description:       htmlText(pickString(productRecord, "description")),
			})
		}

		result = append(result, group)
	}

	return result, nil
}

func (service *Service) legacyGetUpstreamProductTemplate(remoteProductCode string) (UpstreamProductTemplate, error) {
	if strings.TrimSpace(remoteProductCode) == "" {
		return UpstreamProductTemplate{}, fmt.Errorf("缺少远端商品编码")
	}

	raw, err := service.financeRequest(http.MethodGet, "/cart/get_product_config", map[string]any{
		"pid": remoteProductCode,
	})
	if err != nil {
		return UpstreamProductTemplate{}, err
	}

	root, ok := raw.Parsed.(map[string]any)
	if !ok {
		return UpstreamProductTemplate{}, fmt.Errorf("上游商品配置返回格式不正确")
	}

	data, _ := root["data"].(map[string]any)
	productRecord, _ := data["products"].(map[string]any)

	template := UpstreamProductTemplate{
		RemoteProductCode: remoteProductCode,
		Name:              pickString(productRecord, "name"),
		Description:       htmlText(pickString(productRecord, "description")),
		ProductType:       strings.ToUpper(pickString(productRecord, "type")),
		Currency:          firstNonEmpty(pickString(root, "currency"), "CNY"),
		Pricing:           buildUpstreamPricing(data),
		ConfigOptions:     buildUpstreamConfigOptions(data),
	}

	if groups, err := service.ListUpstreamProducts(); err == nil {
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

func buildUpstreamPricing(data map[string]any) []UpstreamCyclePrice {
	priceRows, _ := data["product_pricing"].([]any)
	if len(priceRows) == 0 {
		priceRows, _ = data["product_pricings"].([]any)
	}
	if len(priceRows) == 0 {
		return nil
	}

	priceRecord, _ := priceRows[0].(map[string]any)
	cycles := []struct {
		Code     string
		Label    string
		SetupKey string
	}{
		{Code: "monthly", Label: "月付", SetupKey: "msetupfee"},
		{Code: "quarterly", Label: "季付", SetupKey: "qsetupfee"},
		{Code: "semiannually", Label: "半年付", SetupKey: "ssetupfee"},
		{Code: "annually", Label: "年付", SetupKey: "asetupfee"},
		{Code: "biennially", Label: "两年付", SetupKey: "bsetupfee"},
		{Code: "triennially", Label: "三年付", SetupKey: "tsetupfee"},
		{Code: "onetime", Label: "一次性", SetupKey: "osetupfee"},
	}

	result := make([]UpstreamCyclePrice, 0, len(cycles))
	for _, cycle := range cycles {
		price := pickFloat(priceRecord, cycle.Code)
		if price < 0 {
			continue
		}
		setupFee := pickFloat(priceRecord, cycle.SetupKey)
		if setupFee < 0 {
			setupFee = 0
		}
		result = append(result, UpstreamCyclePrice{
			CycleCode: cycle.Code,
			CycleName: cycle.Label,
			Price:     price,
			SetupFee:  setupFee,
		})
	}
	return result
}

func buildUpstreamConfigOptions(data map[string]any) []UpstreamConfigOption {
	configGroups, _ := data["config_groups"].([]any)
	result := make([]UpstreamConfigOption, 0)

	for _, groupAny := range configGroups {
		groupRecord, ok := groupAny.(map[string]any)
		if !ok {
			continue
		}
		groupDescription := pickString(groupRecord, "description")
		options, _ := groupRecord["options"].([]any)
		for _, optionAny := range options {
			optionRecord, ok := optionAny.(map[string]any)
			if !ok {
				continue
			}
			code, name := splitOptionName(pickString(optionRecord, "option_name"))
			choicesAny, _ := optionRecord["sub"].([]any)
			choices := make([]UpstreamConfigChoice, 0, len(choicesAny))
			for _, choiceAny := range choicesAny {
				choiceRecord, ok := choiceAny.(map[string]any)
				if !ok {
					continue
				}
				value, label := splitChoiceName(pickString(choiceRecord, "option_name"))
				priceDelta := extractChoicePrice(choiceRecord)
				choices = append(choices, UpstreamConfigChoice{
					Value:      value,
					Label:      label,
					PriceDelta: priceDelta,
				})
			}

			defaultValue := ""
			if len(choices) > 0 {
				defaultValue = choices[0].Value
			}
			result = append(result, UpstreamConfigOption{
				Code:         code,
				Name:         name,
				InputType:    mapOptionInputType(pickInt(optionRecord, "option_type")),
				Required:     true,
				DefaultValue: defaultValue,
				Description:  htmlText(groupDescription),
				Choices:      choices,
			})
		}
	}

	return result
}

func splitOptionName(input string) (string, string) {
	parts := strings.SplitN(strings.TrimSpace(input), "|", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	if len(parts) == 1 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[0])
	}
	return "", ""
}

func splitChoiceName(input string) (string, string) {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return "", ""
	}
	mainParts := strings.SplitN(raw, "|", 2)
	if len(mainParts) == 2 {
		value := strings.TrimSpace(mainParts[0])
		label := strings.TrimSpace(mainParts[1])
		labelParts := strings.SplitN(label, "^", 2)
		if len(labelParts) == 2 {
			friendly := strings.TrimSpace(labelParts[1])
			if friendly != "" {
				return value, friendly
			}
		}
		return value, label
	}
	return raw, raw
}

func extractChoicePrice(choiceRecord map[string]any) float64 {
	pricingRows, _ := choiceRecord["pricings"].([]any)
	if len(pricingRows) == 0 {
		return 0
	}
	priceRecord, _ := pricingRows[0].(map[string]any)
	for _, cycle := range []string{"monthly", "quarterly", "semiannually", "annually", "onetime"} {
		value := pickFloat(priceRecord, cycle)
		if value > 0 {
			return value
		}
	}
	return 0
}

func mapOptionInputType(optionType int) string {
	switch optionType {
	case 5, 10, 12, 13, 18:
		return "select"
	default:
		return "select"
	}
}

func htmlText(input string) string {
	value := strings.TrimSpace(input)
	value = strings.ReplaceAll(value, "&lt;br&gt;", "\n")
	value = strings.ReplaceAll(value, "&lt;br/&gt;", "\n")
	value = strings.ReplaceAll(value, "&lt;br /&gt;", "\n")
	value = strings.ReplaceAll(value, "&lt;li&gt;", "- ")
	value = strings.ReplaceAll(value, "&lt;/li&gt;", "\n")
	value = strings.ReplaceAll(value, "&nbsp;", " ")
	value = strings.ReplaceAll(value, "<br>", "\n")
	value = strings.ReplaceAll(value, "<br/>", "\n")
	value = strings.ReplaceAll(value, "<br />", "\n")
	value = strings.ReplaceAll(value, "<li>", "- ")
	value = strings.ReplaceAll(value, "</li>", "\n")
	value = strings.NewReplacer("&lt;", "<", "&gt;", ">", "&amp;", "&", "&quot;", "\"").Replace(value)
	for strings.Contains(value, "  ") {
		value = strings.ReplaceAll(value, "  ", " ")
	}
	return strings.TrimSpace(value)
}

func pickFloat(record map[string]any, keys ...string) float64 {
	for _, key := range keys {
		value, ok := record[key]
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case float64:
			return typed
		case float32:
			return float64(typed)
		case int:
			return float64(typed)
		case int64:
			return float64(typed)
		case string:
			parsed, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
			if err == nil {
				return parsed
			}
		case json.Number:
			parsed, err := typed.Float64()
			if err == nil {
				return parsed
			}
		}
	}
	return 0
}

func (service *Service) getFinanceSessionToken(forceRefresh bool) (string, error) {
	service.financeMu.Lock()
	defer service.financeMu.Unlock()

	if !forceRefresh && strings.TrimSpace(service.financeJWT) != "" {
		return service.financeJWT, nil
	}

	if strings.TrimSpace(service.financeConfig.Username) == "" || strings.TrimSpace(service.financeConfig.Password) == "" {
		return "", fmt.Errorf("未配置上游财务登录账号")
	}

	values := url.Values{}
	values.Set("username", service.financeConfig.Username)
	values.Set("password", service.financeConfig.Password)

	var lastErr error
	for _, baseURL := range service.financeBaseURLCandidates() {
		request, err := http.NewRequest(http.MethodPost, joinURL(baseURL, "/zjmf_api_login"), strings.NewReader(values.Encode()))
		if err != nil {
			lastErr = err
			continue
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

		response, err := service.financeClient.Do(request)
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
			service.financeJWT = token
			service.financeActiveURL = baseURL
			service.financeLastAuthAt = time.Now()
			return token, nil
		}

		lastErr = fmt.Errorf("上游财务登录失败: HTTP %d", response.StatusCode)
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("上游财务登录失败")
	}
	return "", lastErr
}

func (service *Service) financeRequest(method, path string, data map[string]any) (rawResponse, error) {
	token, err := service.getFinanceSessionToken(false)
	if err != nil {
		return rawResponse{}, err
	}

	requestURL := joinURL(service.financeActiveURL, path)
	var bodyString string
	if strings.EqualFold(method, http.MethodGet) && len(data) > 0 {
		values := url.Values{}
		for key, value := range data {
			values.Set(key, fmt.Sprint(value))
		}
		requestURL += "?" + values.Encode()
	} else if len(data) > 0 {
		values := url.Values{}
		for key, value := range data {
			values.Set(key, fmt.Sprint(value))
		}
		bodyString = values.Encode()
	}

	request, err := http.NewRequest(method, requestURL, strings.NewReader(bodyString))
	if err != nil {
		return rawResponse{}, err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	if bodyString != "" {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}

	response, err := service.financeClient.Do(request)
	if err != nil {
		return rawResponse{}, err
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	var parsed any
	_ = json.Unmarshal(body, &parsed)

	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusMethodNotAllowed {
		service.financeMu.Lock()
		service.financeJWT = ""
		service.financeMu.Unlock()
		return service.financeRequest(method, path, data)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return rawResponse{}, fmt.Errorf("上游财务请求失败: HTTP %d", response.StatusCode)
	}

	root, _ := parsed.(map[string]any)
	if status := pickInt(root, "status"); status != 0 && status != 200 {
		return rawResponse{}, fmt.Errorf(firstNonEmpty(pickString(root, "msg", "message"), "上游财务返回业务错误"))
	}

	return rawResponse{
		StatusCode: response.StatusCode,
		Text:       string(body),
		Parsed:     parsed,
	}, nil
}

func (service *Service) financeBaseURLCandidates() []string {
	baseURL := strings.TrimRight(strings.TrimSpace(service.financeConfig.BaseURL), "/")
	if baseURL == "" {
		return nil
	}
	if strings.HasPrefix(baseURL, "https://") {
		return []string{baseURL, "http://" + strings.TrimPrefix(baseURL, "https://")}
	}
	return []string{baseURL}
}

func extractFinanceToken(input any) string {
	record, _ := input.(map[string]any)
	if record == nil {
		return ""
	}
	if token := pickString(record, "jwt", "token", "access_token"); token != "" {
		return token
	}
	if nested, ok := record["data"].(map[string]any); ok {
		return pickString(nested, "jwt", "token", "access_token")
	}
	return ""
}
