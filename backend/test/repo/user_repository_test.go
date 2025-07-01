package repo_test

import (
	"go-backend-v2/internal/models"
	"go-backend-v2/internal/repo"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repository repo.UserRepositoryInterface
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
	// Create in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(suite.T(), err)

	// Auto-migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.UserAuthProvider{},
		&models.Workspace{},
		&models.WorkspaceRole{},
		&models.UserWorkspaceMembership{},
	)
	assert.NoError(suite.T(), err)

	suite.db = db

	// Create repository with test database
	suite.repository = &repo.UserRepository{}
	// We need to inject the test DB into the repository
	// This would normally be done through dependency injection
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	// Clean up data before each test
	suite.db.Exec("DELETE FROM user_workspace_memberships")
	suite.db.Exec("DELETE FROM user_auth_providers")
	suite.db.Exec("DELETE FROM user_profiles")
	suite.db.Exec("DELETE FROM users")
	suite.db.Exec("DELETE FROM workspace_roles")
	suite.db.Exec("DELETE FROM workspaces")
}

func (suite *UserRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *UserRepositoryTestSuite) createTestUser() *models.User {
	user := &models.User{
		ID:         "test-user-123",
		Email:      "test@example.com",
		GlobalRole: "customer",
		Status:     "active",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	profile := &models.UserProfile{
		UserID:    user.ID,
		FirstName: "Test",
		LastName:  "User",
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
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Create user with transaction
	tx := suite.db.Begin()
	err := suite.repository.CreateUserWithAuth(tx, user, profile, authProvider)
	if err != nil {
		tx.Rollback()
		suite.T().Fatalf("Failed to create test user: %v", err)
	}
	tx.Commit()

	return user
}

func stringPtr(s string) *string {
	return &s
}

func (suite *UserRepositoryTestSuite) TestCreateUserWithAuth_Success() {
	user := &models.User{
		Email:      "new@example.com",
		GlobalRole: "customer",
		Status:     "active",
	}

	profile := &models.UserProfile{
		FirstName: "New",
		LastName:  "User",
		Timezone:  "UTC",
		Locale:    "en",
	}

	authProvider := &models.UserAuthProvider{
		Provider:       "local",
		ProviderUserID: user.Email,
		ProviderEmail:  &user.Email,
		PasswordHash:   stringPtr("$2a$10$hashedpassword"),
		IsPrimary:      true,
	}

	tx := suite.db.Begin()
	err := suite.repository.CreateUserWithAuth(tx, user, profile, authProvider)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), user.ID)
	assert.Equal(suite.T(), user.ID, profile.UserID)
	assert.Equal(suite.T(), user.ID, authProvider.UserID)

	tx.Commit()

	// Verify data was saved
	var savedUser models.User
	err = suite.db.First(&savedUser, "email = ?", user.Email).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Email, savedUser.Email)
}

func (suite *UserRepositoryTestSuite) TestCreateUserWithAuth_TransactionRollback() {
	user := &models.User{
		Email:      "test@example.com",
		GlobalRole: "customer",
		Status:     "active",
	}

	profile := &models.UserProfile{
		FirstName: "Test",
		LastName:  "User",
		Timezone:  "UTC",
		Locale:    "en",
	}

	// Invalid auth provider (missing required field)
	authProvider := &models.UserAuthProvider{
		// Missing Provider field - should cause error
		ProviderUserID: user.Email,
	}

	tx := suite.db.Begin()
	err := suite.repository.CreateUserWithAuth(tx, user, profile, authProvider)

	assert.Error(suite.T(), err)
	tx.Rollback()

	// Verify no data was saved
	var count int64
	suite.db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
	assert.Equal(suite.T(), int64(0), count)
}

func (suite *UserRepositoryTestSuite) TestGetUserByEmail_Found() {
	testUser := suite.createTestUser()

	foundUser, err := suite.repository.GetUserByEmail(testUser.Email)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundUser)
	assert.Equal(suite.T(), testUser.ID, foundUser.ID)
	assert.Equal(suite.T(), testUser.Email, foundUser.Email)
}

func (suite *UserRepositoryTestSuite) TestGetUserByEmail_NotFound() {
	foundUser, err := suite.repository.GetUserByEmail("nonexistent@example.com")

	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), foundUser)
}

