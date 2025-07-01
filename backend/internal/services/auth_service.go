package services

import (
	"fmt"
	"go-backend-v2/global"
	"go-backend-v2/internal/common"
	"go-backend-v2/internal/dto"
	"go-backend-v2/internal/models"
	"go-backend-v2/internal/repo"
	"go-backend-v2/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type AuthServiceInterface interface {
	Signup(req *dto.SignupRequest) error
	Login(req *dto.LoginRequest) (string, error) // returns JWT token
	ValidateToken(token string) (string, error)  // returns userID
}

type AuthService struct {
	userRepo repo.UserRepositoryInterface
}

func NewAuthService(userRepo repo.UserRepositoryInterface) AuthServiceInterface {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Signup(req *dto.SignupRequest) error {
	if err := utils.ValidatePassword(req.Password); err != nil {
		return err
	}

	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return common.ErrEmailAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email:      req.Email,
		GlobalRole: common.GlobalRoleCustomer,
		Status:     common.UserStatusActive,
	}

	profile := &models.UserProfile{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Timezone:  "UTC",
		Locale:    "en",
	}

	authProvider := &models.UserAuthProvider{
		Provider:       common.AuthProviderLocal,
		ProviderUserID: req.Email,
		ProviderEmail:  &req.Email,
		ProviderData:   models.ProviderData{}, // Empty for local provider
		PasswordHash:   &hashedPassword,
		IsPrimary:      true,
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		return s.userRepo.CreateUserWithAuth(tx, user, profile, authProvider)
	})
}

func (s *AuthService) Login(req *dto.LoginRequest) (string, error) {
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return "", common.ErrInvalidCredentials
	}

	if user.Status != common.UserStatusActive {
		return "", common.ErrUserInactive
	}

	authProvider, err := s.userRepo.GetUserAuthProvider(user.ID, common.AuthProviderLocal)
	if err != nil {
		return "", fmt.Errorf("failed to get auth provider: %w", err)
	}
	if authProvider == nil || authProvider.PasswordHash == nil {
		return "", common.ErrInvalidCredentials
	}

	if !utils.CheckPassword(req.Password, *authProvider.PasswordHash) {
		return "", common.ErrInvalidCredentials
	}

	err = s.userRepo.UpdateUser(user.ID, map[string]interface{}{
		"last_login_at": time.Now(),
	})
	if err != nil {
		fmt.Printf("Warning: failed to update last login time: %v\n", err)
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *AuthService) ValidateToken(token string) (string, error) {
	userID, err := utils.ValidateToken(token)
	if err != nil {
		return "", common.ErrTokenInvalid
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return "", common.ErrUserNotFound
	}
	if user.Status != common.UserStatusActive {
		return "", common.ErrUserInactive
	}

	return userID, nil
}
