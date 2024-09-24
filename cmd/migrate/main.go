package main

import (
	"api-server/config"
	"api-server/database"
	"flag"
	"log"
	"path/filepath"
)

func main() {
	// Define CLI flags
	migrateUp := flag.Bool("up", false, "Apply all pending migrations")
	migrateDown := flag.Bool("down", false, "Rollback the last migration")
	flag.Parse()

	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	// Initialize the database connection
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// Set the migrations directory
	migrationsDir := filepath.Join(".", "database", "migrations")

	// Apply migrations if the 'up' flag is set
	if *migrateUp {
		log.Println("Applying migrations...")
		if err := database.ApplyMigrations(db, migrationsDir); err != nil {
			log.Fatalf("Error applying migrations: %v", err)
		}
		log.Println("Migrations applied successfully")
		return
	}

	// Rollback the last migration if the 'down' flag is set
	if *migrateDown {
		log.Println("Rolling back the last migration...")
		if err := database.RollbackLastMigration(db, migrationsDir); err != nil {
			log.Fatalf("Error rolling back the migration: %v", err)
		}
		log.Println("Migration rolled back successfully")
		return
	}

	// If no flags are provided, print usage
	flag.Usage()
}
