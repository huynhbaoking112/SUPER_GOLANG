package services

import (
	"go-backend-v2/internal/dto"
	"go-backend-v2/internal/models"
)

type WorkspaceServiceInterface interface {
	CreateWorkspace(userID string, req *dto.CreateWorkspaceRequest) (*dto.WorkspaceResponse, error)
}

type AuthServiceInterface interface {
	Signup(req *dto.SignupRequest) error
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error) // returns login response with tokens
	Logout(userID, encryptedToken string) error              // logout specific token
	ValidateToken(token string) (string, error)              // returns userID

	// Redis token operations
	StoreTokenData(userID, encryptedToken string, tokenData *dto.UserTokenData) error
	GetTokenData(userID, encryptedToken string) (*dto.UserTokenData, error)
	DeleteTokenData(userID, encryptedToken string) error
	InvalidateUserTokens(userID string) error
}

type UserServiceInterface interface {
	GetUserWithWorkspaces(userID string) (*models.User, error)
	GetUserProfile(userID string) (*models.User, error)
	DeleteUser(userID string) error
}
