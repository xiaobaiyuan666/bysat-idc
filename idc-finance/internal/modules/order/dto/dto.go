package dto

import (
	automationDomain "idc-finance/internal/modules/automation/domain"
	"idc-finance/internal/modules/order/domain"
	"idc-finance/internal/platform/audit"
)

type CheckoutRequest struct {
	ProductID       int64                   `json:"productId" binding:"required"`
	CycleCode       string                  `json:"cycleCode" binding:"required"`
	SelectedOptions []CheckoutOptionRequest `json:"selectedOptions"`
}

type CheckoutOptionRequest struct {
	Code  string `json:"code"`
	Value string `json:"value"`
}

type ReceivePaymentRequest struct {
	Channel string `json:"channel"`
	TradeNo string `json:"tradeNo"`
}

type RefundRequest struct {
	Reason string `json:"reason"`
}

type UpdatePendingOrderRequest struct {
	ProductName  string   `json:"productName"`
	BillingCycle string   `json:"billingCycle"`
	Amount       *float64 `json:"amount"`
	DueAt        string   `json:"dueAt"`
	Status       string   `json:"status"`
	Reason       string   `json:"reason"`
}

type UpdateUnpaidInvoiceRequest struct {
	ProductName  string   `json:"productName"`
	BillingCycle string   `json:"billingCycle"`
	Amount       *float64 `json:"amount"`
	DueAt        string   `json:"dueAt"`
	Status       string   `json:"status"`
	Reason       string   `json:"reason"`
}

type UpdateServiceRequest struct {
	ProviderType       string `json:"providerType"`
	ProviderAccountID  *int64 `json:"providerAccountId"`
	ProviderResourceID string `json:"providerResourceId"`
	RegionName         string `json:"regionName"`
	IPAddress          string `json:"ipAddress"`
	NextDueAt          string `json:"nextDueAt"`
	Status             string `json:"status"`
	SyncStatus         string `json:"syncStatus"`
	SyncMessage        string `json:"syncMessage"`
	Reason             string `json:"reason"`
}

type CreateServiceChangeOrderRequest struct {
	ActionName   string         `json:"actionName"`
	Title        string         `json:"title"`
	BillingCycle string         `json:"billingCycle"`
	Amount       float64        `json:"amount"`
	Reason       string         `json:"reason"`
	Payload      map[string]any `json:"payload"`
}

type ServiceActionRequest struct {
	ImageName string `json:"imageName"`
	Password  string `json:"password"`
}

type AdjustWalletRequest struct {
	Target    string  `json:"target"`
	Operation string  `json:"operation"`
	Amount    float64 `json:"amount"`
	Summary   string  `json:"summary"`
	Remark    string  `json:"remark"`
}

type OrderListResponse struct {
	Items []domain.Order `json:"items"`
	Total int            `json:"total"`
}

type InvoiceListResponse struct {
	Items []domain.Invoice `json:"items"`
	Total int              `json:"total"`
}

type ServiceListResponse struct {
	Items []domain.ServiceRecord `json:"items"`
	Total int                    `json:"total"`
}

type AccountTransactionListResponse struct {
	Items []domain.AccountTransaction `json:"items"`
	Total int                         `json:"total"`
}

type PaymentListResponse struct {
	Items []domain.PaymentRecord `json:"items"`
	Total int                    `json:"total"`
}

type RefundListResponse struct {
	Items []domain.RefundRecord `json:"items"`
	Total int                   `json:"total"`
}

type ServiceChangeOrderListResponse struct {
	Items []ServiceChangeOrderRecord `json:"items"`
	Total int                        `json:"total"`
}

type CustomerWalletResponse struct {
	Wallet       domain.CustomerWallet       `json:"wallet"`
	Transactions []domain.AccountTransaction `json:"transactions"`
}

type AdjustWalletResponse struct {
	Wallet      domain.CustomerWallet     `json:"wallet"`
	Transaction domain.AccountTransaction `json:"transaction"`
}

type CheckoutResponse struct {
	Order   domain.Order   `json:"order"`
	Invoice domain.Invoice `json:"invoice"`
}

type OrderDetailResponse struct {
	Order       domain.Order              `json:"order"`
	Invoices    []domain.Invoice          `json:"invoices"`
	Services    []domain.ServiceRecord    `json:"services"`
	ChangeOrder *ServiceChangeOrderRecord `json:"changeOrder,omitempty"`
	AuditLogs   []audit.Entry             `json:"auditLogs"`
}

type InvoiceDetailResponse struct {
	Invoice     domain.Invoice            `json:"invoice"`
	Order       *domain.Order             `json:"order,omitempty"`
	Services    []domain.ServiceRecord    `json:"services"`
	ChangeOrder *ServiceChangeOrderRecord `json:"changeOrder,omitempty"`
	Payments    []domain.PaymentRecord    `json:"payments"`
	Refunds     []domain.RefundRecord     `json:"refunds"`
	AuditLogs   []audit.Entry             `json:"auditLogs"`
}

type ServiceDetailResponse struct {
	Service      domain.ServiceRecord       `json:"service"`
	Order        *domain.Order              `json:"order,omitempty"`
	Invoice      *domain.Invoice            `json:"invoice,omitempty"`
	ChangeOrders []ServiceChangeOrderRecord `json:"changeOrders"`
	AuditLogs    []audit.Entry              `json:"auditLogs"`
}

type CreateServiceChangeOrderResponse struct {
	Service domain.ServiceRecord `json:"service"`
	Order   domain.Order         `json:"order"`
	Invoice domain.Invoice       `json:"invoice"`
}

type ServiceChangeOrderRecord struct {
	ID               int64                  `json:"id"`
	ServiceID        int64                  `json:"serviceId"`
	ServiceNo        string                 `json:"serviceNo"`
	OrderID          int64                  `json:"orderId"`
	OrderNo          string                 `json:"orderNo"`
	InvoiceID        int64                  `json:"invoiceId"`
	InvoiceNo        string                 `json:"invoiceNo"`
	ProductName      string                 `json:"productName"`
	ProviderType     string                 `json:"providerType"`
	ActionName       string                 `json:"actionName"`
	Title            string                 `json:"title"`
	Amount           float64                `json:"amount"`
	Status           string                 `json:"status"`
	ExecutionStatus  string                 `json:"executionStatus"`
	ExecutionMessage string                 `json:"executionMessage"`
	Reason           string                 `json:"reason"`
	BillingCycle     string                 `json:"billingCycle"`
	Payload          map[string]any         `json:"payload"`
	PaidAt           string                 `json:"paidAt"`
	RefundedAt       string                 `json:"refundedAt"`
	CreatedAt        string                 `json:"createdAt"`
	UpdatedAt        string                 `json:"updatedAt"`
	LatestTask       *automationDomain.Task `json:"latestTask,omitempty"`
}
