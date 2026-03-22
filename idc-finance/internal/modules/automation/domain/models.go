package domain

type TaskStatus string

const (
	TaskStatusPending TaskStatus = "PENDING"
	TaskStatusRunning TaskStatus = "RUNNING"
	TaskStatusSuccess TaskStatus = "SUCCESS"
	TaskStatusFailed  TaskStatus = "FAILED"
	TaskStatusBlocked TaskStatus = "BLOCKED"
)

type Task struct {
	ID                 int64      `json:"id"`
	TaskNo             string     `json:"taskNo"`
	TaskType           string     `json:"taskType"`
	Title              string     `json:"title"`
	Channel            string     `json:"channel"`
	Stage              string     `json:"stage"`
	Status             TaskStatus `json:"status"`
	SourceType         string     `json:"sourceType"`
	SourceID           int64      `json:"sourceId"`
	CustomerID         int64      `json:"customerId"`
	CustomerName       string     `json:"customerName"`
	ProductName        string     `json:"productName"`
	OrderID            int64      `json:"orderId"`
	InvoiceID          int64      `json:"invoiceId"`
	ServiceID          int64      `json:"serviceId"`
	ServiceNo          string     `json:"serviceNo"`
	ProviderType       string     `json:"providerType"`
	ProviderResourceID string     `json:"providerResourceId"`
	ActionName         string     `json:"actionName"`
	OperatorType       string     `json:"operatorType"`
	OperatorName       string     `json:"operatorName"`
	RequestPayload     string     `json:"requestPayload"`
	ResultPayload      string     `json:"resultPayload"`
	Message            string     `json:"message"`
	CreatedAt          string     `json:"createdAt"`
	StartedAt          string     `json:"startedAt"`
	FinishedAt         string     `json:"finishedAt"`
	UpdatedAt          string     `json:"updatedAt"`
}

type TaskFilter struct {
	Status     string
	TaskType   string
	Channel    string
	Stage      string
	SourceType string
	SourceID   int64
	OrderID    int64
	InvoiceID  int64
	ServiceID  int64
	Keyword    string
	Limit      int
}

type TaskSummary struct {
	Total   int `json:"total"`
	Running int `json:"running"`
	Success int `json:"success"`
	Failed  int `json:"failed"`
	Blocked int `json:"blocked"`
}

type AutomationSettings struct {
	AutoProvisionEnabled bool   `json:"autoProvisionEnabled"`
	AutoSuspendEnabled   bool   `json:"autoSuspendEnabled"`
	SuspendAfterDays     int    `json:"suspendAfterDays"`
	SuspendNoticeDays    int    `json:"suspendNoticeDays"`
	AutoTerminateEnabled bool   `json:"autoTerminateEnabled"`
	TerminateAfterDays   int    `json:"terminateAfterDays"`
	TerminateNoticeDays  int    `json:"terminateNoticeDays"`
	InvoiceReminderOn    bool   `json:"invoiceReminderOn"`
	InvoiceReminderDays  string `json:"invoiceReminderDays"`
	CreditReminderOn     bool   `json:"creditReminderOn"`
	CreditReminderDays   string `json:"creditReminderDays"`
	TicketAutoCloseOn    bool   `json:"ticketAutoCloseOn"`
	TicketAutoCloseHours int    `json:"ticketAutoCloseHours"`
	LogRetentionDays     int    `json:"logRetentionDays"`
}
