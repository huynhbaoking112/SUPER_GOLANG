package repo

import (
	"fmt"
	"go-backend-v2/global"
	"go-backend-v2/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepositoryInterface {
	return &UserRepository{
		db: global.DB,
	}
}

func (r *UserRepository) CreateUserWithAuth(tx *gorm.DB, user *models.User, profile *models.UserProfile, authProvider *models.UserAuthProvider) error {
	if err := tx.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	profile.UserID = user.ID
	authProvider.UserID = user.ID

	if err := tx.Create(profile).Error; err != nil {
		return fmt.Errorf("failed to create user profile: %w", err)
	}

	if err := tx.Create(authProvider).Error; err != nil {
		return fmt.Errorf("failed to create auth provider: %w", err)
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(userID string) (*models.User, error) {
	var user models.User

	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserWithProfile(userID string) (*models.User, error) {
	var user models.User

	err := r.db.Preload("Profile").Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user with profile: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserWithWorkspaces(userID string) (*models.User, error) {
	var user models.User

	err := r.db.
		Preload("Profile").
		Preload("WorkspaceMemberships").
		Preload("WorkspaceMemberships.Workspace").
		Preload("WorkspaceMemberships.Role").
		Where("id = ?", userID).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user with workspaces: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(userID string, updates map[string]interface{}) error {
	result := r.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// DeleteUser soft deletes a user
func (r *UserRepository) DeleteUser(userID string) error {
	// result := r.db.Delete(&models.User{}, "id = ?", userID)
	result := r.db.Model(&models.User{}).Where("id = ?", userID).Update("status", "deleted")
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserRepository) GetUserAuthProvider(userID, provider string) (*models.UserAuthProvider, error) {
	var authProvider models.UserAuthProvider

	err := r.db.Where("user_id = ? AND provider = ?", userID, provider).First(&authProvider).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user auth provider: %w", err)
	}

	return &authProvider, nil
}

func (r *UserRepository) GetUserAuthProviders(userID string) ([]models.UserAuthProvider, error) {
	var authProviders []models.UserAuthProvider

	err := r.db.Where("user_id = ?", userID).Find(&authProviders).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user auth providers: %w", err)
	}

	return authProviders, nil
}

func (r *UserRepository) CreateAuthProvider(authProvider *models.UserAuthProvider) error {
	if err := r.db.Create(authProvider).Error; err != nil {
		return fmt.Errorf("failed to create auth provider: %w", err)
	}
	return nil
}

func (r *UserRepository) UpdateAuthProvider(providerID string, updates map[string]interface{}) error {
	result := r.db.Model(&models.UserAuthProvider{}).Where("id = ?", providerID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update auth provider: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// DeleteAuthProvider deletes an auth provider
func (r *UserRepository) DeleteAuthProvider(providerID string) error {
	// result := r.db.Delete(&models.UserAuthProvider{}, "id = ?", providerID)
	result := r.db.Model(&models.UserAuthProvider{}).Where("id = ?", providerID).Update("status", "deleted")
	if result.Error != nil {
		return fmt.Errorf("failed to delete auth provider: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserRepository) CreateUserProfile(profile *models.UserProfile) error {
	if err := r.db.Create(profile).Error; err != nil {
		return fmt.Errorf("failed to create user profile: %w", err)
	}
	return nil
}

func (r *UserRepository) UpdateUserProfile(userID string, updates map[string]interface{}) error {
	result := r.db.Model(&models.UserProfile{}).Where("user_id = ?", userID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update user profile: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check user existence by email: %w", err)
	}
	return count > 0, nil
}

func (r *UserRepository) ExistsByID(userID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("id = ?", userID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check user existence by ID: %w", err)
	}
	return count > 0, nil
}
