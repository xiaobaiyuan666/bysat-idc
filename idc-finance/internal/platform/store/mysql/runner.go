package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func RunDirectory(db *sql.DB, dir string) error {
	if err := ensureMigrationsTable(db); err != nil {
		return err
	}

	parentDir := filepath.Base(filepath.Dir(dir))
	dirName := filepath.Base(dir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read directory %s: %w", dir, err)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}
		files = append(files, filepath.Join(dir, entry.Name()))
	}
	slices.Sort(files)

	for _, file := range files {
		key := filepath.ToSlash(filepath.Join(parentDir, dirName, filepath.Base(file)))
		applied, err := isApplied(db, key)
		if err != nil {
			return err
		}
		if applied {
			continue
		}
		if err := RunFile(db, file); err != nil {
			return err
		}
		if err := markApplied(db, key); err != nil {
			return err
		}
	}
	return nil
}

func RunFile(db *sql.DB, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read sql file %s: %w", path, err)
	}

	statements := splitSQLStatements(string(content))
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return fmt.Errorf("exec %s: %w", filepath.Base(path), err)
		}
	}
	return nil
}

func splitSQLStatements(content string) []string {
	content = strings.TrimPrefix(content, "\uFEFF")

	var (
		statements     []string
		builder        strings.Builder
		inSingleQuote  bool
		inDoubleQuote  bool
		inBacktick     bool
		inLineComment  bool
		inBlockComment bool
	)

	for index := 0; index < len(content); index++ {
		char := content[index]
		next := byte(0)
		if index+1 < len(content) {
			next = content[index+1]
		}

		if inLineComment {
			if char == '\n' {
				inLineComment = false
				builder.WriteByte(char)
			}
			continue
		}

		if inBlockComment {
			if char == '*' && next == '/' {
				inBlockComment = false
				index++
			}
			continue
		}

		if !inSingleQuote && !inDoubleQuote && !inBacktick {
			if char == '-' && next == '-' {
				inLineComment = true
				index++
				continue
			}
			if char == '#' {
				inLineComment = true
				continue
			}
			if char == '/' && next == '*' {
				inBlockComment = true
				index++
				continue
			}
		}

		switch char {
		case '\'':
			if !inDoubleQuote && !inBacktick {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if !inSingleQuote && !inBacktick {
				inDoubleQuote = !inDoubleQuote
			}
		case '`':
			if !inSingleQuote && !inDoubleQuote {
				inBacktick = !inBacktick
			}
		case ';':
			if !inSingleQuote && !inDoubleQuote && !inBacktick {
				statement := strings.TrimSpace(builder.String())
				if statement != "" {
					statements = append(statements, statement)
				}
				builder.Reset()
				continue
			}
		}

		builder.WriteByte(char)
	}

	final := strings.TrimSpace(builder.String())
	if final != "" {
		statements = append(statements, final)
	}
	return statements
}

func ensureMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS schema_migrations (
  filename VARCHAR(255) NOT NULL PRIMARY KEY,
  applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	if err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}
	return nil
}

func isApplied(db *sql.DB, filename string) (bool, error) {
	row := db.QueryRow(`SELECT 1 FROM schema_migrations WHERE filename = ? LIMIT 1`, filename)
	var exists int
	if err := row.Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("check schema_migrations for %s: %w", filename, err)
	}
	return true, nil
}

func markApplied(db *sql.DB, filename string) error {
	if _, err := db.Exec(`INSERT INTO schema_migrations (filename) VALUES (?)`, filename); err != nil {
		return fmt.Errorf("mark schema_migrations for %s: %w", filename, err)
	}
	return nil
}
