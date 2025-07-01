package repo

import (
	"go-backend-v2/internal/models"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	CreateUserWithAuth(tx *gorm.DB, user *models.User, profile *models.UserProfile, authProvider *models.UserAuthProvider) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID string) (*models.User, error)
	GetUserWithProfile(userID string) (*models.User, error)
	GetUserWithWorkspaces(userID string) (*models.User, error)
	UpdateUser(userID string, updates map[string]interface{}) error
	DeleteUser(userID string) error

	GetUserAuthProvider(userID, provider string) (*models.UserAuthProvider, error)
	GetUserAuthProviders(userID string) ([]models.UserAuthProvider, error)
	CreateAuthProvider(authProvider *models.UserAuthProvider) error
	UpdateAuthProvider(providerID string, updates map[string]interface{}) error
	DeleteAuthProvider(providerID string) error

	CreateUserProfile(profile *models.UserProfile) error
	UpdateUserProfile(userID string, updates map[string]interface{}) error

	ExistsByEmail(email string) (bool, error)
	ExistsByID(userID string) (bool, error)
}

type WorkspaceRepositoryInterface interface {
	CreateWorkspace(workspace *models.Workspace) error
	GetWorkspaceByID(workspaceID string) (*models.Workspace, error)
	GetWorkspaceBySlug(slug string) (*models.Workspace, error)
	GetWorkspacesByOwnerID(ownerID string) ([]models.Workspace, error)
	UpdateWorkspace(workspaceID string, updates map[string]interface{}) error
	DeleteWorkspace(workspaceID string) error

	CreateMembership(membership *models.UserWorkspaceMembership) error
	GetMembership(userID, workspaceID string) (*models.UserWorkspaceMembership, error)
	GetUserMemberships(userID string) ([]models.UserWorkspaceMembership, error)
	GetWorkspaceMemberships(workspaceID string) ([]models.UserWorkspaceMembership, error)
	UpdateMembership(membershipID string, updates map[string]interface{}) error
	DeleteMembership(membershipID string) error

	CreateWorkspaceRole(role *models.WorkspaceRole) error
	GetWorkspaceRole(roleID string) (*models.WorkspaceRole, error)
	GetWorkspaceRoles(workspaceID string) ([]models.WorkspaceRole, error)
	UpdateWorkspaceRole(roleID string, updates map[string]interface{}) error
	DeleteWorkspaceRole(roleID string) error

	ExistsBySlug(slug string) (bool, error)
	ExistsByID(workspaceID string) (bool, error)
}

type ResourceRepositoryInterface interface {
	CreateResource(resource *models.Resource) error
	GetResourceByID(resourceID string) (*models.Resource, error)
	GetResourceByName(name string) (*models.Resource, error)
	GetResources() ([]models.Resource, error)
	GetResourcesByCategory(category string) ([]models.Resource, error)
	UpdateResource(resourceID string, updates map[string]interface{}) error
	DeleteResource(resourceID string) error

	ExistsByName(name string) (bool, error)
	ExistsByID(resourceID string) (bool, error)
}

type TransactionRepositoryInterface interface {
	BeginTransaction() *gorm.DB
	CommitTransaction(tx *gorm.DB) error
	RollbackTransaction(tx *gorm.DB) error

	WithTransaction(fn func(tx *gorm.DB) error) error
}