func (suite *UserRepositoryTestSuite) TestGetUserByID_Found() {
	testUser := suite.createTestUser()

	foundUser, err := suite.repository.GetUserByID(testUser.ID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundUser)
	assert.Equal(suite.T(), testUser.ID, foundUser.ID)
	assert.Equal(suite.T(), testUser.Email, foundUser.Email)
}

func (suite *UserRepositoryTestSuite) TestGetUserByID_NotFound() {
	foundUser, err := suite.repository.GetUserByID("nonexistent-id")

	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), foundUser)
}

func (suite *UserRepositoryTestSuite) TestGetUserWithProfile_Found() {
	testUser := suite.createTestUser()

	foundUser, err := suite.repository.GetUserWithProfile(testUser.ID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundUser)
	assert.Equal(suite.T(), testUser.ID, foundUser.ID)
	assert.NotNil(suite.T(), foundUser.Profile)
	assert.Equal(suite.T(), "Test", foundUser.Profile.FirstName)
	assert.Equal(suite.T(), "User", foundUser.Profile.LastName)
}

func (suite *UserRepositoryTestSuite) TestUpdateUser_Success() {
	testUser := suite.createTestUser()

	updates := map[string]interface{}{
		"status":      "inactive",
		"global_role": "super_member",
	}

	err := suite.repository.UpdateUser(testUser.ID, updates)

	assert.NoError(suite.T(), err)

	// Verify update
	var updatedUser models.User
	err = suite.db.First(&updatedUser, "id = ?", testUser.ID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "inactive", updatedUser.Status)
	assert.Equal(suite.T(), "super_member", updatedUser.GlobalRole)
}

func (suite *UserRepositoryTestSuite) TestUpdateUser_NotFound() {
	updates := map[string]interface{}{
		"status": "inactive",
	}

	err := suite.repository.UpdateUser("nonexistent-id", updates)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrRecordNotFound, err)
}

func (suite *UserRepositoryTestSuite) TestDeleteUser_Success() {
	testUser := suite.createTestUser()

	err := suite.repository.DeleteUser(testUser.ID)

	assert.NoError(suite.T(), err)

	// Verify deletion (soft delete)
	var deletedUser models.User
	err = suite.db.Unscoped().First(&deletedUser, "id = ?", testUser.ID).Error
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), deletedUser.DeletedAt)
}

func (suite *UserRepositoryTestSuite) TestDeleteUser_NotFound() {
	err := suite.repository.DeleteUser("nonexistent-id")

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrRecordNotFound, err)
}

func (suite *UserRepositoryTestSuite) TestExistsByEmail_Exists() {
	testUser := suite.createTestUser()

	exists, err := suite.repository.ExistsByEmail(testUser.Email)

	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)
}

func (suite *UserRepositoryTestSuite) TestExistsByEmail_NotExists() {
	exists, err := suite.repository.ExistsByEmail("nonexistent@example.com")

	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

func (suite *UserRepositoryTestSuite) TestExistsByID_Exists() {
	testUser := suite.createTestUser()

	exists, err := suite.repository.ExistsByID(testUser.ID)

	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)
}

func (suite *UserRepositoryTestSuite) TestExistsByID_NotExists() {
	exists, err := suite.repository.ExistsByID("nonexistent-id")

	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

func (suite *UserRepositoryTestSuite) TestGetUserAuthProvider_Found() {
	testUser := suite.createTestUser()

	authProvider, err := suite.repository.GetUserAuthProvider(testUser.ID, "local")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), authProvider)
	assert.Equal(suite.T(), testUser.ID, authProvider.UserID)
	assert.Equal(suite.T(), "local", authProvider.Provider)
	assert.True(suite.T(), authProvider.IsPrimary)
}

func (suite *UserRepositoryTestSuite) TestGetUserAuthProvider_NotFound() {
	testUser := suite.createTestUser()

	authProvider, err := suite.repository.GetUserAuthProvider(testUser.ID, "google")

	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), authProvider)
}

func (suite *UserRepositoryTestSuite) TestGetUserAuthProviders_Found() {
	testUser := suite.createTestUser()

	authProviders, err := suite.repository.GetUserAuthProviders(testUser.ID)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), authProviders, 1)
	assert.Equal(suite.T(), "local", authProviders[0].Provider)
}

func (suite *UserRepositoryTestSuite) TestGetUserAuthProviders_NotFound() {
	authProviders, err := suite.repository.GetUserAuthProviders("nonexistent-id")

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), authProviders, 0)
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
