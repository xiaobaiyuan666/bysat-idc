package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"idc-finance/internal/modules/order/domain"
)

type MySQLRepository struct {
	db *sql.DB
}

type serviceChangeLink struct {
	ID          int64
	ServiceID   int64
	OrderID     int64
	InvoiceID   int64
	ActionName  string
	Title       string
	Status      string
	Reason      string
	PayloadJSON string
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func assignNullInt64(target *int64, value sql.NullInt64) {
	if value.Valid {
		*target = value.Int64
		return
	}
	*target = 0
}

func (repository *MySQLRepository) ListOrders(filter domain.OrderListFilter) ([]domain.Order, int) {
	return repository.listOrdersWithFilter(filter)
}

func (repository *MySQLRepository) ListOrdersByCustomer(customerID int64) []domain.Order {
	items, _ := repository.listOrdersWithFilter(domain.OrderListFilter{
		CustomerID: customerID,
		Page:       1,
		Limit:      0,
		Sort:       "create_time",
		Order:      "desc",
	})
	return items
}

func (repository *MySQLRepository) GetOrderByID(id int64) (domain.Order, bool) {
	item, err := repository.loadOrder(context.Background(), repository.db, id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("order mysql get order failed: %v", err)
		}
		return domain.Order{}, false
	}
	return item, true
}

func (repository *MySQLRepository) ListInvoices(filter domain.InvoiceListFilter) ([]domain.Invoice, int) {
	return repository.listInvoicesWithFilter(filter)
}

func (repository *MySQLRepository) ListInvoicesByCustomer(customerID int64) []domain.Invoice {
	return repository.listInvoicesBy("customer", customerID)
}

func (repository *MySQLRepository) ListInvoicesByOrder(orderID int64) []domain.Invoice {
	return repository.listInvoicesBy("order", orderID)
}

func (repository *MySQLRepository) GetInvoiceByID(id int64) (domain.Invoice, bool) {
	item, err := repository.loadInvoice(context.Background(), repository.db, id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("order mysql get invoice failed: %v", err)
		}
		return domain.Invoice{}, false
	}
	return item, true
}

func (repository *MySQLRepository) GetCustomerWallet(customerID int64) (domain.CustomerWallet, bool) {
	item, err := repository.loadCustomerWallet(context.Background(), repository.db, customerID, false)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("order mysql get customer wallet failed: %v", err)
		}
		return domain.CustomerWallet{}, false
	}
	return item, true
}

func (repository *MySQLRepository) ListAccountTransactions(filter domain.AccountTransactionFilter) ([]domain.AccountTransaction, int) {
	whereSQL, args := buildAccountTransactionFilterSQL(filter)

	countQuery := `
SELECT COUNT(*)
FROM account_transactions at
LEFT JOIN customers c ON c.id = at.customer_id
LEFT JOIN orders o ON o.id = at.order_id
LEFT JOIN invoices i ON i.id = at.invoice_id
LEFT JOIN payments p ON p.id = at.payment_id
LEFT JOIN refunds r ON r.id = at.refund_id` + whereSQL

	var total int
	if err := repository.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		log.Printf("order mysql count account transactions failed: %v", err)
		return nil, 0
	}

	sortColumn := resolveAccountTransactionSortColumn(filter.Sort)
	sortDirection := resolveSortDirection(filter.Order)
	query := `
SELECT at.id
FROM account_transactions at
LEFT JOIN customers c ON c.id = at.customer_id
LEFT JOIN orders o ON o.id = at.order_id
LEFT JOIN invoices i ON i.id = at.invoice_id
LEFT JOIN payments p ON p.id = at.payment_id
LEFT JOIN refunds r ON r.id = at.refund_id` + whereSQL + ` ORDER BY ` + sortColumn + ` ` + sortDirection + `, at.id DESC`

	queryArgs := append([]any{}, args...)
	if filter.Limit > 0 {
		page := filter.Page
		if page <= 0 {
			page = 1
		}
		offset := (page - 1) * filter.Limit
		query += ` LIMIT ? OFFSET ?`
		queryArgs = append(queryArgs, filter.Limit, offset)
	}

	rows, err := repository.db.Query(query, queryArgs...)
	if err != nil {
		log.Printf("order mysql list account transactions failed: %v", err)
		return nil, 0
	}
	defer rows.Close()

	items := make([]domain.AccountTransaction, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Printf("order mysql scan account transaction id failed: %v", err)
			return items, total
		}
		item, loadErr := repository.loadAccountTransaction(context.Background(), repository.db, id)
		if loadErr != nil {
			log.Printf("order mysql load account transaction failed: %v", loadErr)
			continue
		}
		items = append(items, item)
	}

	return items, total
}

func (repository *MySQLRepository) ListAccountTransactionsByCustomer(customerID int64, limit int) []domain.AccountTransaction {
	items, _ := repository.ListAccountTransactions(domain.AccountTransactionFilter{
		CustomerID: customerID,
		Page:       1,
		Limit:      limit,
		Sort:       "occurred_at",
		Order:      "desc",
	})
	return items
}

func (repository *MySQLRepository) AdjustCustomerWallet(
	customerID int64,
	input domain.WalletAdjustment,
) (domain.CustomerWallet, domain.AccountTransaction, bool, error) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.CustomerWallet{}, domain.AccountTransaction{}, false, err
	}
	defer rollbackQuietly(tx)

	wallet, err := repository.loadCustomerWallet(ctx, tx, customerID, true)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.CustomerWallet{}, domain.AccountTransaction{}, false, nil
		}
		return domain.CustomerWallet{}, domain.AccountTransaction{}, false, err
	}

	transaction, updatedWallet, err := repository.applyWalletAdjustmentMySQL(ctx, tx, wallet, input, 0, "", 0, "", 0, "", 0, "")
	if err != nil {
		return domain.CustomerWallet{}, domain.AccountTransaction{}, false, err
	}

	if err := tx.Commit(); err != nil {
		return domain.CustomerWallet{}, domain.AccountTransaction{}, false, err
	}

	return updatedWallet, transaction, true, nil
}

func (repository *MySQLRepository) GetServiceChangeOrderByInvoiceID(invoiceID int64) (domain.ServiceChangeOrder, bool) {
	row := repository.db.QueryRow(`
SELECT sco.id, sco.service_id, sco.order_id, sco.invoice_id, sco.action_name, sco.title, sco.amount, sco.status, sco.reason, sco.payload_json,
       IFNULL(i.billing_cycle, ''),
       IFNULL(DATE_FORMAT(sco.paid_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(sco.refunded_at, '%Y-%m-%d %H:%i:%s'), ''),
       DATE_FORMAT(sco.created_at, '%Y-%m-%d %H:%i:%s'),
       DATE_FORMAT(sco.updated_at, '%Y-%m-%d %H:%i:%s')
FROM service_change_orders sco
LEFT JOIN invoices i ON i.id = sco.invoice_id
WHERE sco.invoice_id = ?`, invoiceID)

	var (
		item        domain.ServiceChangeOrder
		payloadJSON sql.NullString
	)
	if err := row.Scan(
		&item.ID,
		&item.ServiceID,
		&item.OrderID,
		&item.InvoiceID,
		&item.ActionName,
		&item.Title,
		&item.Amount,
		&item.Status,
		&item.Reason,
		&payloadJSON,
		&item.BillingCycle,
		&item.PaidAt,
		&item.RefundedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if err != sql.ErrNoRows {
			log.Printf("order mysql get service change failed: %v", err)
		}
		return domain.ServiceChangeOrder{}, false
	}

	if payloadJSON.Valid && strings.TrimSpace(payloadJSON.String) != "" {
		var payload map[string]any
		if err := json.Unmarshal([]byte(payloadJSON.String), &payload); err == nil {
			item.Payload = payload
		}
	}
	return item, true
}

func (repository *MySQLRepository) GetServiceChangeOrderByOrderID(orderID int64) (domain.ServiceChangeOrder, bool) {
	row := repository.db.QueryRow(`
SELECT sco.id, sco.service_id, sco.order_id, sco.invoice_id, sco.action_name, sco.title, sco.amount, sco.status, sco.reason, sco.payload_json,
       IFNULL(i.billing_cycle, ''),
       IFNULL(DATE_FORMAT(sco.paid_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(sco.refunded_at, '%Y-%m-%d %H:%i:%s'), ''),
       DATE_FORMAT(sco.created_at, '%Y-%m-%d %H:%i:%s'),
       DATE_FORMAT(sco.updated_at, '%Y-%m-%d %H:%i:%s')
FROM service_change_orders sco
LEFT JOIN invoices i ON i.id = sco.invoice_id
WHERE sco.order_id = ?`, orderID)

	var (
		item        domain.ServiceChangeOrder
		payloadJSON sql.NullString
	)
	if err := row.Scan(
		&item.ID,
		&item.ServiceID,
		&item.OrderID,
		&item.InvoiceID,
		&item.ActionName,
		&item.Title,
		&item.Amount,
		&item.Status,
		&item.Reason,
		&payloadJSON,
		&item.BillingCycle,
		&item.PaidAt,
		&item.RefundedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if err != sql.ErrNoRows {
			log.Printf("order mysql get service change by order failed: %v", err)
		}
		return domain.ServiceChangeOrder{}, false
	}

	if payloadJSON.Valid && strings.TrimSpace(payloadJSON.String) != "" {
		var payload map[string]any
		if err := json.Unmarshal([]byte(payloadJSON.String), &payload); err == nil {
			item.Payload = payload
		}
	}
	return item, true
}

func (repository *MySQLRepository) ListServiceChangeOrdersByService(serviceID int64) []domain.ServiceChangeOrder {
	rows, err := repository.db.Query(`
SELECT sco.id, sco.service_id, sco.order_id, sco.invoice_id, sco.action_name, sco.title, sco.amount, sco.status, sco.reason, sco.payload_json,
       IFNULL(i.billing_cycle, ''),
       IFNULL(DATE_FORMAT(sco.paid_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(sco.refunded_at, '%Y-%m-%d %H:%i:%s'), ''),
       DATE_FORMAT(sco.created_at, '%Y-%m-%d %H:%i:%s'),
       DATE_FORMAT(sco.updated_at, '%Y-%m-%d %H:%i:%s')
FROM service_change_orders sco
LEFT JOIN invoices i ON i.id = sco.invoice_id
WHERE sco.service_id = ?
ORDER BY sco.id DESC`, serviceID)
	if err != nil {
		log.Printf("order mysql list service changes failed: %v", err)
		return nil
	}
	defer rows.Close()

	items := make([]domain.ServiceChangeOrder, 0)
	for rows.Next() {
		var (
			item        domain.ServiceChangeOrder
			payloadJSON sql.NullString
		)
		if err := rows.Scan(
			&item.ID,
			&item.ServiceID,
			&item.OrderID,
			&item.InvoiceID,
			&item.ActionName,
			&item.Title,
			&item.Amount,
			&item.Status,
			&item.Reason,
			&payloadJSON,
			&item.BillingCycle,
			&item.PaidAt,
			&item.RefundedAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			log.Printf("order mysql scan service change failed: %v", err)
			return items
		}
		if payloadJSON.Valid && strings.TrimSpace(payloadJSON.String) != "" {
			var payload map[string]any
			if err := json.Unmarshal([]byte(payloadJSON.String), &payload); err == nil {
				item.Payload = payload
			}
		}
		items = append(items, item)
	}
	return items
}

