package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"idc-finance/internal/modules/automation/domain"
	"idc-finance/internal/modules/automation/repository"
)

type Service struct {
	repository repository.Repository
	db         *sql.DB
}

type StartTaskRequest struct {
	TaskType           string
	Title              string
	Channel            string
	Stage              string
	SourceType         string
	SourceID           int64
	CustomerID         int64
	CustomerName       string
	ProductName        string
	OrderID            int64
	InvoiceID          int64
	ServiceID          int64
	ServiceNo          string
	ProviderType       string
	ProviderResourceID string
	ActionName         string
	OperatorType       string
	OperatorName       string
	RequestPayload     any
	Message            string
}

type FinishTaskRequest struct {
	Status        domain.TaskStatus
	Message       string
	ResultPayload any
}

func New(repo repository.Repository) *Service {
	return &Service{repository: repo}
}

func NewWithDB(repo repository.Repository, db *sql.DB) *Service {
	return &Service{repository: repo, db: db}
}

func (service *Service) List(filter domain.TaskFilter) ([]domain.Task, domain.TaskSummary) {
	items := service.repository.List(filter)
	summary := domain.TaskSummary{Total: len(items)}
	for _, item := range items {
		switch item.Status {
		case domain.TaskStatusRunning:
			summary.Running++
		case domain.TaskStatusSuccess:
			summary.Success++
		case domain.TaskStatusFailed:
			summary.Failed++
		case domain.TaskStatusBlocked:
			summary.Blocked++
		}
	}
	return items, summary
}

func (service *Service) GetByID(id int64) (domain.Task, bool) {
	return service.repository.GetByID(id)
}

func (service *Service) GetSettings() domain.AutomationSettings {
	settings := defaultSettings()
	if service == nil || service.db == nil {
		return settings
	}

	rows, err := service.db.Query(`
SELECT config_key, config_value
FROM system_configs
WHERE config_key LIKE 'automation.%'`)
	if err != nil {
		return settings
	}
	defer rows.Close()

	values := make(map[string]string)
	for rows.Next() {
		var key string
		var value string
		if err := rows.Scan(&key, &value); err != nil {
			continue
		}
		values[key] = value
	}

	settings.AutoProvisionEnabled = parseBool(values["automation.auto_provision_enabled"], settings.AutoProvisionEnabled)
	settings.AutoSuspendEnabled = parseBool(values["automation.auto_suspend_enabled"], settings.AutoSuspendEnabled)
	settings.SuspendAfterDays = parseInt(values["automation.suspend_after_days"], settings.SuspendAfterDays)
	settings.SuspendNoticeDays = parseInt(values["automation.suspend_notice_days"], settings.SuspendNoticeDays)
	settings.AutoTerminateEnabled = parseBool(values["automation.auto_terminate_enabled"], settings.AutoTerminateEnabled)
	settings.TerminateAfterDays = parseInt(values["automation.terminate_after_days"], settings.TerminateAfterDays)
	settings.TerminateNoticeDays = parseInt(values["automation.terminate_notice_days"], settings.TerminateNoticeDays)
	settings.InvoiceReminderOn = parseBool(values["automation.invoice_reminder_on"], settings.InvoiceReminderOn)
	settings.InvoiceReminderDays = parseString(values["automation.invoice_reminder_days"], settings.InvoiceReminderDays)
	settings.CreditReminderOn = parseBool(values["automation.credit_reminder_on"], settings.CreditReminderOn)
	settings.CreditReminderDays = parseString(values["automation.credit_reminder_days"], settings.CreditReminderDays)
	settings.TicketAutoCloseOn = parseBool(values["automation.ticket_auto_close_on"], settings.TicketAutoCloseOn)
	settings.TicketAutoCloseHours = parseInt(values["automation.ticket_auto_close_hours"], settings.TicketAutoCloseHours)
	settings.LogRetentionDays = parseInt(values["automation.log_retention_days"], settings.LogRetentionDays)
	return settings
}

