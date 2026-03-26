package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"idc-finance/internal/modules/ticket/domain"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (repository *MySQLRepository) List(filter domain.ListFilter) ([]domain.Ticket, int) {
	return repository.list(filter, 0)
}

func (repository *MySQLRepository) ListByCustomer(customerID int64, filter domain.ListFilter) ([]domain.Ticket, int) {
	return repository.list(filter, customerID)
}

func (repository *MySQLRepository) GetByID(id int64) (domain.Ticket, bool) {
	item, err := repository.loadTicket(context.Background(), repository.db, id, 0)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("ticket mysql get failed: %v", err)
		}
		return domain.Ticket{}, false
	}
	return item, true
}

func (repository *MySQLRepository) GetByCustomer(customerID, id int64) (domain.Ticket, bool) {
	item, err := repository.loadTicket(context.Background(), repository.db, id, customerID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("ticket mysql get by customer failed: %v", err)
		}
		return domain.Ticket{}, false
	}
	return item, true
}

func (repository *MySQLRepository) ListReplies(ticketID int64) []domain.Reply {
	rows, err := repository.db.Query(`
SELECT id, ticket_id, author_type, author_id, author_name, content, is_internal,
       IFNULL(DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s'), '')
FROM ticket_replies
WHERE ticket_id = ?
ORDER BY created_at ASC, id ASC`, ticketID)
	if err != nil {
		log.Printf("ticket mysql list replies failed: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]domain.Reply, 0)
	for rows.Next() {
		item, err := scanReply(rows)
		if err != nil {
			log.Printf("ticket mysql scan reply failed: %v", err)
			return result
		}
		result = append(result, item)
	}
	return result
}

func (repository *MySQLRepository) Create(input domain.CreateInput) (domain.Ticket, error) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Ticket{}, err
	}
	defer rollbackTicketTx(tx)

	if !repository.customerExists(ctx, tx, input.CustomerID) {
		return domain.Ticket{}, ErrCustomerNotFound
	}
	if input.ServiceID > 0 && !repository.serviceExists(ctx, tx, input.ServiceID, input.CustomerID) {
		return domain.Ticket{}, ErrServiceNotFound
	}

	tempNo := fmt.Sprintf("TMP-TICKET-%d", time.Now().UnixNano())
	result, err := tx.ExecContext(ctx, `
INSERT INTO tickets (
  customer_id, ticket_no, service_id, title, content, status, priority, source, department_name,
  assigned_admin_id, assigned_admin_name, latest_reply_excerpt, last_reply_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`,
		input.CustomerID,
		tempNo,
		nullInt64(input.ServiceID),
		input.Title,
		input.Content,
		string(input.Status),
		string(input.Priority),
		valueOrDefault(input.Source, "PORTAL"),
		valueOrDefault(input.DepartmentName, "Technical Support"),
		nullInt64(input.AssignedAdminID),
		input.AssignedAdmin,
		excerptContent(input.Content),
	)
	if err != nil {
		return domain.Ticket{}, err
	}

	ticketID, err := result.LastInsertId()
	if err != nil {
		return domain.Ticket{}, err
	}

	ticketNo := fmt.Sprintf("TIC-%08d", ticketID)
	if _, err := tx.ExecContext(ctx, `UPDATE tickets SET ticket_no = ? WHERE id = ?`, ticketNo, ticketID); err != nil {
		return domain.Ticket{}, err
	}

	if err := tx.Commit(); err != nil {
		return domain.Ticket{}, err
	}

	created, ok := repository.GetByID(ticketID)
	if !ok {
		return domain.Ticket{}, ErrTicketNotFound
	}
	return created, nil
}

