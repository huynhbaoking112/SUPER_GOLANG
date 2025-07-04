package services

import (
	"context"
	"encoding/json"
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
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error) // returns login response with tokens
	Logout(userID, encryptedToken string) error              // logout specific token
	ValidateToken(token string) (string, error)              // returns userID

	// Redis token operations
	StoreTokenData(userID, encryptedToken string, tokenData *dto.UserTokenData) error
	GetTokenData(userID, encryptedToken string) (*dto.UserTokenData, error)
	DeleteTokenData(userID, encryptedToken string) error
	InvalidateUserTokens(userID string) error
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

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, common.ErrInvalidCredentials
	}

	if user.Status != common.UserStatusActive {
		return nil, common.ErrUserInactive
	}

	authProvider, err := s.userRepo.GetUserAuthProvider(user.ID, common.AuthProviderLocal)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth provider: %w", err)
	}
	if authProvider == nil || authProvider.PasswordHash == nil {
		return nil, common.ErrInvalidCredentials
	}

	if !utils.CheckPassword(req.Password, *authProvider.PasswordHash) {
		return nil, common.ErrInvalidCredentials
	}

	err = s.userRepo.UpdateUser(user.ID, map[string]interface{}{
		"last_login_at": time.Now(),
	})
	if err != nil {
		fmt.Printf("Warning: failed to update last login time: %v\n", err)
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	encryptedToken, err := utils.EncryptToken(token, global.Config.JWT.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt token: %w", err)
	}

	userWithWorkspaces, err := s.userRepo.GetUserWithWorkspaces(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user with workspaces: %w", err)
	}

	tokenData := s.BuildTokenData(userWithWorkspaces)

	err = s.StoreTokenData(user.ID, encryptedToken, tokenData)
	if err != nil {
		fmt.Printf("Warning: failed to cache token data in Redis: %v\n", err)
	}

	return &dto.LoginResponse{
		AccessToken:    token,
		User:           userWithWorkspaces,
		EncryptedToken: encryptedToken,
	}, nil
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

func (s *AuthService) StoreTokenData(userID, encryptedToken string, tokenData *dto.UserTokenData) error {
	ctx := context.Background()

	jsonData, err := json.Marshal(tokenData)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	tokenKey := fmt.Sprintf("auth:token:%s:%s", userID, encryptedToken)

	expire := global.Config.JWT.ExpirationTime
	if expire == 0 {
		expire = 72 * time.Hour // fallback default
	}

	err = global.RedisClient.Set(ctx, tokenKey, jsonData, expire).Err()
	if err != nil {
		return fmt.Errorf("failed to store token data in Redis: %w", err)
	}

	return nil
}

func (s *AuthService) GetTokenData(userID, encryptedToken string) (*dto.UserTokenData, error) {
	ctx := context.Background()

	tokenKey := fmt.Sprintf("auth:token:%s:%s", userID, encryptedToken)

	jsonData, err := global.RedisClient.Get(ctx, tokenKey).Result()
	if err != nil {
		return nil, fmt.Errorf("token not found or expired: %w", err)
	}

	var tokenData dto.UserTokenData
	err = json.Unmarshal([]byte(jsonData), &tokenData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal token data: %w", err)
	}

	return &tokenData, nil
}

func (s *AuthService) DeleteTokenData(userID, encryptedToken string) error {
	ctx := context.Background()

	tokenKey := fmt.Sprintf("auth:token:%s:%s", userID, encryptedToken)
	err := global.RedisClient.Del(ctx, tokenKey).Err()
	if err != nil {
		return fmt.Errorf("failed to delete token data from Redis: %w", err)
	}

	return nil
}

func (s *AuthService) InvalidateUserTokens(userID string) error {
	ctx := context.Background()

	pattern := fmt.Sprintf("auth:token:%s:*", userID)

	pipe := global.RedisClient.Pipeline()

	keys, err := global.RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get user token keys: %w", err)
	}

	for _, key := range keys {
		pipe.Del(ctx, key)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete user tokens from Redis: %w", err)
	}

	return nil
}

// BuildTokenData creates RBAC data from user with workspaces
func (s *AuthService) BuildTokenData(user *models.User) *dto.UserTokenData {
	tokenData := &dto.UserTokenData{
		GlobalRole: user.GlobalRole,
	}
	for _, membership := range user.WorkspaceMemberships {
		if membership.Status == "active" && membership.RoleID != "" {
			workspaceMembership := dto.WorkspaceMembershipTokenData{
				WorkspaceID: membership.WorkspaceID,
				RoleName:    membership.Role.Name,
				Permissions: membership.Role.Permissions.Permissions,
				Status:      membership.Status,
			}
			tokenData.WorkspaceMemberships = append(tokenData.WorkspaceMemberships, workspaceMembership)
		}
	}
	return tokenData
}

func (s *AuthService) Logout(userID, encryptedToken string) error {
	err := s.DeleteTokenData(userID, encryptedToken)
	if err != nil {
		fmt.Printf("Warning: failed to delete token from Redis: %v\n", err)
	}
	return nil
}
