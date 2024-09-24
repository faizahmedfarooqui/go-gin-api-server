package validators

// RegisterUserValidator holds the fields for validating user registration input.
type RegisterUserValidator struct {
	Username string `json:"username" binding:"required,min=3,max=50" message:"Username is required and must be between 3 and 50 characters"`
	Email    string `json:"email" binding:"required,email" message:"Email is required and must be a valid email address"`
	Password string `json:"password" binding:"required,min=6" message:"Password is required and must be at least 6 characters long"`
}

// LoginUserValidator holds the fields for validating user login input.
type LoginUserValidator struct {
	Email    string `json:"email" binding:"required,email" message:"Email is required and must be a valid email address"`
	Password string `json:"password" binding:"required,min=6" message:"Password is required and must be at least 6 characters long"`
}