func (repository *MySQLRepository) UpdatePendingOrder(
	orderID int64,
	update domain.PendingOrderUpdate,
) (domain.Order, *domain.Invoice, bool, error) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Order{}, nil, false, err
	}
	defer rollbackQuietly(tx)

	order, err := repository.loadOrderForUpdate(ctx, tx, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Order{}, nil, false, nil
		}
		return domain.Order{}, nil, false, err
	}
	_, isServiceChange, err := repository.loadServiceChangeLinkByOrderID(ctx, tx, order.ID)
	if err != nil {
		return domain.Order{}, nil, false, err
	}

	productName := firstNonEmptyString(strings.TrimSpace(update.ProductName), order.ProductName)
	billingCycle := firstNonEmptyString(strings.TrimSpace(update.BillingCycle), order.BillingCycle)
	statusValue, err := normalizeOrderStatusValue(update.Status, string(order.Status))
	if err != nil {
		return domain.Order{}, nil, false, err
	}
	amount := order.Amount
	if update.Amount != nil {
		if *update.Amount < 0 {
			return domain.Order{}, nil, false, fmt.Errorf("订单金额不能小于 0")
		}
		amount = *update.Amount
	}

	var dueAt any = nil
	if strings.TrimSpace(update.DueAt) != "" {
		parsedDueAt, ok := parseMySQLEditableTime(strings.TrimSpace(update.DueAt))
		if !ok {
			return domain.Order{}, nil, false, fmt.Errorf("账单到期时间格式不正确")
		}
		dueAt = parsedDueAt
	}

	invoiceIDs := make([]int64, 0)
	firstPaidInvoiceID := int64(0)
	rows, err := tx.QueryContext(ctx, `
SELECT id, status
FROM invoices
WHERE order_id = ?
FOR UPDATE`, order.ID)
	if err != nil {
		return domain.Order{}, nil, false, err
	}
	for rows.Next() {
		var (
			invoiceID int64
			status    string
		)
		if err := rows.Scan(&invoiceID, &status); err != nil {
			rows.Close()
			return domain.Order{}, nil, false, err
		}
		invoiceIDs = append(invoiceIDs, invoiceID)
		if strings.EqualFold(strings.TrimSpace(status), string(domain.InvoiceStatusPaid)) && firstPaidInvoiceID == 0 {
			firstPaidInvoiceID = invoiceID
		}
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return domain.Order{}, nil, false, err
	}
	rows.Close()

	if _, err := tx.ExecContext(ctx, `
UPDATE orders
SET product_name = ?, billing_cycle = ?, amount = ?, status = ?, updated_at = NOW()
WHERE id = ?`,
		productName,
		billingCycle,
		amount,
		statusValue,
		order.ID,
	); err != nil {
		return domain.Order{}, nil, false, err
	}

	if _, err := tx.ExecContext(ctx, `
UPDATE order_items
SET product_name = ?, billing_cycle = ?, amount = ?
WHERE order_id = ?`,
		productName,
		billingCycle,
		amount,
		order.ID,
	); err != nil {
		return domain.Order{}, nil, false, err
	}

	if len(invoiceIDs) > 0 {
		if dueAt != nil {
			if _, err := tx.ExecContext(ctx, `
UPDATE invoices
SET product_name = ?, billing_cycle = ?, total_amount = ?, due_at = ?, updated_at = NOW()
WHERE order_id = ?`,
				productName,
				billingCycle,
				amount,
				dueAt,
				order.ID,
			); err != nil {
				return domain.Order{}, nil, false, err
			}
		} else {
			if _, err := tx.ExecContext(ctx, `
UPDATE invoices
SET product_name = ?, billing_cycle = ?, total_amount = ?, updated_at = NOW()
WHERE order_id = ?`,
				productName,
				billingCycle,
				amount,
				order.ID,
			); err != nil {
				return domain.Order{}, nil, false, err
			}
		}

		if _, err := tx.ExecContext(ctx, `
UPDATE invoice_items
SET product_name = ?, billing_cycle = ?, amount = ?
WHERE order_id = ?`,
			productName,
			billingCycle,
			amount,
			order.ID,
		); err != nil {
			return domain.Order{}, nil, false, err
		}
	}

	now := time.Now()
	if !isServiceChange {
		switch statusValue {
		case string(domain.OrderStatusCancelled):
			if _, err := tx.ExecContext(ctx, `
UPDATE services
SET status = ?, last_action = ?, sync_status = ?, sync_message = ?, last_sync_at = NOW(), updated_at = NOW()
WHERE order_id = ?`,
				domain.ServiceStatusTerminated,
				"manual-adjust",
				"SUCCESS",
				"后台人工改签订单为已取消",
				order.ID,
			); err != nil {
				return domain.Order{}, nil, false, err
			}
		case string(domain.OrderStatusPending):
			if _, err := tx.ExecContext(ctx, `
UPDATE services
SET status = ?, last_action = ?, sync_status = ?, sync_message = ?, last_sync_at = NOW(), updated_at = NOW()
WHERE order_id = ? AND status <> ?`,
				domain.ServiceStatusSuspended,
				"manual-adjust",
				"SUCCESS",
				"后台人工改签订单为待处理",
				order.ID,
				domain.ServiceStatusTerminated,
			); err != nil {
				return domain.Order{}, nil, false, err
			}
		case string(domain.OrderStatusActive), string(domain.OrderStatusCompleted):
			if firstPaidInvoiceID > 0 {
				paidInvoice, err := repository.loadInvoice(ctx, tx, firstPaidInvoiceID)
				if err != nil {
					return domain.Order{}, nil, false, err
				}
				if _, err := repository.activateOrCreateService(ctx, tx, paidInvoice, now); err != nil {
					return domain.Order{}, nil, false, err
				}
				if _, err := tx.ExecContext(ctx, `
UPDATE services
SET status = ?, last_action = ?, sync_status = ?, sync_message = ?, last_sync_at = NOW(), updated_at = NOW()
WHERE order_id = ? AND status <> ?`,
					domain.ServiceStatusActive,
					"manual-adjust",
					"SUCCESS",
					"后台人工改签订单为生效",
					order.ID,
					domain.ServiceStatusTerminated,
				); err != nil {
					return domain.Order{}, nil, false, err
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return domain.Order{}, nil, false, err
	}

	updatedOrder, ok := repository.GetOrderByID(order.ID)
	if !ok {
		return domain.Order{}, nil, false, fmt.Errorf("订单更新后读取失败")
	}

	var updatedInvoice *domain.Invoice
	if len(invoiceIDs) > 0 {
		if item, exists := repository.GetInvoiceByID(invoiceIDs[0]); exists {
			updatedInvoice = &item
		}
	}

	return updatedOrder, updatedInvoice, true, nil
}

func (repository *MySQLRepository) UpdateUnpaidInvoice(
	invoiceID int64,
	update domain.UnpaidInvoiceUpdate,
) (domain.Invoice, *domain.Order, bool, error) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Invoice{}, nil, false, err
	}
	defer rollbackQuietly(tx)

	var orderID sql.NullInt64
	if err := tx.QueryRowContext(ctx, `SELECT order_id FROM invoices WHERE id = ?`, invoiceID).Scan(&orderID); err != nil {
		if err == sql.ErrNoRows {
			return domain.Invoice{}, nil, false, nil
		}
		return domain.Invoice{}, nil, false, err
	}

	if orderID.Valid && orderID.Int64 > 0 {
		if _, err := repository.loadOrderForUpdate(ctx, tx, orderID.Int64); err != nil && err != sql.ErrNoRows {
			return domain.Invoice{}, nil, false, err
		}
	}

	invoice, err := repository.loadInvoiceForUpdate(ctx, tx, invoiceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Invoice{}, nil, false, nil
		}
		return domain.Invoice{}, nil, false, err
	}

	productName := firstNonEmptyString(strings.TrimSpace(update.ProductName), invoice.ProductName)
	billingCycle := firstNonEmptyString(strings.TrimSpace(update.BillingCycle), invoice.BillingCycle)
	statusValue, err := normalizeInvoiceStatusValue(update.Status, string(invoice.Status))
	if err != nil {
		return domain.Invoice{}, nil, false, err
	}
	amount := invoice.TotalAmount
	if update.Amount != nil {
		if *update.Amount < 0 {
			return domain.Invoice{}, nil, false, fmt.Errorf("调整金额不能小于 0")
		}
		amount = *update.Amount
	}

	var dueAt any = nil
	if strings.TrimSpace(update.DueAt) != "" {
		parsedDueAt, ok := parseMySQLEditableTime(strings.TrimSpace(update.DueAt))
		if !ok {
			return domain.Invoice{}, nil, false, fmt.Errorf("到期时间格式不正确")
		}
		dueAt = parsedDueAt
	}
	previousStatus := string(invoice.Status)
	invoice.ProductName = productName
	invoice.BillingCycle = billingCycle
	invoice.TotalAmount = amount
	invoice.Status = domain.InvoiceStatus(statusValue)
	if dueAt != nil {
		invoice.DueAt = dueAt.(time.Time).Format("2006-01-02 15:04:05")
	}
	switch statusValue {
	case string(domain.InvoiceStatusUnpaid):
		invoice.PaidAt = ""
	case string(domain.InvoiceStatusPaid), string(domain.InvoiceStatusRefunded):
		if parsed, ok := parseMySQLEditableTime(strings.TrimSpace(invoice.PaidAt)); ok {
			invoice.PaidAt = parsed.Format("2006-01-02 15:04:05")
		} else {
			invoice.PaidAt = time.Now().Format("2006-01-02 15:04:05")
		}
	}

	if dueAt != nil {
		if _, err := tx.ExecContext(ctx, `
UPDATE invoices
SET product_name = ?, billing_cycle = ?, total_amount = ?, due_at = ?, status = ?, paid_at = ?, updated_at = NOW()
WHERE id = ?`,
			productName,
			billingCycle,
			amount,
			dueAt,
			statusValue,
			resolveInvoicePaidAt(statusValue, invoice.PaidAt),
			invoice.ID,
		); err != nil {
			return domain.Invoice{}, nil, false, err
		}
	} else {
		if _, err := tx.ExecContext(ctx, `
UPDATE invoices
SET product_name = ?, billing_cycle = ?, total_amount = ?, status = ?, paid_at = ?, updated_at = NOW()
WHERE id = ?`,
			productName,
			billingCycle,
			amount,
			statusValue,
			resolveInvoicePaidAt(statusValue, invoice.PaidAt),
			invoice.ID,
		); err != nil {
			return domain.Invoice{}, nil, false, err
		}
	}

	if _, err := tx.ExecContext(ctx, `
UPDATE invoice_items
SET product_name = ?, billing_cycle = ?, amount = ?
WHERE invoice_id = ?`,
		productName,
		billingCycle,
		amount,
		invoice.ID,
	); err != nil {
		return domain.Invoice{}, nil, false, err
	}

	var updatedOrder *domain.Order
	changeLink, isServiceChange, err := repository.loadServiceChangeLinkByInvoiceID(ctx, tx, invoice.ID)
	if err != nil {
		return domain.Invoice{}, nil, false, err
	}
	if invoice.OrderID > 0 {
		order, err := repository.loadOrder(ctx, tx, invoice.OrderID)
		if err != nil && err != sql.ErrNoRows {
			return domain.Invoice{}, nil, false, err
		}
		if err == nil {
			nextOrderStatus := resolveManualLinkedOrderStatusForInvoice(statusValue, string(order.Status))
			if _, err := tx.ExecContext(ctx, `
UPDATE orders
SET product_name = ?, billing_cycle = ?, amount = ?, status = ?, updated_at = NOW()
WHERE id = ?`,
				productName,
				billingCycle,
				amount,
				nextOrderStatus,
				order.ID,
			); err != nil {
				return domain.Invoice{}, nil, false, err
			}
			if _, err := tx.ExecContext(ctx, `
UPDATE order_items
SET product_name = ?, billing_cycle = ?, amount = ?
WHERE order_id = ?`,
				productName,
				billingCycle,
				amount,
				order.ID,
			); err != nil {
				return domain.Invoice{}, nil, false, err
			}

			now := time.Now()
			if isServiceChange {
				switch statusValue {
				case string(domain.InvoiceStatusPaid):
					if previousStatus != string(domain.InvoiceStatusPaid) {
						if _, ok := repository.loadLatestPayment(ctx, tx, invoice.ID); !ok {
							if _, _, err := repository.insertPayment(ctx, tx, invoice, "MANUAL", "ADMIN", "后台登记收款", "", now); err != nil {
								return domain.Invoice{}, nil, false, err
							}
						}
					}
				case string(domain.InvoiceStatusRefunded):
					if previousStatus != string(domain.InvoiceStatusRefunded) {
						if _, _, err := repository.insertRefund(ctx, tx, invoice, "后台手动退款"); err != nil {
							return domain.Invoice{}, nil, false, err
						}
					}
				}
				if _, err := tx.ExecContext(ctx, `
UPDATE service_change_orders
SET status = ?, updated_at = NOW(),
    paid_at = CASE WHEN ? = 'PAID' THEN NOW() WHEN ? = 'UNPAID' THEN NULL ELSE paid_at END,
    refunded_at = CASE WHEN ? = 'REFUNDED' THEN NOW() WHEN ? = 'UNPAID' THEN NULL ELSE refunded_at END
WHERE id = ?`,
					statusValue,
					statusValue,
					statusValue,
					statusValue,
					statusValue,
					changeLink.ID,
				); err != nil {
					return domain.Invoice{}, nil, false, err
				}
			} else {
				switch statusValue {
				case string(domain.InvoiceStatusPaid):
					if previousStatus != string(domain.InvoiceStatusPaid) {
						if _, ok := repository.loadLatestPayment(ctx, tx, invoice.ID); !ok {
							if _, _, err := repository.insertPayment(ctx, tx, invoice, "MANUAL", "ADMIN", "后台登记收款", "", now); err != nil {
								return domain.Invoice{}, nil, false, err
							}
						}
					}
					if _, err := repository.activateOrCreateService(ctx, tx, invoice, now); err != nil {
						return domain.Invoice{}, nil, false, err
					}
					if _, err := tx.ExecContext(ctx, `
UPDATE services
SET status = ?, last_action = ?, sync_status = ?, sync_message = ?, last_sync_at = NOW(), updated_at = NOW()
WHERE invoice_id = ? AND status <> ?`,
						domain.ServiceStatusActive,
						"manual-adjust",
						"SUCCESS",
						"后台手动标记已支付，服务已激活",
						invoice.ID,
						domain.ServiceStatusTerminated,
					); err != nil {
						return domain.Invoice{}, nil, false, err
					}
				case string(domain.InvoiceStatusUnpaid):
					if _, err := tx.ExecContext(ctx, `
UPDATE services
SET status = ?, last_action = ?, sync_status = ?, sync_message = ?, last_sync_at = NOW(), updated_at = NOW()
WHERE invoice_id = ? AND status <> ?`,
						domain.ServiceStatusSuspended,
						"manual-adjust",
						"SUCCESS",
						"后台手动标记待支付，服务已暂停",
						invoice.ID,
						domain.ServiceStatusTerminated,
					); err != nil {
						return domain.Invoice{}, nil, false, err
					}
				case string(domain.InvoiceStatusRefunded):
					if previousStatus != string(domain.InvoiceStatusRefunded) {
						if _, _, err := repository.insertRefund(ctx, tx, invoice, "后台手动退款"); err != nil {
							return domain.Invoice{}, nil, false, err
						}
					}
					if _, err := tx.ExecContext(ctx, `
UPDATE services
SET status = ?, last_action = ?, sync_status = ?, sync_message = ?, last_sync_at = NOW(), updated_at = NOW()
WHERE invoice_id = ?`,
						domain.ServiceStatusTerminated,
						"refund-terminate",
						"SUCCESS",
						"后台手动退款，服务已终止",
						invoice.ID,
					); err != nil {
						return domain.Invoice{}, nil, false, err
					}
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return domain.Invoice{}, nil, false, err
	}

	updatedInvoice, ok := repository.GetInvoiceByID(invoice.ID)
	if !ok {
		return domain.Invoice{}, nil, false, fmt.Errorf("更新后的账单不存在")
	}

	if updatedInvoice.OrderID > 0 {
		if item, exists := repository.GetOrderByID(updatedInvoice.OrderID); exists {
			updatedOrder = &item
		}
	}

	return updatedInvoice, updatedOrder, true, nil
}

func (repository *MySQLRepository) ListServices(filter domain.ServiceListFilter) ([]domain.ServiceRecord, int) {
	return repository.listServicesWithFilter(filter)
}

func (repository *MySQLRepository) ListServicesByCustomer(customerID int64) []domain.ServiceRecord {
	return repository.listServicesBy("customer", customerID)
}

func (repository *MySQLRepository) ListServicesByOrder(orderID int64) []domain.ServiceRecord {
	return repository.listServicesBy("order", orderID)
}

func (repository *MySQLRepository) GetServiceByID(id int64) (domain.ServiceRecord, bool) {
	item, err := repository.loadService(context.Background(), repository.db, id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("order mysql get service failed: %v", err)
		}
		return domain.ServiceRecord{}, false
	}
	return item, true
}

func (repository *MySQLRepository) UpdateServiceRecord(
	serviceID int64,
	update domain.ManualServiceUpdate,
) (domain.ServiceRecord, bool, error) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.ServiceRecord{}, false, err
	}
	defer rollbackQuietly(tx)

	var lockedID int64
	if err := tx.QueryRowContext(ctx, `SELECT id FROM services WHERE id = ? FOR UPDATE`, serviceID).Scan(&lockedID); err != nil {
		if err == sql.ErrNoRows {
			return domain.ServiceRecord{}, false, nil
		}
		return domain.ServiceRecord{}, false, err
	}

	serviceRecord, err := repository.loadService(ctx, tx, serviceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ServiceRecord{}, false, nil
		}
		return domain.ServiceRecord{}, false, err
	}

	nextProviderType, err := normalizeProviderTypeValue(update.ProviderType, serviceRecord.ProviderType)
	if err != nil {
		return domain.ServiceRecord{}, false, err
	}
	nextStatus, err := normalizeServiceStatusValue(update.Status, string(serviceRecord.Status))
	if err != nil {
		return domain.ServiceRecord{}, false, err
	}
	nextSyncStatus, err := normalizeSyncStatusValue(update.SyncStatus, serviceRecord.SyncStatus)
	if err != nil {
		return domain.ServiceRecord{}, false, err
	}

	if update.ProviderAccountID != nil {
		serviceRecord.ProviderAccountID = *update.ProviderAccountID
	}
	serviceRecord.ProviderType = nextProviderType
	serviceRecord.ProviderResourceID = strings.TrimSpace(update.ProviderResourceID)
	serviceRecord.RegionName = strings.TrimSpace(update.RegionName)
	serviceRecord.IPAddress = strings.TrimSpace(update.IPAddress)
	serviceRecord.Status = domain.ServiceStatus(nextStatus)
	serviceRecord.SyncStatus = nextSyncStatus
	serviceRecord.SyncMessage = strings.TrimSpace(update.SyncMessage)
	serviceRecord.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	serviceRecord.LastAction = "manual-adjust"
	serviceRecord.LastSyncAt = serviceRecord.UpdatedAt

	var nextDueAt any = nil
	if strings.TrimSpace(update.NextDueAt) != "" {
		parsedNextDueAt, ok := parseMySQLEditableTime(strings.TrimSpace(update.NextDueAt))
		if !ok {
			return domain.ServiceRecord{}, false, fmt.Errorf("服务到期时间格式不正确")
		}
		nextDueAt = parsedNextDueAt
		serviceRecord.NextDueAt = parsedNextDueAt.Format("2006-01-02")
	} else if parsedExisting, ok := parseMySQLEditableTime(strings.TrimSpace(serviceRecord.NextDueAt)); ok {
		nextDueAt = parsedExisting
	}

	var linkedInvoiceStatus string
	if serviceRecord.InvoiceID > 0 {
		if linkedInvoice, err := repository.loadInvoice(ctx, tx, serviceRecord.InvoiceID); err == nil {
			linkedInvoiceStatus = string(linkedInvoice.Status)
		}
	}

	var linkedOrderID int64
	if serviceRecord.OrderID > 0 {
		if _, err := repository.loadOrderForUpdate(ctx, tx, serviceRecord.OrderID); err != nil && err != sql.ErrNoRows {
			return domain.ServiceRecord{}, false, err
		}
		linkedOrderID = serviceRecord.OrderID
	}

	if _, err := tx.ExecContext(ctx, `
UPDATE services
SET provider_type = ?, provider_account_id = ?, provider_resource_id = NULLIF(?, ''), region_name = ?, ip_address = ?,
    status = ?, sync_status = ?, sync_message = ?, next_due_at = ?, last_action = ?, last_sync_at = NOW(), updated_at = NOW()
WHERE id = ?`,
		serviceRecord.ProviderType,
		serviceRecord.ProviderAccountID,
		serviceRecord.ProviderResourceID,
		serviceRecord.RegionName,
		serviceRecord.IPAddress,
		serviceRecord.Status,
		serviceRecord.SyncStatus,
		serviceRecord.SyncMessage,
		nextDueAt,
		serviceRecord.LastAction,
		serviceID,
	); err != nil {
		return domain.ServiceRecord{}, false, err
	}

	if linkedOrderID > 0 {
		nextOrderStatus := resolveManualLinkedOrderStatusForService(nextStatus, linkedInvoiceStatus)
		if nextOrderStatus != "" {
			if _, err := tx.ExecContext(ctx, `
UPDATE orders
SET status = ?, updated_at = NOW()
WHERE id = ?`,
				nextOrderStatus,
				linkedOrderID,
			); err != nil {
				return domain.ServiceRecord{}, false, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return domain.ServiceRecord{}, false, err
	}

	reloaded, ok := repository.GetServiceByID(serviceID)
	return reloaded, ok, nil
}

func (repository *MySQLRepository) CreateServiceChangeOrder(
	serviceID int64,
	input domain.ServiceChangeOrderInput,
) (domain.Order, domain.Invoice, *domain.ServiceRecord, bool, error) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Order{}, domain.Invoice{}, nil, false, err
	}
	defer rollbackQuietly(tx)

	var lockedID int64
	if err := tx.QueryRowContext(ctx, `SELECT id FROM services WHERE id = ? FOR UPDATE`, serviceID).Scan(&lockedID); err != nil {
		if err == sql.ErrNoRows {
			return domain.Order{}, domain.Invoice{}, nil, false, nil
		}
		return domain.Order{}, domain.Invoice{}, nil, false, err
	}

	serviceRecord, err := repository.loadService(ctx, tx, serviceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Order{}, domain.Invoice{}, nil, false, nil
		}
		return domain.Order{}, domain.Invoice{}, nil, false, err
	}

	baseOrder := domain.Order{}
	if serviceRecord.OrderID > 0 {
		item, loadErr := repository.loadOrderForUpdate(ctx, tx, serviceRecord.OrderID)
		if loadErr != nil && loadErr != sql.ErrNoRows {
			return domain.Order{}, domain.Invoice{}, nil, false, loadErr
		}
		if loadErr == nil {
			baseOrder = item
		}
	}

	if input.Amount < 0 {
		return domain.Order{}, domain.Invoice{}, nil, false, fmt.Errorf("改配单金额不能小于 0")
	}

	billingCycle := firstNonEmptyString(strings.TrimSpace(input.BillingCycle), firstNonEmptyString(strings.TrimSpace(baseOrder.BillingCycle), "monthly"))
	title := firstNonEmptyString(strings.TrimSpace(input.Title), strings.TrimSpace(serviceRecord.ProductName)+" 改配单")
	productType := firstNonEmptyString(strings.TrimSpace(baseOrder.ProductType), "SERVICE_CHANGE")
	customerName := strings.TrimSpace(baseOrder.CustomerName)
	configJSON, _ := json.Marshal(serviceRecord.Configuration)
	resourceJSON, _ := json.Marshal(serviceRecord.ResourceSnapshot)
	payloadJSON, _ := json.Marshal(input.Payload)

	tempOrderNo := fmt.Sprintf("TMP-ORD-%d", time.Now().UnixNano())
	orderResult, err := tx.ExecContext(ctx, `
INSERT INTO orders (order_no, customer_id, customer_name, product_id, product_name, product_type, automation_type, provider_account_id, billing_cycle, amount, status, configuration_snapshot, resource_snapshot)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		tempOrderNo,
		serviceRecord.CustomerID,
		customerName,
		baseOrder.ProductID,
		title,
		productType,
		normalizeAutomationType("SERVICE_CHANGE"),
		serviceRecord.ProviderAccountID,
		billingCycle,
		input.Amount,
		domain.OrderStatusPending,
		configJSON,
		resourceJSON,
	)
	if err != nil {
		return domain.Order{}, domain.Invoice{}, nil, false, err
	}

	orderID, err := orderResult.LastInsertId()
	if err != nil {
		return domain.Order{}, domain.Invoice{}, nil, false, err
	}
	orderNo := fmt.Sprintf("ORD-%08d", orderID)
	if _, err := tx.ExecContext(ctx, `UPDATE orders SET order_no = ? WHERE id = ?`, orderNo, orderID); err != nil {
		return domain.Order{}, domain.Invoice{}, nil, false, err
	}

	if _, err := tx.ExecContext(ctx, `
INSERT INTO order_items (order_id, product_id, product_name, billing_cycle, quantity, amount)
VALUES (?, ?, ?, ?, 1, ?)`,
		orderID,
		baseOrder.ProductID,
		title,
		billingCycle,
		input.Amount,
	); err != nil {
		return domain.Order{}, domain.Invoice{}, nil, false, err
	}

	tempInvoiceNo := fmt.Sprintf("TMP-INV-%d", time.Now().UnixNano())
	invoiceResult, err := tx.ExecContext(ctx, `
INSERT INTO invoices (customer_id, order_id, order_no, invoice_no, product_name, billing_cycle, status, total_amount, due_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, DATE_ADD(NOW(), INTERVAL 72 HOUR))`,
		serviceRecord.CustomerID,
		orderID,
		orderNo,
		tempInvoiceNo,
		title,
		billingCycle,
		domain.InvoiceStatusUnpaid,
		input.Amount,
	)
	if err != nil {
		return domain.Order{}, domain.Invoice{}, nil, false, err
	}

	invoiceID, err := invoiceResult.LastInsertId()
	if err != nil {
		return domain.Order{}, domain.Invoice{}, nil, false, err
	}
	invoiceNo := fmt.Sprintf("INV-%08d", invoiceID)
	if _, err := tx.ExecContext(ctx, `UPDATE invoices SET invoice_no = ? WHERE id = ?`, invoiceNo, invoiceID); err != nil {
		return domain.Order{}, domain.Invoice{}, nil, false, err
	}

	if _, err := tx.ExecContext(ctx, `
INSERT INTO invoice_items (invoice_id, order_id, product_id, product_name, billing_cycle, amount)
VALUES (?, ?, ?, ?, ?, ?)`,
		invoiceID,
		orderID,
		baseOrder.ProductID,
		title,
		billingCycle,
		input.Amount,
	); err != nil {
		return domain.Order{}, domain.Invoice{}, nil, false, err
	}

	if _, err := tx.ExecContext(ctx, `
INSERT INTO service_change_orders (service_id, order_id, invoice_id, action_name, title, amount, status, reason, payload_json)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		serviceID,
		orderID,
		invoiceID,
		strings.TrimSpace(input.ActionName),
		title,
		input.Amount,
		domain.InvoiceStatusUnpaid,
		firstNonEmptyString(strings.TrimSpace(input.Reason), "后台生成改配单"),
		string(payloadJSON),
	); err != nil {
		return domain.Order{}, domain.Invoice{}, nil, false, err
	}

	if err := tx.Commit(); err != nil {
		return domain.Order{}, domain.Invoice{}, nil, false, err
	}

	order, _ := repository.GetOrderByID(orderID)
	invoice, _ := repository.GetInvoiceByID(invoiceID)
	reloadedService, _ := repository.GetServiceByID(serviceID)
	return order, invoice, &reloadedService, true, nil
}

