package dto

import "go-backend-v2/internal/models"

type SignupRequest struct {
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=6,max=128"`
	FirstName string `json:"first_name" validate:"required,min=1,max=100,alpha_space"`
	LastName  string `json:"last_name" validate:"required,min=1,max=100,alpha_space"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User           *models.User `json:"user"`
	AccessToken    string       `json:"access_token"`
	EncryptedToken string       `json:"encrypted_token"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
