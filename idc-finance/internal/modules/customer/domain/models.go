package domain

type CustomerStatus string

const (
	CustomerStatusActive   CustomerStatus = "ACTIVE"
	CustomerStatusDisabled CustomerStatus = "DISABLED"
)

type IdentityStatus string

const (
	IdentityStatusPending  IdentityStatus = "PENDING"
	IdentityStatusApproved IdentityStatus = "APPROVED"
	IdentityStatusRejected IdentityStatus = "REJECTED"
)

type Contact struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
	RoleName  string `json:"roleName"`
	IsPrimary bool   `json:"isPrimary"`
}

type CustomerGroup struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	CustomerCount int    `json:"customerCount"`
}

type CustomerLevel struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Priority      int    `json:"priority"`
	Description   string `json:"description"`
	CustomerCount int    `json:"customerCount"`
}

type Identity struct {
	ID           int64          `json:"id"`
	IdentityType string         `json:"identityType"`
	VerifyStatus IdentityStatus `json:"verifyStatus"`
	SubjectName  string         `json:"subjectName"`
	CertNo       string         `json:"certNo"`
	CountryCode  string         `json:"countryCode"`
	ReviewRemark string         `json:"reviewRemark"`
	ReviewedAt   string         `json:"reviewedAt"`
	SubmittedAt  string         `json:"submittedAt"`
}

type RelatedItem struct {
	ID                 int64  `json:"id"`
	ServiceID          int64  `json:"serviceId,omitempty"`
	InvoiceID          int64  `json:"invoiceId,omitempty"`
	TicketID           int64  `json:"ticketId,omitempty"`
	No                 string `json:"no"`
	Name               string `json:"name"`
	Status             string `json:"status"`
	Amount             string `json:"amount,omitempty"`
	DueAt              string `json:"dueAt,omitempty"`
	UpdatedAt          string `json:"updatedAt,omitempty"`
	Description        string `json:"description,omitempty"`
	ProviderType       string `json:"providerType,omitempty"`
	ProviderResourceID string `json:"providerResourceId,omitempty"`
	RegionName         string `json:"regionName,omitempty"`
	IPAddress          string `json:"ipAddress,omitempty"`
	BillingCycle       string `json:"billingCycle,omitempty"`
}

type Customer struct {
	ID         int64          `json:"id"`
	CustomerNo string         `json:"customerNo"`
	Name       string         `json:"name"`
	Email      string         `json:"email"`
	Mobile     string         `json:"mobile"`
	Type       string         `json:"type"`
	Status     CustomerStatus `json:"status"`
	GroupName  string         `json:"groupName"`
	LevelName  string         `json:"levelName"`
	SalesOwner string         `json:"salesOwner"`
	Remarks    string         `json:"remarks"`
	Contacts   []Contact      `json:"contacts"`
	Identity   *Identity      `json:"identity"`
}
