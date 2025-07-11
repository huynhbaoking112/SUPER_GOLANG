package repo_test

import (
	"fmt"
	"go-backend-v2/internal/models"
	"go-backend-v2/internal/repo"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type WorkspaceRepositoryTestSuite struct {
	suite.Suite
	db            *gorm.DB
	workspaceRepo repo.WorkspaceRepositoryInterface
	userRepo      repo.UserRepositoryInterface
	testUser      *models.User
	testDBName    string
}

func (suite *WorkspaceRepositoryTestSuite) SetupSuite() {
	// Generate unique test database name
	suite.testDBName = fmt.Sprintf("test_workspace_repo_%d", time.Now().Unix())

	// MySQL test configuration - can be overridden by environment variables
	host := getEnvOrDefault("TEST_MYSQL_HOST", "localhost")
	port := getEnvOrDefault("TEST_MYSQL_PORT", "3307")
	user := getEnvOrDefault("TEST_MYSQL_USER", "root")
	password := getEnvOrDefault("TEST_MYSQL_PASSWORD", "root")

	// Create test database
	rootDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, password, host, port)
	rootDB, err := gorm.Open(mysql.Open(rootDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(suite.T(), err, "Failed to connect to MySQL root")

	// Create test database
	err = rootDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", suite.testDBName)).Error
	assert.NoError(suite.T(), err, "Failed to create test database")

	// Close root connection
	sqlDB, _ := rootDB.DB()
	sqlDB.Close()

	// Connect to test database
	testDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, suite.testDBName)

	testDB, err := gorm.Open(mysql.Open(testDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(suite.T(), err, "Failed to connect to test database")

	suite.db = testDB

	// Auto-migrate the schema
	err = suite.db.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.UserAuthProvider{},
		&models.Workspace{},
		&models.WorkspaceRole{},
		&models.UserWorkspaceMembership{},
		&models.Resource{},
	)
	assert.NoError(suite.T(), err, "Failed to migrate test database schema")

	// Create repositories with dependency injection pattern
	// Note: In real implementation, we'd inject the DB through constructor
	suite.workspaceRepo = &repo.WorkspaceRepository{}
	suite.userRepo = &repo.UserRepository{}
}

func (suite *WorkspaceRepositoryTestSuite) SetupTest() {
	// Clean up data before each test
	// Use TRUNCATE for better performance and to reset auto-increment
	tables := []string{
		"user_workspace_memberships",
		"workspace_roles",
		"workspaces",
		"user_auth_providers",
		"user_profiles",
		"users",
		"resources",
	}

	// Disable foreign key checks temporarily for truncation
	suite.db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	for _, table := range tables {
		suite.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table))
	}
	suite.db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	// Create a test user for workspace operations
	suite.testUser = suite.createTestUser()
}

func (suite *WorkspaceRepositoryTestSuite) TearDownSuite() {
	// Drop test database
	if suite.db != nil {
		suite.db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", suite.testDBName))
		sqlDB, _ := suite.db.DB()
		sqlDB.Close()
	}
}

func (suite *WorkspaceRepositoryTestSuite) createTestUser() *models.User {
	user := &models.User{
		ID:         generateTestID("user"),
		Email:      "admin@example.com",
		GlobalRole: "super_admin",
		Status:     "active",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	profile := &models.UserProfile{
		UserID:    user.ID,
		FirstName: "Super",
		LastName:  "Admin",
		Timezone:  "UTC",
		Locale:    "en",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	authProvider := &models.UserAuthProvider{
		UserID:         user.ID,
		Provider:       "local",
		ProviderUserID: user.Email,
		ProviderEmail:  &user.Email,
		PasswordHash:   stringPtr("$2a$10$hashedpassword"),
		IsPrimary:      true,
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Create user with transaction
	tx := suite.db.Begin()
	err := suite.userRepo.CreateUserWithAuth(tx, user, profile, authProvider)
	if err != nil {
		tx.Rollback()
		suite.T().Fatalf("Failed to create test user: %v", err)
	}
	tx.Commit()

	return user
}

func (suite *WorkspaceRepositoryTestSuite) createTestWorkspace() *models.Workspace {
	workspace := &models.Workspace{
		ID:          generateTestID("workspace"),
		Name:        "Test Workspace",
		Slug:        "test-workspace",
		Description: stringPtr("A test workspace"),
		OwnerID:     suite.testUser.ID,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tx := suite.db.Begin()
	err := suite.workspaceRepo.CreateWorkspace(tx, workspace)
	if err != nil {
		tx.Rollback()
		suite.T().Fatalf("Failed to create test workspace: %v", err)
	}
	tx.Commit()

	return workspace
}

// Test Workspace CRUD Operations

func (suite *WorkspaceRepositoryTestSuite) TestCreateWorkspace_Success() {
	workspace := &models.Workspace{
		Name:        "New Workspace",
		Slug:        "new-workspace",
		Description: stringPtr("A new workspace for testing"),
		OwnerID:     suite.testUser.ID,
		Status:      "active",
	}

	tx := suite.db.Begin()
	err := suite.workspaceRepo.CreateWorkspace(tx, workspace)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), workspace.ID)
	tx.Commit()

	// Verify data was saved
	var savedWorkspace models.Workspace
	err = suite.db.First(&savedWorkspace, "slug = ?", workspace.Slug).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), workspace.Name, savedWorkspace.Name)
	assert.Equal(suite.T(), workspace.Slug, savedWorkspace.Slug)
	assert.Equal(suite.T(), suite.testUser.ID, savedWorkspace.OwnerID)
}

