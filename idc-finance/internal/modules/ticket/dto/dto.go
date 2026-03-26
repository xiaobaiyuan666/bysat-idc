package dto

import (
	"idc-finance/internal/modules/ticket/domain"
	"idc-finance/internal/platform/audit"
)

type ListResponse struct {
	Items []domain.Ticket `json:"items"`
	Total int             `json:"total"`
}

type DetailResponse struct {
	Ticket    domain.Ticket  `json:"ticket"`
	Replies   []domain.Reply `json:"replies"`
	AuditLogs []audit.Entry  `json:"auditLogs"`
}

type CreateTicketRequest struct {
	ServiceID       *int64 `json:"serviceId"`
	Title           string `json:"title" binding:"required"`
	Content         string `json:"content" binding:"required"`
	Priority        string `json:"priority"`
	DepartmentName  string `json:"departmentName"`
	AssignedAdminID *int64 `json:"assignedAdminId"`
}

type CreateAdminTicketRequest struct {
	CustomerID        int64   `json:"customerId" binding:"required"`
	ServiceID         *int64  `json:"serviceId"`
	Title             string  `json:"title" binding:"required"`
	Content           string  `json:"content" binding:"required"`
	Priority          string  `json:"priority"`
	DepartmentName    string  `json:"departmentName"`
	AssignedAdminID   *int64  `json:"assignedAdminId"`
	AssignedAdminName *string `json:"assignedAdminName"`
}

type UpdateTicketRequest struct {
	Title             *string `json:"title"`
	Status            *string `json:"status"`
	Priority          *string `json:"priority"`
	DepartmentName    *string `json:"departmentName"`
	AssignedAdminID   *int64  `json:"assignedAdminId"`
	AssignedAdminName *string `json:"assignedAdminName"`
}

type ReplyRequest struct {
	Content    string  `json:"content" binding:"required"`
	Status     *string `json:"status"`
	IsInternal bool    `json:"isInternal"`
}

type UpdatePresetRepliesRequest struct {
	Items []domain.PresetReply `json:"items"`
}

type UpdateDepartmentsRequest struct {
	Items []domain.Department `json:"items"`
}

type AutoCloseSweepItem struct {
	TicketID     int64  `json:"ticketId"`
	TicketNo     string `json:"ticketNo"`
	CustomerName string `json:"customerName"`
	Department   string `json:"departmentName"`
	AutoCloseAt  string `json:"autoCloseAt"`
	ClosedAt     string `json:"closedAt"`
	Result       string `json:"result"`
	Message      string `json:"message"`
}

type AutoCloseSweepResponse struct {
	Enabled        bool                 `json:"enabled"`
	AutoCloseHours int                  `json:"autoCloseHours"`
	CheckedCount   int                  `json:"checkedCount"`
	SkippedCount   int                  `json:"skippedCount"`
	ClosedCount    int                  `json:"closedCount"`
	FailedCount    int                  `json:"failedCount"`
	TaskID         int64                `json:"taskId"`
	TaskNo         string               `json:"taskNo"`
	Message        string               `json:"message"`
	Items          []AutoCloseSweepItem `json:"items"`
}

type StatisticsResponse struct {
	Summary         domain.SummaryStats      `json:"summary"`
	DepartmentStats []domain.DepartmentStats `json:"departmentStats"`
	AdminStats      []domain.AdminStats      `json:"adminStats"`
}
