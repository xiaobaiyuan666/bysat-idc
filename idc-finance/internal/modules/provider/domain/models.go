package domain

type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "ACTIVE"
	AccountStatusDisabled AccountStatus = "DISABLED"
)

type ProviderType string

const (
	ProviderTypeMofangCloud ProviderType = "MOFANG_CLOUD"
	ProviderTypeZjmfAPI     ProviderType = "ZJMF_API"
	ProviderTypeWHMCS       ProviderType = "WHMCS"
	ProviderTypeManual      ProviderType = "MANUAL"
	ProviderTypeResource    ProviderType = "RESOURCE"
)

type Account struct {
	ID                 int64         `json:"id"`
	ProviderType       string        `json:"providerType"`
	Name               string        `json:"name"`
	BaseURL            string        `json:"baseUrl"`
	Username           string        `json:"username"`
	Password           string        `json:"password"`
	SourceName         string        `json:"sourceName"`
	ContactWay         string        `json:"contactWay"`
	Description        string        `json:"description"`
	AccountMode        string        `json:"accountMode"`
	Lang               string        `json:"lang"`
	ListPath           string        `json:"listPath"`
	DetailPath         string        `json:"detailPath"`
	InsecureSkipVerify bool          `json:"insecureSkipVerify"`
	AutoUpdate         bool          `json:"autoUpdate"`
	ProductCount       int           `json:"productCount"`
	Status             AccountStatus `json:"status"`
	ExtraConfig        string        `json:"extraConfig"`
	CreatedAt          string        `json:"createdAt"`
	UpdatedAt          string        `json:"updatedAt"`
}