func (suite *WorkspaceRepositoryTestSuite) TestCreateWorkspace_TransactionRollback() {
	workspace := &models.Workspace{
		Name:    "Invalid Workspace",
		Slug:    "invalid-workspace",
		OwnerID: "nonexistent-user", // This should cause foreign key constraint error
		Status:  "active",
	}

	tx := suite.db.Begin()
	err := suite.workspaceRepo.CreateWorkspace(tx, workspace)

	// Should get error due to foreign key constraint
	assert.Error(suite.T(), err)
	tx.Rollback()

	// Verify no data was saved
	var count int64
	suite.db.Model(&models.Workspace{}).Where("slug = ?", workspace.Slug).Count(&count)
	assert.Equal(suite.T(), int64(0), count)
}

func (suite *WorkspaceRepositoryTestSuite) TestTransaction_CompleteWorkspaceCreation() {
	// Test complete workspace creation with role and membership in single transaction
	workspace := &models.Workspace{
		Name:        "Complete Workspace",
		Slug:        "complete-workspace",
		Description: stringPtr("Complete workspace with role and membership"),
		OwnerID:     suite.testUser.ID,
		Status:      "active",
	}

	role := &models.WorkspaceRole{
		Name:        "admin",
		Description: stringPtr("Admin role"),
		Permissions: models.RolePermissions{
			Permissions: []string{"all"},
			Metadata: models.PermissionMetadata{
				Version:   "1.0",
				CreatedBy: "system",
				UpdatedBy: "system",
				UpdatedAt: time.Now(),
			},
		},
		Status: "active",
	}

	membership := &models.UserWorkspaceMembership{
		UserID: suite.testUser.ID,
		Status: "active",
	}

	// Execute in transaction
	tx := suite.db.Begin()

	// 1. Create workspace
	err := suite.workspaceRepo.CreateWorkspace(tx, workspace)
	assert.NoError(suite.T(), err)

	// 2. Create role
	role.WorkspaceID = workspace.ID
	err = suite.workspaceRepo.CreateWorkspaceRole(tx, role)
	assert.NoError(suite.T(), err)

	// 3. Create membership
	membership.WorkspaceID = workspace.ID
	membership.RoleID = role.ID
	now := time.Now()
	membership.JoinedAt = &now
	err = suite.workspaceRepo.CreateMembership(tx, membership)
	assert.NoError(suite.T(), err)

	tx.Commit()

	// Verify all data was created
	var savedWorkspace models.Workspace
	err = suite.db.First(&savedWorkspace, "id = ?", workspace.ID).Error
	assert.NoError(suite.T(), err)

	var savedRole models.WorkspaceRole
	err = suite.db.First(&savedRole, "id = ?", role.ID).Error
	assert.NoError(suite.T(), err)

	var savedMembership models.UserWorkspaceMembership
	err = suite.db.First(&savedMembership, "id = ?", membership.ID).Error
	assert.NoError(suite.T(), err)

	// Verify relationships
	assert.Equal(suite.T(), workspace.ID, savedRole.WorkspaceID)
	assert.Equal(suite.T(), workspace.ID, savedMembership.WorkspaceID)
	assert.Equal(suite.T(), role.ID, savedMembership.RoleID)
	assert.Equal(suite.T(), suite.testUser.ID, savedMembership.UserID)
}

// Helper functions

func stringPtr(s string) *string {
	return &s
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func generateTestID(prefix string) string {
	return fmt.Sprintf("%s_%d_%d", prefix, time.Now().UnixNano(), time.Now().Nanosecond()%1000)
}

func TestWorkspaceRepositorySuite(t *testing.T) {
	suite.Run(t, new(WorkspaceRepositoryTestSuite))
}
