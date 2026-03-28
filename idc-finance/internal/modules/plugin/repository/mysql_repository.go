package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"idc-finance/internal/modules/plugin/domain"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (repository *MySQLRepository) ListDefinitions() []domain.PluginDefinition {
	return defaultDefinitions()
}

func (repository *MySQLRepository) GetConfig(code string) (domain.PluginConfig, bool) {
	normalized := normalizeCode(code)
	if normalized == "" {
		return domain.PluginConfig{}, false
	}
	row := repository.db.QueryRow(`
SELECT code, name, enabled, config_json,
       IFNULL(updated_by, ''),
       IFNULL(DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s'), '')
FROM plugin_configs
WHERE code = ?`, normalized)

	var (
		item       domain.PluginConfig
		enabledInt int
		configJSON string
	)
	if err := row.Scan(&item.Code, &item.Name, &enabledInt, &configJSON, &item.UpdatedBy, &item.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			definition, ok := lookupDefinition(normalized)
			if !ok {
				return domain.PluginConfig{}, false
			}
			return domain.PluginConfig{
				Code:    definition.Code,
				Name:    definition.Name,
				Enabled: true,
				Config:  map[string]any{},
			}, true
		}
		return domain.PluginConfig{}, false
	}

	item.Enabled = enabledInt == 1
	item.Config = decodeJSONMap(configJSON)
	return item, true
}

func (repository *MySQLRepository) ListConfigs() []domain.PluginConfig {
	rows, err := repository.db.Query(`
SELECT code, name, enabled, config_json,
       IFNULL(updated_by, ''),
       IFNULL(DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s'), '')
FROM plugin_configs
ORDER BY code ASC`)
	if err != nil {
		return []domain.PluginConfig{}
	}
	defer rows.Close()

	items := make([]domain.PluginConfig, 0)
	for rows.Next() {
		var (
			item       domain.PluginConfig
			enabledInt int
			configJSON string
		)
		if scanErr := rows.Scan(&item.Code, &item.Name, &enabledInt, &configJSON, &item.UpdatedBy, &item.UpdatedAt); scanErr != nil {
			continue
		}
		item.Enabled = enabledInt == 1
		item.Config = decodeJSONMap(configJSON)
		items = append(items, item)
	}
	for _, definition := range defaultDefinitions() {
		exists := false
		for _, item := range items {
			if item.Code == definition.Code {
				exists = true
				break
			}
		}
		if !exists {
			items = append(items, domain.PluginConfig{
				Code:    definition.Code,
				Name:    definition.Name,
				Enabled: true,
				Config:  map[string]any{},
			})
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Code < items[j].Code
	})
	return items
}

func (repository *MySQLRepository) SaveConfig(code string, enabled bool, config map[string]any, updatedBy string) (domain.PluginConfig, error) {
	normalized := normalizeCode(code)
	if normalized == "" {
		return domain.PluginConfig{}, fmt.Errorf("plugin not found: %s", code)
	}
	definition, ok := lookupDefinition(normalized)
	if !ok {
		return domain.PluginConfig{}, fmt.Errorf("plugin not found: %s", code)
	}

	payload, _ := json.Marshal(cloneMap(config))
	enabledInt := 0
	if enabled {
		enabledInt = 1
	}

	_, err := repository.db.Exec(`
INSERT INTO plugin_configs(code, name, enabled, config_json, updated_by, updated_at)
VALUES(?, ?, ?, ?, ?, NOW())
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  enabled = VALUES(enabled),
  config_json = VALUES(config_json),
  updated_by = VALUES(updated_by),
  updated_at = NOW()`,
		normalized,
		definition.Name,
		enabledInt,
		string(payload),
		strings.TrimSpace(updatedBy),
	)
	if err != nil {
		return domain.PluginConfig{}, err
	}

	item, exists := repository.GetConfig(normalized)
	if !exists {
		return domain.PluginConfig{}, fmt.Errorf("plugin config save failed")
	}
	return item, nil
}

