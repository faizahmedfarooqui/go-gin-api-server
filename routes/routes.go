package routes

import (
	"database/sql"
	"net/http"

	"api-server/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	authController := controllers.NewAuthController(db)
	itemController := controllers.NewItemController(db)

	// Routes: Auth
	router.POST("/register", authController.Register)
	router.POST("/login", authController.Login)

	// Routes: Items
	router.GET("/items", itemController.GetItems)
	router.GET("/items/:id", itemController.GetItem)
	router.POST("/items", itemController.CreateItem)

	// Routes: Custom 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": http.StatusNotFound,
			"status":  "Sorry, the requested URL was not found on this server.",
		})
	})
}
