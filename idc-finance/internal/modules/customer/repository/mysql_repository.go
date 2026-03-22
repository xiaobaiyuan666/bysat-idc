package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"idc-finance/internal/modules/customer/domain"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (repository *MySQLRepository) List() []domain.Customer {
	rows, err := repository.db.Query(`SELECT id FROM customers ORDER BY id DESC`)
	if err != nil {
		log.Printf("customer mysql list failed: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]domain.Customer, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Printf("customer mysql scan id failed: %v", err)
			return result
		}
		item, err := repository.loadCustomer(context.Background(), repository.db, id)
		if err != nil {
			log.Printf("customer mysql load customer failed: %v", err)
			continue
		}
		result = append(result, item)
	}
	return result
}

func (repository *MySQLRepository) ListGroups() []domain.CustomerGroup {
	rows, err := repository.db.Query(`
SELECT g.id, g.name, g.description, COUNT(c.id) AS customer_count
FROM customer_groups g
LEFT JOIN customers c ON c.group_id = g.id
GROUP BY g.id, g.name, g.description
ORDER BY g.id`)
	if err != nil {
		log.Printf("customer mysql list groups failed: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]domain.CustomerGroup, 0)
	for rows.Next() {
		var item domain.CustomerGroup
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.CustomerCount); err != nil {
			log.Printf("customer mysql scan group failed: %v", err)
			return result
		}
		result = append(result, item)
	}
	return result
}

func (repository *MySQLRepository) ListLevels() []domain.CustomerLevel {
	rows, err := repository.db.Query(`
SELECT l.id, l.name, l.priority, l.description, COUNT(c.id) AS customer_count
FROM customer_levels l
LEFT JOIN customers c ON c.level_id = l.id
GROUP BY l.id, l.name, l.priority, l.description
ORDER BY l.priority DESC, l.id`)
	if err != nil {
		log.Printf("customer mysql list levels failed: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]domain.CustomerLevel, 0)
	for rows.Next() {
		var item domain.CustomerLevel
		if err := rows.Scan(&item.ID, &item.Name, &item.Priority, &item.Description, &item.CustomerCount); err != nil {
			log.Printf("customer mysql scan level failed: %v", err)
			return result
		}
		result = append(result, item)
	}
	return result
}

func (repository *MySQLRepository) CreateGroup(group domain.CustomerGroup) (domain.CustomerGroup, error) {
	result, err := repository.db.Exec(`
INSERT INTO customer_groups (name, description)
VALUES (?, ?)`,
		group.Name,
		group.Description,
	)
	if err != nil {
		return domain.CustomerGroup{}, err
	}

	groupID, err := result.LastInsertId()
	if err != nil {
		return domain.CustomerGroup{}, err
	}

	row := repository.db.QueryRow(`
SELECT g.id, g.name, g.description, COUNT(c.id) AS customer_count
FROM customer_groups g
LEFT JOIN customers c ON c.group_id = g.id
WHERE g.id = ?
GROUP BY g.id, g.name, g.description`, groupID)

	var created domain.CustomerGroup
	if err := row.Scan(&created.ID, &created.Name, &created.Description, &created.CustomerCount); err != nil {
		return domain.CustomerGroup{}, err
	}
	return created, nil
}

func (repository *MySQLRepository) UpdateGroup(id int64, group domain.CustomerGroup) (domain.CustomerGroup, error) {
	result, err := repository.db.Exec(`
UPDATE customer_groups
SET name = ?, description = ?, updated_at = NOW()
WHERE id = ?`,
		group.Name,
		group.Description,
		id,
	)
	if err != nil {
		return domain.CustomerGroup{}, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return domain.CustomerGroup{}, err
	}
	if affected == 0 {
		return domain.CustomerGroup{}, ErrGroupNotFound
	}

	row := repository.db.QueryRow(`
SELECT g.id, g.name, g.description, COUNT(c.id) AS customer_count
FROM customer_groups g
LEFT JOIN customers c ON c.group_id = g.id
WHERE g.id = ?
GROUP BY g.id, g.name, g.description`, id)

	var updated domain.CustomerGroup
	if err := row.Scan(&updated.ID, &updated.Name, &updated.Description, &updated.CustomerCount); err != nil {
		return domain.CustomerGroup{}, err
	}
	return updated, nil
}

func (repository *MySQLRepository) DeleteGroup(id int64) error {
	row := repository.db.QueryRow(`SELECT COUNT(1) FROM customers WHERE group_id = ?`, id)
	var customerCount int
	if err := row.Scan(&customerCount); err != nil {
		return err
	}
	if customerCount > 0 {
		return ErrGroupInUse
	}

	result, err := repository.db.Exec(`DELETE FROM customer_groups WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrGroupNotFound
	}
	return nil
}

func (repository *MySQLRepository) CreateLevel(level domain.CustomerLevel) (domain.CustomerLevel, error) {
	result, err := repository.db.Exec(`
INSERT INTO customer_levels (name, priority, description)
VALUES (?, ?, ?)`,
		level.Name,
		level.Priority,
		level.Description,
	)
	if err != nil {
		return domain.CustomerLevel{}, err
	}

	levelID, err := result.LastInsertId()
	if err != nil {
		return domain.CustomerLevel{}, err
	}

	row := repository.db.QueryRow(`
SELECT l.id, l.name, l.priority, l.description, COUNT(c.id) AS customer_count
FROM customer_levels l
LEFT JOIN customers c ON c.level_id = l.id
WHERE l.id = ?
GROUP BY l.id, l.name, l.priority, l.description`, levelID)

	var created domain.CustomerLevel
	if err := row.Scan(&created.ID, &created.Name, &created.Priority, &created.Description, &created.CustomerCount); err != nil {
		return domain.CustomerLevel{}, err
	}
	return created, nil
}

func (repository *MySQLRepository) UpdateLevel(id int64, level domain.CustomerLevel) (domain.CustomerLevel, error) {
	result, err := repository.db.Exec(`
UPDATE customer_levels
SET name = ?, priority = ?, description = ?, updated_at = NOW()
WHERE id = ?`,
		level.Name,
		level.Priority,
		level.Description,
		id,
	)
	if err != nil {
		return domain.CustomerLevel{}, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return domain.CustomerLevel{}, err
	}
	if affected == 0 {
		return domain.CustomerLevel{}, ErrLevelNotFound
	}

	row := repository.db.QueryRow(`
SELECT l.id, l.name, l.priority, l.description, COUNT(c.id) AS customer_count
FROM customer_levels l
LEFT JOIN customers c ON c.level_id = l.id
WHERE l.id = ?
GROUP BY l.id, l.name, l.priority, l.description`, id)

	var updated domain.CustomerLevel
	if err := row.Scan(&updated.ID, &updated.Name, &updated.Priority, &updated.Description, &updated.CustomerCount); err != nil {
		return domain.CustomerLevel{}, err
	}
	return updated, nil
}

func (repository *MySQLRepository) DeleteLevel(id int64) error {
	row := repository.db.QueryRow(`SELECT COUNT(1) FROM customers WHERE level_id = ?`, id)
	var customerCount int
	if err := row.Scan(&customerCount); err != nil {
		return err
	}
	if customerCount > 0 {
		return ErrLevelInUse
	}

	result, err := repository.db.Exec(`DELETE FROM customer_levels WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrLevelNotFound
	}
	return nil
}

func (repository *MySQLRepository) GetByID(id int64) (domain.Customer, bool) {
	item, err := repository.loadCustomer(context.Background(), repository.db, id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("customer mysql get customer failed: %v", err)
		}
		return domain.Customer{}, false
	}
	return item, true
}

func (repository *MySQLRepository) Create(customer domain.Customer) domain.Customer {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("customer mysql begin create failed: %v", err)
		return domain.Customer{}
	}
	defer rollbackCustomerTx(tx)

	groupID, err := ensureCustomerGroup(ctx, tx, customer.GroupName)
	if err != nil {
		log.Printf("customer mysql ensure group failed: %v", err)
		return domain.Customer{}
	}
	levelID, err := ensureCustomerLevel(ctx, tx, customer.LevelName)
	if err != nil {
		log.Printf("customer mysql ensure level failed: %v", err)
		return domain.Customer{}
	}

	customerType := firstCustomerValue(customer.Type, "PERSONAL")
	tempNo := fmt.Sprintf("TMP-CUST-%d", time.Now().UnixNano())
	result, err := tx.ExecContext(ctx, `
INSERT INTO customers (customer_no, customer_type, name, email, mobile, status, group_id, level_id, sales_owner, remarks)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		tempNo,
		customerType,
		customer.Name,
		customer.Email,
		customer.Mobile,
		firstCustomerValue(string(customer.Status), string(domain.CustomerStatusActive)),
		groupID,
		levelID,
		customer.SalesOwner,
		customer.Remarks,
	)
	if err != nil {
		log.Printf("customer mysql insert customer failed: %v", err)
		return domain.Customer{}
	}

	customerID, err := result.LastInsertId()
	if err != nil {
		log.Printf("customer mysql last insert id failed: %v", err)
		return domain.Customer{}
	}

	customerNo := fmt.Sprintf("CUST-%08d", customerID)
	if _, err := tx.ExecContext(ctx, `UPDATE customers SET customer_no = ? WHERE id = ?`, customerNo, customerID); err != nil {
		log.Printf("customer mysql update customer no failed: %v", err)
		return domain.Customer{}
	}

	if _, err := tx.ExecContext(ctx, `
INSERT INTO customer_identities (customer_id, identity_type, verify_status, subject_name, cert_no, country_code, review_remark, submitted_at)
VALUES (?, ?, ?, ?, ?, ?, ?, NOW())`,
		customerID,
		customerType,
		domain.IdentityStatusPending,
		customer.Name,
		"",
		"CN",
		"",
	); err != nil {
		log.Printf("customer mysql insert identity failed: %v", err)
		return domain.Customer{}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("customer mysql commit create failed: %v", err)
		return domain.Customer{}
	}

	created, ok := repository.GetByID(customerID)
	if !ok {
		return domain.Customer{}
	}
	return created
}

func (repository *MySQLRepository) Update(id int64, customer domain.Customer) (domain.Customer, bool) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("customer mysql begin update failed: %v", err)
		return domain.Customer{}, false
	}
	defer rollbackCustomerTx(tx)

	current, err := repository.loadCustomer(ctx, tx, id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("customer mysql load update customer failed: %v", err)
		}
		return domain.Customer{}, false
	}

	groupID, err := ensureCustomerGroup(ctx, tx, customer.GroupName)
	if err != nil {
		log.Printf("customer mysql ensure group failed: %v", err)
		return domain.Customer{}, false
	}
	levelID, err := ensureCustomerLevel(ctx, tx, customer.LevelName)
	if err != nil {
		log.Printf("customer mysql ensure level failed: %v", err)
		return domain.Customer{}, false
	}

	result, err := tx.ExecContext(ctx, `
UPDATE customers
SET name = ?, email = ?, mobile = ?, status = ?, group_id = ?, level_id = ?, sales_owner = ?, remarks = ?, updated_at = NOW()
WHERE id = ?`,
		customer.Name,
		customer.Email,
		customer.Mobile,
		firstCustomerValue(string(customer.Status), string(current.Status)),
		groupID,
		levelID,
		customer.SalesOwner,
		customer.Remarks,
		id,
	)
	if err != nil {
		log.Printf("customer mysql update customer failed: %v", err)
		return domain.Customer{}, false
	}

	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		return domain.Customer{}, false
	}

	if err := tx.Commit(); err != nil {
		log.Printf("customer mysql commit update failed: %v", err)
		return domain.Customer{}, false
	}

	updated, ok := repository.GetByID(id)
	return updated, ok
}

func (repository *MySQLRepository) AddContact(customerID int64, contact domain.Contact) (domain.Customer, domain.Contact, bool) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("customer mysql begin add contact failed: %v", err)
		return domain.Customer{}, domain.Contact{}, false
	}
	defer rollbackCustomerTx(tx)

	if !repository.customerExists(ctx, tx, customerID) {
		return domain.Customer{}, domain.Contact{}, false
	}

	if contact.IsPrimary {
		if _, err := tx.ExecContext(ctx, `UPDATE customer_contacts SET is_primary = 0 WHERE customer_id = ?`, customerID); err != nil {
			log.Printf("customer mysql reset primary contact failed: %v", err)
			return domain.Customer{}, domain.Contact{}, false
		}
	}

	result, err := tx.ExecContext(ctx, `
INSERT INTO customer_contacts (customer_id, name, email, mobile, role_name, is_primary)
VALUES (?, ?, ?, ?, ?, ?)`,
		customerID,
		contact.Name,
		contact.Email,
		contact.Mobile,
		contact.RoleName,
		contact.IsPrimary,
	)
	if err != nil {
		log.Printf("customer mysql insert contact failed: %v", err)
		return domain.Customer{}, domain.Contact{}, false
	}

	contactID, err := result.LastInsertId()
	if err != nil {
		log.Printf("customer mysql last contact id failed: %v", err)
		return domain.Customer{}, domain.Contact{}, false
	}

	if err := tx.Commit(); err != nil {
		log.Printf("customer mysql commit add contact failed: %v", err)
		return domain.Customer{}, domain.Contact{}, false
	}

	customer, ok := repository.GetByID(customerID)
	if !ok {
		return domain.Customer{}, domain.Contact{}, false
	}
	createdContact, ok := repository.loadContact(context.Background(), repository.db, customerID, contactID)
	if !ok {
		return domain.Customer{}, domain.Contact{}, false
	}
	return customer, createdContact, true
}

func (repository *MySQLRepository) UpdateContact(customerID, contactID int64, contact domain.Contact) (domain.Customer, domain.Contact, bool) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("customer mysql begin update contact failed: %v", err)
		return domain.Customer{}, domain.Contact{}, false
	}
	defer rollbackCustomerTx(tx)

	existing, ok := repository.loadContact(ctx, tx, customerID, contactID)
	if !ok {
		return domain.Customer{}, domain.Contact{}, false
	}

	if contact.IsPrimary {
		if _, err := tx.ExecContext(ctx, `UPDATE customer_contacts SET is_primary = 0 WHERE customer_id = ?`, customerID); err != nil {
			log.Printf("customer mysql reset primary contact failed: %v", err)
			return domain.Customer{}, domain.Contact{}, false
		}
	} else {
		contact.IsPrimary = existing.IsPrimary
	}

	result, err := tx.ExecContext(ctx, `
UPDATE customer_contacts
SET name = ?, email = ?, mobile = ?, role_name = ?, is_primary = ?, updated_at = NOW()
WHERE id = ? AND customer_id = ?`,
		contact.Name,
		contact.Email,
		contact.Mobile,
		contact.RoleName,
		contact.IsPrimary,
		contactID,
		customerID,
	)
	if err != nil {
		log.Printf("customer mysql update contact failed: %v", err)
		return domain.Customer{}, domain.Contact{}, false
	}

	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		return domain.Customer{}, domain.Contact{}, false
	}

	if err := tx.Commit(); err != nil {
		log.Printf("customer mysql commit update contact failed: %v", err)
		return domain.Customer{}, domain.Contact{}, false
	}

	customer, customerOK := repository.GetByID(customerID)
	updatedContact, contactOK := repository.loadContact(context.Background(), repository.db, customerID, contactID)
	return customer, updatedContact, customerOK && contactOK
}

func (repository *MySQLRepository) DeleteContact(customerID, contactID int64) (domain.Customer, bool) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("customer mysql begin delete contact failed: %v", err)
		return domain.Customer{}, false
	}
	defer rollbackCustomerTx(tx)

	existing, ok := repository.loadContact(ctx, tx, customerID, contactID)
	if !ok {
		return domain.Customer{}, false
	}

	result, err := tx.ExecContext(ctx, `DELETE FROM customer_contacts WHERE id = ? AND customer_id = ?`, contactID, customerID)
	if err != nil {
		log.Printf("customer mysql delete contact failed: %v", err)
		return domain.Customer{}, false
	}

	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		return domain.Customer{}, false
	}

	if existing.IsPrimary {
		row := tx.QueryRowContext(ctx, `SELECT id FROM customer_contacts WHERE customer_id = ? ORDER BY id LIMIT 1`, customerID)
		var newPrimaryID int64
		if err := row.Scan(&newPrimaryID); err == nil {
			if _, err := tx.ExecContext(ctx, `UPDATE customer_contacts SET is_primary = 1 WHERE id = ?`, newPrimaryID); err != nil {
				log.Printf("customer mysql promote primary contact failed: %v", err)
				return domain.Customer{}, false
			}
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("customer mysql commit delete contact failed: %v", err)
		return domain.Customer{}, false
	}

	customer, ok := repository.GetByID(customerID)
	return customer, ok
}

func (repository *MySQLRepository) ReviewIdentity(customerID int64, status domain.IdentityStatus, reviewRemark string) (domain.Customer, bool) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("customer mysql begin review identity failed: %v", err)
		return domain.Customer{}, false
	}
	defer rollbackCustomerTx(tx)

	customer, err := repository.loadCustomer(ctx, tx, customerID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("customer mysql load customer for identity review failed: %v", err)
		}
		return domain.Customer{}, false
	}

	if _, err := tx.ExecContext(ctx, `
INSERT INTO customer_identities (customer_id, identity_type, verify_status, subject_name, cert_no, country_code, review_remark, reviewed_at, submitted_at)
VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
ON DUPLICATE KEY UPDATE
  identity_type = VALUES(identity_type),
  verify_status = VALUES(verify_status),
  subject_name = VALUES(subject_name),
  cert_no = VALUES(cert_no),
  country_code = VALUES(country_code),
  review_remark = VALUES(review_remark),
  reviewed_at = VALUES(reviewed_at),
  updated_at = NOW()`,
		customerID,
		firstCustomerValue(customer.Type, "PERSONAL"),
		status,
		customerIdentitySubject(customer),
		customerIdentityCertNo(customer),
		firstCustomerValue(customerIdentityCountryCode(customer), "CN"),
		reviewRemark,
	); err != nil {
		log.Printf("customer mysql upsert identity failed: %v", err)
		return domain.Customer{}, false
	}

	if err := tx.Commit(); err != nil {
		log.Printf("customer mysql commit review identity failed: %v", err)
		return domain.Customer{}, false
	}

	updatedCustomer, ok := repository.GetByID(customerID)
	return updatedCustomer, ok
}

func (repository *MySQLRepository) ListServiceItems(customerID int64) []domain.RelatedItem {
	return repository.listRelatedItems(`
SELECT service_no, product_name, status,
       '',
       IFNULL(DATE_FORMAT(next_due_at, '%Y-%m-%d'), ''),
       '',
       CONCAT('操作类型：', IFNULL(last_action, ''))
FROM services
WHERE customer_id = ?
ORDER BY id DESC`, customerID)
}

func (repository *MySQLRepository) ListInvoiceItems(customerID int64) []domain.RelatedItem {
	return repository.listRelatedItems(`
SELECT invoice_no, product_name, status,
       CAST(total_amount AS CHAR),
       IFNULL(DATE_FORMAT(due_at, '%Y-%m-%d'), ''),
       '',
       ''
FROM invoices
WHERE customer_id = ?
ORDER BY id DESC`, customerID)
}

func (repository *MySQLRepository) ListTicketItems(customerID int64) []domain.RelatedItem {
	return repository.listRelatedItems(`
SELECT ticket_no, title, status,
       '',
       '',
       IFNULL(DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s'), ''),
       ''
FROM tickets
WHERE customer_id = ?
ORDER BY id DESC`, customerID)
}

func (repository *MySQLRepository) listRelatedItems(query string, customerID int64) []domain.RelatedItem {
	rows, err := repository.db.Query(query, customerID)
	if err != nil {
		log.Printf("customer mysql list related items failed: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]domain.RelatedItem, 0)
	for rows.Next() {
		var item domain.RelatedItem
		if err := rows.Scan(&item.No, &item.Name, &item.Status, &item.Amount, &item.DueAt, &item.UpdatedAt, &item.Description); err != nil {
			log.Printf("customer mysql scan related item failed: %v", err)
			return result
		}
		result = append(result, item)
	}
	return result
}

func (repository *MySQLRepository) loadCustomer(ctx context.Context, queryer customerQueryer, id int64) (domain.Customer, error) {
	row := queryer.QueryRowContext(ctx, `
SELECT
  c.id,
  c.customer_no,
  c.customer_type,
  c.name,
  c.email,
  c.mobile,
  c.status,
  COALESCE(g.name, ''),
  COALESCE(l.name, ''),
  c.sales_owner,
  c.remarks
FROM customers c
LEFT JOIN customer_groups g ON g.id = c.group_id
LEFT JOIN customer_levels l ON l.id = c.level_id
WHERE c.id = ?`, id)

	var (
		item     domain.Customer
		statusDB string
	)
	if err := row.Scan(
		&item.ID,
		&item.CustomerNo,
		&item.Type,
		&item.Name,
		&item.Email,
		&item.Mobile,
		&statusDB,
		&item.GroupName,
		&item.LevelName,
		&item.SalesOwner,
		&item.Remarks,
	); err != nil {
		return domain.Customer{}, err
	}

	item.Status = domain.CustomerStatus(statusDB)
	item.Contacts = repository.loadContacts(ctx, queryer, item.ID)
	identity, ok := repository.loadIdentity(ctx, queryer, item.ID)
	if ok {
		item.Identity = identity
	}
	return item, nil
}

func (repository *MySQLRepository) loadContacts(ctx context.Context, queryer customerQueryer, customerID int64) []domain.Contact {
	rows, err := queryer.QueryContext(ctx, `
SELECT id, name, email, mobile, role_name, is_primary
FROM customer_contacts
WHERE customer_id = ?
ORDER BY is_primary DESC, id`, customerID)
	if err != nil {
		log.Printf("customer mysql load contacts failed: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]domain.Contact, 0)
	for rows.Next() {
		var (
			item        domain.Contact
			isPrimaryDB bool
		)
		if err := rows.Scan(&item.ID, &item.Name, &item.Email, &item.Mobile, &item.RoleName, &isPrimaryDB); err != nil {
			log.Printf("customer mysql scan contact failed: %v", err)
			return result
		}
		item.IsPrimary = isPrimaryDB
		result = append(result, item)
	}
	return result
}

func (repository *MySQLRepository) loadIdentity(ctx context.Context, queryer customerQueryer, customerID int64) (*domain.Identity, bool) {
	row := queryer.QueryRowContext(ctx, `
SELECT id, identity_type, verify_status, subject_name, cert_no, country_code,
       COALESCE(review_remark, ''),
       IFNULL(DATE_FORMAT(reviewed_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(submitted_at, '%Y-%m-%d %H:%i:%s'), '')
FROM customer_identities
WHERE customer_id = ?`, customerID)

	var (
		item     domain.Identity
		statusDB string
	)
	if err := row.Scan(
		&item.ID,
		&item.IdentityType,
		&statusDB,
		&item.SubjectName,
		&item.CertNo,
		&item.CountryCode,
		&item.ReviewRemark,
		&item.ReviewedAt,
		&item.SubmittedAt,
	); err != nil {
		if err != sql.ErrNoRows {
			log.Printf("customer mysql load identity failed: %v", err)
		}
		return nil, false
	}
	item.VerifyStatus = domain.IdentityStatus(statusDB)
	return &item, true
}

func (repository *MySQLRepository) loadContact(ctx context.Context, queryer customerQueryer, customerID, contactID int64) (domain.Contact, bool) {
	row := queryer.QueryRowContext(ctx, `
SELECT id, name, email, mobile, role_name, is_primary
FROM customer_contacts
WHERE customer_id = ? AND id = ?`, customerID, contactID)

	var (
		item        domain.Contact
		isPrimaryDB bool
	)
	if err := row.Scan(&item.ID, &item.Name, &item.Email, &item.Mobile, &item.RoleName, &isPrimaryDB); err != nil {
		return domain.Contact{}, false
	}
	item.IsPrimary = isPrimaryDB
	return item, true
}

func (repository *MySQLRepository) customerExists(ctx context.Context, queryer customerQueryer, customerID int64) bool {
	row := queryer.QueryRowContext(ctx, `SELECT id FROM customers WHERE id = ?`, customerID)
	var id int64
	return row.Scan(&id) == nil
}

func ensureCustomerGroup(ctx context.Context, tx *sql.Tx, groupName string) (sql.NullInt64, error) {
	groupName = strings.TrimSpace(groupName)
	if groupName == "" {
		return sql.NullInt64{}, nil
	}

	result, err := tx.ExecContext(ctx, `
INSERT INTO customer_groups (name, description)
VALUES (?, ?)
ON DUPLICATE KEY UPDATE id = LAST_INSERT_ID(id), description = VALUES(description)`,
		groupName,
		groupName+" 自动维护",
	)
	if err != nil {
		return sql.NullInt64{}, err
	}

	groupID, err := result.LastInsertId()
	if err != nil {
		return sql.NullInt64{}, err
	}
	return sql.NullInt64{Int64: groupID, Valid: true}, nil
}

func ensureCustomerLevel(ctx context.Context, tx *sql.Tx, levelName string) (sql.NullInt64, error) {
	levelName = strings.TrimSpace(levelName)
	if levelName == "" {
		return sql.NullInt64{}, nil
	}

	result, err := tx.ExecContext(ctx, `
INSERT INTO customer_levels (name, priority, description)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE id = LAST_INSERT_ID(id), description = VALUES(description)`,
		levelName,
		50,
		levelName+" 自动维护",
	)
	if err != nil {
		return sql.NullInt64{}, err
	}

	levelID, err := result.LastInsertId()
	if err != nil {
		return sql.NullInt64{}, err
	}
	return sql.NullInt64{Int64: levelID, Valid: true}, nil
}

func customerIdentitySubject(customer domain.Customer) string {
	if customer.Identity != nil && strings.TrimSpace(customer.Identity.SubjectName) != "" {
		return customer.Identity.SubjectName
	}
	return customer.Name
}

func customerIdentityCertNo(customer domain.Customer) string {
	if customer.Identity != nil {
		return customer.Identity.CertNo
	}
	return ""
}

func customerIdentityCountryCode(customer domain.Customer) string {
	if customer.Identity != nil && strings.TrimSpace(customer.Identity.CountryCode) != "" {
		return customer.Identity.CountryCode
	}
	return "CN"
}

func firstCustomerValue(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

type customerQueryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func rollbackCustomerTx(tx *sql.Tx) {
	_ = tx.Rollback()
}