func (repository *MySQLRepository) ListPaymentsByInvoice(invoiceID int64) []domain.PaymentRecord {
	rows, err := repository.db.Query(`
SELECT id, payment_no, invoice_id, order_id, customer_id, channel, trade_no, amount, source, status, operator_name,
       DATE_FORMAT(paid_at, '%Y-%m-%d %H:%i:%s')
FROM payments
WHERE invoice_id = ?
ORDER BY id DESC`, invoiceID)
	if err != nil {
		log.Printf("order mysql list payments failed: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]domain.PaymentRecord, 0)
	for rows.Next() {
		var (
			item    domain.PaymentRecord
			orderID sql.NullInt64
		)
		if err := rows.Scan(
			&item.ID,
			&item.PaymentNo,
			&item.InvoiceID,
			&orderID,
			&item.CustomerID,
			&item.Channel,
			&item.TradeNo,
			&item.Amount,
			&item.Source,
			&item.Status,
			&item.Operator,
			&item.PaidAt,
		); err != nil {
			log.Printf("order mysql scan payment failed: %v", err)
			return result
		}
		assignNullInt64(&item.OrderID, orderID)
		result = append(result, item)
	}
	return result
}

func (repository *MySQLRepository) ListPayments(filter domain.PaymentListFilter) ([]domain.PaymentRecord, int) {
	whereSQL, args := buildPaymentFilterSQL(filter)

	countQuery := `SELECT COUNT(*) FROM payments p` + whereSQL
	var total int
	if err := repository.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		log.Printf("order mysql count payments failed: %v", err)
		return nil, 0
	}

	sortColumn := resolvePaymentSortColumn(filter.Sort)
	sortDirection := resolveSortDirection(filter.Order)
	query := `
SELECT p.id, p.payment_no, p.invoice_id, p.order_id, p.customer_id, p.channel, p.trade_no, p.amount, p.source, p.status, p.operator_name,
       DATE_FORMAT(p.paid_at, '%Y-%m-%d %H:%i:%s')
FROM payments p` + whereSQL + ` ORDER BY ` + sortColumn + ` ` + sortDirection + `, p.id DESC`

	queryArgs := append([]any{}, args...)
	if filter.Limit > 0 {
		page := filter.Page
		if page <= 0 {
			page = 1
		}
		offset := (page - 1) * filter.Limit
		query += ` LIMIT ? OFFSET ?`
		queryArgs = append(queryArgs, filter.Limit, offset)
	}

	rows, err := repository.db.Query(query, queryArgs...)
	if err != nil {
		log.Printf("order mysql list payments failed: %v", err)
		return nil, 0
	}
	defer rows.Close()

	result := make([]domain.PaymentRecord, 0)
	for rows.Next() {
		var (
			item    domain.PaymentRecord
			orderID sql.NullInt64
		)
		if err := rows.Scan(
			&item.ID,
			&item.PaymentNo,
			&item.InvoiceID,
			&orderID,
			&item.CustomerID,
			&item.Channel,
			&item.TradeNo,
			&item.Amount,
			&item.Source,
			&item.Status,
			&item.Operator,
			&item.PaidAt,
		); err != nil {
			log.Printf("order mysql scan payment failed: %v", err)
			return result, total
		}
		assignNullInt64(&item.OrderID, orderID)
		result = append(result, item)
	}
	return result, total
}

func (repository *MySQLRepository) ListRefundsByInvoice(invoiceID int64) []domain.RefundRecord {
	rows, err := repository.db.Query(`
SELECT id, refund_no, invoice_id, order_id, customer_id, amount, reason, status,
       DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s')
FROM refunds
WHERE invoice_id = ?
ORDER BY id DESC`, invoiceID)
	if err != nil {
		log.Printf("order mysql list refunds failed: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]domain.RefundRecord, 0)
	for rows.Next() {
		var (
			item    domain.RefundRecord
			orderID sql.NullInt64
		)
		if err := rows.Scan(
			&item.ID,
			&item.RefundNo,
			&item.InvoiceID,
			&orderID,
			&item.CustomerID,
			&item.Amount,
			&item.Reason,
			&item.Status,
			&item.CreatedAt,
		); err != nil {
			log.Printf("order mysql scan refund failed: %v", err)
			return result
		}
		assignNullInt64(&item.OrderID, orderID)
		result = append(result, item)
	}
	return result
}

func (repository *MySQLRepository) ListRefunds(filter domain.RefundListFilter) ([]domain.RefundRecord, int) {
	whereSQL, args := buildRefundFilterSQL(filter)

	countQuery := `SELECT COUNT(*) FROM refunds r` + whereSQL
	var total int
	if err := repository.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		log.Printf("order mysql count refunds failed: %v", err)
		return nil, 0
	}

	sortColumn := resolveRefundSortColumn(filter.Sort)
	sortDirection := resolveSortDirection(filter.Order)
	query := `
SELECT r.id, r.refund_no, r.invoice_id, r.order_id, r.customer_id, r.amount, r.reason, r.status,
       DATE_FORMAT(r.created_at, '%Y-%m-%d %H:%i:%s')
FROM refunds r` + whereSQL + ` ORDER BY ` + sortColumn + ` ` + sortDirection + `, r.id DESC`

	queryArgs := append([]any{}, args...)
	if filter.Limit > 0 {
		page := filter.Page
		if page <= 0 {
			page = 1
		}
		offset := (page - 1) * filter.Limit
		query += ` LIMIT ? OFFSET ?`
		queryArgs = append(queryArgs, filter.Limit, offset)
	}

	rows, err := repository.db.Query(query, queryArgs...)
	if err != nil {
		log.Printf("order mysql list refunds failed: %v", err)
		return nil, 0
	}
	defer rows.Close()

	result := make([]domain.RefundRecord, 0)
	for rows.Next() {
		var (
			item    domain.RefundRecord
			orderID sql.NullInt64
		)
		if err := rows.Scan(
			&item.ID,
			&item.RefundNo,
			&item.InvoiceID,
			&orderID,
			&item.CustomerID,
			&item.Amount,
			&item.Reason,
			&item.Status,
			&item.CreatedAt,
		); err != nil {
			log.Printf("order mysql scan refund failed: %v", err)
			return result, total
		}
		assignNullInt64(&item.OrderID, orderID)
		result = append(result, item)
	}
	return result, total
}

func (repository *MySQLRepository) ListOrderRequestsByOrder(orderID int64) []domain.OrderRequest {
	items, _ := repository.ListOrderRequests(domain.OrderRequestListFilter{
		OrderID: orderID,
		Page:    1,
		Limit:   0,
		Sort:    "created_at",
		Order:   "desc",
	})
	return items
}

func (repository *MySQLRepository) ListOrderRequests(filter domain.OrderRequestListFilter) ([]domain.OrderRequest, int) {
	whereSQL, args := buildOrderRequestFilterSQL(filter)

	countQuery := `SELECT COUNT(*) FROM order_requests orq` + whereSQL
	var total int
	if err := repository.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		log.Printf("order mysql count order requests failed: %v", err)
		return nil, 0
	}

	sortColumn := resolveOrderRequestSortColumn(filter.Sort)
	sortDirection := resolveSortDirection(filter.Order)
	query := `SELECT orq.id FROM order_requests orq` + whereSQL + ` ORDER BY ` + sortColumn + ` ` + sortDirection + `, orq.id DESC`
	queryArgs := append([]any{}, args...)
	if filter.Limit > 0 {
		page := filter.Page
		if page <= 0 {
			page = 1
		}
		offset := (page - 1) * filter.Limit
		query += ` LIMIT ? OFFSET ?`
		queryArgs = append(queryArgs, filter.Limit, offset)
	}

	rows, err := repository.db.Query(query, queryArgs...)
	if err != nil {
		log.Printf("order mysql list order requests failed: %v", err)
		return nil, 0
	}
	defer rows.Close()

	items := make([]domain.OrderRequest, 0)
	ctx := context.Background()
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Printf("order mysql scan order request id failed: %v", err)
			return items, total
		}
		item, err := repository.loadOrderRequest(ctx, repository.db, id)
		if err != nil {
			log.Printf("order mysql load order request failed: %v", err)
			return items, total
		}
		items = append(items, item)
	}
	return items, total
}