func (service *Service) SaveSettings(settings domain.AutomationSettings) (domain.AutomationSettings, error) {
	normalized := normalizeSettings(settings)
	if service == nil || service.db == nil {
		return normalized, nil
	}

	tx, err := service.db.Begin()
	if err != nil {
		return domain.AutomationSettings{}, err
	}
	defer func() { _ = tx.Rollback() }()

	entries := map[string]string{
		"automation.auto_provision_enabled":  strconv.FormatBool(normalized.AutoProvisionEnabled),
		"automation.auto_suspend_enabled":    strconv.FormatBool(normalized.AutoSuspendEnabled),
		"automation.suspend_after_days":      strconv.Itoa(normalized.SuspendAfterDays),
		"automation.suspend_notice_days":     strconv.Itoa(normalized.SuspendNoticeDays),
		"automation.auto_terminate_enabled":  strconv.FormatBool(normalized.AutoTerminateEnabled),
		"automation.terminate_after_days":    strconv.Itoa(normalized.TerminateAfterDays),
		"automation.terminate_notice_days":   strconv.Itoa(normalized.TerminateNoticeDays),
		"automation.invoice_reminder_on":     strconv.FormatBool(normalized.InvoiceReminderOn),
		"automation.invoice_reminder_days":   normalized.InvoiceReminderDays,
		"automation.credit_reminder_on":      strconv.FormatBool(normalized.CreditReminderOn),
		"automation.credit_reminder_days":    normalized.CreditReminderDays,
		"automation.ticket_auto_close_on":    strconv.FormatBool(normalized.TicketAutoCloseOn),
		"automation.ticket_auto_close_hours": strconv.Itoa(normalized.TicketAutoCloseHours),
		"automation.log_retention_days":      strconv.Itoa(normalized.LogRetentionDays),
	}

	for key, value := range entries {
		if _, err := tx.Exec(`
INSERT INTO system_configs (config_key, config_value, description)
VALUES (?, ?, 'automation setting')
ON DUPLICATE KEY UPDATE config_value = VALUES(config_value), updated_at = NOW()`,
			key, value,
		); err != nil {
			return domain.AutomationSettings{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return domain.AutomationSettings{}, err
	}
	return normalized, nil
}

func (service *Service) GetSystemConfig(key, fallback string) string {
	if service == nil || service.db == nil {
		return fallback
	}

	var value string
	err := service.db.QueryRow(`
SELECT config_value
FROM system_configs
WHERE config_key = ?
LIMIT 1`, strings.TrimSpace(key)).Scan(&value)
	if err != nil {
		return fallback
	}

	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func (service *Service) SaveSystemConfig(key, value, description string) error {
	if service == nil || service.db == nil {
		return nil
	}

	_, err := service.db.Exec(`
INSERT INTO system_configs (config_key, config_value, description)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE config_value = VALUES(config_value), description = VALUES(description), updated_at = NOW()`,
		strings.TrimSpace(key),
		value,
		valueOrFallback(description, "system config"),
	)
	return err
}

func (service *Service) Start(request StartTaskRequest) domain.Task {
	if service == nil || service.repository == nil {
		return domain.Task{}
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	task := domain.Task{
		TaskNo:             buildTaskNo(),
		TaskType:           strings.TrimSpace(request.TaskType),
		Title:              strings.TrimSpace(request.Title),
		Channel:            strings.TrimSpace(request.Channel),
		Stage:              strings.TrimSpace(request.Stage),
		Status:             domain.TaskStatusRunning,
		SourceType:         strings.TrimSpace(request.SourceType),
		SourceID:           request.SourceID,
		CustomerID:         request.CustomerID,
		CustomerName:       strings.TrimSpace(request.CustomerName),
		ProductName:        strings.TrimSpace(request.ProductName),
		OrderID:            request.OrderID,
		InvoiceID:          request.InvoiceID,
		ServiceID:          request.ServiceID,
		ServiceNo:          strings.TrimSpace(request.ServiceNo),
		ProviderType:       strings.TrimSpace(request.ProviderType),
		ProviderResourceID: strings.TrimSpace(request.ProviderResourceID),
		ActionName:         strings.TrimSpace(request.ActionName),
		OperatorType:       strings.TrimSpace(request.OperatorType),
		OperatorName:       strings.TrimSpace(request.OperatorName),
		RequestPayload:     marshalPayload(request.RequestPayload),
		Message:            strings.TrimSpace(request.Message),
		CreatedAt:          now,
		StartedAt:          now,
		UpdatedAt:          now,
	}
	return service.repository.Create(task)
}

func (service *Service) Complete(taskID int64, request FinishTaskRequest) (domain.Task, bool) {
	if service == nil || service.repository == nil || taskID == 0 {
		return domain.Task{}, false
	}

	task, ok := service.repository.GetByID(taskID)
	if !ok {
		return domain.Task{}, false
	}

	task.Status = request.Status
	task.Message = strings.TrimSpace(request.Message)
	task.ResultPayload = marshalPayload(request.ResultPayload)
	task.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	task.FinishedAt = task.UpdatedAt
	return service.repository.Update(task)
}

func (service *Service) MarkSuccess(taskID int64, message string, result any) {
	_, _ = service.Complete(taskID, FinishTaskRequest{
		Status:        domain.TaskStatusSuccess,
		Message:       message,
		ResultPayload: result,
	})
}

func (service *Service) MarkFailed(taskID int64, message string, result any) {
	_, _ = service.Complete(taskID, FinishTaskRequest{
		Status:        domain.TaskStatusFailed,
		Message:       message,
		ResultPayload: result,
	})
}

func (service *Service) MarkBlocked(taskID int64, message string, result any) {
	_, _ = service.Complete(taskID, FinishTaskRequest{
		Status:        domain.TaskStatusBlocked,
		Message:       message,
		ResultPayload: result,
	})
}

func buildTaskNo() string {
	return fmt.Sprintf("TASK-%s-%04d", time.Now().Format("20060102150405"), time.Now().Nanosecond()%10000)
}

func marshalPayload(value any) string {
	if value == nil {
		return ""
	}
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return ""
	}
	return string(encoded)
}

func defaultSettings() domain.AutomationSettings {
	return domain.AutomationSettings{
		AutoProvisionEnabled: true,
		AutoSuspendEnabled:   true,
		SuspendAfterDays:     3,
		SuspendNoticeDays:    1,
		AutoTerminateEnabled: false,
		TerminateAfterDays:   7,
		TerminateNoticeDays:  3,
		InvoiceReminderOn:    true,
		InvoiceReminderDays:  "7,3,1",
		CreditReminderOn:     true,
		CreditReminderDays:   "3,1",
		TicketAutoCloseOn:    true,
		TicketAutoCloseHours: 72,
		LogRetentionDays:     90,
	}
}

func normalizeSettings(settings domain.AutomationSettings) domain.AutomationSettings {
	normalized := settings
	if normalized.SuspendAfterDays < 0 {
		normalized.SuspendAfterDays = 0
	}
	if normalized.SuspendNoticeDays < 0 {
		normalized.SuspendNoticeDays = 0
	}
	if normalized.TerminateAfterDays < 0 {
		normalized.TerminateAfterDays = 0
	}
	if normalized.TerminateNoticeDays < 0 {
		normalized.TerminateNoticeDays = 0
	}
	if normalized.TicketAutoCloseHours < 0 {
		normalized.TicketAutoCloseHours = 0
	}
	if normalized.LogRetentionDays <= 0 {
		normalized.LogRetentionDays = 90
	}
	normalized.InvoiceReminderDays = strings.TrimSpace(normalized.InvoiceReminderDays)
	if normalized.InvoiceReminderDays == "" {
		normalized.InvoiceReminderDays = "7,3,1"
	}
	normalized.CreditReminderDays = strings.TrimSpace(normalized.CreditReminderDays)
	if normalized.CreditReminderDays == "" {
		normalized.CreditReminderDays = "3,1"
	}
	return normalized
}

func parseBool(value string, fallback bool) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func parseInt(value string, fallback int) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseString(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func valueOrFallback(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}
