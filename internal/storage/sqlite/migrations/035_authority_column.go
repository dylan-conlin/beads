package migrations

import (
	"database/sql"
	"fmt"
)

// MigrateAuthorityColumn adds the authority column to the dependencies table.
// This column stores who can traverse/resolve the dependency edge:
// - daemon: automated processes can resolve
// - orchestrator: requires orchestrator judgment
// - human: requires human decision
// Default is 'daemon' for backward compatibility with existing dependencies.
func MigrateAuthorityColumn(db *sql.DB) error {
	var columnExists bool
	err := db.QueryRow(`
		SELECT COUNT(*) > 0
		FROM pragma_table_info('dependencies')
		WHERE name = 'authority'
	`).Scan(&columnExists)
	if err != nil {
		return fmt.Errorf("failed to check authority column: %w", err)
	}

	if columnExists {
		return nil
	}

	_, err = db.Exec(`ALTER TABLE dependencies ADD COLUMN authority TEXT DEFAULT 'daemon'`)
	if err != nil {
		return fmt.Errorf("failed to add authority column: %w", err)
	}

	return nil
}
