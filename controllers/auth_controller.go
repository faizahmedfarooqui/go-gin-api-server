package controllers

import (
	"database/sql"
	"net/http"
	"time"

	"api-server/models"
	"api-server/utils"
	"api-server/validators"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	db *sql.DB
}

// NewAuthController initializes AuthController with DB connection.
func NewAuthController(db *sql.DB) *AuthController {
	return &AuthController{db: db}
}

// Register handles user registration
func (ac *AuthController) Register(c *gin.Context) {
	var input validators.RegisterUserValidator

	// Bind and validate the input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		// When there's a validation error, it will automatically caught by the middleware
		c.Error(err).SetType(gin.ErrorTypeBind) // Set the error type to bind so it triggers the middleware
		return
	}

	// Check if the email already exists
	var existingUser models.User
	err := ac.db.QueryRow(`SELECT id FROM users WHERE email = $1`, input.Email).Scan(&existingUser.ID)
	if err != sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already in use"})
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Create the user
	var newUser models.User
	err = ac.db.QueryRow(
		`INSERT INTO users (username, email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, username, email, created_at`,
		input.Username, input.Email, hashedPassword, time.Now(), time.Now(),
	).Scan(&newUser.ID, &newUser.Username, &newUser.Email, &newUser.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// Login handles user login
func (ac *AuthController) Login(c *gin.Context) {
	var input validators.LoginUserValidator

	// Bind and validate the input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		// When there's a validation error, it will automatically caught by the middleware
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Set the error type to bind so it triggers the middleware
		return
	}

	// Get the user by email
	var user models.User
	err := ac.db.QueryRow(`SELECT id, username, email, password_hash FROM users WHERE email = $1`, input.Email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// Check the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successful", "user": user})
}
