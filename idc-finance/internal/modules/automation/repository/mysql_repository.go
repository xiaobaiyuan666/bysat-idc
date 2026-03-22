package repository

import (
	"database/sql"
	"log"
	"strings"

	"idc-finance/internal/modules/automation/domain"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (repository *MySQLRepository) List(filter domain.TaskFilter) []domain.Task {
	query := `
SELECT id, task_no, task_type, title, channel, stage_name, status, source_type, IFNULL(source_id, 0),
       IFNULL(customer_id, 0), IFNULL(customer_name, ''), IFNULL(product_name, ''), IFNULL(order_id, 0),
       IFNULL(invoice_id, 0), IFNULL(service_id, 0), IFNULL(service_no, ''), IFNULL(provider_type, ''),
       IFNULL(provider_resource_id, ''), IFNULL(action_name, ''), IFNULL(operator_type, ''), IFNULL(operator_name, ''),
       IFNULL(request_payload, ''), IFNULL(result_payload, ''), IFNULL(message, ''),
       DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'),
       IFNULL(DATE_FORMAT(started_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(finished_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s'), '')
FROM automation_tasks
WHERE 1 = 1`

	args := make([]any, 0, 4)
	if strings.TrimSpace(filter.Status) != "" {
		query += ` AND status = ?`
		args = append(args, filter.Status)
	}
	if strings.TrimSpace(filter.TaskType) != "" {
		query += ` AND task_type = ?`
		args = append(args, filter.TaskType)
	}
	if strings.TrimSpace(filter.Channel) != "" {
		query += ` AND channel = ?`
		args = append(args, filter.Channel)
	}
	if strings.TrimSpace(filter.Stage) != "" {
		query += ` AND stage_name = ?`
		args = append(args, filter.Stage)
	}
	if strings.TrimSpace(filter.SourceType) != "" {
		query += ` AND source_type = ?`
		args = append(args, filter.SourceType)
	}
	if filter.SourceID > 0 {
		query += ` AND source_id = ?`
		args = append(args, filter.SourceID)
	}
	if filter.OrderID > 0 {
		query += ` AND order_id = ?`
		args = append(args, filter.OrderID)
	}
	if filter.InvoiceID > 0 {
		query += ` AND invoice_id = ?`
		args = append(args, filter.InvoiceID)
	}
	if filter.ServiceID > 0 {
		query += ` AND service_id = ?`
		args = append(args, filter.ServiceID)
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		query += ` AND (task_no LIKE ? OR title LIKE ? OR service_no LIKE ? OR customer_name LIKE ?)`
		like := "%" + keyword + "%"
		args = append(args, like, like, like, like)
	}
	query += ` ORDER BY id DESC`
	if filter.Limit > 0 {
		query += ` LIMIT ?`
		args = append(args, filter.Limit)
	}

	rows, err := repository.db.Query(query, args...)
	if err != nil {
		log.Printf("automation mysql list tasks failed: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]domain.Task, 0)
	for rows.Next() {
		var item domain.Task
		var status string
		if err := rows.Scan(
			&item.ID,
			&item.TaskNo,
			&item.TaskType,
			&item.Title,
			&item.Channel,
			&item.Stage,
			&status,
			&item.SourceType,
			&item.SourceID,
			&item.CustomerID,
			&item.CustomerName,
			&item.ProductName,
			&item.OrderID,
			&item.InvoiceID,
			&item.ServiceID,
			&item.ServiceNo,
			&item.ProviderType,
			&item.ProviderResourceID,
			&item.ActionName,
			&item.OperatorType,
			&item.OperatorName,
			&item.RequestPayload,
			&item.ResultPayload,
			&item.Message,
			&item.CreatedAt,
			&item.StartedAt,
			&item.FinishedAt,
			&item.UpdatedAt,
		); err != nil {
			log.Printf("automation mysql scan task failed: %v", err)
			return result
		}
		item.Status = domain.TaskStatus(status)
		result = append(result, item)
	}
	return result
}

func (repository *MySQLRepository) GetByID(id int64) (domain.Task, bool) {
	row := repository.db.QueryRow(`
SELECT id, task_no, task_type, title, channel, stage_name, status, source_type, IFNULL(source_id, 0),
       IFNULL(customer_id, 0), IFNULL(customer_name, ''), IFNULL(product_name, ''), IFNULL(order_id, 0),
       IFNULL(invoice_id, 0), IFNULL(service_id, 0), IFNULL(service_no, ''), IFNULL(provider_type, ''),
       IFNULL(provider_resource_id, ''), IFNULL(action_name, ''), IFNULL(operator_type, ''), IFNULL(operator_name, ''),
       IFNULL(request_payload, ''), IFNULL(result_payload, ''), IFNULL(message, ''),
       DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'),
       IFNULL(DATE_FORMAT(started_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(finished_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s'), '')
FROM automation_tasks
WHERE id = ?`, id)

	var item domain.Task
	var status string
	if err := row.Scan(
		&item.ID,
		&item.TaskNo,
		&item.TaskType,
		&item.Title,
		&item.Channel,
		&item.Stage,
		&status,
		&item.SourceType,
		&item.SourceID,
		&item.CustomerID,
		&item.CustomerName,
		&item.ProductName,
		&item.OrderID,
		&item.InvoiceID,
		&item.ServiceID,
		&item.ServiceNo,
		&item.ProviderType,
		&item.ProviderResourceID,
		&item.ActionName,
		&item.OperatorType,
		&item.OperatorName,
		&item.RequestPayload,
		&item.ResultPayload,
		&item.Message,
		&item.CreatedAt,
		&item.StartedAt,
		&item.FinishedAt,
		&item.UpdatedAt,
	); err != nil {
		return domain.Task{}, false
	}
	item.Status = domain.TaskStatus(status)
	return item, true
}

func (repository *MySQLRepository) Create(task domain.Task) domain.Task {
	result, err := repository.db.Exec(`
INSERT INTO automation_tasks (
  task_no, task_type, title, channel, stage_name, status, source_type, source_id,
  customer_id, customer_name, product_name, order_id, invoice_id, service_id, service_no,
  provider_type, provider_resource_id, action_name, operator_type, operator_name,
  request_payload, result_payload, message, started_at, finished_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		task.TaskNo,
		task.TaskType,
		task.Title,
		task.Channel,
		task.Stage,
		task.Status,
		task.SourceType,
		nullIfZero(task.SourceID),
		nullIfZero(task.CustomerID),
		task.CustomerName,
		task.ProductName,
		nullIfZero(task.OrderID),
		nullIfZero(task.InvoiceID),
		nullIfZero(task.ServiceID),
		task.ServiceNo,
		task.ProviderType,
		task.ProviderResourceID,
		task.ActionName,
		task.OperatorType,
		task.OperatorName,
		nullIfEmpty(task.RequestPayload),
		nullIfEmpty(task.ResultPayload),
		task.Message,
		nullTime(task.StartedAt),
		nullTime(task.FinishedAt),
	)
	if err != nil {
		log.Printf("automation mysql create task failed: %v", err)
		return domain.Task{}
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		return domain.Task{}
	}
	created, _ := repository.GetByID(taskID)
	return created
}

func (repository *MySQLRepository) Update(task domain.Task) (domain.Task, bool) {
	_, err := repository.db.Exec(`
UPDATE automation_tasks
SET task_type = ?, title = ?, channel = ?, stage_name = ?, status = ?, source_type = ?, source_id = ?,
    customer_id = ?, customer_name = ?, product_name = ?, order_id = ?, invoice_id = ?, service_id = ?, service_no = ?,
    provider_type = ?, provider_resource_id = ?, action_name = ?, operator_type = ?, operator_name = ?,
    request_payload = ?, result_payload = ?, message = ?, started_at = ?, finished_at = ?, updated_at = NOW()
WHERE id = ?`,
		task.TaskType,
		task.Title,
		task.Channel,
		task.Stage,
		task.Status,
		task.SourceType,
		nullIfZero(task.SourceID),
		nullIfZero(task.CustomerID),
		task.CustomerName,
		task.ProductName,
		nullIfZero(task.OrderID),
		nullIfZero(task.InvoiceID),
		nullIfZero(task.ServiceID),
		task.ServiceNo,
		task.ProviderType,
		task.ProviderResourceID,
		task.ActionName,
		task.OperatorType,
		task.OperatorName,
		nullIfEmpty(task.RequestPayload),
		nullIfEmpty(task.ResultPayload),
		task.Message,
		nullTime(task.StartedAt),
		nullTime(task.FinishedAt),
		task.ID,
	)
	if err != nil {
		log.Printf("automation mysql update task failed: %v", err)
		return domain.Task{}, false
	}
	updated, ok := repository.GetByID(task.ID)
	return updated, ok
}

func nullIfZero(value int64) any {
	if value == 0 {
		return nil
	}
	return value
}

func nullIfEmpty(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}

func nullTime(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}
