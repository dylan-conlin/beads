package migrations

import (
	"database/sql"
	"fmt"
)

// MigrateCloseOutcomeColumn adds the close_outcome column to the issues table.
// This column stores the outcome category when closing an issue (e.g., could-not-reproduce).
func MigrateCloseOutcomeColumn(db *sql.DB) error {
	var columnExists bool
	err := db.QueryRow(`
		SELECT COUNT(*) > 0
		FROM pragma_table_info('issues')
		WHERE name = 'close_outcome'
	`).Scan(&columnExists)
	if err != nil {
		return fmt.Errorf("failed to check close_outcome column: %w", err)
	}

	if columnExists {
		return nil
	}

	_, err = db.Exec(`ALTER TABLE issues ADD COLUMN close_outcome TEXT DEFAULT ''`)
	if err != nil {
		return fmt.Errorf("failed to add close_outcome column: %w", err)
	}

	return nil
}
