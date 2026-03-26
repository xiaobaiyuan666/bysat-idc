package domain

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusActive    OrderStatus = "ACTIVE"
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type InvoiceStatus string

const (
	InvoiceStatusUnpaid   InvoiceStatus = "UNPAID"
	InvoiceStatusPaid     InvoiceStatus = "PAID"
	InvoiceStatusRefunded InvoiceStatus = "REFUNDED"
)

type ServiceStatus string

const (
	ServiceStatusPending    ServiceStatus = "PENDING"
	ServiceStatusActive     ServiceStatus = "ACTIVE"
	ServiceStatusSuspended  ServiceStatus = "SUSPENDED"
	ServiceStatusTerminated ServiceStatus = "TERMINATED"
)

type RefundStatus string

const (
	RefundStatusCompleted RefundStatus = "COMPLETED"
)

type PaymentStatus string

const (
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
)

type OrderRequestType string

const (
	OrderRequestTypeCancel      OrderRequestType = "CANCEL"
	OrderRequestTypeRenew       OrderRequestType = "RENEW"
	OrderRequestTypePriceAdjust OrderRequestType = "PRICE_ADJUST"
)

type OrderRequestStatus string

const (
	OrderRequestStatusPending   OrderRequestStatus = "PENDING"
	OrderRequestStatusApproved  OrderRequestStatus = "APPROVED"
	OrderRequestStatusRejected  OrderRequestStatus = "REJECTED"
	OrderRequestStatusCompleted OrderRequestStatus = "COMPLETED"
	OrderRequestStatusCancelled OrderRequestStatus = "CANCELLED"
)

type AccountTransactionType string

const (
	AccountTransactionTypeRecharge    AccountTransactionType = "RECHARGE"
	AccountTransactionTypeConsume     AccountTransactionType = "CONSUME"
	AccountTransactionTypeRefund      AccountTransactionType = "REFUND"
	AccountTransactionTypeAdjustment  AccountTransactionType = "ADJUSTMENT"
	AccountTransactionTypeCreditLimit AccountTransactionType = "CREDIT_LIMIT"
)

type AccountTransactionDirection string

const (
	AccountTransactionDirectionIn   AccountTransactionDirection = "IN"
	AccountTransactionDirectionOut  AccountTransactionDirection = "OUT"
	AccountTransactionDirectionFlat AccountTransactionDirection = "FLAT"
)

type OrderListFilter struct {
	Page        int     `json:"page"`
	Limit       int     `json:"limit"`
	Sort        string  `json:"sort"`
	Order       string  `json:"order"`
	Status      string  `json:"status"`
	OrderNo     string  `json:"orderNo"`
	StartTime   string  `json:"startTime"`
	EndTime     string  `json:"endTime"`
	Amount      float64 `json:"amount"`
	HasAmount   bool    `json:"hasAmount"`
	CustomerID  int64   `json:"customerId"`
	Payment     string  `json:"payment"`
	PayStatus   string  `json:"payStatus"`
	ProductName string  `json:"productName"`
}

type InvoiceListFilter struct {
	Page         int    `json:"page"`
	Limit        int    `json:"limit"`
	Sort         string `json:"sort"`
	Order        string `json:"order"`
	Status       string `json:"status"`
	InvoiceNo    string `json:"invoiceNo"`
	OrderNo      string `json:"orderNo"`
	ProductName  string `json:"productName"`
	BillingCycle string `json:"billingCycle"`
	CustomerID   int64  `json:"customerId"`
}

type ServiceListFilter struct {
	Page              int    `json:"page"`
	Limit             int    `json:"limit"`
	Sort              string `json:"sort"`
	Order             string `json:"order"`
	Status            string `json:"status"`
	CustomerID        int64  `json:"customerId"`
	OrderID           int64  `json:"orderId"`
	ProviderType      string `json:"providerType"`
	ProviderAccountID int64  `json:"providerAccountId"`
	SyncStatus        string `json:"syncStatus"`
	Keyword           string `json:"keyword"`
}

type AccountTransactionFilter struct {
	Page            int    `json:"page"`
	Limit           int    `json:"limit"`
	Sort            string `json:"sort"`
	Order           string `json:"order"`
	CustomerID      int64  `json:"customerId"`
	TransactionType string `json:"transactionType"`
	Direction       string `json:"direction"`
	Channel         string `json:"channel"`
	Keyword         string `json:"keyword"`
	StartTime       string `json:"startTime"`
	EndTime         string `json:"endTime"`
}

