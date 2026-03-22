package service

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	automationService "idc-finance/internal/modules/automation/service"
	providerDomain "idc-finance/internal/modules/provider/domain"
	providerRepository "idc-finance/internal/modules/provider/repository"
	"idc-finance/internal/platform/audit"
	"idc-finance/internal/platform/config"
)

type Service struct {
	config            config.MofangCloudConfig
	financeConfig     config.FinanceUpstreamConfig
	accounts          providerRepository.Repository
	db                *sql.DB
	audit             *audit.Service
	tasks             *automationService.Service
	client            *http.Client
	financeClient     *http.Client
	mu                sync.Mutex
	sessionJWT        string
	activeURL         string
	lastAuthAt        time.Time
	financeMu         sync.Mutex
	financeJWT        string
	financeActiveURL  string
	financeLastAuthAt time.Time
	sessionStoreMu    sync.Mutex
	sessionStore      map[string]providerSessionState
}

type HealthResponse struct {
	Enabled    bool      `json:"enabled"`
	Connected  bool      `json:"connected"`
	BaseURL    string    `json:"baseUrl"`
	ActiveURL  string    `json:"activeUrl,omitempty"`
	AuthMode   string    `json:"authMode"`
	Message    string    `json:"message"`
	LastAuthAt time.Time `json:"lastAuthAt,omitempty"`
}

type InstanceSummary struct {
	RemoteID   string         `json:"remoteId"`
	Name       string         `json:"name"`
	Status     string         `json:"status"`
	Region     string         `json:"region,omitempty"`
	IPAddress  string         `json:"ipAddress,omitempty"`
	ConsoleURL string         `json:"consoleUrl,omitempty"`
	ExpiresAt  string         `json:"expiresAt,omitempty"`
	Raw        map[string]any `json:"raw,omitempty"`
}

type ActionRequest struct {
	ImageName string `json:"imageName"`
	Password  string `json:"password"`
}

type ActionResponse struct {
	OK       bool           `json:"ok"`
	Action   string         `json:"action"`
	RemoteID string         `json:"remoteId"`
	Status   string         `json:"status"`
	Message  string         `json:"message"`
	Response map[string]any `json:"response,omitempty"`
}

type PullSyncOptions struct {
	ProviderAccountID int64    `json:"providerAccountId"`
	Limit            int      `json:"limit"`
	IncludeResources bool     `json:"includeResources"`
	RemoteIDs        []string `json:"remoteIds"`
}

type PullSyncItem struct {
	RemoteID     string `json:"remoteId"`
	ServiceID    int64  `json:"serviceId,omitempty"`
	ServiceNo    string `json:"serviceNo,omitempty"`
	CustomerID   int64  `json:"customerId,omitempty"`
	CustomerName string `json:"customerName,omitempty"`
	Operation    string `json:"operation"`
	Status       string `json:"status,omitempty"`
	Message      string `json:"message"`
}

type PullSyncSummary struct {
	RemoteCount             int    `json:"remoteCount"`
	ProcessedCount          int    `json:"processedCount"`
	CreatedCustomers        int    `json:"createdCustomers"`
	CreatedServices         int    `json:"createdServices"`
	UpdatedServices         int    `json:"updatedServices"`
	FailedServices          int    `json:"failedServices"`
	SyncedNetworkInterfaces int    `json:"syncedNetworkInterfaces"`
	SyncedIps               int    `json:"syncedIps"`
	SyncedDisks             int    `json:"syncedDisks"`
	SyncedSnapshots         int    `json:"syncedSnapshots"`
	SyncedBackups           int    `json:"syncedBackups"`
	SyncedVpcs              int    `json:"syncedVpcs"`
	SyncedSecurityGroups    int    `json:"syncedSecurityGroups"`
	SyncedSecurityRules     int    `json:"syncedSecurityRules"`
	Message                 string `json:"message"`
}

type PullSyncResult struct {
	Summary PullSyncSummary `json:"summary"`
	Items   []PullSyncItem  `json:"items"`
}

