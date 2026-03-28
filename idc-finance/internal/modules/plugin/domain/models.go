package domain

type PluginCode string

const (
	PluginCodeZhimaKYC PluginCode = "zhima_kyc"
	PluginCodeSMSBao   PluginCode = "smsbao_sms"
	PluginCodeEPay     PluginCode = "epay"
)

type PluginField struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Required    bool   `json:"required"`
	Placeholder string `json:"placeholder,omitempty"`
	Secret      bool   `json:"secret"`
}

type PluginDefinition struct {
	Code        string        `json:"code"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Fields      []PluginField `json:"fields"`
}

type PluginConfig struct {
	Code      string         `json:"code"`
	Name      string         `json:"name"`
	Enabled   bool           `json:"enabled"`
	Config    map[string]any `json:"config"`
	UpdatedBy string         `json:"updatedBy,omitempty"`
	UpdatedAt string         `json:"updatedAt,omitempty"`
}

type KYCRecord struct {
	ID         int64          `json:"id"`
	PluginCode string         `json:"pluginCode"`
	CustomerID int64          `json:"customerId"`
	RealName   string         `json:"realName"`
	IDCardNo   string         `json:"idCardNo"`
	Mobile     string         `json:"mobile"`
	Status     string         `json:"status"`
	Message    string         `json:"message"`
	Raw        map[string]any `json:"raw,omitempty"`
	CreatedAt  string         `json:"createdAt"`
}

type SMSRecord struct {
	ID           int64          `json:"id"`
	PluginCode   string         `json:"pluginCode"`
	CustomerID   int64          `json:"customerId"`
	Mobile       string         `json:"mobile"`
	TemplateCode string         `json:"templateCode"`
	Content      string         `json:"content"`
	Status       string         `json:"status"`
	Message      string         `json:"message"`
	BizNo        string         `json:"bizNo"`
	Raw          map[string]any `json:"raw,omitempty"`
	CreatedAt    string         `json:"createdAt"`
}

type PaymentOrder struct {
	ID         int64          `json:"id"`
	PluginCode string         `json:"pluginCode"`
	OrderNo    string         `json:"orderNo"`
	CustomerID int64          `json:"customerId"`
	Amount     float64        `json:"amount"`
	Currency   string         `json:"currency"`
	Subject    string         `json:"subject"`
	Status     string         `json:"status"`
	PayURL     string         `json:"payUrl"`
	NotifyURL  string         `json:"notifyUrl"`
	ReturnURL  string         `json:"returnUrl"`
	TradeNo    string         `json:"tradeNo"`
	Raw        map[string]any `json:"raw,omitempty"`
	PaidAt     string         `json:"paidAt,omitempty"`
	CreatedAt  string         `json:"createdAt"`
	UpdatedAt  string         `json:"updatedAt"`
}