func (repository *MySQLRepository) CreateOrderRequest(
	orderID int64,
	input domain.OrderRequestCreateInput,
) (domain.OrderRequest, bool, error) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.OrderRequest{}, false, err
	}
	defer rollbackQuietly(tx)

	order, err := repository.loadOrderForUpdate(ctx, tx, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.OrderRequest{}, false, nil
		}
		return domain.OrderRequest{}, false, err
	}

	requestedAmount := order.Amount
	if input.RequestedAmount != nil {
		requestedAmount = *input.RequestedAmount
	} else if strings.EqualFold(strings.TrimSpace(input.Type), string(domain.OrderRequestTypeCancel)) {
		requestedAmount = 0
	}

	requestedBillingCycle := firstNonEmptyMySQL(input.RequestedBillingCycle, order.BillingCycle)
	payloadJSON := []byte("{}")
	if input.Payload != nil {
		if encoded, err := json.Marshal(input.Payload); err == nil {
			payloadJSON = encoded
		}
	}

	tempRequestNo := fmt.Sprintf("TMP-REQ-%d", time.Now().UnixNano())
	result, err := tx.ExecContext(ctx, `
INSERT INTO order_requests (
	request_no, order_id, order_no, customer_id, customer_name, product_name, type, status, summary, reason,
	current_amount, requested_amount, current_billing_cycle, requested_billing_cycle,
	source_type, source_id, source_name, payload_json
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		tempRequestNo,
		order.ID,
		order.OrderNo,
		order.CustomerID,
		order.CustomerName,
		order.ProductName,
		strings.ToUpper(strings.TrimSpace(input.Type)),
		string(domain.OrderRequestStatusPending),
		firstNonEmptyMySQL(input.Summary, order.ProductName+"业务申请"),
		strings.TrimSpace(input.Reason),
		order.Amount,
		requestedAmount,
		order.BillingCycle,
		requestedBillingCycle,
		firstNonEmptyMySQL(strings.ToUpper(strings.TrimSpace(input.SourceType)), "ADMIN"),
		nullInt64(input.SourceID),
		firstNonEmptyMySQL(input.SourceName, "系统"),
		payloadJSON,
	)
	if err != nil {
		return domain.OrderRequest{}, false, err
	}

	requestRowID, err := result.LastInsertId()
	if err != nil {
		return domain.OrderRequest{}, false, err
	}
	requestNo := fmt.Sprintf("REQ-%08d", requestRowID)
	if _, err := tx.ExecContext(ctx, `UPDATE order_requests SET request_no = ? WHERE id = ?`, requestNo, requestRowID); err != nil {
		return domain.OrderRequest{}, false, err
	}

	item, err := repository.loadOrderRequest(ctx, tx, requestRowID)
	if err != nil {
		return domain.OrderRequest{}, false, err
	}
	if err := tx.Commit(); err != nil {
		return domain.OrderRequest{}, false, err
	}
	return item, true, nil
}

func (repository *MySQLRepository) ProcessOrderRequest(
	requestID int64,
	input domain.OrderRequestProcessInput,
) (domain.OrderRequest, bool, error) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.OrderRequest{}, false, err
	}
	defer rollbackQuietly(tx)

	item, err := repository.loadOrderRequestForUpdate(ctx, tx, requestID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.OrderRequest{}, false, nil
		}
		return domain.OrderRequest{}, false, err
	}
	if !canTransitionOrderRequestStatusMySQL(string(item.Status), input.Status) {
		return domain.OrderRequest{}, false, fmt.Errorf("当前申请状态不允许继续处理")
	}

	if _, err := tx.ExecContext(ctx, `
UPDATE order_requests
SET status = ?, process_note = ?, processor_type = ?, processor_id = ?, processor_name = ?, processed_at = NOW()
WHERE id = ?`,
		strings.ToUpper(strings.TrimSpace(input.Status)),
		strings.TrimSpace(input.ProcessNote),
		firstNonEmptyMySQL(strings.ToUpper(strings.TrimSpace(input.ProcessorType)), "ADMIN"),
		nullInt64(input.ProcessorID),
		firstNonEmptyMySQL(input.ProcessorName, "系统"),
		requestID,
	); err != nil {
		return domain.OrderRequest{}, false, err
	}

	updated, err := repository.loadOrderRequest(ctx, tx, requestID)
	if err != nil {
		return domain.OrderRequest{}, false, err
	}
	if err := tx.Commit(); err != nil {
		return domain.OrderRequest{}, false, err
	}
	return updated, true, nil
}

func (repository *MySQLRepository) ListServiceChangeOrders(filter domain.ServiceChangeOrderListFilter) ([]domain.ServiceChangeOrder, int) {
	whereSQL, args := buildServiceChangeOrderFilterSQL(filter)

	countQuery := `SELECT COUNT(*) FROM service_change_orders sco
LEFT JOIN orders o ON o.id = sco.order_id
LEFT JOIN invoices i ON i.id = sco.invoice_id
LEFT JOIN services s ON s.id = sco.service_id` + whereSQL
	var total int
	if err := repository.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		log.Printf("order mysql count service change orders failed: %v", err)
		return nil, 0
	}

	sortColumn := resolveServiceChangeOrderSortColumn(filter.Sort)
	sortDirection := resolveSortDirection(filter.Order)
	query := `
SELECT sco.id, sco.service_id, sco.order_id, sco.invoice_id, sco.action_name, sco.title, sco.amount, sco.status, sco.reason, sco.payload_json,
       IFNULL(i.billing_cycle, ''),
       IFNULL(DATE_FORMAT(sco.paid_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(sco.refunded_at, '%Y-%m-%d %H:%i:%s'), ''),
       DATE_FORMAT(sco.created_at, '%Y-%m-%d %H:%i:%s'),
       DATE_FORMAT(sco.updated_at, '%Y-%m-%d %H:%i:%s')
FROM service_change_orders sco
LEFT JOIN orders o ON o.id = sco.order_id
LEFT JOIN invoices i ON i.id = sco.invoice_id
LEFT JOIN services s ON s.id = sco.service_id` + whereSQL + ` ORDER BY ` + sortColumn + ` ` + sortDirection + `, sco.id DESC`

	queryArgs := append([]any{}, args...)
	if filter.Limit > 0 {
		page := filter.Page
		if page <= 0 {
			page = 1
		}
		offset := (page - 1) * filter.Limit
		query += ` LIMIT ? OFFSET ?`
		queryArgs = append(queryArgs, filter.Limit, offset)
	}

	rows, err := repository.db.Query(query, queryArgs...)
	if err != nil {
		log.Printf("order mysql list service change orders failed: %v", err)
		return nil, 0
	}
	defer rows.Close()

	items := make([]domain.ServiceChangeOrder, 0)
	for rows.Next() {
		var (
			item        domain.ServiceChangeOrder
			payloadJSON sql.NullString
		)
		if err := rows.Scan(
			&item.ID,
			&item.ServiceID,
			&item.OrderID,
			&item.InvoiceID,
			&item.ActionName,
			&item.Title,
			&item.Amount,
			&item.Status,
			&item.Reason,
			&payloadJSON,
			&item.BillingCycle,
			&item.PaidAt,
			&item.RefundedAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			log.Printf("order mysql scan service change order failed: %v", err)
			return items, total
		}
		if payloadJSON.Valid && strings.TrimSpace(payloadJSON.String) != "" {
			var payload map[string]any
			if err := json.Unmarshal([]byte(payloadJSON.String), &payload); err == nil {
				item.Payload = payload
			}
		}
		items = append(items, item)
	}
	return items, total
}

func (repository *MySQLRepository) Checkout(
	customerID int64,
	customerName string,
	productID int64,
	productName, productType, automationType string,
	providerAccountID int64,
	cycleCode string,
	amount float64,
	configuration []domain.ServiceConfigSelection,
	resourceSnapshot domain.ServiceResourceSnapshot,
) (domain.Order, domain.Invoice) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("order mysql begin checkout failed: %v", err)
		return domain.Order{}, domain.Invoice{}
	}
	defer rollbackQuietly(tx)

	configJSON, _ := json.Marshal(configuration)
	resourceJSON, _ := json.Marshal(resourceSnapshot)
	tempOrderNo := fmt.Sprintf("TMP-ORD-%d", time.Now().UnixNano())

	orderResult, err := tx.ExecContext(ctx, `
INSERT INTO orders (order_no, customer_id, customer_name, product_id, product_name, product_type, automation_type, provider_account_id, billing_cycle, amount, status, configuration_snapshot, resource_snapshot)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		tempOrderNo,
		customerID,
		customerName,
		productID,
		productName,
		productType,
		normalizeAutomationType(automationType),
		providerAccountID,
		cycleCode,
		amount,
		domain.OrderStatusPending,
		configJSON,
		resourceJSON,
	)
	if err != nil {
		log.Printf("order mysql insert checkout order failed: %v", err)
		return domain.Order{}, domain.Invoice{}
	}

	orderID, err := orderResult.LastInsertId()
	if err != nil {
		log.Printf("order mysql get checkout order id failed: %v", err)
		return domain.Order{}, domain.Invoice{}
	}
	orderNo := fmt.Sprintf("ORD-%08d", orderID)
	if _, err := tx.ExecContext(ctx, `UPDATE orders SET order_no = ? WHERE id = ?`, orderNo, orderID); err != nil {
		log.Printf("order mysql update order no failed: %v", err)
		return domain.Order{}, domain.Invoice{}
	}

	if _, err := tx.ExecContext(ctx, `
INSERT INTO order_items (order_id, product_id, product_name, billing_cycle, quantity, amount)
VALUES (?, ?, ?, ?, 1, ?)`,
		orderID,
		productID,
		productName,
		cycleCode,
		amount,
	); err != nil {
		log.Printf("order mysql insert order item failed: %v", err)
		return domain.Order{}, domain.Invoice{}
	}

	tempInvoiceNo := fmt.Sprintf("TMP-INV-%d", time.Now().UnixNano())
	invoiceResult, err := tx.ExecContext(ctx, `
INSERT INTO invoices (customer_id, order_id, order_no, invoice_no, product_name, billing_cycle, status, total_amount, due_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, DATE_ADD(NOW(), INTERVAL 72 HOUR))`,
		customerID,
		orderID,
		orderNo,
		tempInvoiceNo,
		productName,
		cycleCode,
		domain.InvoiceStatusUnpaid,
		amount,
	)
	if err != nil {
		log.Printf("order mysql insert invoice failed: %v", err)
		return domain.Order{}, domain.Invoice{}
	}

	invoiceID, err := invoiceResult.LastInsertId()
	if err != nil {
		log.Printf("order mysql get invoice id failed: %v", err)
		return domain.Order{}, domain.Invoice{}
	}
	invoiceNo := fmt.Sprintf("INV-%08d", invoiceID)
	if _, err := tx.ExecContext(ctx, `UPDATE invoices SET invoice_no = ? WHERE id = ?`, invoiceNo, invoiceID); err != nil {
		log.Printf("order mysql update invoice no failed: %v", err)
		return domain.Order{}, domain.Invoice{}
	}

	if _, err := tx.ExecContext(ctx, `
INSERT INTO invoice_items (invoice_id, order_id, product_id, product_name, billing_cycle, amount)
VALUES (?, ?, ?, ?, ?, ?)`,
		invoiceID,
		orderID,
		productID,
		productName,
		cycleCode,
		amount,
	); err != nil {
		log.Printf("order mysql insert invoice item failed: %v", err)
		return domain.Order{}, domain.Invoice{}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("order mysql commit checkout failed: %v", err)
		return domain.Order{}, domain.Invoice{}
	}

	order, _ := repository.GetOrderByID(orderID)
	invoice, _ := repository.GetInvoiceByID(invoiceID)
	return order, invoice
}

func (repository *MySQLRepository) PayInvoice(invoiceID int64, channel, source, operator, tradeNo string) (domain.Invoice, *domain.ServiceRecord, domain.PaymentRecord, bool) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("order mysql begin pay failed: %v", err)
		return domain.Invoice{}, nil, domain.PaymentRecord{}, false
	}
	defer rollbackQuietly(tx)

	invoice, err := repository.loadInvoiceForUpdate(ctx, tx, invoiceID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("order mysql lock invoice failed: %v", err)
		}
		return domain.Invoice{}, nil, domain.PaymentRecord{}, false
	}

	if invoice.Status == domain.InvoiceStatusRefunded {
		return domain.Invoice{}, nil, domain.PaymentRecord{}, false
	}

	if invoice.Status == domain.InvoiceStatusPaid {
		payment, _ := repository.loadLatestPayment(ctx, tx, invoiceID)
		serviceRecord, _ := repository.findServiceByInvoice(ctx, tx, invoiceID)
		if serviceRecord == nil {
			if changeLink, exists, linkErr := repository.loadServiceChangeLinkByInvoiceID(ctx, tx, invoiceID); linkErr == nil && exists {
				if linkedService, loadErr := repository.loadService(ctx, tx, changeLink.ServiceID); loadErr == nil {
					serviceRecord = &linkedService
				}
			}
		}
		if err := tx.Commit(); err != nil {
			log.Printf("order mysql commit pay noop failed: %v", err)
		}
		return invoice, serviceRecord, payment, true
	}

	now := time.Now()
	paymentChannel := strings.ToUpper(strings.TrimSpace(firstNonEmptyString(channel, "ONLINE")))
	var wallet *domain.CustomerWallet
	if paymentChannel == "BALANCE" {
		item, walletErr := repository.loadCustomerWallet(ctx, tx, invoice.CustomerID, true)
		if walletErr != nil || item.Balance+0.00001 < invoice.TotalAmount {
			return domain.Invoice{}, nil, domain.PaymentRecord{}, false
		}
		wallet = &item
	}
	if _, err := tx.ExecContext(ctx, `
UPDATE invoices
SET status = ?, paid_at = ?, updated_at = NOW()
WHERE id = ?`,
		domain.InvoiceStatusPaid,
		now,
		invoiceID,
	); err != nil {
		log.Printf("order mysql update invoice paid failed: %v", err)
		return domain.Invoice{}, nil, domain.PaymentRecord{}, false
	}

	paymentNo, paymentID, err := repository.insertPayment(ctx, tx, invoice, paymentChannel, source, operator, tradeNo, now)
	if err != nil {
		log.Printf("order mysql insert payment failed: %v", err)
		return domain.Invoice{}, nil, domain.PaymentRecord{}, false
	}
	if wallet != nil {
		if _, _, err := repository.applyWalletAdjustmentMySQL(ctx, tx, *wallet, domain.WalletAdjustment{
			Target:       "BALANCE",
			Operation:    "DECREASE",
			Amount:       invoice.TotalAmount,
			Summary:      fmt.Sprintf("浣欓鏀粯璐﹀崟 %s", invoice.InvoiceNo),
			Remark:       fmt.Sprintf("payment:%s", paymentNo),
			OperatorType: source,
			OperatorName: operator,
		}, invoice.OrderID, invoice.OrderNo, invoice.ID, invoice.InvoiceNo, paymentID, paymentNo, 0, ""); err != nil {
			log.Printf("order mysql apply wallet payment failed: %v", err)
			return domain.Invoice{}, nil, domain.PaymentRecord{}, false
		}
	}

	if _, err := tx.ExecContext(ctx, `
UPDATE orders
SET status = ?, updated_at = NOW()
WHERE id = ?`, domain.OrderStatusActive, invoice.OrderID); err != nil {
		log.Printf("order mysql update order status failed: %v", err)
		return domain.Invoice{}, nil, domain.PaymentRecord{}, false
	}

	changeLink, isServiceChange, err := repository.loadServiceChangeLinkByInvoiceID(ctx, tx, invoiceID)
	if err != nil {
		log.Printf("order mysql load service change failed: %v", err)
		return domain.Invoice{}, nil, domain.PaymentRecord{}, false
	}

	var serviceRecord *domain.ServiceRecord
	if isServiceChange {
		if _, err := tx.ExecContext(ctx, `
UPDATE service_change_orders
SET status = ?, paid_at = NOW(), updated_at = NOW()
WHERE id = ?`,
			domain.InvoiceStatusPaid,
			changeLink.ID,
		); err != nil {
			log.Printf("order mysql update service change paid failed: %v", err)
			return domain.Invoice{}, nil, domain.PaymentRecord{}, false
		}
		if linkedService, loadErr := repository.loadService(ctx, tx, changeLink.ServiceID); loadErr == nil {
			serviceRecord = &linkedService
		}
	} else {
		serviceRecord, err = repository.activateOrCreateService(ctx, tx, invoice, now)
		if err != nil {
			log.Printf("order mysql activate service failed: %v", err)
			return domain.Invoice{}, nil, domain.PaymentRecord{}, false
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("order mysql commit pay failed: %v", err)
		return domain.Invoice{}, nil, domain.PaymentRecord{}, false
	}

	updatedInvoice, _ := repository.GetInvoiceByID(invoiceID)
	payment, _ := repository.loadPaymentByID(context.Background(), repository.db, paymentID)
	_ = paymentNo
	if serviceRecord == nil {
		return updatedInvoice, nil, payment, true
	}
	reloadedService, ok := repository.GetServiceByID(serviceRecord.ID)
	if !ok {
		return updatedInvoice, nil, payment, true
	}
	return updatedInvoice, &reloadedService, payment, true
}

func (repository *MySQLRepository) ExecuteServiceAction(serviceID int64, action string, params domain.ServiceActionParams) (domain.ServiceRecord, bool) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("order mysql begin service action failed: %v", err)
		return domain.ServiceRecord{}, false
	}
	defer rollbackQuietly(tx)

	serviceRecord, err := repository.loadService(ctx, tx, serviceID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("order mysql load service failed: %v", err)
		}
		return domain.ServiceRecord{}, false
	}

	switch action {
	case "activate":
		serviceRecord.Status = domain.ServiceStatusActive
	case "suspend":
		serviceRecord.Status = domain.ServiceStatusSuspended
	case "terminate":
		serviceRecord.Status = domain.ServiceStatusTerminated
	case "reboot":
		if serviceRecord.Status == domain.ServiceStatusSuspended || serviceRecord.Status == domain.ServiceStatusTerminated {
			return domain.ServiceRecord{}, false
		}
	case "reset-password":
		if strings.TrimSpace(params.Password) != "" {
			serviceRecord.ResourceSnapshot.PasswordHint = fmt.Sprintf("最近一次密码重置：%s", time.Now().Format("2006-01-02 15:04:05"))
		} else {
			serviceRecord.ResourceSnapshot.PasswordHint = fmt.Sprintf("最近一次密码重置：%s（系统生成密码）", time.Now().Format("2006-01-02 15:04:05"))
		}
	case "reinstall":
		imageName := strings.TrimSpace(params.ImageName)
		if imageName == "" {
			imageName = "Rocky Linux 9"
		}
		serviceRecord.ResourceSnapshot.OperatingSystem = imageName
		serviceRecord.Status = domain.ServiceStatusActive
	default:
		return domain.ServiceRecord{}, false
	}

	serviceRecord.LastAction = action
	serviceRecord.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	configJSON, _ := json.Marshal(serviceRecord.Configuration)
	resourceJSON, _ := json.Marshal(serviceRecord.ResourceSnapshot)
	if _, err := tx.ExecContext(ctx, `
UPDATE services
SET status = ?, last_action = ?, sync_status = ?, sync_message = ?, last_sync_at = NOW(), configuration_snapshot = ?, resource_snapshot = ?, updated_at = NOW()
WHERE id = ?`,
		serviceRecord.Status,
		serviceRecord.LastAction,
		"SUCCESS",
		"",
		configJSON,
		resourceJSON,
		serviceID,
	); err != nil {
		log.Printf("order mysql update service action failed: %v", err)
		return domain.ServiceRecord{}, false
	}

	if err := tx.Commit(); err != nil {
		log.Printf("order mysql commit service action failed: %v", err)
		return domain.ServiceRecord{}, false
	}

	reloaded, ok := repository.GetServiceByID(serviceID)
	return reloaded, ok
}

