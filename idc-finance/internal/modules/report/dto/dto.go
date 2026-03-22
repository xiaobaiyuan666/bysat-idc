package dto

type MetricCard struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Value string `json:"value"`
	Hint  string `json:"hint"`
	Tone  string `json:"tone"`
}

type StatusBucket struct {
	Name   string  `json:"name"`
	Count  int64   `json:"count"`
	Amount float64 `json:"amount"`
}

type TrendPoint struct {
	Label  string  `json:"label"`
	Amount float64 `json:"amount"`
	Count  int64   `json:"count"`
}

type RankItem struct {
	Name   string  `json:"name"`
	Count  int64   `json:"count"`
	Amount float64 `json:"amount"`
	Extra  string  `json:"extra"`
}

type PendingIdentityItem struct {
	CustomerID   int64  `json:"customerId"`
	CustomerNo   string `json:"customerNo"`
	CustomerName string `json:"customerName"`
	SubjectName  string `json:"subjectName"`
	SubmittedAt  string `json:"submittedAt"`
}

type OverdueInvoiceItem struct {
	InvoiceID    int64   `json:"invoiceId"`
	InvoiceNo    string  `json:"invoiceNo"`
	CustomerName string  `json:"customerName"`
	Amount       float64 `json:"amount"`
	DueAt        string  `json:"dueAt"`
	DaysOverdue  int64   `json:"daysOverdue"`
}

type ExpiringServiceItem struct {
	ServiceID     int64  `json:"serviceId"`
	ServiceNo     string `json:"serviceNo"`
	CustomerName  string `json:"customerName"`
	ProductName   string `json:"productName"`
	Status        string `json:"status"`
	NextDueAt     string `json:"nextDueAt"`
	DaysRemaining int64  `json:"daysRemaining"`
}

type TicketTodoItem struct {
	TicketID     int64  `json:"ticketId"`
	TicketNo     string `json:"ticketNo"`
	CustomerName string `json:"customerName"`
	Title        string `json:"title"`
	Status       string `json:"status"`
	UpdatedAt    string `json:"updatedAt"`
}

type AuditTimelineItem struct {
	ID          int64  `json:"id"`
	Actor       string `json:"actor"`
	Action      string `json:"action"`
	Target      string `json:"target"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

type WorkbenchResponse struct {
	SummaryCards      []MetricCard         `json:"summaryCards"`
	RiskCards         []MetricCard         `json:"riskCards"`
	RevenueTrend      []TrendPoint         `json:"revenueTrend"`
	ServiceStatus     []StatusBucket       `json:"serviceStatus"`
	PendingIdentities []PendingIdentityItem `json:"pendingIdentities"`
	OverdueInvoices   []OverdueInvoiceItem `json:"overdueInvoices"`
	ExpiringServices  []ExpiringServiceItem `json:"expiringServices"`
	OpenTickets       []TicketTodoItem     `json:"openTickets"`
	RecentAudits      []AuditTimelineItem  `json:"recentAudits"`
}

type ReportOverviewResponse struct {
	Headline        []MetricCard   `json:"headline"`
	RevenueTrend    []TrendPoint   `json:"revenueTrend"`
	RefundTrend     []TrendPoint   `json:"refundTrend"`
	InvoiceStatus   []StatusBucket `json:"invoiceStatus"`
	ServiceStatus   []StatusBucket `json:"serviceStatus"`
	PaymentChannels []StatusBucket `json:"paymentChannels"`
	BillingCycles   []StatusBucket `json:"billingCycles"`
	CustomerGroups  []StatusBucket `json:"customerGroups"`
	TopProducts     []RankItem     `json:"topProducts"`
	TopReceivables  []RankItem     `json:"topReceivables"`
}