func (repository *MySQLRepository) CreateKYCRecord(record domain.KYCRecord) (domain.KYCRecord, error) {
	rawJSON, _ := json.Marshal(cloneMap(record.Raw))
	result, err := repository.db.Exec(`
INSERT INTO plugin_kyc_records(plugin_code, customer_id, real_name, id_card_no, mobile, status, message, raw_json, created_at)
VALUES(?, ?, ?, ?, ?, ?, ?, ?, NOW())`,
		normalizeCode(record.PluginCode),
		record.CustomerID,
		strings.TrimSpace(record.RealName),
		strings.TrimSpace(record.IDCardNo),
		strings.TrimSpace(record.Mobile),
		strings.TrimSpace(record.Status),
		strings.TrimSpace(record.Message),
		string(rawJSON),
	)
	if err != nil {
		return domain.KYCRecord{}, err
	}
	id, _ := result.LastInsertId()
	return repository.getKYCRecordByID(id)
}

func (repository *MySQLRepository) ListKYCRecords(customerID int64, limit int) []domain.KYCRecord {
	if limit <= 0 {
		limit = 20
	}
	rows, err := repository.db.Query(`
SELECT id, plugin_code, customer_id, real_name, id_card_no, mobile, status, message,
       IFNULL(raw_json, '{}'),
       IFNULL(DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), '')
FROM plugin_kyc_records
WHERE customer_id = ?
ORDER BY id DESC
LIMIT ?`, customerID, limit)
	if err != nil {
		return []domain.KYCRecord{}
	}
	defer rows.Close()

	items := make([]domain.KYCRecord, 0)
	for rows.Next() {
		var (
			item    domain.KYCRecord
			rawJSON string
		)
		if scanErr := rows.Scan(
			&item.ID,
			&item.PluginCode,
			&item.CustomerID,
			&item.RealName,
			&item.IDCardNo,
			&item.Mobile,
			&item.Status,
			&item.Message,
			&rawJSON,
			&item.CreatedAt,
		); scanErr != nil {
			continue
		}
		item.Raw = decodeJSONMap(rawJSON)
		items = append(items, item)
	}
	return items
}

func (repository *MySQLRepository) CreateSMSRecord(record domain.SMSRecord) (domain.SMSRecord, error) {
	rawJSON, _ := json.Marshal(cloneMap(record.Raw))
	result, err := repository.db.Exec(`
INSERT INTO plugin_sms_records(plugin_code, customer_id, mobile, template_code, content, status, message, biz_no, raw_json, created_at)
VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`,
		normalizeCode(record.PluginCode),
		record.CustomerID,
		strings.TrimSpace(record.Mobile),
		strings.TrimSpace(record.TemplateCode),
		strings.TrimSpace(record.Content),
		strings.TrimSpace(record.Status),
		strings.TrimSpace(record.Message),
		strings.TrimSpace(record.BizNo),
		string(rawJSON),
	)
	if err != nil {
		return domain.SMSRecord{}, err
	}
	id, _ := result.LastInsertId()
	return repository.getSMSRecordByID(id)
}

func (repository *MySQLRepository) ListSMSRecords(customerID int64, limit int) []domain.SMSRecord {
	if limit <= 0 {
		limit = 20
	}
	rows, err := repository.db.Query(`
SELECT id, plugin_code, customer_id, mobile, template_code, content, status, message, biz_no,
       IFNULL(raw_json, '{}'),
       IFNULL(DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), '')
FROM plugin_sms_records
WHERE customer_id = ?
ORDER BY id DESC
LIMIT ?`, customerID, limit)
	if err != nil {
		return []domain.SMSRecord{}
	}
	defer rows.Close()

	items := make([]domain.SMSRecord, 0)
	for rows.Next() {
		var (
			item    domain.SMSRecord
			rawJSON string
		)
		if scanErr := rows.Scan(
			&item.ID,
			&item.PluginCode,
			&item.CustomerID,
			&item.Mobile,
			&item.TemplateCode,
			&item.Content,
			&item.Status,
			&item.Message,
			&item.BizNo,
			&rawJSON,
			&item.CreatedAt,
		); scanErr != nil {
			continue
		}
		item.Raw = decodeJSONMap(rawJSON)
		items = append(items, item)
	}
	return items
}