type PaymentListFilter struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Sort       string `json:"sort"`
	Order      string `json:"order"`
	CustomerID int64  `json:"customerId"`
	InvoiceID  int64  `json:"invoiceId"`
	Keyword    string `json:"keyword"`
	Channel    string `json:"channel"`
	Status     string `json:"status"`
}

type RefundListFilter struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Sort       string `json:"sort"`
	Order      string `json:"order"`
	CustomerID int64  `json:"customerId"`
	InvoiceID  int64  `json:"invoiceId"`
	Keyword    string `json:"keyword"`
	Status     string `json:"status"`
}

type ServiceChangeOrderListFilter struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Sort      string `json:"sort"`
	Order     string `json:"order"`
	Status    string `json:"status"`
	ServiceID int64  `json:"serviceId"`
	OrderID   int64  `json:"orderId"`
	InvoiceID int64  `json:"invoiceId"`
	Action    string `json:"action"`
	Keyword   string `json:"keyword"`
}

type OrderRequestListFilter struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Sort       string `json:"sort"`
	Order      string `json:"order"`
	OrderID    int64  `json:"orderId"`
	CustomerID int64  `json:"customerId"`
	Type       string `json:"type"`
	Status     string `json:"status"`
	Keyword    string `json:"keyword"`
}

type WalletAdjustment struct {
	Target       string  `json:"target"`
	Operation    string  `json:"operation"`
	Amount       float64 `json:"amount"`
	Summary      string  `json:"summary"`
	Remark       string  `json:"remark"`
	OperatorType string  `json:"operatorType"`
	OperatorID   int64   `json:"operatorId"`
	OperatorName string  `json:"operatorName"`
}

type ServiceConfigSelection struct {
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	Value      string  `json:"value"`
	ValueLabel string  `json:"valueLabel"`
	PriceDelta float64 `json:"priceDelta"`
}

type ServiceResourceSnapshot struct {
	RegionName      string `json:"regionName"`
	ZoneName        string `json:"zoneName"`
	Hostname        string `json:"hostname"`
	OperatingSystem string `json:"operatingSystem"`
	LoginUsername   string `json:"loginUsername"`
	PasswordHint    string `json:"passwordHint"`
	SecurityGroup   string `json:"securityGroup"`
	CPUCores        int    `json:"cpuCores"`
	MemoryGB        int    `json:"memoryGB"`
	SystemDiskGB    int    `json:"systemDiskGB"`
	DataDiskGB      int    `json:"dataDiskGB"`
	BandwidthMbps   int    `json:"bandwidthMbps"`
	PublicIPv4      string `json:"publicIpv4"`
	PublicIPv6      string `json:"publicIpv6"`
}

type ServiceActionParams struct {
	ImageName string
	Password  string
}

type PendingOrderUpdate struct {
	ProductName  string   `json:"productName"`
	BillingCycle string   `json:"billingCycle"`
	Amount       *float64 `json:"amount"`
	DueAt        string   `json:"dueAt"`
	Status       string   `json:"status"`
}

type UnpaidInvoiceUpdate struct {
	ProductName  string   `json:"productName"`
	BillingCycle string   `json:"billingCycle"`
	Amount       *float64 `json:"amount"`
	DueAt        string   `json:"dueAt"`
	Status       string   `json:"status"`
}

type ManualServiceUpdate struct {
	ProviderType       string `json:"providerType"`
	ProviderAccountID  *int64 `json:"providerAccountId"`
	ProviderResourceID string `json:"providerResourceId"`
	RegionName         string `json:"regionName"`
	IPAddress          string `json:"ipAddress"`
	NextDueAt          string `json:"nextDueAt"`
	Status             string `json:"status"`
	SyncStatus         string `json:"syncStatus"`
	SyncMessage        string `json:"syncMessage"`
}

type ServiceChangeOrderInput struct {
	ActionName   string         `json:"actionName"`
	Title        string         `json:"title"`
	BillingCycle string         `json:"billingCycle"`
	Amount       float64        `json:"amount"`
	Reason       string         `json:"reason"`
	Payload      map[string]any `json:"payload"`
}

type ServiceChangeOrder struct {
	ID           int64          `json:"id"`
	ServiceID    int64          `json:"serviceId"`
	OrderID      int64          `json:"orderId"`
	InvoiceID    int64          `json:"invoiceId"`
	ActionName   string         `json:"actionName"`
	Title        string         `json:"title"`
	Amount       float64        `json:"amount"`
	Status       string         `json:"status"`
	Reason       string         `json:"reason"`
	BillingCycle string         `json:"billingCycle"`
	Payload      map[string]any `json:"payload"`
	PaidAt       string         `json:"paidAt"`
	RefundedAt   string         `json:"refundedAt"`
	CreatedAt    string         `json:"createdAt"`
	UpdatedAt    string         `json:"updatedAt"`
}

