package main

import (
	"log"
	"os"
	"strings"

	"api-server/config"
	"api-server/middlewares"
	"api-server/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	// Set GIN mode
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.ReleaseMode // Defaults to release mode if not set
	} else if ginMode == gin.DebugMode {
		ginMode = gin.DebugMode // Set to debug mode
	} else if ginMode == gin.ReleaseMode {
		ginMode = gin.ReleaseMode // Set to release mode
	} else if ginMode == gin.TestMode {
		ginMode = gin.TestMode // Set to test mode
	} else {
		log.Fatalf("Invalid GIN_MODE: %s", ginMode)
	}
	gin.SetMode(ginMode)

	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// Router: Initialize router
	router := gin.Default()

	// Router: Configure trusted proxies
	trustedProxies := os.Getenv("TRUSTED_PROXIES")
	if trustedProxies == "" {
		// Default to loopback addresses if not set
		router.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	} else {
		proxyList := strings.Split(trustedProxies, ",")
		router.SetTrustedProxies(proxyList)
	}

	// Register middleware
	router.Use(middlewares.ErrorHandler)

	// Router: Setup routes
	routes.SetupRoutes(router, db)

	// Get port from environment variables or use default
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // default port
	}

	// Start the server
	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
