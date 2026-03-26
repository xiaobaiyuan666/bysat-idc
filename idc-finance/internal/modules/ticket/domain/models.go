package domain

type Status string

const (
	StatusOpen            Status = "OPEN"
	StatusProcessing      Status = "PROCESSING"
	StatusWaitingCustomer Status = "WAITING_CUSTOMER"
	StatusClosed          Status = "CLOSED"
)

type Priority string

const (
	PriorityLow    Priority = "LOW"
	PriorityNormal Priority = "NORMAL"
	PriorityHigh   Priority = "HIGH"
	PriorityUrgent Priority = "URGENT"
)

type AuthorType string

const (
	AuthorTypeAdmin    AuthorType = "ADMIN"
	AuthorTypeCustomer AuthorType = "CUSTOMER"
)

type ListFilter struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Keyword    string `json:"keyword"`
	Status     string `json:"status"`
	Priority   string `json:"priority"`
	CustomerID int64  `json:"customerId"`
	ServiceID  int64  `json:"serviceId"`
	Department string `json:"departmentName"`
	AdminID    int64  `json:"assignedAdminId"`
}

type Ticket struct {
	ID                 int64    `json:"id"`
	TicketNo           string   `json:"ticketNo"`
	CustomerID         int64    `json:"customerId"`
	CustomerNo         string   `json:"customerNo"`
	CustomerName       string   `json:"customerName"`
	ServiceID          int64    `json:"serviceId"`
	ServiceNo          string   `json:"serviceNo"`
	ProductName        string   `json:"productName"`
	Title              string   `json:"title"`
	Content            string   `json:"content"`
	Status             Status   `json:"status"`
	Priority           Priority `json:"priority"`
	Source             string   `json:"source"`
	DepartmentName     string   `json:"departmentName"`
	AssignedAdminID    int64    `json:"assignedAdminId"`
	AssignedAdminName  string   `json:"assignedAdminName"`
	LatestReplyExcerpt string   `json:"latestReplyExcerpt"`
	LastReplyAt        string   `json:"lastReplyAt"`
	ClosedAt           string   `json:"closedAt"`
	CreatedAt          string   `json:"createdAt"`
	UpdatedAt          string   `json:"updatedAt"`
	SLAStatus          string   `json:"slaStatus"`
	SLADeadlineAt      string   `json:"slaDeadlineAt"`
	SLARemainingMins   int      `json:"slaRemainingMins"`
	SLAPaused          bool     `json:"slaPaused"`
	AutoCloseAt        string   `json:"autoCloseAt"`
	AutoCloseMins      int      `json:"autoCloseMins"`
}

type Reply struct {
	ID         int64      `json:"id"`
	TicketID   int64      `json:"ticketId"`
	AuthorType AuthorType `json:"authorType"`
	AuthorID   int64      `json:"authorId"`
	AuthorName string     `json:"authorName"`
	Content    string     `json:"content"`
	IsInternal bool       `json:"isInternal"`
	CreatedAt  string     `json:"createdAt"`
	UpdatedAt  string     `json:"updatedAt"`
}

type CreateInput struct {
	CustomerID      int64    `json:"customerId"`
	CustomerNo      string   `json:"customerNo"`
	CustomerName    string   `json:"customerName"`
	ServiceID       int64    `json:"serviceId"`
	Title           string   `json:"title"`
	Content         string   `json:"content"`
	Status          Status   `json:"status"`
	Priority        Priority `json:"priority"`
	Source          string   `json:"source"`
	DepartmentName  string   `json:"departmentName"`
	AssignedAdminID int64    `json:"assignedAdminId"`
	AssignedAdmin   string   `json:"assignedAdminName"`
}

type UpdateInput struct {
	Title             *string   `json:"title"`
	Status            *Status   `json:"status"`
	Priority          *Priority `json:"priority"`
	DepartmentName    *string   `json:"departmentName"`
	AssignedAdminID   *int64    `json:"assignedAdminId"`
	AssignedAdminName *string   `json:"assignedAdminName"`
}

type ReplyInput struct {
	AuthorType AuthorType `json:"authorType"`
	AuthorID   int64      `json:"authorId"`
	AuthorName string     `json:"authorName"`
	Content    string     `json:"content"`
	Status     Status     `json:"status"`
	IsInternal bool       `json:"isInternal"`
}

type PresetReply struct {
	Key            string `json:"key"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	DepartmentName string `json:"departmentName"`
	Status         string `json:"status"`
}

type Department struct {
	Key         string  `json:"key"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Enabled     bool    `json:"enabled"`
	IsDefault   bool    `json:"isDefault"`
	Sort        int     `json:"sort"`
	AdminIDs    []int64 `json:"adminIds"`
}

type SummaryStats struct {
	Total           int `json:"total"`
	Unassigned      int `json:"unassigned"`
	WaitingCustomer int `json:"waitingCustomer"`
	Breached        int `json:"breached"`
	Closed          int `json:"closed"`
}

type DepartmentStats struct {
	Key             string `json:"key"`
	Name            string `json:"name"`
	Total           int    `json:"total"`
	Open            int    `json:"open"`
	Processing      int    `json:"processing"`
	WaitingCustomer int    `json:"waitingCustomer"`
	Closed          int    `json:"closed"`
	Breached        int    `json:"breached"`
}

type AdminStats struct {
	AdminID         int64  `json:"adminId"`
	AdminName       string `json:"adminName"`
	Total           int    `json:"total"`
	Open            int    `json:"open"`
	Processing      int    `json:"processing"`
	WaitingCustomer int    `json:"waitingCustomer"`
	Closed          int    `json:"closed"`
	Breached        int    `json:"breached"`
}