type OrderRequestCreateInput struct {
	Type                 string         `json:"type"`
	Summary              string         `json:"summary"`
	Reason               string         `json:"reason"`
	RequestedAmount      *float64       `json:"requestedAmount"`
	RequestedBillingCycle string        `json:"requestedBillingCycle"`
	SourceType           string         `json:"sourceType"`
	SourceID             int64          `json:"sourceId"`
	SourceName           string         `json:"sourceName"`
	Payload              map[string]any `json:"payload"`
}

type OrderRequestProcessInput struct {
	Status        string `json:"status"`
	ProcessNote   string `json:"processNote"`
	ProcessorType string `json:"processorType"`
	ProcessorID   int64  `json:"processorId"`
	ProcessorName string `json:"processorName"`
}

type OrderRequest struct {
	ID                   int64              `json:"id"`
	RequestNo            string             `json:"requestNo"`
	OrderID              int64              `json:"orderId"`
	OrderNo              string             `json:"orderNo"`
	CustomerID           int64              `json:"customerId"`
	CustomerName         string             `json:"customerName"`
	ProductName          string             `json:"productName"`
	Type                 OrderRequestType   `json:"type"`
	Status               OrderRequestStatus `json:"status"`
	Summary              string             `json:"summary"`
	Reason               string             `json:"reason"`
	CurrentAmount        float64            `json:"currentAmount"`
	RequestedAmount      float64            `json:"requestedAmount"`
	CurrentBillingCycle  string             `json:"currentBillingCycle"`
	RequestedBillingCycle string            `json:"requestedBillingCycle"`
	SourceType           string             `json:"sourceType"`
	SourceID             int64              `json:"sourceId"`
	SourceName           string             `json:"sourceName"`
	ProcessorType        string             `json:"processorType"`
	ProcessorID          int64              `json:"processorId"`
	ProcessorName        string             `json:"processorName"`
	ProcessNote          string             `json:"processNote"`
	Payload              map[string]any     `json:"payload"`
	CreatedAt            string             `json:"createdAt"`
	UpdatedAt            string             `json:"updatedAt"`
	ProcessedAt          string             `json:"processedAt"`
}

type ServiceProvisioningUpdate struct {
	ProviderAccountID  int64                    `json:"providerAccountId"`
	ProviderType       string                   `json:"providerType"`
	ProviderResourceID string                   `json:"providerResourceId"`
	RegionName         string                   `json:"regionName"`
	IPAddress          string                   `json:"ipAddress"`
	Status             ServiceStatus            `json:"status"`
	SyncStatus         string                   `json:"syncStatus"`
	SyncMessage        string                   `json:"syncMessage"`
	LastAction         string                   `json:"lastAction"`
	ResourceSnapshot   ServiceResourceSnapshot  `json:"resourceSnapshot"`
	Configuration      []ServiceConfigSelection `json:"configuration"`
}

type Order struct {
	ID                int64                    `json:"id"`
	OrderNo           string                   `json:"orderNo"`
	CustomerID        int64                    `json:"customerId"`
	CustomerName      string                   `json:"customerName"`
	ProductID         int64                    `json:"productId"`
	ProductName       string                   `json:"productName"`
	ProductType       string                   `json:"productType"`
	AutomationType    string                   `json:"automationType"`
	ProviderAccountID int64                    `json:"providerAccountId"`
	BillingCycle      string                   `json:"billingCycle"`
	Amount            float64                  `json:"amount"`
	Status            OrderStatus              `json:"status"`
	Payment           string                   `json:"payment"`
	PayStatus         string                   `json:"payStatus"`
	Configuration     []ServiceConfigSelection `json:"configuration"`
	ResourceSnapshot  ServiceResourceSnapshot  `json:"resourceSnapshot"`
	CreatedAt         string                   `json:"createdAt"`
}

type Invoice struct {
	ID           int64         `json:"id"`
	InvoiceNo    string        `json:"invoiceNo"`
	OrderID      int64         `json:"orderId"`
	OrderNo      string        `json:"orderNo"`
	CustomerID   int64         `json:"customerId"`
	ProductName  string        `json:"productName"`
	TotalAmount  float64       `json:"totalAmount"`
	Status       InvoiceStatus `json:"status"`
	DueAt        string        `json:"dueAt"`
	PaidAt       string        `json:"paidAt,omitempty"`
	BillingCycle string        `json:"billingCycle"`
}

