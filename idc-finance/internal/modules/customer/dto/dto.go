package dto

import (
	"idc-finance/internal/modules/customer/domain"
	"idc-finance/internal/platform/audit"
)

type CreateCustomerRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Mobile     string `json:"mobile"`
	Type       string `json:"type"`
	GroupName  string `json:"groupName"`
	LevelName  string `json:"levelName"`
	SalesOwner string `json:"salesOwner"`
	Remarks    string `json:"remarks"`
}

type CustomerListResponse struct {
	Items []domain.Customer `json:"items"`
	Total int               `json:"total"`
}

type UpdateCustomerRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Mobile     string `json:"mobile"`
	GroupName  string `json:"groupName"`
	LevelName  string `json:"levelName"`
	SalesOwner string `json:"salesOwner"`
	Remarks    string `json:"remarks"`
	Status     string `json:"status"`
}

type CreateContactRequest struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"omitempty,email"`
	Mobile    string `json:"mobile"`
	RoleName  string `json:"roleName"`
	IsPrimary bool   `json:"isPrimary"`
}

type UpdateContactRequest struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"omitempty,email"`
	Mobile    string `json:"mobile"`
	RoleName  string `json:"roleName"`
	IsPrimary bool   `json:"isPrimary"`
}

type ReviewIdentityRequest struct {
	Status       string `json:"status" binding:"required"`
	ReviewRemark string `json:"reviewRemark"`
	Remark       string `json:"remark"`
}

type SaveCustomerGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type SaveCustomerLevelRequest struct {
	Name        string `json:"name" binding:"required"`
	Priority    int    `json:"priority"`
	Description string `json:"description"`
}

type IdentityOverviewItem struct {
	ID           int64                 `json:"id"`
	CustomerID   int64                 `json:"customerId"`
	CustomerNo   string                `json:"customerNo"`
	CustomerName string                `json:"customerName"`
	CustomerType string                `json:"customerType"`
	IdentityType string                `json:"identityType"`
	VerifyStatus domain.IdentityStatus `json:"verifyStatus"`
	SubjectName  string                `json:"subjectName"`
	CertNo       string                `json:"certNo"`
	CountryCode  string                `json:"countryCode"`
	ReviewRemark string                `json:"reviewRemark"`
	ReviewedAt   string                `json:"reviewedAt"`
	SubmittedAt  string                `json:"submittedAt"`
}

type RelatedItem = domain.RelatedItem

type CustomerWorkbenchResponse struct {
	Customer  domain.Customer      `json:"customer"`
	Services  []domain.RelatedItem `json:"services"`
	Invoices  []domain.RelatedItem `json:"invoices"`
	Tickets   []domain.RelatedItem `json:"tickets"`
	AuditLogs []audit.Entry        `json:"auditLogs"`
}
