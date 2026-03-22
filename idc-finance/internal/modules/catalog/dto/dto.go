package dto

import (
	"idc-finance/internal/modules/catalog/domain"
	"idc-finance/internal/platform/audit"
)

type CreateProductRequest struct {
	GroupName        string                `json:"groupName" binding:"required"`
	Name             string                `json:"name" binding:"required"`
	Description      string                `json:"description"`
	ProductType      string                `json:"productType" binding:"required"`
	Status           string                `json:"status"`
	Pricing          []PriceOptionInput    `json:"pricing"`
	ConfigOptions    []ConfigOptionInput   `json:"configOptions"`
	ResourceTemplate ResourceTemplateInput `json:"resourceTemplate"`
	AutomationConfig AutomationConfigInput `json:"automationConfig"`
	UpstreamMapping  UpstreamMappingInput  `json:"upstreamMapping"`
}

type UpdateProductRequest struct {
	GroupName        string                `json:"groupName" binding:"required"`
	Name             string                `json:"name" binding:"required"`
	Description      string                `json:"description"`
	ProductType      string                `json:"productType" binding:"required"`
	Status           string                `json:"status"`
	Pricing          []PriceOptionInput    `json:"pricing"`
	ConfigOptions    []ConfigOptionInput   `json:"configOptions"`
	ResourceTemplate ResourceTemplateInput `json:"resourceTemplate"`
	AutomationConfig AutomationConfigInput `json:"automationConfig"`
	UpstreamMapping  UpstreamMappingInput  `json:"upstreamMapping"`
}

type AutomationConfigInput struct {
	ProviderAccountID int64  `json:"providerAccountId"`
	Channel           string `json:"channel"`
	ModuleType        string `json:"moduleType"`
	ProvisionStage    string `json:"provisionStage"`
	AutoProvision     bool   `json:"autoProvision"`
	ServerGroup       string `json:"serverGroup"`
	ProviderNode      string `json:"providerNode"`
}

type UpstreamMappingInput struct {
	ProviderAccountID int64  `json:"providerAccountId"`
	ProviderType      string `json:"providerType"`
	SourceName        string `json:"sourceName"`
	RemoteProductCode string `json:"remoteProductCode"`
	RemoteProductName string `json:"remoteProductName"`
	PricePolicy       string `json:"pricePolicy"`
	AutoSyncPricing   bool   `json:"autoSyncPricing"`
	AutoSyncConfig    bool   `json:"autoSyncConfig"`
	AutoSyncTemplate  bool   `json:"autoSyncTemplate"`
}

type SyncProductUpstreamRequest struct {
	ProviderAccountID int64  `json:"providerAccountId"`
	ProviderType      string `json:"providerType"`
	SourceName        string `json:"sourceName"`
	RemoteProductCode string `json:"remoteProductCode"`
	RemoteProductName string `json:"remoteProductName"`
	PricePolicy       string `json:"pricePolicy"`
	AutoSyncPricing   bool   `json:"autoSyncPricing"`
	AutoSyncConfig    bool   `json:"autoSyncConfig"`
	AutoSyncTemplate  bool   `json:"autoSyncTemplate"`
}

type ImportUpstreamProductsRequest struct {
	ProviderAccountID  int64    `json:"providerAccountId"`
	RemoteProductCodes []string `json:"remoteProductCodes"`
	ImportAll          bool     `json:"importAll"`
	AutoSyncPricing    bool     `json:"autoSyncPricing"`
	AutoSyncConfig     bool     `json:"autoSyncConfig"`
	AutoSyncTemplate   bool     `json:"autoSyncTemplate"`
}

type ImportUpstreamProductItem struct {
	RemoteProductCode string `json:"remoteProductCode"`
	RemoteProductName string `json:"remoteProductName"`
	GroupName         string `json:"groupName"`
	ProductID         int64  `json:"productId"`
	ProductNo         string `json:"productNo"`
	Operation         string `json:"operation"`
	Status            string `json:"status"`
	Message           string `json:"message"`
}

type ImportUpstreamProductsResponse struct {
	ProviderAccountID int64                       `json:"providerAccountId"`
	ImportedCount     int                         `json:"importedCount"`
	UpdatedCount      int                         `json:"updatedCount"`
	FailedCount       int                         `json:"failedCount"`
	Items             []ImportUpstreamProductItem `json:"items"`
	Total             int                         `json:"total"`
	Created           int                         `json:"created"`
	Updated           int                         `json:"updated"`
	Skipped           int                         `json:"skipped"`
	Failed            int                         `json:"failed"`
	Message           string                      `json:"message"`
}

type PriceOptionInput struct {
	CycleCode string  `json:"cycleCode"`
	CycleName string  `json:"cycleName"`
	Price     float64 `json:"price"`
	SetupFee  float64 `json:"setupFee"`
}

type ConfigOptionChoiceInput struct {
	Value      string  `json:"value"`
	Label      string  `json:"label"`
	PriceDelta float64 `json:"priceDelta"`
}

type ConfigOptionInput struct {
	Code         string                    `json:"code"`
	Name         string                    `json:"name"`
	InputType    string                    `json:"inputType"`
	Required     bool                      `json:"required"`
	DefaultValue string                    `json:"defaultValue"`
	Description  string                    `json:"description"`
	Choices      []ConfigOptionChoiceInput `json:"choices"`
}

type ResourceTemplateInput struct {
	RegionName      string `json:"regionName"`
	ZoneName        string `json:"zoneName"`
	OperatingSystem string `json:"operatingSystem"`
	LoginUsername   string `json:"loginUsername"`
	SecurityGroup   string `json:"securityGroup"`
	CPUCores        int    `json:"cpuCores"`
	MemoryGB        int    `json:"memoryGB"`
	SystemDiskGB    int    `json:"systemDiskGB"`
	DataDiskGB      int    `json:"dataDiskGB"`
	BandwidthMbps   int    `json:"bandwidthMbps"`
	PublicIPCount   int    `json:"publicIpCount"`
}

type ProductListResponse struct {
	Items []domain.Product `json:"items"`
	Total int              `json:"total"`
}

type ProductDetailResponse struct {
	Product   domain.Product `json:"product"`
	AuditLogs []audit.Entry  `json:"auditLogs"`
}
