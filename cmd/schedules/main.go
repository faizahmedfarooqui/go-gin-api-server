package main

import (
	"api-server/config"
	"api-server/schedules"
	"database/sql"
	"log"

	"github.com/robfig/cron/v3"
)

func main() {
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	// Initialize the database
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// Initialize the cron scheduler
	c := cron.New(cron.WithSeconds())

	// Register cron jobs
	registerSchedules(c, db)

	// Start the cron scheduler
	c.Start()
	defer c.Stop()

	// Keep the cron job running
	select {}
}

// registerSchedules registers all cron jobs
func registerSchedules(c *cron.Cron, db *sql.DB) {
	// Run ExampleTask every minute
	_, err := c.AddFunc("*/1 * * * * *", func() {
		schedules.ExampleTask(db) // Call the decoupled task from schedules package
	})
	if err != nil {
		log.Fatalf("Error scheduling ExampleTask: %v", err)
	}

	// Run DailyCleanupTask every day at midnight
	_, err = c.AddFunc("0 0 * * * *", func() {
		schedules.DailyCleanupTask(db)
	})
	if err != nil {
		log.Fatalf("Error scheduling DailyCleanupTask: %v", err)
	}
}