func (repository *MySQLRepository) UpdateServiceProvisioning(serviceID int64, update domain.ServiceProvisioningUpdate) (domain.ServiceRecord, bool) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("order mysql begin provision update failed: %v", err)
		return domain.ServiceRecord{}, false
	}
	defer rollbackQuietly(tx)

	serviceRecord, err := repository.loadService(ctx, tx, serviceID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("order mysql load provision service failed: %v", err)
		}
		return domain.ServiceRecord{}, false
	}

	if strings.TrimSpace(update.ProviderType) != "" {
		serviceRecord.ProviderType = update.ProviderType
	}
	if update.ProviderAccountID > 0 {
		serviceRecord.ProviderAccountID = update.ProviderAccountID
	}
	serviceRecord.ProviderResourceID = strings.TrimSpace(update.ProviderResourceID)
	serviceRecord.RegionName = strings.TrimSpace(update.RegionName)
	serviceRecord.IPAddress = strings.TrimSpace(update.IPAddress)
	serviceRecord.Status = update.Status
	serviceRecord.SyncStatus = strings.TrimSpace(update.SyncStatus)
	serviceRecord.SyncMessage = strings.TrimSpace(update.SyncMessage)
	serviceRecord.LastAction = strings.TrimSpace(update.LastAction)
	serviceRecord.LastSyncAt = time.Now().Format("2006-01-02 15:04:05")
	serviceRecord.UpdatedAt = serviceRecord.LastSyncAt

	if len(update.Configuration) > 0 {
		serviceRecord.Configuration = cloneConfigSelections(update.Configuration)
	}
	if update.ResourceSnapshot != (domain.ServiceResourceSnapshot{}) {
		serviceRecord.ResourceSnapshot = cloneSnapshot(update.ResourceSnapshot)
	}

	configJSON, _ := json.Marshal(serviceRecord.Configuration)
	resourceJSON, _ := json.Marshal(serviceRecord.ResourceSnapshot)
	if _, err := tx.ExecContext(ctx, `
UPDATE services
SET provider_type = ?, provider_account_id = ?, provider_resource_id = NULLIF(?, ''), region_name = ?, ip_address = ?, status = ?, sync_status = ?, sync_message = ?, last_action = ?, last_sync_at = NOW(), configuration_snapshot = ?, resource_snapshot = ?, updated_at = NOW()
WHERE id = ?`,
		serviceRecord.ProviderType,
		serviceRecord.ProviderAccountID,
		serviceRecord.ProviderResourceID,
		serviceRecord.RegionName,
		serviceRecord.IPAddress,
		serviceRecord.Status,
		serviceRecord.SyncStatus,
		serviceRecord.SyncMessage,
		serviceRecord.LastAction,
		configJSON,
		resourceJSON,
		serviceID,
	); err != nil {
		log.Printf("order mysql update provision service failed: %v", err)
		return domain.ServiceRecord{}, false
	}

	if err := tx.Commit(); err != nil {
		log.Printf("order mysql commit provision update failed: %v", err)
		return domain.ServiceRecord{}, false
	}

	reloaded, ok := repository.GetServiceByID(serviceID)
	return reloaded, ok
}

func (repository *MySQLRepository) RefundInvoice(invoiceID int64, reason string) (domain.Invoice, *domain.ServiceRecord, domain.RefundRecord, bool) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("order mysql begin refund failed: %v", err)
		return domain.Invoice{}, nil, domain.RefundRecord{}, false
	}
	defer rollbackQuietly(tx)

	invoice, err := repository.loadInvoiceForUpdate(ctx, tx, invoiceID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("order mysql lock refund invoice failed: %v", err)
		}
		return domain.Invoice{}, nil, domain.RefundRecord{}, false
	}
	if invoice.Status != domain.InvoiceStatusPaid {
		return domain.Invoice{}, nil, domain.RefundRecord{}, false
	}
	payment, hasPayment := repository.loadLatestPayment(ctx, tx, invoiceID)
	var wallet *domain.CustomerWallet
	if hasPayment && strings.EqualFold(payment.Channel, "BALANCE") {
		item, walletErr := repository.loadCustomerWallet(ctx, tx, invoice.CustomerID, true)
		if walletErr != nil {
			return domain.Invoice{}, nil, domain.RefundRecord{}, false
		}
		wallet = &item
	}

	if _, err := tx.ExecContext(ctx, `
UPDATE invoices
SET status = ?, updated_at = NOW()
WHERE id = ?`,
		domain.InvoiceStatusRefunded,
		invoiceID,
	); err != nil {
		log.Printf("order mysql update refund invoice failed: %v", err)
		return domain.Invoice{}, nil, domain.RefundRecord{}, false
	}

	refundNo, refundID, err := repository.insertRefund(ctx, tx, invoice, reason)
	if err != nil {
		log.Printf("order mysql insert refund failed: %v", err)
		return domain.Invoice{}, nil, domain.RefundRecord{}, false
	}
	_ = refundNo
	if wallet != nil {
		if _, _, err := repository.applyWalletAdjustmentMySQL(ctx, tx, *wallet, domain.WalletAdjustment{
			Target:       "BALANCE",
			Operation:    "INCREASE",
			Amount:       invoice.TotalAmount,
			Summary:      fmt.Sprintf("閫€娆惧洖閫€浣欓 %s", refundNo),
			Remark:       reason,
			OperatorType: "ADMIN",
			OperatorName: "绯荤粺",
		}, invoice.OrderID, invoice.OrderNo, invoice.ID, invoice.InvoiceNo, payment.ID, payment.PaymentNo, refundID, refundNo); err != nil {
			log.Printf("order mysql apply wallet refund failed: %v", err)
			return domain.Invoice{}, nil, domain.RefundRecord{}, false
		}
	}

	if _, err := tx.ExecContext(ctx, `
UPDATE orders
SET status = ?, updated_at = NOW()
WHERE id = ?`,
		domain.OrderStatusPending,
		invoice.OrderID,
	); err != nil {
		log.Printf("order mysql rollback order status failed: %v", err)
		return domain.Invoice{}, nil, domain.RefundRecord{}, false
	}

	serviceRecord, _ := repository.findServiceByInvoice(ctx, tx, invoiceID)
	changeLink, isServiceChange, err := repository.loadServiceChangeLinkByInvoiceID(ctx, tx, invoiceID)
	if err != nil {
		log.Printf("order mysql load service change for refund failed: %v", err)
		return domain.Invoice{}, nil, domain.RefundRecord{}, false
	}
	if isServiceChange {
		if _, err := tx.ExecContext(ctx, `
UPDATE service_change_orders
SET status = ?, refunded_at = NOW(), updated_at = NOW()
WHERE id = ?`,
			domain.InvoiceStatusRefunded,
			changeLink.ID,
		); err != nil {
			log.Printf("order mysql update service change refund failed: %v", err)
			return domain.Invoice{}, nil, domain.RefundRecord{}, false
		}
		if linkedService, loadErr := repository.loadService(ctx, tx, changeLink.ServiceID); loadErr == nil {
			serviceRecord = &linkedService
		}
	} else if serviceRecord != nil {
		if _, err := tx.ExecContext(ctx, `
UPDATE services
SET status = ?, last_action = ?, sync_status = ?, sync_message = ?, last_sync_at = NOW(), updated_at = NOW()
WHERE id = ?`,
			domain.ServiceStatusTerminated,
			"refund-terminate",
			"SUCCESS",
			"退款后终止",
			serviceRecord.ID,
		); err != nil {
			log.Printf("order mysql terminate service after refund failed: %v", err)
			return domain.Invoice{}, nil, domain.RefundRecord{}, false
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("order mysql commit refund failed: %v", err)
		return domain.Invoice{}, nil, domain.RefundRecord{}, false
	}

	updatedInvoice, _ := repository.GetInvoiceByID(invoiceID)
	refund, _ := repository.loadRefundByID(context.Background(), repository.db, refundID)
	if serviceRecord == nil {
		return updatedInvoice, nil, refund, true
	}
	reloadedService, ok := repository.GetServiceByID(serviceRecord.ID)
	if !ok {
		return updatedInvoice, nil, refund, true
	}
	return updatedInvoice, &reloadedService, refund, true
}

func (repository *MySQLRepository) listOrdersWithFilter(filter domain.OrderListFilter) ([]domain.Order, int) {
	whereSQL, args := buildOrderFilterSQL(filter)

	countQuery := `SELECT COUNT(*) FROM orders o` + whereSQL
	var total int
	if err := repository.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		log.Printf("order mysql count orders failed: %v", err)
		return nil, 0
	}

	sortColumn := resolveOrderSortColumn(filter.Sort)
	sortDirection := resolveSortDirection(filter.Order)

	query := `SELECT o.id FROM orders o` + whereSQL + ` ORDER BY ` + sortColumn + ` ` + sortDirection + `, o.id DESC`
	queryArgs := append([]any{}, args...)
	if filter.Limit > 0 {
		page := filter.Page
		if page <= 0 {
			page = 1
		}
		offset := (page - 1) * filter.Limit
		query += ` LIMIT ? OFFSET ?`
		queryArgs = append(queryArgs, filter.Limit, offset)
	}

	rows, err := repository.db.Query(query, queryArgs...)
	if err != nil {
		log.Printf("order mysql list orders failed: %v", err)
		return nil, 0
	}
	defer rows.Close()

	result := make([]domain.Order, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Printf("order mysql scan order id failed: %v", err)
			return result, total
		}
		item, err := repository.loadOrder(context.Background(), repository.db, id)
		if err != nil {
			log.Printf("order mysql load order failed: %v", err)
			continue
		}
		result = append(result, item)
	}
	return result, total
}

func (repository *MySQLRepository) listInvoicesBy(filter string, value int64) []domain.Invoice {
	query := `SELECT id FROM invoices`
	args := make([]any, 0)
	switch filter {
	case "customer":
		query += ` WHERE customer_id = ?`
		args = append(args, value)
	case "order":
		query += ` WHERE order_id = ?`
		args = append(args, value)
	}
	query += ` ORDER BY id DESC`

	rows, err := repository.db.Query(query, args...)
	if err != nil {
		log.Printf("order mysql list invoices failed: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]domain.Invoice, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Printf("order mysql scan invoice id failed: %v", err)
			return result
		}
		item, err := repository.loadInvoice(context.Background(), repository.db, id)
		if err != nil {
			log.Printf("order mysql load invoice failed: %v", err)
			continue
		}
		result = append(result, item)
	}
	return result
}

func (repository *MySQLRepository) listInvoicesWithFilter(filter domain.InvoiceListFilter) ([]domain.Invoice, int) {
	whereSQL, args := buildInvoiceFilterSQL(filter)

	countQuery := `SELECT COUNT(*) FROM invoices i` + whereSQL
	var total int
	if err := repository.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		log.Printf("order mysql count invoices failed: %v", err)
		return nil, 0
	}

	sortColumn := resolveInvoiceSortColumn(filter.Sort)
	sortDirection := resolveSortDirection(filter.Order)
	query := `SELECT i.id FROM invoices i` + whereSQL + ` ORDER BY ` + sortColumn + ` ` + sortDirection + `, i.id DESC`
	queryArgs := append([]any{}, args...)
	if filter.Limit > 0 {
		page := filter.Page
		if page <= 0 {
			page = 1
		}
		offset := (page - 1) * filter.Limit
		query += ` LIMIT ? OFFSET ?`
		queryArgs = append(queryArgs, filter.Limit, offset)
	}

	rows, err := repository.db.Query(query, queryArgs...)
	if err != nil {
		log.Printf("order mysql list invoices failed: %v", err)
		return nil, 0
	}
	defer rows.Close()

	result := make([]domain.Invoice, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Printf("order mysql scan invoice id failed: %v", err)
			return result, total
		}
		item, err := repository.loadInvoice(context.Background(), repository.db, id)
		if err != nil {
			log.Printf("order mysql load invoice failed: %v", err)
			continue
		}
		result = append(result, item)
	}
	return result, total
}

func (repository *MySQLRepository) listServicesBy(filter string, value int64) []domain.ServiceRecord {
	query := `SELECT id FROM services`
	args := make([]any, 0)
	switch filter {
	case "customer":
		query += ` WHERE customer_id = ?`
		args = append(args, value)
	case "order":
		query = `
SELECT id
FROM services
WHERE order_id = ?
UNION
SELECT service_id AS id
FROM service_change_orders
WHERE order_id = ?
ORDER BY id DESC`
		args = append(args, value, value)
	}
	if filter != "order" {
		query += ` ORDER BY id DESC`
	}

	rows, err := repository.db.Query(query, args...)
	if err != nil {
		log.Printf("order mysql list services failed: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]domain.ServiceRecord, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Printf("order mysql scan service id failed: %v", err)
			return result
		}
		item, err := repository.loadService(context.Background(), repository.db, id)
		if err != nil {
			log.Printf("order mysql load service failed: %v", err)
			continue
		}
		result = append(result, item)
	}
	return result
}

func (repository *MySQLRepository) listServicesWithFilter(filter domain.ServiceListFilter) ([]domain.ServiceRecord, int) {
	whereSQL, args := buildServiceFilterSQL(filter)

	countQuery := `SELECT COUNT(*) FROM services s` + whereSQL
	var total int
	if err := repository.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		log.Printf("order mysql count services failed: %v", err)
		return nil, 0
	}

	sortColumn := resolveServiceSortColumn(filter.Sort)
	sortDirection := resolveSortDirection(filter.Order)
	query := `SELECT s.id FROM services s` + whereSQL + ` ORDER BY ` + sortColumn + ` ` + sortDirection + `, s.id DESC`
	queryArgs := append([]any{}, args...)
	if filter.Limit > 0 {
		page := filter.Page
		if page <= 0 {
			page = 1
		}
		offset := (page - 1) * filter.Limit
		query += ` LIMIT ? OFFSET ?`
		queryArgs = append(queryArgs, filter.Limit, offset)
	}

	rows, err := repository.db.Query(query, queryArgs...)
	if err != nil {
		log.Printf("order mysql list services failed: %v", err)
		return nil, 0
	}
	defer rows.Close()

	result := make([]domain.ServiceRecord, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Printf("order mysql scan service id failed: %v", err)
			return result, total
		}
		item, err := repository.loadService(context.Background(), repository.db, id)
		if err != nil {
			log.Printf("order mysql load service failed: %v", err)
			continue
		}
		result = append(result, item)
	}
	return result, total
}

func (repository *MySQLRepository) loadOrder(ctx context.Context, queryer queryer, id int64) (domain.Order, error) {
	row := queryer.QueryRowContext(ctx, `
SELECT id, order_no, customer_id, customer_name, product_id, product_name, product_type, automation_type, provider_account_id, billing_cycle, amount, status,
       configuration_snapshot, resource_snapshot, DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s')
FROM orders
WHERE id = ?`, id)

	var (
		item         domain.Order
		configJSON   []byte
		resourceJSON []byte
		statusDB     string
	)
	if err := row.Scan(
		&item.ID,
		&item.OrderNo,
		&item.CustomerID,
		&item.CustomerName,
		&item.ProductID,
		&item.ProductName,
		&item.ProductType,
		&item.AutomationType,
		&item.ProviderAccountID,
		&item.BillingCycle,
		&item.Amount,
		&statusDB,
		&configJSON,
		&resourceJSON,
		&item.CreatedAt,
	); err != nil {
		return domain.Order{}, err
	}
	item.Status = domain.OrderStatus(statusDB)
	_ = json.Unmarshal(configJSON, &item.Configuration)
	_ = json.Unmarshal(resourceJSON, &item.ResourceSnapshot)
	repository.fillOrderBillingMeta(ctx, queryer, &item)
	return item, nil
}

func (repository *MySQLRepository) loadOrderForUpdate(ctx context.Context, tx *sql.Tx, id int64) (domain.Order, error) {
	row := tx.QueryRowContext(ctx, `
SELECT id, order_no, customer_id, customer_name, product_id, product_name, product_type, automation_type, provider_account_id, billing_cycle, amount, status,
       configuration_snapshot, resource_snapshot, DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s')
FROM orders
WHERE id = ?
FOR UPDATE`, id)

	var (
		item              domain.Order
		configJSON        []byte
		resourceJSON      []byte
		statusDB          string
		providerAccountID sql.NullInt64
	)
	if err := row.Scan(
		&item.ID,
		&item.OrderNo,
		&item.CustomerID,
		&item.CustomerName,
		&item.ProductID,
		&item.ProductName,
		&item.ProductType,
		&item.AutomationType,
		&providerAccountID,
		&item.BillingCycle,
		&item.Amount,
		&statusDB,
		&configJSON,
		&resourceJSON,
		&item.CreatedAt,
	); err != nil {
		return domain.Order{}, err
	}
	assignNullInt64(&item.ProviderAccountID, providerAccountID)
	item.Status = domain.OrderStatus(statusDB)
	_ = json.Unmarshal(configJSON, &item.Configuration)
	_ = json.Unmarshal(resourceJSON, &item.ResourceSnapshot)
	repository.fillOrderBillingMeta(ctx, tx, &item)
	return item, nil
}

