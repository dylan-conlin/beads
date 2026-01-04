package migrations

import (
	"database/sql"
	"fmt"
)

// MigrateReproColumns adds repro and no_repro_reason columns to the issues table.
// These fields support bug reproduction requirements for type=bug issues.
// - repro: Steps/evidence to reproduce the bug
// - no_repro_reason: Explanation why reproduction is not possible (if --no-repro was used)
func MigrateReproColumns(db *sql.DB) error {
	// Check if repro column already exists
	var reproExists bool
	err := db.QueryRow(`
		SELECT COUNT(*) > 0
		FROM pragma_table_info('issues')
		WHERE name = 'repro'
	`).Scan(&reproExists)
	if err != nil {
		return fmt.Errorf("failed to check repro column: %w", err)
	}

	if !reproExists {
		_, err = db.Exec(`ALTER TABLE issues ADD COLUMN repro TEXT DEFAULT ''`)
		if err != nil {
			return fmt.Errorf("failed to add repro column: %w", err)
		}
	}

	// Check if no_repro_reason column already exists
	var noReproExists bool
	err = db.QueryRow(`
		SELECT COUNT(*) > 0
		FROM pragma_table_info('issues')
		WHERE name = 'no_repro_reason'
	`).Scan(&noReproExists)
	if err != nil {
		return fmt.Errorf("failed to check no_repro_reason column: %w", err)
	}

	if !noReproExists {
		_, err = db.Exec(`ALTER TABLE issues ADD COLUMN no_repro_reason TEXT DEFAULT ''`)
		if err != nil {
			return fmt.Errorf("failed to add no_repro_reason column: %w", err)
		}
	}

	return nil
}
