package dto

type SavePluginConfigRequest struct {
	Enabled bool           `json:"enabled"`
	Config  map[string]any `json:"config"`
}

type KYCVerifyRequest struct {
	RealName string `json:"realName" binding:"required"`
	IDCardNo string `json:"idCardNo" binding:"required"`
	Mobile   string `json:"mobile"`
}

type SMSDispatchRequest struct {
	Mobile       string `json:"mobile" binding:"required"`
	TemplateCode string `json:"templateCode"`
	Content      string `json:"content" binding:"required"`
}

type CreatePaymentOrderRequest struct {
	OrderNo   string  `json:"orderNo" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
	Currency  string  `json:"currency"`
	Subject   string  `json:"subject"`
	NotifyURL string  `json:"notifyUrl"`
	ReturnURL string  `json:"returnUrl"`
}

type PaymentNotifyRequest struct {
	OrderNo string `json:"orderNo" binding:"required"`
	TradeNo string `json:"tradeNo"`
	Status  string `json:"status" binding:"required"`
	Sign    string `json:"sign"`
}