func (repository *MySQLRepository) loadInvoice(ctx context.Context, queryer queryer, id int64) (domain.Invoice, error) {
	row := queryer.QueryRowContext(ctx, `
SELECT id, invoice_no, order_id, order_no, customer_id, product_name, total_amount, status,
       IFNULL(DATE_FORMAT(due_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(paid_at, '%Y-%m-%d %H:%i:%s'), ''),
       billing_cycle
FROM invoices
WHERE id = ?`, id)

	var (
		item     domain.Invoice
		statusDB string
		orderID  sql.NullInt64
	)
	if err := row.Scan(
		&item.ID,
		&item.InvoiceNo,
		&orderID,
		&item.OrderNo,
		&item.CustomerID,
		&item.ProductName,
		&item.TotalAmount,
		&statusDB,
		&item.DueAt,
		&item.PaidAt,
		&item.BillingCycle,
	); err != nil {
		return domain.Invoice{}, err
	}
	assignNullInt64(&item.OrderID, orderID)
	item.Status = domain.InvoiceStatus(statusDB)
	return item, nil
}

func (repository *MySQLRepository) loadInvoiceForUpdate(ctx context.Context, tx *sql.Tx, id int64) (domain.Invoice, error) {
	row := tx.QueryRowContext(ctx, `
SELECT id, invoice_no, order_id, order_no, customer_id, product_name, total_amount, status,
       IFNULL(DATE_FORMAT(due_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(paid_at, '%Y-%m-%d %H:%i:%s'), ''),
       billing_cycle
FROM invoices
WHERE id = ?
FOR UPDATE`, id)

	var (
		item     domain.Invoice
		statusDB string
	)
	if err := row.Scan(
		&item.ID,
		&item.InvoiceNo,
		&item.OrderID,
		&item.OrderNo,
		&item.CustomerID,
		&item.ProductName,
		&item.TotalAmount,
		&statusDB,
		&item.DueAt,
		&item.PaidAt,
		&item.BillingCycle,
	); err != nil {
		return domain.Invoice{}, err
	}
	item.Status = domain.InvoiceStatus(statusDB)
	return item, nil
}

func (repository *MySQLRepository) loadService(ctx context.Context, queryer queryer, id int64) (domain.ServiceRecord, error) {
	row := queryer.QueryRowContext(ctx, `
SELECT id, service_no, order_id, invoice_id, customer_id, product_name, provider_type, provider_account_id,
       IFNULL(provider_resource_id, ''),
       IFNULL(region_name, ''),
       IFNULL(ip_address, ''),
       status,
       IFNULL(sync_status, ''),
       IFNULL(sync_message, ''),
       IFNULL(DATE_FORMAT(next_due_at, '%Y-%m-%d'), ''),
       IFNULL(last_action, ''),
       IFNULL(DATE_FORMAT(last_sync_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s'), ''),
       configuration_snapshot,
       resource_snapshot,
       DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s')
FROM services
WHERE id = ?`, id)

	var (
		item              domain.ServiceRecord
		statusDB          string
		configJSON        []byte
		resourceJSON      []byte
		orderID           sql.NullInt64
		invoiceID         sql.NullInt64
		providerAccountID sql.NullInt64
	)
	if err := row.Scan(
		&item.ID,
		&item.ServiceNo,
		&orderID,
		&invoiceID,
		&item.CustomerID,
		&item.ProductName,
		&item.ProviderType,
		&providerAccountID,
		&item.ProviderResourceID,
		&item.RegionName,
		&item.IPAddress,
		&statusDB,
		&item.SyncStatus,
		&item.SyncMessage,
		&item.NextDueAt,
		&item.LastAction,
		&item.LastSyncAt,
		&item.UpdatedAt,
		&configJSON,
		&resourceJSON,
		&item.CreatedAt,
	); err != nil {
		return domain.ServiceRecord{}, err
	}
	assignNullInt64(&item.OrderID, orderID)
	assignNullInt64(&item.InvoiceID, invoiceID)
	assignNullInt64(&item.ProviderAccountID, providerAccountID)
	item.Status = domain.ServiceStatus(statusDB)
	_ = json.Unmarshal(configJSON, &item.Configuration)
	_ = json.Unmarshal(resourceJSON, &item.ResourceSnapshot)
	return item, nil
}

func (repository *MySQLRepository) loadPaymentByID(ctx context.Context, queryer queryer, id int64) (domain.PaymentRecord, error) {
	row := queryer.QueryRowContext(ctx, `
SELECT id, payment_no, invoice_id, order_id, customer_id, channel, trade_no, amount, source, status, operator_name,
       DATE_FORMAT(paid_at, '%Y-%m-%d %H:%i:%s')
FROM payments
WHERE id = ?`, id)
	var (
		item    domain.PaymentRecord
		orderID sql.NullInt64
	)
	if err := row.Scan(
		&item.ID,
		&item.PaymentNo,
		&item.InvoiceID,
		&orderID,
		&item.CustomerID,
		&item.Channel,
		&item.TradeNo,
		&item.Amount,
		&item.Source,
		&item.Status,
		&item.Operator,
		&item.PaidAt,
	); err != nil {
		return domain.PaymentRecord{}, err
	}
	assignNullInt64(&item.OrderID, orderID)
	return item, nil
}

func (repository *MySQLRepository) loadRefundByID(ctx context.Context, queryer queryer, id int64) (domain.RefundRecord, error) {
	row := queryer.QueryRowContext(ctx, `
SELECT id, refund_no, invoice_id, order_id, customer_id, amount, reason, status,
       DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s')
FROM refunds
WHERE id = ?`, id)
	var (
		item    domain.RefundRecord
		orderID sql.NullInt64
	)
	if err := row.Scan(
		&item.ID,
		&item.RefundNo,
		&item.InvoiceID,
		&orderID,
		&item.CustomerID,
		&item.Amount,
		&item.Reason,
		&item.Status,
		&item.CreatedAt,
	); err != nil {
		return domain.RefundRecord{}, err
	}
	assignNullInt64(&item.OrderID, orderID)
	return item, nil
}

func (repository *MySQLRepository) loadOrderRequest(ctx context.Context, queryer queryer, id int64) (domain.OrderRequest, error) {
	row := queryer.QueryRowContext(ctx, `
SELECT id, request_no, order_id, order_no, customer_id, customer_name, product_name, type, status, summary, reason,
       current_amount, requested_amount, IFNULL(current_billing_cycle, ''), IFNULL(requested_billing_cycle, ''),
       IFNULL(source_type, ''), source_id, IFNULL(source_name, ''),
       IFNULL(processor_type, ''), processor_id, IFNULL(processor_name, ''), IFNULL(process_note, ''),
       payload_json,
       DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'),
       DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s'),
       IFNULL(DATE_FORMAT(processed_at, '%Y-%m-%d %H:%i:%s'), '')
FROM order_requests
WHERE id = ?`, id)

	var (
		item        domain.OrderRequest
		typeDB      string
		statusDB    string
		sourceID    sql.NullInt64
		processorID sql.NullInt64
		payloadJSON []byte
	)
	if err := row.Scan(
		&item.ID,
		&item.RequestNo,
		&item.OrderID,
		&item.OrderNo,
		&item.CustomerID,
		&item.CustomerName,
		&item.ProductName,
		&typeDB,
		&statusDB,
		&item.Summary,
		&item.Reason,
		&item.CurrentAmount,
		&item.RequestedAmount,
		&item.CurrentBillingCycle,
		&item.RequestedBillingCycle,
		&item.SourceType,
		&sourceID,
		&item.SourceName,
		&item.ProcessorType,
		&processorID,
		&item.ProcessorName,
		&item.ProcessNote,
		&payloadJSON,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.ProcessedAt,
	); err != nil {
		return domain.OrderRequest{}, err
	}
	item.Type = domain.OrderRequestType(typeDB)
	item.Status = domain.OrderRequestStatus(statusDB)
	assignNullInt64(&item.SourceID, sourceID)
	assignNullInt64(&item.ProcessorID, processorID)
	_ = json.Unmarshal(payloadJSON, &item.Payload)
	return item, nil
}

func (repository *MySQLRepository) loadOrderRequestForUpdate(ctx context.Context, tx *sql.Tx, id int64) (domain.OrderRequest, error) {
	row := tx.QueryRowContext(ctx, `
SELECT id, request_no, order_id, order_no, customer_id, customer_name, product_name, type, status, summary, reason,
       current_amount, requested_amount, IFNULL(current_billing_cycle, ''), IFNULL(requested_billing_cycle, ''),
       IFNULL(source_type, ''), source_id, IFNULL(source_name, ''),
       IFNULL(processor_type, ''), processor_id, IFNULL(processor_name, ''), IFNULL(process_note, ''),
       payload_json,
       DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'),
       DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s'),
       IFNULL(DATE_FORMAT(processed_at, '%Y-%m-%d %H:%i:%s'), '')
FROM order_requests
WHERE id = ?
FOR UPDATE`, id)

	var (
		item        domain.OrderRequest
		typeDB      string
		statusDB    string
		sourceID    sql.NullInt64
		processorID sql.NullInt64
		payloadJSON []byte
	)
	if err := row.Scan(
		&item.ID,
		&item.RequestNo,
		&item.OrderID,
		&item.OrderNo,
		&item.CustomerID,
		&item.CustomerName,
		&item.ProductName,
		&typeDB,
		&statusDB,
		&item.Summary,
		&item.Reason,
		&item.CurrentAmount,
		&item.RequestedAmount,
		&item.CurrentBillingCycle,
		&item.RequestedBillingCycle,
		&item.SourceType,
		&sourceID,
		&item.SourceName,
		&item.ProcessorType,
		&processorID,
		&item.ProcessorName,
		&item.ProcessNote,
		&payloadJSON,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.ProcessedAt,
	); err != nil {
		return domain.OrderRequest{}, err
	}
	item.Type = domain.OrderRequestType(typeDB)
	item.Status = domain.OrderRequestStatus(statusDB)
	assignNullInt64(&item.SourceID, sourceID)
	assignNullInt64(&item.ProcessorID, processorID)
	_ = json.Unmarshal(payloadJSON, &item.Payload)
	return item, nil
}

func (repository *MySQLRepository) loadLatestPayment(ctx context.Context, queryer queryer, invoiceID int64) (domain.PaymentRecord, bool) {
	row := queryer.QueryRowContext(ctx, `
SELECT id, payment_no, invoice_id, order_id, customer_id, channel, trade_no, amount, source, status, operator_name,
       DATE_FORMAT(paid_at, '%Y-%m-%d %H:%i:%s')
FROM payments
WHERE invoice_id = ?
ORDER BY id DESC
LIMIT 1`, invoiceID)

	var (
		item    domain.PaymentRecord
		orderID sql.NullInt64
	)
	if err := row.Scan(
		&item.ID,
		&item.PaymentNo,
		&item.InvoiceID,
		&orderID,
		&item.CustomerID,
		&item.Channel,
		&item.TradeNo,
		&item.Amount,
		&item.Source,
		&item.Status,
		&item.Operator,
		&item.PaidAt,
	); err != nil {
		return domain.PaymentRecord{}, false
	}
	assignNullInt64(&item.OrderID, orderID)
	return item, true
}

func (repository *MySQLRepository) loadCustomerWallet(
	ctx context.Context,
	queryer queryer,
	customerID int64,
	forUpdate bool,
) (domain.CustomerWallet, error) {
	query := `
SELECT c.id, c.customer_no, c.name,
       CAST(c.balance AS DECIMAL(12,2)),
       CAST(c.credit_limit AS DECIMAL(12,2)),
       CAST(c.credit_used AS DECIMAL(12,2)),
       IFNULL(DATE_FORMAT(c.updated_at, '%Y-%m-%d %H:%i:%s'), '')
FROM customers c
WHERE c.id = ?`
	if forUpdate {
		query += ` FOR UPDATE`
	}

	row := queryer.QueryRowContext(ctx, query, customerID)
	var item domain.CustomerWallet
	if err := row.Scan(
		&item.CustomerID,
		&item.CustomerNo,
		&item.CustomerName,
		&item.Balance,
		&item.CreditLimit,
		&item.CreditUsed,
		&item.UpdatedAt,
	); err != nil {
		return domain.CustomerWallet{}, err
	}
	item.AvailableCredit = item.CreditLimit - item.CreditUsed
	return item, nil
}

