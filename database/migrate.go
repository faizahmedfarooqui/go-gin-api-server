package database

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ensureSchemaMigrationsTable ensures the schema_migrations table exists.
func ensureSchemaMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id TEXT PRIMARY KEY
		);
	`)
	return err
}

// getAppliedMigrations returns the list of migrations already applied in the database.
func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	appliedMigrations := make(map[string]bool)
	rows, err := db.Query(`SELECT id FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		appliedMigrations[id] = true
	}

	return appliedMigrations, nil
}

// getMigrationFiles reads all the migration files from the directory.
func getMigrationFiles(migrationsDir string) ([]string, error) {
	var migrations []string

	err := filepath.WalkDir(migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(d.Name(), ".up.sql") {
			migrationID := strings.TrimSuffix(d.Name(), ".up.sql")
			migrations = append(migrations, migrationID)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return migrations, nil
}

// ApplyMigrations applies any pending migrations and checks for missing migration files.
func ApplyMigrations(db *sql.DB, migrationsDir string) error {
	// Ensure the schema_migrations table exists
	if err := ensureSchemaMigrationsTable(db); err != nil {
		return fmt.Errorf("error ensuring schema_migrations table: %w", err)
	}

	// Get applied migrations from the database
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("error fetching applied migrations: %w", err)
	}

	// Get migration files from the directory
	migrationFiles, err := getMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("error reading migration files: %w", err)
	}

	// Check if any applied migration is missing from the directory
	for appliedMigration := range appliedMigrations {
		found := false
		for _, migration := range migrationFiles {
			if appliedMigration == migration {
				found = true
				break
			}
		}
		if !found {
			// If an applied migration is missing from the filesystem, throw an error and exit
			return fmt.Errorf("critical error: migration %s is missing from the migrations directory", appliedMigration)
		}
	}

	// Apply pending migrations
	for _, migrationID := range migrationFiles {
		if !appliedMigrations[migrationID] {
			upSQLPath := filepath.Join(migrationsDir, fmt.Sprintf("%s.up.sql", migrationID))

			// Read the up migration SQL file
			upSQL, err := os.ReadFile(upSQLPath)
			if err != nil {
				return fmt.Errorf("error reading migration file %s: %w", upSQLPath, err)
			}

			// Apply the migration
			log.Printf("Applying migration: %s", migrationID)
			if _, err := db.Exec(string(upSQL)); err != nil {
				return fmt.Errorf("error applying migration %s: %w", migrationID, err)
			}

			// Record the applied migration
			if _, err := db.Exec(`INSERT INTO schema_migrations (id) VALUES ($1)`, migrationID); err != nil {
				return fmt.Errorf("error recording migration %s: %w", migrationID, err)
			}
		}
	}

	log.Println("All migrations applied successfully")
	return nil
}

// RollbackLastMigration rolls back the last applied migration.
func RollbackLastMigration(db *sql.DB, migrationsDir string) error {
	// Get the last applied migration
	row := db.QueryRow(`SELECT id FROM schema_migrations ORDER BY id DESC LIMIT 1`)
	var lastMigrationID string
	if err := row.Scan(&lastMigrationID); err != nil {
		if err == sql.ErrNoRows {
			log.Println("No migrations found to rollback")
			return nil
		}
		return err
	}

	// Check if the migration file still exists
	downSQLPath := filepath.Join(migrationsDir, fmt.Sprintf("%s.down.sql", lastMigrationID))
	if _, err := os.Stat(downSQLPath); os.IsNotExist(err) {
		return fmt.Errorf("critical error: migration file %s is missing", downSQLPath)
	}

	// Read the down migration SQL file
	downSQL, err := os.ReadFile(downSQLPath)
	if err != nil {
		return fmt.Errorf("error reading migration file %s: %w", downSQLPath, err)
	}

	// Rollback the migration
	log.Printf("Rolling back migration: %s", lastMigrationID)
	if _, err := db.Exec(string(downSQL)); err != nil {
		return fmt.Errorf("error rolling back migration %s: %w", lastMigrationID, err)
	}

	// Remove the migration from the schema_migrations table
	if _, err := db.Exec(`DELETE FROM schema_migrations WHERE id = $1`, lastMigrationID); err != nil {
		return fmt.Errorf("error deleting migration record %s: %w", lastMigrationID, err)
	}

	log.Printf("Migration %s rolled back successfully", lastMigrationID)
	return nil
}