func (repository *MySQLRepository) Update(id int64, input domain.UpdateInput) (domain.Ticket, bool, error) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Ticket{}, false, err
	}
	defer rollbackTicketTx(tx)

	current, err := repository.loadTicket(ctx, tx, id, 0)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Ticket{}, false, nil
		}
		return domain.Ticket{}, false, err
	}

	title := current.Title
	if input.Title != nil {
		title = strings.TrimSpace(*input.Title)
	}

	status := current.Status
	if input.Status != nil {
		status = *input.Status
	}

	priority := current.Priority
	if input.Priority != nil {
		priority = *input.Priority
	}

	departmentName := current.DepartmentName
	if input.DepartmentName != nil {
		departmentName = strings.TrimSpace(*input.DepartmentName)
	}

	assignedAdminID := current.AssignedAdminID
	if input.AssignedAdminID != nil {
		assignedAdminID = *input.AssignedAdminID
	}

	assignedAdminName := current.AssignedAdminName
	if input.AssignedAdminID != nil && *input.AssignedAdminID == 0 {
		assignedAdminName = ""
	}
	if input.AssignedAdminName != nil {
		assignedAdminName = strings.TrimSpace(*input.AssignedAdminName)
	}

	var closedAt any
	if status == domain.StatusClosed {
		if strings.TrimSpace(current.ClosedAt) != "" {
			closedAt = current.ClosedAt
		} else {
			closedAt = time.Now()
		}
	}

	_, err = tx.ExecContext(ctx, `
UPDATE tickets
SET title = ?, status = ?, priority = ?, department_name = ?,
    assigned_admin_id = ?, assigned_admin_name = ?, closed_at = ?, updated_at = NOW()
WHERE id = ?`,
		title,
		string(status),
		string(priority),
		departmentName,
		nullInt64(assignedAdminID),
		assignedAdminName,
		closedAt,
		id,
	)
	if err != nil {
		return domain.Ticket{}, false, err
	}

	if err := tx.Commit(); err != nil {
		return domain.Ticket{}, false, err
	}

	updated, ok := repository.GetByID(id)
	return updated, ok, nil
}

func (repository *MySQLRepository) AddReply(ticketID int64, input domain.ReplyInput) (domain.Reply, domain.Ticket, bool, error) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Reply{}, domain.Ticket{}, false, err
	}
	defer rollbackTicketTx(tx)

	current, err := repository.loadTicket(ctx, tx, ticketID, 0)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Reply{}, domain.Ticket{}, false, nil
		}
		return domain.Reply{}, domain.Ticket{}, false, err
	}
	if current.Status == domain.StatusClosed && input.Status == "" {
		return domain.Reply{}, domain.Ticket{}, false, ErrTicketClosed
	}

	result, err := tx.ExecContext(ctx, `
INSERT INTO ticket_replies (ticket_id, author_type, author_id, author_name, content, is_internal)
VALUES (?, ?, ?, ?, ?, ?)`,
		ticketID,
		string(input.AuthorType),
		input.AuthorID,
		input.AuthorName,
		input.Content,
		input.IsInternal,
	)
	if err != nil {
		return domain.Reply{}, domain.Ticket{}, false, err
	}

	replyID, err := result.LastInsertId()
	if err != nil {
		return domain.Reply{}, domain.Ticket{}, false, err
	}

	status := current.Status
	if input.Status != "" {
		status = input.Status
	}

	var closedAt any
	if status == domain.StatusClosed {
		closedAt = time.Now()
	}

	if _, err := tx.ExecContext(ctx, `
UPDATE tickets
SET status = ?, latest_reply_excerpt = ?, last_reply_at = NOW(), closed_at = ?, updated_at = NOW()
WHERE id = ?`,
		string(status),
		excerptContent(input.Content),
		closedAt,
		ticketID,
	); err != nil {
		return domain.Reply{}, domain.Ticket{}, false, err
	}

	if err := tx.Commit(); err != nil {
		return domain.Reply{}, domain.Ticket{}, false, err
	}

	reply, err := repository.loadReply(context.Background(), repository.db, replyID)
	if err != nil {
		return domain.Reply{}, domain.Ticket{}, false, err
	}
	ticket, ok := repository.GetByID(ticketID)
	return reply, ticket, ok, nil
}

func (repository *MySQLRepository) Close(ticketID int64) (domain.Ticket, bool, error) {
	result, err := repository.db.Exec(`
UPDATE tickets
SET status = ?, closed_at = NOW(), updated_at = NOW()
WHERE id = ?`,
		string(domain.StatusClosed),
		ticketID,
	)
	if err != nil {
		return domain.Ticket{}, false, err
	}
	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		return domain.Ticket{}, false, rowsErr
	}
	if affected == 0 {
		return domain.Ticket{}, false, nil
	}
	updated, ok := repository.GetByID(ticketID)
	return updated, ok, nil
}

