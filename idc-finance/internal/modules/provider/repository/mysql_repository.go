package repository

import (
	"database/sql"
	"log"
	"strings"

	"idc-finance/internal/modules/provider/domain"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (repository *MySQLRepository) ListAccounts(providerType string) []domain.Account {
	query := `
SELECT id, provider_type, account_name, base_url, username, password, source_name, contact_way, description,
       account_mode, lang, list_path, detail_path, insecure_skip_verify, auto_update, product_count, status,
       IFNULL(extra_config, ''), DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s')
FROM provider_accounts`
	args := make([]any, 0, 1)
	if strings.TrimSpace(providerType) != "" {
		query += ` WHERE provider_type = ?`
		args = append(args, strings.ToUpper(strings.TrimSpace(providerType)))
	}
	query += ` ORDER BY provider_type, account_name, id`

	rows, err := repository.db.Query(query, args...)
	if err != nil {
		log.Printf("provider mysql list accounts failed: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]domain.Account, 0)
	for rows.Next() {
		item, err := scanAccount(rows)
		if err != nil {
			log.Printf("provider mysql scan account failed: %v", err)
			return result
		}
		result = append(result, item)
	}
	return result
}

func (repository *MySQLRepository) GetAccountByID(id int64) (domain.Account, bool) {
	row := repository.db.QueryRow(`
SELECT id, provider_type, account_name, base_url, username, password, source_name, contact_way, description,
       account_mode, lang, list_path, detail_path, insecure_skip_verify, auto_update, product_count, status,
       IFNULL(extra_config, ''), DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s')
FROM provider_accounts
WHERE id = ?`, id)

	item, err := scanAccount(row)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("provider mysql get account failed: %v", err)
		}
		return domain.Account{}, false
	}
	return item, true
}

func (repository *MySQLRepository) CreateAccount(account domain.Account) domain.Account {
	result, err := repository.db.Exec(`
INSERT INTO provider_accounts (
  provider_type, account_name, base_url, username, password, source_name, contact_way, description,
  account_mode, lang, list_path, detail_path, insecure_skip_verify, auto_update, product_count, status, extra_config
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		strings.ToUpper(strings.TrimSpace(account.ProviderType)),
		strings.TrimSpace(account.Name),
		strings.TrimSpace(account.BaseURL),
		strings.TrimSpace(account.Username),
		strings.TrimSpace(account.Password),
		strings.TrimSpace(account.SourceName),
		strings.TrimSpace(account.ContactWay),
		strings.TrimSpace(account.Description),
		strings.TrimSpace(account.AccountMode),
		strings.TrimSpace(account.Lang),
		strings.TrimSpace(account.ListPath),
		strings.TrimSpace(account.DetailPath),
		account.InsecureSkipVerify,
		account.AutoUpdate,
		account.ProductCount,
		firstAccountStatus(account.Status),
		normalizeExtraConfig(account.ExtraConfig),
	)
	if err != nil {
		log.Printf("provider mysql create account failed: %v", err)
		return domain.Account{}
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("provider mysql create account id failed: %v", err)
		return domain.Account{}
	}

	created, ok := repository.GetAccountByID(id)
	if !ok {
		return domain.Account{}
	}
	return created
}

func (repository *MySQLRepository) UpdateAccount(id int64, account domain.Account) (domain.Account, bool) {
	result, err := repository.db.Exec(`
UPDATE provider_accounts
SET provider_type = ?, account_name = ?, base_url = ?, username = ?, password = ?, source_name = ?, contact_way = ?,
    description = ?, account_mode = ?, lang = ?, list_path = ?, detail_path = ?, insecure_skip_verify = ?,
    auto_update = ?, product_count = ?, status = ?, extra_config = ?
WHERE id = ?`,
		strings.ToUpper(strings.TrimSpace(account.ProviderType)),
		strings.TrimSpace(account.Name),
		strings.TrimSpace(account.BaseURL),
		strings.TrimSpace(account.Username),
		strings.TrimSpace(account.Password),
		strings.TrimSpace(account.SourceName),
		strings.TrimSpace(account.ContactWay),
		strings.TrimSpace(account.Description),
		strings.TrimSpace(account.AccountMode),
		strings.TrimSpace(account.Lang),
		strings.TrimSpace(account.ListPath),
		strings.TrimSpace(account.DetailPath),
		account.InsecureSkipVerify,
		account.AutoUpdate,
		account.ProductCount,
		firstAccountStatus(account.Status),
		normalizeExtraConfig(account.ExtraConfig),
		id,
	)
	if err != nil {
		log.Printf("provider mysql update account failed: %v", err)
		return domain.Account{}, false
	}

	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		return domain.Account{}, false
	}

	updated, ok := repository.GetAccountByID(id)
	if !ok {
		return domain.Account{}, false
	}
	return updated, true
}

func (repository *MySQLRepository) DeleteAccount(id int64) bool {
	result, err := repository.db.Exec(`DELETE FROM provider_accounts WHERE id = ?`, id)
	if err != nil {
		log.Printf("provider mysql delete account failed: %v", err)
		return false
	}
	affected, err := result.RowsAffected()
	return err == nil && affected > 0
}

func (repository *MySQLRepository) FindByProviderAndBaseURL(providerType, baseURL string) (domain.Account, bool) {
	row := repository.db.QueryRow(`
SELECT id, provider_type, account_name, base_url, username, password, source_name, contact_way, description,
       account_mode, lang, list_path, detail_path, insecure_skip_verify, auto_update, product_count, status,
       IFNULL(extra_config, ''), DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s')
FROM provider_accounts
WHERE provider_type = ? AND REPLACE(TRIM(TRAILING '/' FROM base_url), ' ', '') = REPLACE(TRIM(TRAILING '/' FROM ?), ' ', '')
LIMIT 1`,
		strings.ToUpper(strings.TrimSpace(providerType)),
		normalizeBaseURL(baseURL),
	)

	item, err := scanAccount(row)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("provider mysql find account failed: %v", err)
		}
		return domain.Account{}, false
	}
	return item, true
}

type accountScanner interface {
	Scan(dest ...any) error
}

func scanAccount(scanner accountScanner) (domain.Account, error) {
	var item domain.Account
	err := scanner.Scan(
		&item.ID,
		&item.ProviderType,
		&item.Name,
		&item.BaseURL,
		&item.Username,
		&item.Password,
		&item.SourceName,
		&item.ContactWay,
		&item.Description,
		&item.AccountMode,
		&item.Lang,
		&item.ListPath,
		&item.DetailPath,
		&item.InsecureSkipVerify,
		&item.AutoUpdate,
		&item.ProductCount,
		&item.Status,
		&item.ExtraConfig,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return domain.Account{}, err
	}
	return item, nil
}

func firstAccountStatus(status domain.AccountStatus) domain.AccountStatus {
	if status == "" {
		return domain.AccountStatusActive
	}
	return status
}

func normalizeExtraConfig(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "{}"
	}
	return trimmed
}