func (repository *MySQLRepository) applyWalletAdjustmentMySQL(
	ctx context.Context,
	tx *sql.Tx,
	wallet domain.CustomerWallet,
	input domain.WalletAdjustment,
	orderID int64,
	orderNo string,
	invoiceID int64,
	invoiceNo string,
	paymentID int64,
	paymentNo string,
	refundID int64,
	refundNo string,
) (domain.AccountTransaction, domain.CustomerWallet, error) {
	target := strings.ToUpper(strings.TrimSpace(input.Target))
	operation := strings.ToUpper(strings.TrimSpace(input.Operation))
	if input.Amount < 0 {
		return domain.AccountTransaction{}, domain.CustomerWallet{}, fmt.Errorf("wallet amount cannot be negative")
	}

	transaction := domain.AccountTransaction{
		CustomerID:   wallet.CustomerID,
		CustomerNo:   wallet.CustomerNo,
		CustomerName: wallet.CustomerName,
		OrderID:      orderID,
		OrderNo:      orderNo,
		InvoiceID:    invoiceID,
		InvoiceNo:    invoiceNo,
		PaymentID:    paymentID,
		PaymentNo:    paymentNo,
		RefundID:     refundID,
		RefundNo:     refundNo,
		Channel:      firstNonEmptyString(strings.ToUpper(strings.TrimSpace(input.OperatorType)), "SYSTEM"),
		Summary:      strings.TrimSpace(input.Summary),
		Remark:       strings.TrimSpace(input.Remark),
		OperatorType: firstNonEmptyString(strings.ToUpper(strings.TrimSpace(input.OperatorType)), "ADMIN"),
		OperatorID:   input.OperatorID,
		OperatorName: firstNonEmptyString(input.OperatorName, "System"),
		OccurredAt:   time.Now().Format("2006-01-02 15:04:05"),
	}

	switch target {
	case "BALANCE":
		transaction.TransactionType = domain.AccountTransactionTypeAdjustment
		transaction.BalanceBefore = wallet.Balance
		transaction.BalanceAfter = wallet.Balance
		transaction.CreditBefore = wallet.CreditLimit
		transaction.CreditAfter = wallet.CreditLimit
		switch operation {
		case "INCREASE":
			transaction.Direction = domain.AccountTransactionDirectionIn
			transaction.BalanceAfter = wallet.Balance + input.Amount
		case "DECREASE":
			if wallet.Balance+0.00001 < input.Amount {
				return domain.AccountTransaction{}, domain.CustomerWallet{}, fmt.Errorf("wallet balance is insufficient")
			}
			transaction.Direction = domain.AccountTransactionDirectionOut
			transaction.BalanceAfter = wallet.Balance - input.Amount
		case "SET":
			transaction.BalanceAfter = input.Amount
			switch {
			case transaction.BalanceAfter > transaction.BalanceBefore:
				transaction.Direction = domain.AccountTransactionDirectionIn
			case transaction.BalanceAfter < transaction.BalanceBefore:
				transaction.Direction = domain.AccountTransactionDirectionOut
			default:
				transaction.Direction = domain.AccountTransactionDirectionFlat
			}
			input.Amount = absDecimal(transaction.BalanceAfter - transaction.BalanceBefore)
		default:
			return domain.AccountTransaction{}, domain.CustomerWallet{}, fmt.Errorf("wallet operation is invalid")
		}
		transaction.Amount = input.Amount
		if strings.TrimSpace(transaction.Summary) == "" {
			transaction.Summary = "Balance adjustment"
		}
		wallet.Balance = transaction.BalanceAfter

	case "CREDIT_LIMIT":
		transaction.TransactionType = domain.AccountTransactionTypeCreditLimit
		transaction.BalanceBefore = wallet.Balance
		transaction.BalanceAfter = wallet.Balance
		transaction.CreditBefore = wallet.CreditLimit
		transaction.CreditAfter = wallet.CreditLimit
		switch operation {
		case "INCREASE":
			transaction.Direction = domain.AccountTransactionDirectionIn
			transaction.CreditAfter = wallet.CreditLimit + input.Amount
		case "DECREASE":
			if wallet.CreditLimit+0.00001 < input.Amount {
				return domain.AccountTransaction{}, domain.CustomerWallet{}, fmt.Errorf("credit limit is insufficient")
			}
			transaction.Direction = domain.AccountTransactionDirectionOut
			transaction.CreditAfter = wallet.CreditLimit - input.Amount
		case "SET":
			transaction.CreditAfter = input.Amount
			switch {
			case transaction.CreditAfter > transaction.CreditBefore:
				transaction.Direction = domain.AccountTransactionDirectionIn
			case transaction.CreditAfter < transaction.CreditBefore:
				transaction.Direction = domain.AccountTransactionDirectionOut
			default:
				transaction.Direction = domain.AccountTransactionDirectionFlat
			}
			input.Amount = absDecimal(transaction.CreditAfter - transaction.CreditBefore)
		default:
			return domain.AccountTransaction{}, domain.CustomerWallet{}, fmt.Errorf("credit operation is invalid")
		}
		if transaction.CreditAfter+0.00001 < wallet.CreditUsed {
			return domain.AccountTransaction{}, domain.CustomerWallet{}, fmt.Errorf("credit limit cannot be lower than used credit")
		}
		transaction.Amount = input.Amount
		if strings.TrimSpace(transaction.Summary) == "" {
			transaction.Summary = "Credit limit adjustment"
		}
		wallet.CreditLimit = transaction.CreditAfter

	default:
		return domain.AccountTransaction{}, domain.CustomerWallet{}, fmt.Errorf("wallet target is invalid")
	}

	if _, err := tx.ExecContext(ctx, `
UPDATE customers
SET balance = ?, credit_limit = ?, updated_at = NOW()
WHERE id = ?`,
		wallet.Balance,
		wallet.CreditLimit,
		wallet.CustomerID,
	); err != nil {
		return domain.AccountTransaction{}, domain.CustomerWallet{}, err
	}

	tempNo := fmt.Sprintf("TMP-ACT-%d", time.Now().UnixNano())
	result, err := tx.ExecContext(ctx, `
INSERT INTO account_transactions (
  transaction_no, customer_id, order_id, invoice_id, payment_id, refund_id,
  transaction_type, direction, amount,
  balance_before, balance_after, credit_before, credit_after,
  channel, summary, remark, operator_type, operator_id, operator_name, occurred_at
) VALUES (?, ?, NULLIF(?, 0), NULLIF(?, 0), NULLIF(?, 0), NULLIF(?, 0), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		tempNo,
		wallet.CustomerID,
		orderID,
		invoiceID,
		paymentID,
		refundID,
		transaction.TransactionType,
		transaction.Direction,
		transaction.Amount,
		transaction.BalanceBefore,
		transaction.BalanceAfter,
		transaction.CreditBefore,
		transaction.CreditAfter,
		transaction.Channel,
		transaction.Summary,
		transaction.Remark,
		transaction.OperatorType,
		transaction.OperatorID,
		transaction.OperatorName,
		transaction.OccurredAt,
	)
	if err != nil {
		return domain.AccountTransaction{}, domain.CustomerWallet{}, err
	}

	transactionID, err := result.LastInsertId()
	if err != nil {
		return domain.AccountTransaction{}, domain.CustomerWallet{}, err
	}
	transactionNo := fmt.Sprintf("ACT-%08d", transactionID)
	if _, err := tx.ExecContext(ctx, `UPDATE account_transactions SET transaction_no = ? WHERE id = ?`, transactionNo, transactionID); err != nil {
		return domain.AccountTransaction{}, domain.CustomerWallet{}, err
	}

	transaction.ID = transactionID
	transaction.TransactionNo = transactionNo
	wallet.AvailableCredit = wallet.CreditLimit - wallet.CreditUsed
	wallet.UpdatedAt = transaction.OccurredAt
	return transaction, wallet, nil
}

func (repository *MySQLRepository) loadAccountTransaction(ctx context.Context, queryer queryer, id int64) (domain.AccountTransaction, error) {
	row := queryer.QueryRowContext(ctx, `
SELECT at.id, at.transaction_no, at.customer_id,
       IFNULL(c.customer_no, ''), IFNULL(c.name, ''),
       at.order_id, IFNULL(o.order_no, ''),
       at.invoice_id, IFNULL(i.invoice_no, ''),
       at.payment_id, IFNULL(p.payment_no, ''),
       at.refund_id, IFNULL(r.refund_no, ''),
       at.transaction_type, at.direction, at.amount,
       at.balance_before, at.balance_after, at.credit_before, at.credit_after,
       at.channel, at.summary, IFNULL(at.remark, ''), at.operator_type, at.operator_id, at.operator_name,
       DATE_FORMAT(at.occurred_at, '%Y-%m-%d %H:%i:%s')
FROM account_transactions at
LEFT JOIN customers c ON c.id = at.customer_id
LEFT JOIN orders o ON o.id = at.order_id
LEFT JOIN invoices i ON i.id = at.invoice_id
LEFT JOIN payments p ON p.id = at.payment_id
LEFT JOIN refunds r ON r.id = at.refund_id
WHERE at.id = ?`, id)

	var (
		item      domain.AccountTransaction
		orderID   sql.NullInt64
		invoiceID sql.NullInt64
		paymentID sql.NullInt64
		refundID  sql.NullInt64
	)
	if err := row.Scan(
		&item.ID,
		&item.TransactionNo,
		&item.CustomerID,
		&item.CustomerNo,
		&item.CustomerName,
		&orderID,
		&item.OrderNo,
		&invoiceID,
		&item.InvoiceNo,
		&paymentID,
		&item.PaymentNo,
		&refundID,
		&item.RefundNo,
		&item.TransactionType,
		&item.Direction,
		&item.Amount,
		&item.BalanceBefore,
		&item.BalanceAfter,
		&item.CreditBefore,
		&item.CreditAfter,
		&item.Channel,
		&item.Summary,
		&item.Remark,
		&item.OperatorType,
		&item.OperatorID,
		&item.OperatorName,
		&item.OccurredAt,
	); err != nil {
		return domain.AccountTransaction{}, err
	}

	assignNullInt64(&item.OrderID, orderID)
	assignNullInt64(&item.InvoiceID, invoiceID)
	assignNullInt64(&item.PaymentID, paymentID)
	assignNullInt64(&item.RefundID, refundID)
	return item, nil
}

func (repository *MySQLRepository) findServiceByInvoice(ctx context.Context, queryer queryer, invoiceID int64) (*domain.ServiceRecord, bool) {
	row := queryer.QueryRowContext(ctx, `SELECT id FROM services WHERE invoice_id = ? ORDER BY id DESC LIMIT 1`, invoiceID)
	var id int64
	if err := row.Scan(&id); err != nil {
		return nil, false
	}
	item, err := repository.loadService(ctx, queryer, id)
	if err != nil {
		return nil, false
	}
	return &item, true
}

func (repository *MySQLRepository) insertPayment(ctx context.Context, tx *sql.Tx, invoice domain.Invoice, channel, source, operator, tradeNo string, now time.Time) (string, int64, error) {
	tempNo := fmt.Sprintf("TMP-PAY-%d", time.Now().UnixNano())
	if tradeNo == "" {
		tradeNo = fmt.Sprintf("TRADE-%d", time.Now().UnixNano())
	}
	result, err := tx.ExecContext(ctx, `
INSERT INTO payments (payment_no, invoice_id, order_id, customer_id, channel, trade_no, amount, source, status, operator_name, paid_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		tempNo,
		invoice.ID,
		invoice.OrderID,
		invoice.CustomerID,
		firstNonEmptyString(channel, "ONLINE"),
		tradeNo,
		invoice.TotalAmount,
		firstNonEmptyString(source, "PORTAL"),
		domain.PaymentStatusCompleted,
		firstNonEmptyString(operator, "系统"),
		now,
	)
	if err != nil {
		return "", 0, err
	}

	paymentID, err := result.LastInsertId()
	if err != nil {
		return "", 0, err
	}
	paymentNo := fmt.Sprintf("PAY-%08d", paymentID)
	if _, err := tx.ExecContext(ctx, `UPDATE payments SET payment_no = ? WHERE id = ?`, paymentNo, paymentID); err != nil {
		return "", 0, err
	}
	return paymentNo, paymentID, nil
}

func (repository *MySQLRepository) activateOrCreateService(ctx context.Context, tx *sql.Tx, invoice domain.Invoice, now time.Time) (*domain.ServiceRecord, error) {
	if existing, ok := repository.findServiceByInvoice(ctx, tx, invoice.ID); ok {
		configJSON, _ := json.Marshal(existing.Configuration)
		resourceJSON, _ := json.Marshal(existing.ResourceSnapshot)
		if _, err := tx.ExecContext(ctx, `
UPDATE services
SET status = ?, last_action = ?, sync_status = ?, sync_message = ?, last_sync_at = NOW(), configuration_snapshot = ?, resource_snapshot = ?, updated_at = NOW()
WHERE id = ?`,
			domain.ServiceStatusActive,
			"activate",
			"SUCCESS",
			"",
			configJSON,
			resourceJSON,
			existing.ID,
		); err != nil {
			return nil, err
		}
		return existing, nil
	}

	order, err := repository.loadOrder(ctx, tx, invoice.OrderID)
	if err != nil {
		return nil, err
	}

	configJSON, _ := json.Marshal(order.Configuration)
	resourceJSON, _ := json.Marshal(order.ResourceSnapshot)
	tempServiceNo := fmt.Sprintf("TMP-SRV-%d", time.Now().UnixNano())
	result, err := tx.ExecContext(ctx, `
INSERT INTO services (service_no, customer_id, order_id, invoice_id, product_name, provider_type, provider_account_id, status, sync_status, sync_message, last_action, configuration_snapshot, resource_snapshot, next_due_at, last_sync_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`,
		tempServiceNo,
		order.CustomerID,
		order.ID,
		invoice.ID,
		order.ProductName,
		resolveServiceProviderType(order.AutomationType),
		order.ProviderAccountID,
		resolveProvisionedServiceStatus(order.AutomationType),
		resolveProvisionedSyncStatus(order.AutomationType),
		resolveProvisionedSyncMessage(order.AutomationType),
		resolveProvisionedLastAction(order.AutomationType),
		configJSON,
		resourceJSON,
		parseNextDue(order.BillingCycle, now),
	)
	if err != nil {
		return nil, err
	}

	serviceID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	serviceNo := fmt.Sprintf("SRV-%08d", serviceID)
	if _, err := tx.ExecContext(ctx, `UPDATE services SET service_no = ? WHERE id = ?`, serviceNo, serviceID); err != nil {
		return nil, err
	}

	item, err := repository.loadService(ctx, tx, serviceID)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (repository *MySQLRepository) insertRefund(ctx context.Context, tx *sql.Tx, invoice domain.Invoice, reason string) (string, int64, error) {
	tempNo := fmt.Sprintf("TMP-RFD-%d", time.Now().UnixNano())
	result, err := tx.ExecContext(ctx, `
INSERT INTO refunds (refund_no, invoice_id, order_id, customer_id, amount, reason, status)
VALUES (?, ?, ?, ?, ?, ?, ?)`,
		tempNo,
		invoice.ID,
		invoice.OrderID,
		invoice.CustomerID,
		invoice.TotalAmount,
		firstNonEmptyString(reason, "后台手工退款"),
		domain.RefundStatusCompleted,
	)
	if err != nil {
		return "", 0, err
	}

	refundID, err := result.LastInsertId()
	if err != nil {
		return "", 0, err
	}
	refundNo := fmt.Sprintf("RFD-%08d", refundID)
	if _, err := tx.ExecContext(ctx, `UPDATE refunds SET refund_no = ? WHERE id = ?`, refundNo, refundID); err != nil {
		return "", 0, err
	}
	return refundNo, refundID, nil
}

func (repository *MySQLRepository) loadServiceChangeLinkByInvoiceID(ctx context.Context, tx *sql.Tx, invoiceID int64) (serviceChangeLink, bool, error) {
	var item serviceChangeLink
	err := tx.QueryRowContext(ctx, `
SELECT id, service_id, order_id, invoice_id, action_name, title, status, reason, payload_json
FROM service_change_orders
WHERE invoice_id = ?
FOR UPDATE`, invoiceID).Scan(
		&item.ID,
		&item.ServiceID,
		&item.OrderID,
		&item.InvoiceID,
		&item.ActionName,
		&item.Title,
		&item.Status,
		&item.Reason,
		&item.PayloadJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return serviceChangeLink{}, false, nil
		}
		return serviceChangeLink{}, false, err
	}
	return item, true, nil
}

func (repository *MySQLRepository) loadServiceChangeLinkByOrderID(ctx context.Context, tx *sql.Tx, orderID int64) (serviceChangeLink, bool, error) {
	var item serviceChangeLink
	err := tx.QueryRowContext(ctx, `
SELECT id, service_id, order_id, invoice_id, action_name, title, status, reason, payload_json
FROM service_change_orders
WHERE order_id = ?
FOR UPDATE`, orderID).Scan(
		&item.ID,
		&item.ServiceID,
		&item.OrderID,
		&item.InvoiceID,
		&item.ActionName,
		&item.Title,
		&item.Status,
		&item.Reason,
		&item.PayloadJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return serviceChangeLink{}, false, nil
		}
		return serviceChangeLink{}, false, err
	}
	return item, true, nil
}

func parseNextDue(cycle string, now time.Time) time.Time {
	switch cycle {
	case "annual":
		return now.AddDate(1, 0, 0)
	case "quarterly":
		return now.AddDate(0, 3, 0)
	default:
		return now.AddDate(0, 1, 0)
	}
}