type ServiceRecord struct {
	ID                 int64                    `json:"id"`
	ServiceNo          string                   `json:"serviceNo"`
	OrderID            int64                    `json:"orderId"`
	InvoiceID          int64                    `json:"invoiceId"`
	CustomerID         int64                    `json:"customerId"`
	ProductName        string                   `json:"productName"`
	ProviderType       string                   `json:"providerType"`
	ProviderAccountID  int64                    `json:"providerAccountId"`
	ProviderResourceID string                   `json:"providerResourceId"`
	RegionName         string                   `json:"regionName"`
	IPAddress          string                   `json:"ipAddress"`
	Status             ServiceStatus            `json:"status"`
	SyncStatus         string                   `json:"syncStatus"`
	SyncMessage        string                   `json:"syncMessage"`
	NextDueAt          string                   `json:"nextDueAt"`
	LastAction         string                   `json:"lastAction"`
	LastSyncAt         string                   `json:"lastSyncAt"`
	UpdatedAt          string                   `json:"updatedAt"`
	Configuration      []ServiceConfigSelection `json:"configuration"`
	ResourceSnapshot   ServiceResourceSnapshot  `json:"resourceSnapshot"`
	CreatedAt          string                   `json:"createdAt"`
}

type RefundRecord struct {
	ID         int64        `json:"id"`
	RefundNo   string       `json:"refundNo"`
	InvoiceID  int64        `json:"invoiceId"`
	OrderID    int64        `json:"orderId"`
	CustomerID int64        `json:"customerId"`
	Amount     float64      `json:"amount"`
	Reason     string       `json:"reason"`
	Status     RefundStatus `json:"status"`
	CreatedAt  string       `json:"createdAt"`
}

type PaymentRecord struct {
	ID         int64         `json:"id"`
	PaymentNo  string        `json:"paymentNo"`
	InvoiceID  int64         `json:"invoiceId"`
	OrderID    int64         `json:"orderId"`
	CustomerID int64         `json:"customerId"`
	Channel    string        `json:"channel"`
	TradeNo    string        `json:"tradeNo"`
	Amount     float64       `json:"amount"`
	Source     string        `json:"source"`
	Status     PaymentStatus `json:"status"`
	Operator   string        `json:"operator"`
	PaidAt     string        `json:"paidAt"`
}

type CustomerWallet struct {
	CustomerID      int64   `json:"customerId"`
	CustomerNo      string  `json:"customerNo"`
	CustomerName    string  `json:"customerName"`
	Balance         float64 `json:"balance"`
	CreditLimit     float64 `json:"creditLimit"`
	CreditUsed      float64 `json:"creditUsed"`
	AvailableCredit float64 `json:"availableCredit"`
	UpdatedAt       string  `json:"updatedAt"`
}

type AccountTransaction struct {
	ID              int64                       `json:"id"`
	TransactionNo   string                      `json:"transactionNo"`
	CustomerID      int64                       `json:"customerId"`
	CustomerNo      string                      `json:"customerNo"`
	CustomerName    string                      `json:"customerName"`
	OrderID         int64                       `json:"orderId"`
	OrderNo         string                      `json:"orderNo"`
	InvoiceID       int64                       `json:"invoiceId"`
	InvoiceNo       string                      `json:"invoiceNo"`
	PaymentID       int64                       `json:"paymentId"`
	PaymentNo       string                      `json:"paymentNo"`
	RefundID        int64                       `json:"refundId"`
	RefundNo        string                      `json:"refundNo"`
	TransactionType AccountTransactionType      `json:"transactionType"`
	Direction       AccountTransactionDirection `json:"direction"`
	Amount          float64                     `json:"amount"`
	BalanceBefore   float64                     `json:"balanceBefore"`
	BalanceAfter    float64                     `json:"balanceAfter"`
	CreditBefore    float64                     `json:"creditBefore"`
	CreditAfter     float64                     `json:"creditAfter"`
	Channel         string                      `json:"channel"`
	Summary         string                      `json:"summary"`
	Remark          string                      `json:"remark"`
	OperatorType    string                      `json:"operatorType"`
	OperatorID      int64                       `json:"operatorId"`
	OperatorName    string                      `json:"operatorName"`
	OccurredAt      string                      `json:"occurredAt"`
}
