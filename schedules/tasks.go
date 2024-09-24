package schedules

import (
	"database/sql"
	"log"
	"time"
)

// Task: Example task that runs every minute
func ExampleTask(db *sql.DB) {
	log.Println("Running cron task: Current time:", time.Now())
	performDatabaseTask(db)
}

// Task: Example daily cleanup task
func DailyCleanupTask(db *sql.DB) {
	log.Println("Running daily cleanup task")

	// Perform cleanup
	err := cleanupOldRecords(db)
	if err != nil {
		log.Printf("Error during cleanup: %v", err)
	} else {
		log.Println("Cleanup completed successfully")
	}
}

// Helper function: Database interaction for task
func performDatabaseTask(db *sql.DB) {
	var result string
	err := db.QueryRow("SELECT 'Hello from schedules!'").Scan(&result)
	if err != nil {
		log.Println("Database error:", err)
		return
	}
	log.Println("Database result:", result)
}

// Helper function: Cleanup old records
func cleanupOldRecords(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM some_table WHERE created_at < NOW() - INTERVAL '30 days'")
	if err != nil {
		return err
	}
	return nil
}
