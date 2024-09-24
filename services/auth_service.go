package services

import (
	"api-server/models"
	"api-server/repositories"
	"api-server/utils"
	"errors"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// RegisterUser registers a new user by hashing the password and saving the user.
func (s *AuthService) RegisterUser(username, email, password string) (*models.User, error) {
	// Check if email is already in use
	existingUser, _ := s.userRepo.GetUserByEmail(email)
	if existingUser != nil {
		return nil, errors.New("email already in use")
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create and save the user
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}
	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

// AuthenticateUser authenticates a user by checking the password and returning the user if valid.
func (s *AuthService) AuthenticateUser(email, password string) (*models.User, error) {
	// Get user by email
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if err := utils.CheckPasswordHash(password, user.PasswordHash); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