func (repository *MySQLRepository) CreatePaymentOrder(order domain.PaymentOrder) (domain.PaymentOrder, error) {
	rawJSON, _ := json.Marshal(cloneMap(order.Raw))
	result, err := repository.db.Exec(`
INSERT INTO plugin_payment_orders(plugin_code, order_no, customer_id, amount, currency, subject, status, pay_url, notify_url, return_url, trade_no, raw_json, paid_at, created_at, updated_at)
VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULL, NOW(), NOW())`,
		normalizeCode(order.PluginCode),
		strings.TrimSpace(order.OrderNo),
		order.CustomerID,
		order.Amount,
		firstNonEmpty(strings.TrimSpace(order.Currency), "CNY"),
		strings.TrimSpace(order.Subject),
		firstNonEmpty(strings.TrimSpace(order.Status), "UNPAID"),
		strings.TrimSpace(order.PayURL),
		strings.TrimSpace(order.NotifyURL),
		strings.TrimSpace(order.ReturnURL),
		strings.TrimSpace(order.TradeNo),
		string(rawJSON),
	)
	if err != nil {
		return domain.PaymentOrder{}, err
	}
	id, _ := result.LastInsertId()
	return repository.getPaymentOrderByID(id)
}

func (repository *MySQLRepository) GetPaymentOrderByOrderNo(orderNo string) (domain.PaymentOrder, bool) {
	row := repository.db.QueryRow(`
SELECT id, plugin_code, order_no, customer_id, amount, currency, subject, status, pay_url, notify_url, return_url, trade_no,
       IFNULL(raw_json, '{}'),
       IFNULL(DATE_FORMAT(paid_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s'), '')
FROM plugin_payment_orders
WHERE order_no = ?`, strings.TrimSpace(orderNo))

	item, err := scanPaymentOrder(row)
	if err != nil {
		return domain.PaymentOrder{}, false
	}
	return item, true
}

func (repository *MySQLRepository) UpdatePaymentOrder(order domain.PaymentOrder) (domain.PaymentOrder, bool, error) {
	rawJSON, _ := json.Marshal(cloneMap(order.Raw))
	paidAt := strings.TrimSpace(order.PaidAt)
	if paidAt == "" {
		_, err := repository.db.Exec(`
UPDATE plugin_payment_orders
SET plugin_code=?, customer_id=?, amount=?, currency=?, subject=?, status=?, pay_url=?, notify_url=?, return_url=?, trade_no=?, raw_json=?, paid_at=NULL, updated_at=NOW()
WHERE order_no=?`,
			normalizeCode(order.PluginCode),
			order.CustomerID,
			order.Amount,
			firstNonEmpty(strings.TrimSpace(order.Currency), "CNY"),
			strings.TrimSpace(order.Subject),
			strings.TrimSpace(order.Status),
			strings.TrimSpace(order.PayURL),
			strings.TrimSpace(order.NotifyURL),
			strings.TrimSpace(order.ReturnURL),
			strings.TrimSpace(order.TradeNo),
			string(rawJSON),
			strings.TrimSpace(order.OrderNo),
		)
		if err != nil {
			return domain.PaymentOrder{}, false, err
		}
	} else {
		_, err := repository.db.Exec(`
UPDATE plugin_payment_orders
SET plugin_code=?, customer_id=?, amount=?, currency=?, subject=?, status=?, pay_url=?, notify_url=?, return_url=?, trade_no=?, raw_json=?, paid_at=?, updated_at=NOW()
WHERE order_no=?`,
			normalizeCode(order.PluginCode),
			order.CustomerID,
			order.Amount,
			firstNonEmpty(strings.TrimSpace(order.Currency), "CNY"),
			strings.TrimSpace(order.Subject),
			strings.TrimSpace(order.Status),
			strings.TrimSpace(order.PayURL),
			strings.TrimSpace(order.NotifyURL),
			strings.TrimSpace(order.ReturnURL),
			strings.TrimSpace(order.TradeNo),
			string(rawJSON),
			paidAt,
			strings.TrimSpace(order.OrderNo),
		)
		if err != nil {
			return domain.PaymentOrder{}, false, err
		}
	}

	updated, exists := repository.GetPaymentOrderByOrderNo(order.OrderNo)
	if !exists {
		return domain.PaymentOrder{}, false, nil
	}
	return updated, true, nil
}

