package migrations

import (
	"database/sql"
	"fmt"
)

// MigrateResolutionTypeDomainColumns adds the resolution_type and domain columns
// to the issues table for decidability substrate support.
//
// resolution_type: How questions should be resolved (factual, judgment, framing)
// domain: Categorizes decisions for frontier queries (e.g., model-selection, spawn-architecture)
func MigrateResolutionTypeDomainColumns(db *sql.DB) error {
	// Check and add resolution_type column
	var resolutionTypeExists bool
	err := db.QueryRow(`
		SELECT COUNT(*) > 0
		FROM pragma_table_info('issues')
		WHERE name = 'resolution_type'
	`).Scan(&resolutionTypeExists)
	if err != nil {
		return fmt.Errorf("failed to check resolution_type column: %w", err)
	}

	if !resolutionTypeExists {
		_, err = db.Exec(`ALTER TABLE issues ADD COLUMN resolution_type TEXT DEFAULT ''`)
		if err != nil {
			return fmt.Errorf("failed to add resolution_type column: %w", err)
		}
	}

	// Check and add domain column
	var domainExists bool
	err = db.QueryRow(`
		SELECT COUNT(*) > 0
		FROM pragma_table_info('issues')
		WHERE name = 'domain'
	`).Scan(&domainExists)
	if err != nil {
		return fmt.Errorf("failed to check domain column: %w", err)
	}

	if !domainExists {
		_, err = db.Exec(`ALTER TABLE issues ADD COLUMN domain TEXT DEFAULT ''`)
		if err != nil {
			return fmt.Errorf("failed to add domain column: %w", err)
		}
	}

	return nil
}