type SyncLogItem struct {
	ID           int64  `json:"id"`
	ProviderType string `json:"providerType"`
	Action       string `json:"action"`
	ResourceType string `json:"resourceType"`
	ResourceID   string `json:"resourceId"`
	ServiceID    int64  `json:"serviceId,omitempty"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	CreatedAt    string `json:"createdAt"`
}

type NetworkInterface struct {
	ID                  int64  `json:"id"`
	ProviderInterfaceID string `json:"providerInterfaceId,omitempty"`
	Name                string `json:"name"`
	MACAddress          string `json:"macAddress"`
	Bridge              string `json:"bridge"`
	NetworkName         string `json:"networkName"`
	InterfaceName       string `json:"interfaceName"`
	NICModel            string `json:"nicModel"`
	InboundMbps         int    `json:"inboundMbps"`
	OutboundMbps        int    `json:"outboundMbps"`
	VPCName             string `json:"vpcName"`
	Status              string `json:"status"`
}

type IPAddress struct {
	ID                 int64  `json:"id"`
	NetworkInterfaceID int64  `json:"networkInterfaceId,omitempty"`
	ProviderIPID       string `json:"providerIpId,omitempty"`
	Address            string `json:"address"`
	Version            string `json:"version"`
	IPRole             string `json:"ipRole"`
	SubnetMask         string `json:"subnetMask"`
	Gateway            string `json:"gateway"`
	BandwidthMbps      int    `json:"bandwidthMbps"`
	IsPrimary          bool   `json:"isPrimary"`
	Status             string `json:"status"`
}

type Disk struct {
	ID             int64  `json:"id"`
	ProviderDiskID string `json:"providerDiskId,omitempty"`
	Name           string `json:"name"`
	DiskType       string `json:"diskType"`
	SizeGB         int    `json:"sizeGb"`
	DeviceName     string `json:"deviceName"`
	DriverName     string `json:"driverName"`
	FSType         string `json:"fsType"`
	MountPoint     string `json:"mountPoint"`
	IsSystem       bool   `json:"isSystem"`
	Status         string `json:"status"`
}

type Snapshot struct {
	ID                 int64  `json:"id"`
	DiskID             int64  `json:"diskId,omitempty"`
	ProviderSnapshotID string `json:"providerSnapshotId,omitempty"`
	Name               string `json:"name"`
	SizeGB             int    `json:"sizeGb"`
	Status             string `json:"status"`
}

type Backup struct {
	ID               int64  `json:"id"`
	ProviderBackupID string `json:"providerBackupId,omitempty"`
	Name             string `json:"name"`
	SizeGB           int    `json:"sizeGb"`
	Status           string `json:"status"`
	ExpiresAt        string `json:"expiresAt,omitempty"`
}

type VPCNetwork struct {
	ID            int64  `json:"id"`
	ProviderVPCID string `json:"providerVpcId,omitempty"`
	Name          string `json:"name"`
	CIDR          string `json:"cidr"`
	Gateway       string `json:"gateway"`
	InterfaceName string `json:"interfaceName"`
	Status        string `json:"status"`
}

type SecurityGroupRule struct {
	ID          int64  `json:"id"`
	Direction   string `json:"direction"`
	Protocol    string `json:"protocol"`
	PortRange   string `json:"portRange"`
	SourceCIDR  string `json:"sourceCidr"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

type SecurityGroup struct {
	ID                      int64               `json:"id"`
	ProviderSecurityGroupID string              `json:"providerSecurityGroupId,omitempty"`
	Name                    string              `json:"name"`
	Status                  string              `json:"status"`
	Rules                   []SecurityGroupRule `json:"rules"`
}

type ServiceResourcesResponse struct {
	ServiceID          int64              `json:"serviceId"`
	ServiceNo          string             `json:"serviceNo"`
	ProviderType       string             `json:"providerType"`
	ProviderResourceID string             `json:"providerResourceId"`
	Status             string             `json:"status"`
	RegionName         string             `json:"regionName"`
	IPAddress          string             `json:"ipAddress"`
	SyncStatus         string             `json:"syncStatus"`
	SyncMessage        string             `json:"syncMessage"`
	LastSyncAt         string             `json:"lastSyncAt"`
	NetworkInterfaces  []NetworkInterface `json:"networkInterfaces"`
	IPAddresses        []IPAddress        `json:"ipAddresses"`
	Disks              []Disk             `json:"disks"`
	Snapshots          []Snapshot         `json:"snapshots"`
	Backups            []Backup           `json:"backups"`
	VPCNetworks        []VPCNetwork       `json:"vpcNetworks"`
	SecurityGroups     []SecurityGroup    `json:"securityGroups"`
	SyncLogs           []SyncLogItem      `json:"syncLogs"`
}

type rawResponse struct {
	StatusCode int
	Text       string
	Parsed     any
}

func New(
	cfg config.MofangCloudConfig,
	financeCfg config.FinanceUpstreamConfig,
	accountRepo providerRepository.Repository,
	db *sql.DB,
	auditService *audit.Service,
	taskService *automationService.Service,
) *Service {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.InsecureSkipVerify}, //nolint:gosec
	}
	financeTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: financeCfg.InsecureSkipVerify}, //nolint:gosec
	}

	instance := &Service{
		config:        cfg,
		financeConfig: financeCfg,
		accounts:      accountRepo,
		db:            db,
		audit:         auditService,
		tasks:         taskService,
		client: &http.Client{
			Timeout:   20 * time.Second,
			Transport: transport,
		},
		financeClient: &http.Client{
			Timeout:   20 * time.Second,
			Transport: financeTransport,
		},
		sessionStore: make(map[string]providerSessionState),
	}
	instance.ensureDefaultAccounts()
	return instance
}

