package repo

import (
	"fmt"
	"go-backend-v2/global"
	"go-backend-v2/internal/models"

	"gorm.io/gorm"
)

type WorkspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository() WorkspaceRepositoryInterface {
	return &WorkspaceRepository{
		db: global.DB,
	}
}

func (r *WorkspaceRepository) CreateWorkspace(tx *gorm.DB, workspace *models.Workspace) error {
	if err := tx.Create(workspace).Error; err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}
	return nil
}

func (r *WorkspaceRepository) GetWorkspaceByID(workspaceID string) (*models.Workspace, error) {
	var workspace models.Workspace

	err := r.db.Where("id = ?", workspaceID).First(&workspace).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get workspace by ID: %w", err)
	}

	return &workspace, nil
}

// GetWorkspaceBySlug retrieves a workspace by its slug
func (r *WorkspaceRepository) GetWorkspaceBySlug(slug string) (*models.Workspace, error) {
	var workspace models.Workspace

	err := r.db.Where("slug = ?", slug).First(&workspace).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get workspace by slug: %w", err)
	}

	return &workspace, nil
}

// GetWorkspacesByOwnerID retrieves all workspaces owned by a user
func (r *WorkspaceRepository) GetWorkspacesByOwnerID(ownerID string) ([]models.Workspace, error) {
	var workspaces []models.Workspace

	err := r.db.Where("owner_id = ?", ownerID).Find(&workspaces).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get workspaces by owner ID: %w", err)
	}

	return workspaces, nil
}

// UpdateWorkspace updates a workspace with the given updates
func (r *WorkspaceRepository) UpdateWorkspace(workspaceID string, updates map[string]interface{}) error {
	result := r.db.Model(&models.Workspace{}).Where("id = ?", workspaceID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update workspace: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// DeleteWorkspace soft deletes a workspace by updating its status
func (r *WorkspaceRepository) DeleteWorkspace(workspaceID string) error {
	result := r.db.Model(&models.Workspace{}).Where("id = ?", workspaceID).Update("status", "deleted")
	if result.Error != nil {
		return fmt.Errorf("failed to delete workspace: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// CreateMembership creates a new workspace membership within a transaction
func (r *WorkspaceRepository) CreateMembership(tx *gorm.DB, membership *models.UserWorkspaceMembership) error {
	if err := tx.Create(membership).Error; err != nil {
		return fmt.Errorf("failed to create membership: %w", err)
	}
	return nil
}

// GetMembership retrieves a membership by user ID and workspace ID
func (r *WorkspaceRepository) GetMembership(userID, workspaceID string) (*models.UserWorkspaceMembership, error) {
	var membership models.UserWorkspaceMembership

	err := r.db.Where("user_id = ? AND workspace_id = ?", userID, workspaceID).First(&membership).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get membership: %w", err)
	}

	return &membership, nil
}

// GetUserMemberships retrieves all memberships for a user
func (r *WorkspaceRepository) GetUserMemberships(userID string) ([]models.UserWorkspaceMembership, error) {
	var memberships []models.UserWorkspaceMembership

	err := r.db.Where("user_id = ?", userID).Find(&memberships).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user memberships: %w", err)
	}

	return memberships, nil
}

// GetWorkspaceMemberships retrieves all memberships for a workspace
func (r *WorkspaceRepository) GetWorkspaceMemberships(workspaceID string) ([]models.UserWorkspaceMembership, error) {
	var memberships []models.UserWorkspaceMembership

	err := r.db.Where("workspace_id = ?", workspaceID).Find(&memberships).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace memberships: %w", err)
	}

	return memberships, nil
}

// UpdateMembership updates a membership with the given updates
func (r *WorkspaceRepository) UpdateMembership(membershipID string, updates map[string]interface{}) error {
	result := r.db.Model(&models.UserWorkspaceMembership{}).Where("id = ?", membershipID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update membership: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// DeleteMembership soft deletes a membership by updating its status
func (r *WorkspaceRepository) DeleteMembership(membershipID string) error {
	result := r.db.Model(&models.UserWorkspaceMembership{}).Where("id = ?", membershipID).Update("status", "inactive")
	if result.Error != nil {
		return fmt.Errorf("failed to delete membership: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// CreateWorkspaceRole creates a new workspace role within a transaction
func (r *WorkspaceRepository) CreateWorkspaceRole(tx *gorm.DB, role *models.WorkspaceRole) error {
	if err := tx.Create(role).Error; err != nil {
		return fmt.Errorf("failed to create workspace role: %w", err)
	}
	return nil
}

// GetWorkspaceRole retrieves a workspace role by its ID
func (r *WorkspaceRepository) GetWorkspaceRole(roleID string) (*models.WorkspaceRole, error) {
	var role models.WorkspaceRole

	err := r.db.Where("id = ?", roleID).First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get workspace role: %w", err)
	}

	return &role, nil
}

// GetWorkspaceRoles retrieves all roles for a workspace
func (r *WorkspaceRepository) GetWorkspaceRoles(workspaceID string) ([]models.WorkspaceRole, error) {
	var roles []models.WorkspaceRole

	err := r.db.Where("workspace_id = ?", workspaceID).Find(&roles).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace roles: %w", err)
	}

	return roles, nil
}

// UpdateWorkspaceRole updates a workspace role with the given updates
func (r *WorkspaceRepository) UpdateWorkspaceRole(roleID string, updates map[string]interface{}) error {
	result := r.db.Model(&models.WorkspaceRole{}).Where("id = ?", roleID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update workspace role: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// DeleteWorkspaceRole soft deletes a workspace role by updating its status
func (r *WorkspaceRepository) DeleteWorkspaceRole(roleID string) error {
	result := r.db.Model(&models.WorkspaceRole{}).Where("id = ?", roleID).Update("status", "inactive")
	if result.Error != nil {
		return fmt.Errorf("failed to delete workspace role: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// ExistsBySlug checks if a workspace with the given slug exists
func (r *WorkspaceRepository) ExistsBySlug(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Workspace{}).Where("slug = ?", slug).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check workspace existence by slug: %w", err)
	}
	return count > 0, nil
}

// ExistsByID checks if a workspace with the given ID exists
func (r *WorkspaceRepository) ExistsByID(workspaceID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Workspace{}).Where("id = ?", workspaceID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check workspace existence by ID: %w", err)
	}
	return count > 0, nil
}
