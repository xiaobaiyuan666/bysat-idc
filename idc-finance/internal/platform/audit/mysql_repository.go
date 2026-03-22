package audit

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (repository *MySQLRepository) Record(entry Entry) Entry {
	payloadJSON, err := json.Marshal(entry.Payload)
	if err != nil {
		log.Printf("audit mysql marshal payload failed: %v", err)
		payloadJSON = []byte("{}")
	}

	result, err := repository.db.Exec(`
INSERT INTO audit_logs (actor_type, actor_id, actor, action, target_type, target_id, target, request_id, description, payload)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entry.ActorType,
		entry.ActorID,
		entry.Actor,
		entry.Action,
		entry.TargetType,
		entry.TargetID,
		entry.Target,
		entry.RequestID,
		entry.Description,
		payloadJSON,
	)
	if err != nil {
		log.Printf("audit mysql insert failed: %v", err)
		return entry
	}

	entryID, err := result.LastInsertId()
	if err != nil {
		log.Printf("audit mysql last insert id failed: %v", err)
		return entry
	}

	recorded, ok := repository.getByID(context.Background(), entryID)
	if !ok {
		entry.ID = entryID
		return entry
	}
	return recorded
}

func (repository *MySQLRepository) List() []Entry {
	return repository.listByTarget("", 0)
}

func (repository *MySQLRepository) ListByTarget(targetType string, targetID int64) []Entry {
	return repository.listByTarget(targetType, targetID)
}

func (repository *MySQLRepository) listByTarget(targetType string, targetID int64) []Entry {
	query := `
SELECT id, actor_type, actor_id, actor, action, target_type, target_id, target, request_id, description,
       payload, DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s')
FROM audit_logs`
	args := make([]any, 0)
	if targetType != "" {
		query += ` WHERE target_type = ? AND target_id = ?`
		args = append(args, targetType, targetID)
	}
	query += ` ORDER BY id DESC`

	rows, err := repository.db.Query(query, args...)
	if err != nil {
		log.Printf("audit mysql list failed: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]Entry, 0)
	for rows.Next() {
		item, err := scanAuditEntry(rows)
		if err != nil {
			log.Printf("audit mysql scan failed: %v", err)
			return result
		}
		result = append(result, item)
	}
	return result
}

func (repository *MySQLRepository) getByID(ctx context.Context, id int64) (Entry, bool) {
	row := repository.db.QueryRowContext(ctx, `
SELECT id, actor_type, actor_id, actor, action, target_type, target_id, target, request_id, description,
       payload, DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s')
FROM audit_logs
WHERE id = ?`, id)

	item, err := scanAuditRow(row)
	if err != nil {
		return Entry{}, false
	}
	return item, true
}

type auditScanner interface {
	Scan(dest ...any) error
}

func scanAuditEntry(scanner auditScanner) (Entry, error) {
	var (
		item        Entry
		payloadJSON []byte
	)
	if err := scanner.Scan(
		&item.ID,
		&item.ActorType,
		&item.ActorID,
		&item.Actor,
		&item.Action,
		&item.TargetType,
		&item.TargetID,
		&item.Target,
		&item.RequestID,
		&item.Description,
		&payloadJSON,
		&item.CreatedAt,
	); err != nil {
		return Entry{}, err
	}
	if len(payloadJSON) > 0 {
		_ = json.Unmarshal(payloadJSON, &item.Payload)
	}
	return item, nil
}

func scanAuditRow(row *sql.Row) (Entry, error) {
	return scanAuditEntry(row)
}