func (repository *MySQLRepository) list(filter domain.ListFilter, onlyCustomerID int64) ([]domain.Ticket, int) {
	whereSQL, args := buildTicketWhere(filter, onlyCustomerID)
	countQuery := `
SELECT COUNT(1)
FROM tickets t
LEFT JOIN customers c ON c.id = t.customer_id
LEFT JOIN services s ON s.id = t.service_id
` + whereSQL

	var total int
	if err := repository.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		log.Printf("ticket mysql count failed: %v", err)
		return nil, 0
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := `
SELECT
  t.id,
  t.ticket_no,
  t.customer_id,
  COALESCE(c.customer_no, ''),
  COALESCE(c.name, ''),
  t.service_id,
  COALESCE(s.service_no, ''),
  COALESCE(s.product_name, ''),
  t.title,
  COALESCE(t.content, ''),
  t.status,
  t.priority,
  COALESCE(t.source, ''),
  COALESCE(t.department_name, ''),
  t.assigned_admin_id,
  COALESCE(t.assigned_admin_name, ''),
  COALESCE(t.latest_reply_excerpt, ''),
  IFNULL(DATE_FORMAT(t.last_reply_at, '%Y-%m-%d %H:%i:%s'), ''),
  IFNULL(DATE_FORMAT(t.closed_at, '%Y-%m-%d %H:%i:%s'), ''),
  IFNULL(DATE_FORMAT(t.created_at, '%Y-%m-%d %H:%i:%s'), ''),
  IFNULL(DATE_FORMAT(t.updated_at, '%Y-%m-%d %H:%i:%s'), '')
FROM tickets t
LEFT JOIN customers c ON c.id = t.customer_id
LEFT JOIN services s ON s.id = t.service_id
` + whereSQL + `
ORDER BY t.updated_at DESC, t.id DESC
LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := repository.db.Query(query, args...)
	if err != nil {
		log.Printf("ticket mysql list failed: %v", err)
		return nil, total
	}
	defer rows.Close()

	result := make([]domain.Ticket, 0)
	for rows.Next() {
		item, err := scanTicket(rows)
		if err != nil {
			log.Printf("ticket mysql scan list failed: %v", err)
			return result, total
		}
		result = append(result, item)
	}
	return result, total
}

func (repository *MySQLRepository) loadTicket(ctx context.Context, queryer ticketQueryer, id, customerID int64) (domain.Ticket, error) {
	query := `
SELECT
  t.id,
  t.ticket_no,
  t.customer_id,
  COALESCE(c.customer_no, ''),
  COALESCE(c.name, ''),
  t.service_id,
  COALESCE(s.service_no, ''),
  COALESCE(s.product_name, ''),
  t.title,
  COALESCE(t.content, ''),
  t.status,
  t.priority,
  COALESCE(t.source, ''),
  COALESCE(t.department_name, ''),
  t.assigned_admin_id,
  COALESCE(t.assigned_admin_name, ''),
  COALESCE(t.latest_reply_excerpt, ''),
  IFNULL(DATE_FORMAT(t.last_reply_at, '%Y-%m-%d %H:%i:%s'), ''),
  IFNULL(DATE_FORMAT(t.closed_at, '%Y-%m-%d %H:%i:%s'), ''),
  IFNULL(DATE_FORMAT(t.created_at, '%Y-%m-%d %H:%i:%s'), ''),
  IFNULL(DATE_FORMAT(t.updated_at, '%Y-%m-%d %H:%i:%s'), '')
FROM tickets t
LEFT JOIN customers c ON c.id = t.customer_id
LEFT JOIN services s ON s.id = t.service_id
WHERE t.id = ?`

	args := []any{id}
	if customerID > 0 {
		query += ` AND t.customer_id = ?`
		args = append(args, customerID)
	}

	row := queryer.QueryRowContext(ctx, query, args...)
	return scanTicketRow(row)
}

func (repository *MySQLRepository) loadReply(ctx context.Context, queryer ticketQueryer, id int64) (domain.Reply, error) {
	row := queryer.QueryRowContext(ctx, `
SELECT id, ticket_id, author_type, author_id, author_name, content, is_internal,
       IFNULL(DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s'), '')
FROM ticket_replies
WHERE id = ?`, id)

	reply, err := scanReplyRow(row)
	if err != nil {
		return domain.Reply{}, err
	}
	return reply, nil
}

func (repository *MySQLRepository) customerExists(ctx context.Context, queryer ticketQueryer, customerID int64) bool {
	row := queryer.QueryRowContext(ctx, `SELECT id FROM customers WHERE id = ?`, customerID)
	var id int64
	return row.Scan(&id) == nil
}

func (repository *MySQLRepository) serviceExists(ctx context.Context, queryer ticketQueryer, serviceID, customerID int64) bool {
	row := queryer.QueryRowContext(ctx, `SELECT id FROM services WHERE id = ? AND customer_id = ?`, serviceID, customerID)
	var id int64
	return row.Scan(&id) == nil
}

func buildTicketWhere(filter domain.ListFilter, onlyCustomerID int64) (string, []any) {
	conditions := []string{"WHERE 1=1"}
	args := make([]any, 0)

	if onlyCustomerID > 0 {
		conditions = append(conditions, "AND t.customer_id = ?")
		args = append(args, onlyCustomerID)
	}
	if filter.CustomerID > 0 {
		conditions = append(conditions, "AND t.customer_id = ?")
		args = append(args, filter.CustomerID)
	}
	if filter.ServiceID > 0 {
		conditions = append(conditions, "AND t.service_id = ?")
		args = append(args, filter.ServiceID)
	}
	if value := strings.TrimSpace(filter.Department); value != "" {
		conditions = append(conditions, "AND t.department_name = ?")
		args = append(args, value)
	}
	if filter.AdminID > 0 {
		conditions = append(conditions, "AND t.assigned_admin_id = ?")
		args = append(args, filter.AdminID)
	}
	if value := strings.TrimSpace(filter.Status); value != "" {
		conditions = append(conditions, "AND t.status = ?")
		args = append(args, value)
	}
	if value := strings.TrimSpace(filter.Priority); value != "" {
		conditions = append(conditions, "AND t.priority = ?")
		args = append(args, value)
	}
	if value := strings.TrimSpace(filter.Keyword); value != "" {
		keyword := "%" + value + "%"
		conditions = append(conditions, `AND (
  t.ticket_no LIKE ? OR
  t.title LIKE ? OR
  c.name LIKE ? OR
  s.service_no LIKE ? OR
  s.product_name LIKE ?
)`)
		args = append(args, keyword, keyword, keyword, keyword, keyword)
	}

	return strings.Join(conditions, "\n"), args
}

func scanTicket(rows *sql.Rows) (domain.Ticket, error) {
	var (
		item            domain.Ticket
		status          string
		priority        string
		serviceID       sql.NullInt64
		assignedAdminID sql.NullInt64
	)
	err := rows.Scan(
		&item.ID,
		&item.TicketNo,
		&item.CustomerID,
		&item.CustomerNo,
		&item.CustomerName,
		&serviceID,
		&item.ServiceNo,
		&item.ProductName,
		&item.Title,
		&item.Content,
		&status,
		&priority,
		&item.Source,
		&item.DepartmentName,
		&assignedAdminID,
		&item.AssignedAdminName,
		&item.LatestReplyExcerpt,
		&item.LastReplyAt,
		&item.ClosedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return domain.Ticket{}, err
	}
	item.ServiceID = serviceID.Int64
	item.AssignedAdminID = assignedAdminID.Int64
	item.Status = domain.Status(status)
	item.Priority = domain.Priority(priority)
	return item, nil
}

func scanTicketRow(row *sql.Row) (domain.Ticket, error) {
	var (
		item            domain.Ticket
		status          string
		priority        string
		serviceID       sql.NullInt64
		assignedAdminID sql.NullInt64
	)
	err := row.Scan(
		&item.ID,
		&item.TicketNo,
		&item.CustomerID,
		&item.CustomerNo,
		&item.CustomerName,
		&serviceID,
		&item.ServiceNo,
		&item.ProductName,
		&item.Title,
		&item.Content,
		&status,
		&priority,
		&item.Source,
		&item.DepartmentName,
		&assignedAdminID,
		&item.AssignedAdminName,
		&item.LatestReplyExcerpt,
		&item.LastReplyAt,
		&item.ClosedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return domain.Ticket{}, err
	}
	item.ServiceID = serviceID.Int64
	item.AssignedAdminID = assignedAdminID.Int64
	item.Status = domain.Status(status)
	item.Priority = domain.Priority(priority)
	return item, nil
}

func scanReply(rows *sql.Rows) (domain.Reply, error) {
	var (
		item       domain.Reply
		authorType string
		isInternal bool
	)
	err := rows.Scan(
		&item.ID,
		&item.TicketID,
		&authorType,
		&item.AuthorID,
		&item.AuthorName,
		&item.Content,
		&isInternal,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return domain.Reply{}, err
	}
	item.AuthorType = domain.AuthorType(authorType)
	item.IsInternal = isInternal
	return item, nil
}

func scanReplyRow(row *sql.Row) (domain.Reply, error) {
	var (
		item       domain.Reply
		authorType string
		isInternal bool
	)
	err := row.Scan(
		&item.ID,
		&item.TicketID,
		&authorType,
		&item.AuthorID,
		&item.AuthorName,
		&item.Content,
		&isInternal,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return domain.Reply{}, err
	}
	item.AuthorType = domain.AuthorType(authorType)
	item.IsInternal = isInternal
	return item, nil
}

func nullInt64(value int64) any {
	if value <= 0 {
		return nil
	}
	return value
}

func valueOrDefault(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func excerptContent(content string) string {
	content = strings.TrimSpace(content)
	if len(content) <= 255 {
		return content
	}
	return content[:255]
}

type ticketQueryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func rollbackTicketTx(tx *sql.Tx) {
	_ = tx.Rollback()
}