type AccountUpsertRequest struct {
	ProviderType       string `json:"providerType"`
	Name               string `json:"name"`
	BaseURL            string `json:"baseUrl"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	SourceName         string `json:"sourceName"`
	ContactWay         string `json:"contactWay"`
	Description        string `json:"description"`
	AccountMode        string `json:"accountMode"`
	Lang               string `json:"lang"`
	ListPath           string `json:"listPath"`
	DetailPath         string `json:"detailPath"`
	InsecureSkipVerify bool   `json:"insecureSkipVerify"`
	AutoUpdate         bool   `json:"autoUpdate"`
	Status             string `json:"status"`
	ExtraConfig        string `json:"extraConfig"`
}

type providerSessionState struct {
	Token      string
	ActiveURL  string
	LastAuthAt time.Time
}

type resolvedProviderAccount struct {
	AccountID           int64
	ProviderType        string
	Name                string
	BaseURL             string
	Username            string
	Password            string
	SourceName          string
	Lang                string
	ListPath            string
	DetailPath          string
	InsecureSkipVerify  bool
	AuthMode            string
	Status              providerDomain.AccountStatus
}

func (service *Service) legacyCheckHealth() HealthResponse {
	if strings.TrimSpace(service.config.BaseURL) == "" {
		return HealthResponse{
			Enabled:   false,
			Connected: false,
			AuthMode:  "disabled",
			Message:   "未配置魔方云地址",
		}
	}

	_, err := service.getSessionToken(false)
	if err != nil {
		return HealthResponse{
			Enabled:    true,
			Connected:  false,
			BaseURL:    service.config.BaseURL,
			ActiveURL:  service.activeURL,
			AuthMode:   "jwt",
			Message:    err.Error(),
			LastAuthAt: service.lastAuthAt,
		}
	}

	return HealthResponse{
		Enabled:    true,
		Connected:  true,
		BaseURL:    service.config.BaseURL,
		ActiveURL:  service.activeURL,
		AuthMode:   "jwt",
		Message:    "魔方云连接正常",
		LastAuthAt: service.lastAuthAt,
	}
}

func (service *Service) legacyListInstances() ([]InstanceSummary, error) {
	if strings.TrimSpace(service.config.BaseURL) == "" {
		return nil, fmt.Errorf("未配置魔方云地址")
	}

	page := 1
	totalPages := 1
	result := make([]InstanceSummary, 0)
	for page <= totalPages {
		raw, err := service.request(http.MethodGet, service.buildListPath(page, 100), nil)
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
		page += 1
	}

	return result, nil
}

func (service *Service) legacyGetInstanceDetail(remoteID string) (InstanceSummary, error) {
	path := strings.ReplaceAll(service.config.DetailPath, ":id", remoteID)
	raw, err := service.request(http.MethodGet, path, nil)
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

func (service *Service) legacyExecuteAction(remoteID, action string, request ActionRequest) (ActionResponse, error) {
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

	path := strings.ReplaceAll(spec.Path, ":id", remoteID)
	raw, err := service.request(spec.Method, path, spec.Body)
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

func (service *Service) buildListPath(page, limit int) string {
	basePath := strings.TrimSpace(service.config.ListPath)
	if basePath == "" {
		basePath = "/v1/clouds"
	}

	separator := "?"
	if strings.Contains(basePath, "?") {
		separator = "&"
	}
	return fmt.Sprintf("%s%spage=%d&limit=%d", basePath, separator, page, limit)
}

func (service *Service) getSessionToken(forceRefresh bool) (string, error) {
	service.mu.Lock()
	defer service.mu.Unlock()

	if !forceRefresh && strings.TrimSpace(service.sessionJWT) != "" {
		return service.sessionJWT, nil
	}

	if strings.TrimSpace(service.config.Username) == "" || strings.TrimSpace(service.config.Password) == "" {
		return "", fmt.Errorf("未配置魔方云登录账号")
	}

	values := url.Values{}
	values.Set("username", service.config.Username)
	values.Set("password", service.config.Password)
	values.Set("customfield[google_code]", "")

	var lastErr error
	for _, baseURL := range service.baseURLCandidates() {
		request, err := http.NewRequest(http.MethodPost, joinURL(baseURL, "/v1/login"), strings.NewReader(values.Encode()))
		if err != nil {
			lastErr = err
			continue
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		request.Header.Set("think-lang", service.config.Lang)

		response, err := service.client.Do(request)
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
			service.sessionJWT = token
			service.activeURL = baseURL
			service.lastAuthAt = time.Now()
			return token, nil
		}

		lastErr = fmt.Errorf("魔方云登录失败: HTTP %d", response.StatusCode)
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("魔方云登录失败")
	}
	return "", lastErr
}

func (service *Service) request(method, path string, payload map[string]any) (rawResponse, error) {
	token, err := service.getSessionToken(false)
	if err != nil {
		return rawResponse{}, err
	}

	bodyString := ""
	requestURL := joinURL(service.activeURL, path)
	if strings.EqualFold(method, http.MethodGet) && payload != nil {
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
	} else if payload != nil {
		values := url.Values{}
		for key, value := range payload {
			if strings.TrimSpace(fmt.Sprint(value)) == "" {
				continue
			}
			values.Set(key, fmt.Sprint(value))
		}
		bodyString = values.Encode()
	}

	request, err := http.NewRequest(method, requestURL, bytes.NewBufferString(bodyString))
	if err != nil {
		return rawResponse{}, err
	}
	request.Header.Set("think-lang", service.config.Lang)
	request.Header.Set("Authorization", "JWT "+token)
	request.Header.Set("access-token", token)
	if bodyString != "" {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}

	response, err := service.client.Do(request)
	if err != nil {
		return rawResponse{}, err
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	var parsed any
	_ = json.Unmarshal(body, &parsed)

	if response.StatusCode == http.StatusUnauthorized {
		service.mu.Lock()
		service.sessionJWT = ""
		service.mu.Unlock()
		return service.request(method, path, payload)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		detail := strings.TrimSpace(pickString(extractRecord(parsed), "error", "msg", "message"))
		if detail == "" {
			detail = strings.TrimSpace(string(body))
		}
		if detail != "" {
			return rawResponse{}, fmt.Errorf("魔方云请求失败: HTTP %d - %s", response.StatusCode, detail)
		}
		return rawResponse{}, fmt.Errorf("魔方云请求失败: HTTP %d", response.StatusCode)
	}

	return rawResponse{
		StatusCode: response.StatusCode,
		Text:       string(body),
		Parsed:     parsed,
	}, nil
}

func (service *Service) baseURLCandidates() []string {
	baseURL := strings.TrimRight(strings.TrimSpace(service.config.BaseURL), "/")
	if baseURL == "" {
		return nil
	}
	if strings.HasPrefix(baseURL, "https://") {
		return []string{baseURL, "http://" + strings.TrimPrefix(baseURL, "https://")}
	}
	return []string{baseURL}
}

func joinURL(baseURL, path string) string {
	baseURL = strings.TrimRight(baseURL, "/")
	if strings.HasPrefix(path, "/") {
		return baseURL + path
	}
	return baseURL + "/" + path
}

type listMeta struct {
	totalPages int
}

func extractMeta(input any) listMeta {
	record := extractRecord(input)
	if record == nil {
		return listMeta{}
	}
	meta := map[string]any{}
	if nested, ok := record["meta"].(map[string]any); ok {
		meta = nested
	}
	totalPages := pickInt(meta, "total_page", "totalPages")
	if totalPages == 0 {
		totalPages = pickInt(record, "total_page", "totalPages")
	}
	return listMeta{totalPages: totalPages}
}

func extractToken(input any) string {
	record := extractRecord(input)
	if record == nil {
		if value, ok := input.(string); ok {
			return strings.TrimSpace(value)
		}
		return ""
	}

	if token := pickString(record, "token", "access_token"); token != "" {
		return token
	}

	if data, ok := record["data"].(map[string]any); ok {
		return pickString(data, "token", "access_token")
	}

	return ""
}

func extractRecord(input any) map[string]any {
	switch value := input.(type) {
	case map[string]any:
		if nested, ok := value["data"].(map[string]any); ok {
			return nested
		}
		if nested, ok := value["result"].(map[string]any); ok {
			return nested
		}
		return value
	default:
		return nil
	}
}

func extractRecords(input any) []map[string]any {
	switch value := input.(type) {
	case map[string]any:
		if items, ok := value["data"].([]any); ok {
			return normalizeRecordArray(items)
		}
		if items, ok := value["result"].([]any); ok {
			return normalizeRecordArray(items)
		}
	case []any:
		return normalizeRecordArray(value)
	}

	record := extractRecord(input)
	if record == nil {
		return nil
	}

	for _, key := range []string{"list", "items", "rows", "clouds", "servers", "instances", "records"} {
		if items, ok := record[key].([]any); ok {
			return normalizeRecordArray(items)
		}
	}

	return []map[string]any{record}
}

func normalizeRecordArray(items []any) []map[string]any {
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if record, ok := item.(map[string]any); ok {
			result = append(result, record)
		}
	}
	return result
}

func mapInstance(record map[string]any) InstanceSummary {
	return InstanceSummary{
		RemoteID:   pickString(record, "id", "cloudid", "cloud_id", "instance_id", "hostid", "uuid"),
		Name:       firstNonEmpty(pickString(record, "name", "remarks", "remark", "title", "instance_name", "hostname"), "未命名实例"),
		Status:     normalizeStatus(pickString(record, "status", "instance_status", "cloud_status", "power_status", "state")),
		Region:     pickRegion(record),
		IPAddress:  pickString(record, "ip", "ip_address", "main_ip", "mainip", "public_ip", "primary_ip"),
		ConsoleURL: pickString(record, "url", "vnc_url", "ws_url", "remote_url"),
		ExpiresAt:  pickString(record, "expire_time", "expired_at", "end_time", "renew_time", "due_date"),
		Raw:        record,
	}
}

func pickRegion(record map[string]any) string {
	for _, key := range []string{"region_name", "region", "area_name", "area", "zone_name", "node_name"} {
		if value := pickString(record, key); value != "" {
			return value
		}
	}
	if area, ok := record["area"].(map[string]any); ok {
		return pickString(area, "name", "short_name")
	}
	return ""
}

func pickString(record map[string]any, keys ...string) string {
	for _, key := range keys {
		value, ok := record[key]
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case string:
			if strings.TrimSpace(typed) != "" {
				return strings.TrimSpace(typed)
			}
		case float64:
			return fmt.Sprintf("%.0f", typed)
		case int:
			return fmt.Sprintf("%d", typed)
		case int64:
			return strconv.FormatInt(typed, 10)
		case json.Number:
			return typed.String()
		}
	}
	return ""
}

func pickInt(record map[string]any, keys ...string) int {
	for _, key := range keys {
		value, ok := record[key]
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case int:
			return typed
		case int64:
			return int(typed)
		case float64:
			return int(typed)
		case string:
			parsed, err := strconv.Atoi(strings.TrimSpace(typed))
			if err == nil {
				return parsed
			}
		case json.Number:
			parsed, err := typed.Int64()
			if err == nil {
				return int(parsed)
			}
		}
	}
	return 0
}

func normalizeStatus(input string) string {
	value := strings.ToUpper(strings.TrimSpace(input))
	switch value {
	case "", "RUNNING", "ON", "NORMAL", "ACTIVE":
		return "ACTIVE"
	case "SUSPENDED", "SUSPEND", "OFF", "SHUTOFF", "STOPPED", "PAUSED":
		return "SUSPENDED"
	case "PROVISIONING", "BUILDING", "CREATING", "INSTALLING", "REINSTALLING", "PENDING":
		return "PROVISIONING"
	case "DELETED", "DESTROYED", "TERMINATED", "CANCELLED":
		return "TERMINATED"
	case "FAILED", "ERROR":
		return "FAILED"
	default:
		return value
	}
}

func firstNonEmpty(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
