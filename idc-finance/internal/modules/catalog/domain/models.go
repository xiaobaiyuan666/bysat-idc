package domain

type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "ACTIVE"
	ProductStatusInactive ProductStatus = "INACTIVE"
)

type PriceOption struct {
	CycleCode string  `json:"cycleCode"`
	CycleName string  `json:"cycleName"`
	Price     float64 `json:"price"`
	SetupFee  float64 `json:"setupFee"`
}

type ConfigOptionChoice struct {
	Value      string  `json:"value"`
	Label      string  `json:"label"`
	PriceDelta float64 `json:"priceDelta"`
}

type ConfigOption struct {
	Code         string               `json:"code"`
	Name         string               `json:"name"`
	InputType    string               `json:"inputType"`
	Required     bool                 `json:"required"`
	DefaultValue string               `json:"defaultValue"`
	Description  string               `json:"description"`
	Choices      []ConfigOptionChoice `json:"choices"`
}

type ResourceTemplate struct {
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

type AutomationConfig struct {
	ProviderAccountID int64  `json:"providerAccountId"`
	Channel        string `json:"channel"`
	ModuleType     string `json:"moduleType"`
	ProvisionStage string `json:"provisionStage"`
	AutoProvision  bool   `json:"autoProvision"`
	ServerGroup    string `json:"serverGroup"`
	ProviderNode   string `json:"providerNode"`
}

type UpstreamMapping struct {
	ProviderAccountID int64  `json:"providerAccountId"`
	ProviderType      string `json:"providerType"`
	SourceName        string `json:"sourceName"`
	RemoteProductCode string `json:"remoteProductCode"`
	RemoteProductName string `json:"remoteProductName"`
	PricePolicy       string `json:"pricePolicy"`
	AutoSyncPricing   bool   `json:"autoSyncPricing"`
	AutoSyncConfig    bool   `json:"autoSyncConfig"`
	AutoSyncTemplate  bool   `json:"autoSyncTemplate"`
	SyncStatus        string `json:"syncStatus"`
	SyncMessage       string `json:"syncMessage"`
	LastSyncedAt      string `json:"lastSyncedAt"`
}

type Product struct {
	ID               int64            `json:"id"`
	ProductNo        string           `json:"productNo"`
	GroupName        string           `json:"groupName"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	ProductType      string           `json:"productType"`
	Status           ProductStatus    `json:"status"`
	Pricing          []PriceOption    `json:"pricing"`
	ConfigOptions    []ConfigOption   `json:"configOptions"`
	ResourceTemplate ResourceTemplate `json:"resourceTemplate"`
	AutomationConfig AutomationConfig `json:"automationConfig"`
	UpstreamMapping  UpstreamMapping  `json:"upstreamMapping"`
}