func (repository *MySQLRepository) ListPaymentOrders(customerID int64, limit int) []domain.PaymentOrder {
	if limit <= 0 {
		limit = 20
	}
	rows, err := repository.db.Query(`
SELECT id, plugin_code, order_no, customer_id, amount, currency, subject, status, pay_url, notify_url, return_url, trade_no,
       IFNULL(raw_json, '{}'),
       IFNULL(DATE_FORMAT(paid_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s'), '')
FROM plugin_payment_orders
WHERE customer_id = ?
ORDER BY id DESC
LIMIT ?`, customerID, limit)
	if err != nil {
		return []domain.PaymentOrder{}
	}
	defer rows.Close()

	items := make([]domain.PaymentOrder, 0)
	for rows.Next() {
		item, err := scanPaymentOrder(rows)
		if err != nil {
			continue
		}
		items = append(items, item)
	}
	return items
}

func (repository *MySQLRepository) getKYCRecordByID(id int64) (domain.KYCRecord, error) {
	row := repository.db.QueryRow(`
SELECT id, plugin_code, customer_id, real_name, id_card_no, mobile, status, message,
       IFNULL(raw_json, '{}'),
       IFNULL(DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), '')
FROM plugin_kyc_records
WHERE id = ?`, id)

	var (
		item    domain.KYCRecord
		rawJSON string
	)
	if err := row.Scan(
		&item.ID,
		&item.PluginCode,
		&item.CustomerID,
		&item.RealName,
		&item.IDCardNo,
		&item.Mobile,
		&item.Status,
		&item.Message,
		&rawJSON,
		&item.CreatedAt,
	); err != nil {
		return domain.KYCRecord{}, err
	}
	item.Raw = decodeJSONMap(rawJSON)
	return item, nil
}

func (repository *MySQLRepository) getSMSRecordByID(id int64) (domain.SMSRecord, error) {
	row := repository.db.QueryRow(`
SELECT id, plugin_code, customer_id, mobile, template_code, content, status, message, biz_no,
       IFNULL(raw_json, '{}'),
       IFNULL(DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), '')
FROM plugin_sms_records
WHERE id = ?`, id)

	var (
		item    domain.SMSRecord
		rawJSON string
	)
	if err := row.Scan(
		&item.ID,
		&item.PluginCode,
		&item.CustomerID,
		&item.Mobile,
		&item.TemplateCode,
		&item.Content,
		&item.Status,
		&item.Message,
		&item.BizNo,
		&rawJSON,
		&item.CreatedAt,
	); err != nil {
		return domain.SMSRecord{}, err
	}
	item.Raw = decodeJSONMap(rawJSON)
	return item, nil
}

func (repository *MySQLRepository) getPaymentOrderByID(id int64) (domain.PaymentOrder, error) {
	row := repository.db.QueryRow(`
SELECT id, plugin_code, order_no, customer_id, amount, currency, subject, status, pay_url, notify_url, return_url, trade_no,
       IFNULL(raw_json, '{}'),
       IFNULL(DATE_FORMAT(paid_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), ''),
       IFNULL(DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s'), '')
FROM plugin_payment_orders
WHERE id = ?`, id)
	return scanPaymentOrder(row)
}

type scanner interface {
	Scan(dest ...any) error
}

func scanPaymentOrder(sc scanner) (domain.PaymentOrder, error) {
	var (
		item    domain.PaymentOrder
		rawJSON string
	)
	if err := sc.Scan(
		&item.ID,
		&item.PluginCode,
		&item.OrderNo,
		&item.CustomerID,
		&item.Amount,
		&item.Currency,
		&item.Subject,
		&item.Status,
		&item.PayURL,
		&item.NotifyURL,
		&item.ReturnURL,
		&item.TradeNo,
		&rawJSON,
		&item.PaidAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return domain.PaymentOrder{}, err
	}
	item.Raw = decodeJSONMap(rawJSON)
	return item, nil
}

func decodeJSONMap(payload string) map[string]any {
	text := strings.TrimSpace(payload)
	if text == "" {
		return map[string]any{}
	}
	result := map[string]any{}
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return map[string]any{}
	}
	return result
}

func lookupDefinition(code string) (domain.PluginDefinition, bool) {
	normalized := normalizeCode(code)
	for _, item := range defaultDefinitions() {
		if normalizeCode(item.Code) == normalized {
			return item, true
		}
	}
	return domain.PluginDefinition{}, false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