func firstNonEmptyString(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func resolveServiceProviderType(automationType string) string {
	return normalizeAutomationType(automationType)
}

func resolveProvisionedServiceStatus(automationType string) domain.ServiceStatus {
	if normalizeAutomationType(automationType) == "LOCAL" {
		return domain.ServiceStatusActive
	}
	return domain.ServiceStatusPending
}

func resolveProvisionedSyncStatus(automationType string) string {
	if normalizeAutomationType(automationType) == "LOCAL" {
		return "SUCCESS"
	}
	return "PENDING"
}

func resolveProvisionedSyncMessage(automationType string) string {
	providerType := normalizeAutomationType(automationType)
	if providerType == "LOCAL" {
		return ""
	}
	return fmt.Sprintf("待按自动化渠道 %s 开通", providerType)
}

func resolveProvisionedLastAction(automationType string) string {
	if normalizeAutomationType(automationType) == "LOCAL" {
		return "activate"
	}
	return "pending-provision"
}

type queryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func rollbackQuietly(tx *sql.Tx) {
	_ = tx.Rollback()
}

func (repository *MySQLRepository) fillOrderBillingMeta(ctx context.Context, queryer queryer, item *domain.Order) {
	row := queryer.QueryRowContext(ctx, `
SELECT
  IFNULL(
    (
      SELECT p.channel
      FROM payments p
      WHERE p.order_id = o.id
      ORDER BY p.id DESC
      LIMIT 1
    ),
    ''
  ) AS payment_channel,
  IFNULL(
    (
      SELECT i.status
      FROM invoices i
      WHERE i.order_id = o.id
      ORDER BY i.id DESC
      LIMIT 1
    ),
    'UNPAID'
  ) AS invoice_status
FROM orders o
WHERE o.id = ?`, item.ID)

	var (
		paymentChannel string
		invoiceStatus  string
	)
	if err := row.Scan(&paymentChannel, &invoiceStatus); err != nil {
		return
	}

	item.Payment = strings.TrimSpace(paymentChannel)
	switch strings.ToUpper(strings.TrimSpace(invoiceStatus)) {
	case "REFUNDED":
		item.PayStatus = "REFUNDED"
	case "PAID":
		item.PayStatus = "PAID"
	default:
		item.PayStatus = "UNPAID"
	}
}

func buildOrderFilterSQL(filter domain.OrderListFilter) (string, []any) {
	conditions := []string{"1 = 1"}
	args := make([]any, 0)

	if value := strings.TrimSpace(filter.Status); value != "" {
		conditions = append(conditions, "o.status = ?")
		args = append(args, strings.ToUpper(value))
	}
	if value := strings.TrimSpace(filter.OrderNo); value != "" {
		conditions = append(conditions, "o.order_no LIKE ?")
		args = append(args, "%"+value+"%")
	}
	if value := strings.TrimSpace(filter.ProductName); value != "" {
		conditions = append(conditions, "o.product_name LIKE ?")
		args = append(args, "%"+value+"%")
	}
	if filter.CustomerID > 0 {
		conditions = append(conditions, "o.customer_id = ?")
		args = append(args, filter.CustomerID)
	}
	if filter.HasAmount {
		conditions = append(conditions, "o.amount = ?")
		args = append(args, filter.Amount)
	}
	if value := strings.TrimSpace(filter.StartTime); value != "" {
		if parsed, ok := parseMySQLFilterTime(value, false); ok {
			conditions = append(conditions, "o.created_at >= ?")
			args = append(args, parsed)
		}
	}
	if value := strings.TrimSpace(filter.EndTime); value != "" {
		if parsed, ok := parseMySQLFilterTime(value, true); ok {
			conditions = append(conditions, "o.created_at <= ?")
			args = append(args, parsed)
		}
	}
	if value := strings.TrimSpace(filter.Payment); value != "" {
		conditions = append(conditions, `EXISTS (
SELECT 1
FROM payments p
WHERE p.order_id = o.id AND p.channel = ?
)`)
		args = append(args, strings.ToUpper(value))
	}
	if value := strings.TrimSpace(filter.PayStatus); value != "" {
		if strings.EqualFold(value, "PAID") {
			conditions = append(conditions, `EXISTS (
SELECT 1
FROM invoices i
WHERE i.order_id = o.id AND i.status = 'PAID'
)`)
		} else if strings.EqualFold(value, "REFUNDED") {
			conditions = append(conditions, `EXISTS (
SELECT 1
FROM invoices i
WHERE i.order_id = o.id AND i.status = 'REFUNDED'
)`)
		} else if strings.EqualFold(value, "UNPAID") {
			conditions = append(conditions, `EXISTS (
SELECT 1
FROM invoices i
WHERE i.order_id = o.id AND i.status = 'UNPAID'
)`)
		}
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}

func buildInvoiceFilterSQL(filter domain.InvoiceListFilter) (string, []any) {
	conditions := []string{"1 = 1"}
	args := make([]any, 0)

	if value := strings.TrimSpace(filter.Status); value != "" {
		conditions = append(conditions, "i.status = ?")
		args = append(args, strings.ToUpper(value))
	}
	if value := strings.TrimSpace(filter.InvoiceNo); value != "" {
		conditions = append(conditions, "i.invoice_no LIKE ?")
		args = append(args, "%"+value+"%")
	}
	if value := strings.TrimSpace(filter.OrderNo); value != "" {
		conditions = append(conditions, "i.order_no LIKE ?")
		args = append(args, "%"+value+"%")
	}
	if value := strings.TrimSpace(filter.ProductName); value != "" {
		conditions = append(conditions, "i.product_name LIKE ?")
		args = append(args, "%"+value+"%")
	}
	if value := strings.TrimSpace(filter.BillingCycle); value != "" {
		conditions = append(conditions, "i.billing_cycle = ?")
		args = append(args, value)
	}
	if filter.CustomerID > 0 {
		conditions = append(conditions, "i.customer_id = ?")
		args = append(args, filter.CustomerID)
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}

func buildServiceFilterSQL(filter domain.ServiceListFilter) (string, []any) {
	conditions := []string{"1 = 1"}
	args := make([]any, 0)

	if value := strings.TrimSpace(filter.Status); value != "" {
		conditions = append(conditions, "s.status = ?")
		args = append(args, strings.ToUpper(value))
	}
	if filter.CustomerID > 0 {
		conditions = append(conditions, "s.customer_id = ?")
		args = append(args, filter.CustomerID)
	}
	if filter.OrderID > 0 {
		conditions = append(conditions, "s.order_id = ?")
		args = append(args, filter.OrderID)
	}
	if value := strings.TrimSpace(filter.ProviderType); value != "" {
		conditions = append(conditions, "s.provider_type = ?")
		args = append(args, strings.ToUpper(value))
	}
	if filter.ProviderAccountID > 0 {
		conditions = append(conditions, "s.provider_account_id = ?")
		args = append(args, filter.ProviderAccountID)
	}
	if value := strings.TrimSpace(filter.SyncStatus); value != "" {
		conditions = append(conditions, "s.sync_status = ?")
		args = append(args, strings.ToUpper(value))
	}
	if value := strings.TrimSpace(filter.Keyword); value != "" {
		conditions = append(conditions, "(s.service_no LIKE ? OR s.product_name LIKE ? OR IFNULL(s.ip_address, '') LIKE ? OR IFNULL(s.provider_resource_id, '') LIKE ?)")
		pattern := "%" + value + "%"
		args = append(args, pattern, pattern, pattern, pattern)
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}

func buildAccountTransactionFilterSQL(filter domain.AccountTransactionFilter) (string, []any) {
	conditions := []string{"1 = 1"}
	args := make([]any, 0)

	if filter.CustomerID > 0 {
		conditions = append(conditions, "at.customer_id = ?")
		args = append(args, filter.CustomerID)
	}
	if value := strings.TrimSpace(filter.TransactionType); value != "" {
		tokens := make([]string, 0)
		for _, item := range strings.Split(value, ",") {
			token := strings.ToUpper(strings.TrimSpace(item))
			if token == "" {
				continue
			}
			tokens = append(tokens, token)
		}
		switch len(tokens) {
		case 1:
			conditions = append(conditions, "at.transaction_type = ?")
			args = append(args, tokens[0])
		default:
			if len(tokens) > 1 {
				placeholders := make([]string, len(tokens))
				for index, token := range tokens {
					placeholders[index] = "?"
					args = append(args, token)
				}
				conditions = append(conditions, "at.transaction_type IN ("+strings.Join(placeholders, ", ")+")")
			}
		}
	}
	if value := strings.TrimSpace(filter.Direction); value != "" {
		conditions = append(conditions, "at.direction = ?")
		args = append(args, strings.ToUpper(value))
	}
	if value := strings.TrimSpace(filter.Channel); value != "" {
		conditions = append(conditions, "at.channel = ?")
		args = append(args, strings.ToUpper(value))
	}
	if value := strings.TrimSpace(filter.Keyword); value != "" {
		pattern := "%" + value + "%"
		conditions = append(conditions, `(at.transaction_no LIKE ? OR c.name LIKE ? OR at.summary LIKE ? OR IFNULL(i.invoice_no, '') LIKE ? OR IFNULL(o.order_no, '') LIKE ? OR IFNULL(p.payment_no, '') LIKE ? OR IFNULL(r.refund_no, '') LIKE ?)`)
		args = append(args, pattern, pattern, pattern, pattern, pattern, pattern, pattern)
	}
	if value := strings.TrimSpace(filter.StartTime); value != "" {
		if parsed, ok := parseMySQLFilterTime(value, false); ok {
			conditions = append(conditions, "at.occurred_at >= ?")
			args = append(args, parsed)
		}
	}
	if value := strings.TrimSpace(filter.EndTime); value != "" {
		if parsed, ok := parseMySQLFilterTime(value, true); ok {
			conditions = append(conditions, "at.occurred_at <= ?")
			args = append(args, parsed)
		}
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}

func buildPaymentFilterSQL(filter domain.PaymentListFilter) (string, []any) {
	conditions := []string{"1 = 1"}
	args := make([]any, 0)

	if filter.CustomerID > 0 {
		conditions = append(conditions, "p.customer_id = ?")
		args = append(args, filter.CustomerID)
	}
	if filter.InvoiceID > 0 {
		conditions = append(conditions, "p.invoice_id = ?")
		args = append(args, filter.InvoiceID)
	}
	if value := strings.TrimSpace(filter.Keyword); value != "" {
		pattern := "%" + value + "%"
		conditions = append(conditions, "(p.payment_no LIKE ? OR IFNULL(p.trade_no, '') LIKE ?)")
		args = append(args, pattern, pattern)
	}
	if value := strings.TrimSpace(filter.Channel); value != "" {
		conditions = append(conditions, "p.channel = ?")
		args = append(args, strings.ToUpper(value))
	}
	if value := strings.TrimSpace(filter.Status); value != "" {
		conditions = append(conditions, "p.status = ?")
		args = append(args, strings.ToUpper(value))
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}

func buildRefundFilterSQL(filter domain.RefundListFilter) (string, []any) {
	conditions := []string{"1 = 1"}
	args := make([]any, 0)

	if filter.CustomerID > 0 {
		conditions = append(conditions, "r.customer_id = ?")
		args = append(args, filter.CustomerID)
	}
	if filter.InvoiceID > 0 {
		conditions = append(conditions, "r.invoice_id = ?")
		args = append(args, filter.InvoiceID)
	}
	if value := strings.TrimSpace(filter.Keyword); value != "" {
		pattern := "%" + value + "%"
		conditions = append(conditions, "(r.refund_no LIKE ? OR IFNULL(r.reason, '') LIKE ?)")
		args = append(args, pattern, pattern)
	}
	if value := strings.TrimSpace(filter.Status); value != "" {
		conditions = append(conditions, "r.status = ?")
		args = append(args, strings.ToUpper(value))
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}

func buildOrderRequestFilterSQL(filter domain.OrderRequestListFilter) (string, []any) {
	conditions := []string{"1 = 1"}
	args := make([]any, 0)

	if filter.OrderID > 0 {
		conditions = append(conditions, "orq.order_id = ?")
		args = append(args, filter.OrderID)
	}
	if filter.CustomerID > 0 {
		conditions = append(conditions, "orq.customer_id = ?")
		args = append(args, filter.CustomerID)
	}
	if value := strings.TrimSpace(filter.Type); value != "" {
		conditions = append(conditions, "orq.type = ?")
		args = append(args, strings.ToUpper(value))
	}
	if value := strings.TrimSpace(filter.Status); value != "" {
		conditions = append(conditions, "orq.status = ?")
		args = append(args, strings.ToUpper(value))
	}
	if value := strings.TrimSpace(filter.Keyword); value != "" {
		pattern := "%" + value + "%"
		conditions = append(conditions, "(orq.request_no LIKE ? OR orq.order_no LIKE ? OR orq.customer_name LIKE ? OR orq.product_name LIKE ? OR IFNULL(orq.summary, '') LIKE ? OR IFNULL(orq.reason, '') LIKE ?)")
		args = append(args, pattern, pattern, pattern, pattern, pattern, pattern)
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}

func buildServiceChangeOrderFilterSQL(filter domain.ServiceChangeOrderListFilter) (string, []any) {
	conditions := []string{"1 = 1"}
	args := make([]any, 0)

	if filter.ServiceID > 0 {
		conditions = append(conditions, "sco.service_id = ?")
		args = append(args, filter.ServiceID)
	}
	if filter.OrderID > 0 {
		conditions = append(conditions, "sco.order_id = ?")
		args = append(args, filter.OrderID)
	}
	if filter.InvoiceID > 0 {
		conditions = append(conditions, "sco.invoice_id = ?")
		args = append(args, filter.InvoiceID)
	}
	if value := strings.TrimSpace(filter.Status); value != "" {
		conditions = append(conditions, "sco.status = ?")
		args = append(args, strings.ToUpper(value))
	}
	if value := strings.TrimSpace(filter.Action); value != "" {
		conditions = append(conditions, "sco.action_name = ?")
		args = append(args, value)
	}
	if value := strings.TrimSpace(filter.Keyword); value != "" {
		pattern := "%" + value + "%"
		conditions = append(conditions, "(sco.title LIKE ? OR sco.action_name LIKE ? OR IFNULL(sco.reason, '') LIKE ? OR IFNULL(o.order_no, '') LIKE ? OR IFNULL(i.invoice_no, '') LIKE ? OR IFNULL(s.service_no, '') LIKE ?)")
		args = append(args, pattern, pattern, pattern, pattern, pattern, pattern)
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}

func resolveOrderSortColumn(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "amount":
		return "o.amount"
	case "ordernum":
		return "o.order_no"
	case "status":
		return "o.status"
	default:
		return "o.created_at"
	}
}

func resolveInvoiceSortColumn(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "invoice_no":
		return "i.invoice_no"
	case "due_at":
		return "i.due_at"
	case "amount", "total_amount":
		return "i.total_amount"
	case "paid_at":
		return "i.paid_at"
	default:
		return "i.created_at"
	}
}

func resolveServiceSortColumn(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "service_no":
		return "s.service_no"
	case "next_due_at":
		return "s.next_due_at"
	case "last_sync_at":
		return "s.last_sync_at"
	default:
		return "s.created_at"
	}
}

func resolveAccountTransactionSortColumn(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "amount":
		return "at.amount"
	case "transaction_no":
		return "at.transaction_no"
	default:
		return "at.occurred_at"
	}
}

func resolvePaymentSortColumn(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "amount":
		return "p.amount"
	case "payment_no":
		return "p.payment_no"
	default:
		return "p.paid_at"
	}
}

func resolveRefundSortColumn(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "amount":
		return "r.amount"
	case "refund_no":
		return "r.refund_no"
	default:
		return "r.created_at"
	}
}

func resolveOrderRequestSortColumn(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "request_no":
		return "orq.request_no"
	case "updated_at":
		return "orq.updated_at"
	case "processed_at":
		return "orq.processed_at"
	case "requested_amount":
		return "orq.requested_amount"
	default:
		return "orq.created_at"
	}
}

func resolveServiceChangeOrderSortColumn(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "amount":
		return "sco.amount"
	case "status":
		return "sco.status"
	case "paid_at":
		return "sco.paid_at"
	case "refunded_at":
		return "sco.refunded_at"
	default:
		return "sco.created_at"
	}
}

func resolveSortDirection(value string) string {
	if strings.EqualFold(strings.TrimSpace(value), "asc") {
		return "ASC"
	}
	return "DESC"
}

func firstNonEmptyMySQL(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func nullInt64(value int64) any {
	if value <= 0 {
		return nil
	}
	return value
}

func canTransitionOrderRequestStatusMySQL(current, next string) bool {
	current = strings.ToUpper(strings.TrimSpace(current))
	next = strings.ToUpper(strings.TrimSpace(next))
	if current == "" || next == "" || current == next {
		return false
	}

	switch current {
	case string(domain.OrderRequestStatusPending):
		return next == string(domain.OrderRequestStatusApproved) ||
			next == string(domain.OrderRequestStatusRejected) ||
			next == string(domain.OrderRequestStatusCancelled)
	case string(domain.OrderRequestStatusApproved):
		return next == string(domain.OrderRequestStatusCompleted) ||
			next == string(domain.OrderRequestStatusCancelled)
	default:
		return false
	}
}

func parseMySQLFilterTime(value string, endOfDay bool) (time.Time, bool) {
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		parsed, err := time.ParseInLocation(layout, value, time.Local)
		if err != nil {
			continue
		}
		if layout == "2006-01-02" && endOfDay {
			parsed = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
		return parsed, true
	}
	return time.Time{}, false
}

func parseMySQLEditableTime(value string) (time.Time, bool) {
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	for _, layout := range layouts {
		parsed, err := time.ParseInLocation(layout, value, time.Local)
		if err != nil {
			continue
		}
		return parsed, true
	}
	return time.Time{}, false
}

func absDecimal(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}

func normalizeOrderStatusValue(value, fallback string) (string, error) {
	status := strings.ToUpper(strings.TrimSpace(value))
	if status == "" {
		status = strings.ToUpper(strings.TrimSpace(fallback))
	}
	switch status {
	case string(domain.OrderStatusPending), string(domain.OrderStatusActive), string(domain.OrderStatusCompleted), string(domain.OrderStatusCancelled):
		return status, nil
	default:
		return "", fmt.Errorf("订单状态不合法")
	}
}

func normalizeInvoiceStatusValue(value, fallback string) (string, error) {
	status := strings.ToUpper(strings.TrimSpace(value))
	if status == "" {
		status = strings.ToUpper(strings.TrimSpace(fallback))
	}
	switch status {
	case string(domain.InvoiceStatusUnpaid), string(domain.InvoiceStatusPaid), string(domain.InvoiceStatusRefunded):
		return status, nil
	default:
		return "", fmt.Errorf("账单状态不合法")
	}
}

func normalizeServiceStatusValue(value, fallback string) (string, error) {
	status := strings.ToUpper(strings.TrimSpace(value))
	if status == "" {
		status = strings.ToUpper(strings.TrimSpace(fallback))
	}
	switch status {
	case string(domain.ServiceStatusPending), string(domain.ServiceStatusActive), string(domain.ServiceStatusSuspended), string(domain.ServiceStatusTerminated):
		return status, nil
	default:
		return "", fmt.Errorf("服务状态不合法")
	}
}

func normalizeProviderTypeValue(value, fallback string) (string, error) {
	providerType := strings.ToUpper(strings.TrimSpace(value))
	if providerType == "" {
		providerType = strings.ToUpper(strings.TrimSpace(fallback))
	}
	switch providerType {
	case "LOCAL", "MOFANG_CLOUD", "ZJMF_API", "RESOURCE", "MANUAL":
		return providerType, nil
	default:
		return "", fmt.Errorf("自动化渠道不合法")
	}
}

func normalizeSyncStatusValue(value, fallback string) (string, error) {
	syncStatus := strings.ToUpper(strings.TrimSpace(value))
	if syncStatus == "" {
		syncStatus = strings.ToUpper(strings.TrimSpace(fallback))
	}
	switch syncStatus {
	case "", "PENDING", "RUNNING", "SUCCESS", "FAILED":
		return syncStatus, nil
	default:
		return "", fmt.Errorf("同步状态不合法")
	}
}

func resolveManualLinkedOrderStatusForInvoice(invoiceStatus, fallback string) string {
	switch invoiceStatus {
	case string(domain.InvoiceStatusPaid):
		return string(domain.OrderStatusActive)
	case string(domain.InvoiceStatusUnpaid), string(domain.InvoiceStatusRefunded):
		return string(domain.OrderStatusPending)
	default:
		return fallback
	}
}

func resolveManualLinkedOrderStatusForService(serviceStatus, linkedInvoiceStatus string) string {
	switch serviceStatus {
	case string(domain.ServiceStatusActive):
		if linkedInvoiceStatus == string(domain.InvoiceStatusPaid) {
			return string(domain.OrderStatusActive)
		}
	case string(domain.ServiceStatusTerminated):
		if linkedInvoiceStatus == string(domain.InvoiceStatusRefunded) {
			return string(domain.OrderStatusPending)
		}
	}
	return ""
}

func resolveInvoicePaidAt(status string, existing string) any {
	switch status {
	case string(domain.InvoiceStatusUnpaid):
		return nil
	case string(domain.InvoiceStatusPaid), string(domain.InvoiceStatusRefunded):
		if parsed, ok := parseMySQLEditableTime(strings.TrimSpace(existing)); ok {
			return parsed
		}
		return time.Now()
	default:
		return nil
	}
}
