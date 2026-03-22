package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"idc-finance/internal/modules/catalog/domain"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (repository *MySQLRepository) List() []domain.Product {
	return repository.listByStatus("")
}

func (repository *MySQLRepository) ListActive() []domain.Product {
	return repository.listByStatus(string(domain.ProductStatusActive))
}

func (repository *MySQLRepository) ListGroups() []string {
	rows, err := repository.db.Query(`
SELECT name
FROM product_groups
ORDER BY name`)
	if err != nil {
		log.Printf("catalog mysql list groups failed: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Printf("catalog mysql scan group failed: %v", err)
			return result
		}
		result = append(result, name)
	}
	return result
}

func (repository *MySQLRepository) GetByID(id int64) (domain.Product, bool) {
	product, err := repository.loadProduct(context.Background(), repository.db, id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("catalog mysql get product failed: %v", err)
		}
		return domain.Product{}, false
	}
	return product, true
}

func (repository *MySQLRepository) Create(product domain.Product) domain.Product {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("catalog mysql begin create failed: %v", err)
		return domain.Product{}
	}
	defer rollbackQuietly(tx)

	groupID, err := ensureProductGroup(ctx, tx, product.GroupName)
	if err != nil {
		log.Printf("catalog mysql ensure group failed: %v", err)
		return domain.Product{}
	}

	resourceJSON, err := json.Marshal(product.ResourceTemplate)
	if err != nil {
		log.Printf("catalog mysql marshal template failed: %v", err)
		return domain.Product{}
	}
	automationJSON, err := json.Marshal(product.AutomationConfig)
	if err != nil {
		log.Printf("catalog mysql marshal automation config failed: %v", err)
		return domain.Product{}
	}
	upstreamJSON, err := json.Marshal(product.UpstreamMapping)
	if err != nil {
		log.Printf("catalog mysql marshal upstream mapping failed: %v", err)
		return domain.Product{}
	}

	tempNo := fmt.Sprintf("TMP-%d", time.Now().UnixNano())
	result, err := tx.ExecContext(ctx, `
INSERT INTO products (product_no, group_id, name, description, product_type, status, resource_template, automation_config, upstream_mapping)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		tempNo,
		groupID,
		product.Name,
		product.Description,
		product.ProductType,
		product.Status,
		resourceJSON,
		automationJSON,
		upstreamJSON,
	)
	if err != nil {
		log.Printf("catalog mysql insert product failed: %v", err)
		return domain.Product{}
	}

	productID, err := result.LastInsertId()
	if err != nil {
		log.Printf("catalog mysql get product id failed: %v", err)
		return domain.Product{}
	}

	productNo := fmt.Sprintf("PROD-%06d", productID)
	if _, err := tx.ExecContext(ctx, `UPDATE products SET product_no = ? WHERE id = ?`, productNo, productID); err != nil {
		log.Printf("catalog mysql update product no failed: %v", err)
		return domain.Product{}
	}

	if err := repository.replaceProductRelations(ctx, tx, productID, product); err != nil {
		log.Printf("catalog mysql insert relations failed: %v", err)
		return domain.Product{}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("catalog mysql commit create failed: %v", err)
		return domain.Product{}
	}

	created, ok := repository.GetByID(productID)
	if !ok {
		return domain.Product{}
	}
	return created
}

func (repository *MySQLRepository) Update(id int64, product domain.Product) (domain.Product, bool) {
	ctx := context.Background()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("catalog mysql begin update failed: %v", err)
		return domain.Product{}, false
	}
	defer rollbackQuietly(tx)

	groupID, err := ensureProductGroup(ctx, tx, product.GroupName)
	if err != nil {
		log.Printf("catalog mysql ensure group failed: %v", err)
		return domain.Product{}, false
	}

	resourceJSON, err := json.Marshal(product.ResourceTemplate)
	if err != nil {
		log.Printf("catalog mysql marshal template failed: %v", err)
		return domain.Product{}, false
	}
	automationJSON, err := json.Marshal(product.AutomationConfig)
	if err != nil {
		log.Printf("catalog mysql marshal automation config failed: %v", err)
		return domain.Product{}, false
	}
	upstreamJSON, err := json.Marshal(product.UpstreamMapping)
	if err != nil {
		log.Printf("catalog mysql marshal upstream mapping failed: %v", err)
		return domain.Product{}, false
	}

	result, err := tx.ExecContext(ctx, `
UPDATE products
SET group_id = ?, name = ?, description = ?, product_type = ?, status = ?, resource_template = ?, automation_config = ?, upstream_mapping = ?
WHERE id = ?`,
		groupID,
		product.Name,
		product.Description,
		product.ProductType,
		product.Status,
		resourceJSON,
		automationJSON,
		upstreamJSON,
		id,
	)
	if err != nil {
		log.Printf("catalog mysql update product failed: %v", err)
		return domain.Product{}, false
	}

	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		return domain.Product{}, false
	}

	if err := repository.replaceProductRelations(ctx, tx, id, product); err != nil {
		log.Printf("catalog mysql replace relations failed: %v", err)
		return domain.Product{}, false
	}

	if err := tx.Commit(); err != nil {
		log.Printf("catalog mysql commit update failed: %v", err)
		return domain.Product{}, false
	}

	updated, ok := repository.GetByID(id)
	if !ok {
		return domain.Product{}, false
	}
	return updated, true
}

func (repository *MySQLRepository) listByStatus(status string) []domain.Product {
	query := `
SELECT p.id
FROM products p
WHERE (? = '' OR p.status = ?)
ORDER BY p.id`

	rows, err := repository.db.Query(query, status, status)
	if err != nil {
		log.Printf("catalog mysql list products failed: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]domain.Product, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Printf("catalog mysql scan product id failed: %v", err)
			return result
		}
		item, err := repository.loadProduct(context.Background(), repository.db, id)
		if err != nil {
			log.Printf("catalog mysql load product failed: %v", err)
			continue
		}
		result = append(result, item)
	}
	return result
}

func (repository *MySQLRepository) loadProduct(ctx context.Context, queryer queryer, id int64) (domain.Product, error) {
	row := queryer.QueryRowContext(ctx, `
SELECT
  p.id,
  p.product_no,
  COALESCE(g.name, '') AS group_name,
  p.name,
  p.description,
  p.product_type,
  p.status,
  p.resource_template,
  p.automation_config,
  p.upstream_mapping
FROM products p
LEFT JOIN product_groups g ON g.id = p.group_id
WHERE p.id = ?`, id)

	var (
		product            domain.Product
		resourceTemplateDB []byte
		automationConfigDB []byte
		upstreamMappingDB  []byte
	)
	if err := row.Scan(
		&product.ID,
		&product.ProductNo,
		&product.GroupName,
		&product.Name,
		&product.Description,
		&product.ProductType,
		&product.Status,
		&resourceTemplateDB,
		&automationConfigDB,
		&upstreamMappingDB,
	); err != nil {
		return domain.Product{}, err
	}

	if len(resourceTemplateDB) > 0 {
		_ = json.Unmarshal(resourceTemplateDB, &product.ResourceTemplate)
	}
	if len(automationConfigDB) > 0 {
		_ = json.Unmarshal(automationConfigDB, &product.AutomationConfig)
	}
	if len(upstreamMappingDB) > 0 {
		_ = json.Unmarshal(upstreamMappingDB, &product.UpstreamMapping)
	}

	product.Pricing = repository.loadPricing(ctx, queryer, id)
	product.ConfigOptions = repository.loadConfigOptions(ctx, queryer, id)
	return product, nil
}

func (repository *MySQLRepository) loadPricing(ctx context.Context, queryer queryer, productID int64) []domain.PriceOption {
	rows, err := queryer.QueryContext(ctx, `
SELECT cycle_code, cycle_name, price, setup_fee
FROM product_prices
WHERE product_id = ?
ORDER BY id`, productID)
	if err != nil {
		log.Printf("catalog mysql load pricing failed: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]domain.PriceOption, 0)
	for rows.Next() {
		var item domain.PriceOption
		if err := rows.Scan(&item.CycleCode, &item.CycleName, &item.Price, &item.SetupFee); err != nil {
			log.Printf("catalog mysql scan pricing failed: %v", err)
			return result
		}
		result = append(result, item)
	}
	return result
}

func (repository *MySQLRepository) loadConfigOptions(ctx context.Context, queryer queryer, productID int64) []domain.ConfigOption {
	rows, err := queryer.QueryContext(ctx, `
SELECT id, code, name, input_type, is_required, default_value, description
FROM product_config_options
WHERE product_id = ?
ORDER BY sort_index, id`, productID)
	if err != nil {
		log.Printf("catalog mysql load config options failed: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]domain.ConfigOption, 0)
	for rows.Next() {
		var (
			optionID   int64
			requiredDB bool
			item       domain.ConfigOption
		)
		if err := rows.Scan(&optionID, &item.Code, &item.Name, &item.InputType, &requiredDB, &item.DefaultValue, &item.Description); err != nil {
			log.Printf("catalog mysql scan config option failed: %v", err)
			return result
		}
		item.Required = requiredDB
		item.Choices = repository.loadConfigChoices(ctx, queryer, optionID)
		result = append(result, item)
	}
	return result
}

func (repository *MySQLRepository) loadConfigChoices(ctx context.Context, queryer queryer, optionID int64) []domain.ConfigOptionChoice {
	rows, err := queryer.QueryContext(ctx, `
SELECT option_value, option_label, price_delta
FROM product_config_choices
WHERE option_id = ?
ORDER BY sort_index, id`, optionID)
	if err != nil {
		log.Printf("catalog mysql load config choices failed: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]domain.ConfigOptionChoice, 0)
	for rows.Next() {
		var item domain.ConfigOptionChoice
		if err := rows.Scan(&item.Value, &item.Label, &item.PriceDelta); err != nil {
			log.Printf("catalog mysql scan config choice failed: %v", err)
			return result
		}
		result = append(result, item)
	}
	return result
}

func (repository *MySQLRepository) replaceProductRelations(ctx context.Context, tx *sql.Tx, productID int64, product domain.Product) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM product_prices WHERE product_id = ?`, productID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
DELETE c
FROM product_config_choices c
INNER JOIN product_config_options o ON o.id = c.option_id
WHERE o.product_id = ?`, productID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM product_config_options WHERE product_id = ?`, productID); err != nil {
		return err
	}

	for _, item := range product.Pricing {
		if _, err := tx.ExecContext(ctx, `
INSERT INTO product_prices (product_id, cycle_code, cycle_name, price, setup_fee)
VALUES (?, ?, ?, ?, ?)`,
			productID,
			item.CycleCode,
			item.CycleName,
			item.Price,
			item.SetupFee,
		); err != nil {
			return err
		}
	}

	for optionIndex, option := range product.ConfigOptions {
		result, err := tx.ExecContext(ctx, `
INSERT INTO product_config_options (product_id, code, name, input_type, is_required, default_value, description, sort_index)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			productID,
			option.Code,
			option.Name,
			option.InputType,
			option.Required,
			option.DefaultValue,
			option.Description,
			optionIndex,
		)
		if err != nil {
			return err
		}

		optionID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		for choiceIndex, choice := range option.Choices {
			if _, err := tx.ExecContext(ctx, `
INSERT INTO product_config_choices (option_id, option_value, option_label, price_delta, sort_index)
VALUES (?, ?, ?, ?, ?)`,
				optionID,
				choice.Value,
				choice.Label,
				choice.PriceDelta,
				choiceIndex,
			); err != nil {
				return err
			}
		}
	}

	return nil
}

func ensureProductGroup(ctx context.Context, tx *sql.Tx, groupName string) (sql.NullInt64, error) {
	if groupName == "" {
		return sql.NullInt64{}, nil
	}

	result, err := tx.ExecContext(ctx, `
INSERT INTO product_groups (name, description)
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

type queryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func rollbackQuietly(tx *sql.Tx) {
	_ = tx.Rollback()
}
