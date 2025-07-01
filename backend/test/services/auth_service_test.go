package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"go-backend-v2/internal/common"
	"go-backend-v2/internal/dto"
	"go-backend-v2/internal/models"
	"go-backend-v2/internal/services"
	"go-backend-v2/pkg/utils"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUserWithAuth(tx *gorm.DB, user *models.User, profile *models.UserProfile, authProvider *models.UserAuthProvider) error {
	args := m.Called(tx, user, profile, authProvider)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(userID string) (*models.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserWithProfile(userID string) (*models.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserWithWorkspaces(userID string) (*models.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(userID string, updates map[string]interface{}) error {
	args := m.Called(userID, updates)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserAuthProvider(userID, provider string) (*models.UserAuthProvider, error) {
	args := m.Called(userID, provider)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserAuthProvider), args.Error(1)
}

func (m *MockUserRepository) GetUserAuthProviders(userID string) ([]models.UserAuthProvider, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.UserAuthProvider), args.Error(1)
}

func (m *MockUserRepository) CreateAuthProvider(authProvider *models.UserAuthProvider) error {
	args := m.Called(authProvider)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateAuthProvider(providerID string, updates map[string]interface{}) error {
	args := m.Called(providerID, updates)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteAuthProvider(providerID string) error {
	args := m.Called(providerID)
	return args.Error(0)
}

func (m *MockUserRepository) CreateUserProfile(profile *models.UserProfile) error {
	args := m.Called(profile)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserProfile(userID string, updates map[string]interface{}) error {
	args := m.Called(userID, updates)
	return args.Error(0)
}

func (m *MockUserRepository) ExistsByEmail(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByID(userID string) (bool, error) {
	args := m.Called(userID)
	return args.Bool(0), args.Error(1)
}

func TestNewAuthService(t *testing.T) {
	mockRepo := &MockUserRepository{}

	authService := services.NewAuthService(mockRepo)

	assert.NotNil(t, authService)
	assert.Implements(t, (*services.AuthServiceInterface)(nil), authService)
}

func TestAuthService_Signup(t *testing.T) {
	tests := []struct {
		name          string
		request       *dto.SignupRequest
		setupMock     func(*MockUserRepository)
		expectedError error
	}{
		{
			name: "successful signup",
			request: &dto.SignupRequest{
				Email:     "test@example.com",
				Password:  "Password123!",
				FirstName: "John",
				LastName:  "Doe",
			},
			setupMock: func(mockRepo *MockUserRepository) {
				mockRepo.On("ExistsByEmail", "test@example.com").Return(false, nil)
				mockRepo.On("CreateUserWithAuth", mock.Anything, mock.AnythingOfType("*models.User"), mock.AnythingOfType("*models.UserProfile"), mock.AnythingOfType("*models.UserAuthProvider")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "user already exists",
			request: &dto.SignupRequest{
				Email:     "existing@example.com",
				Password:  "Password123!",
				FirstName: "Jane",
				LastName:  "Doe",
			},
			setupMock: func(mockRepo *MockUserRepository) {
				mockRepo.On("ExistsByEmail", "existing@example.com").Return(true, nil)
			},
			expectedError: common.ErrEmailAlreadyExists,
		},
		{
			name: "weak password",
			request: &dto.SignupRequest{
				Email:     "test@example.com",
				Password:  "123",
				FirstName: "John",
				LastName:  "Doe",
			},
			setupMock:     func(mockRepo *MockUserRepository) {},
			expectedError: common.ErrPasswordTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			tt.setupMock(mockRepo)

			authService := services.NewAuthService(mockRepo)
			err := authService.Signup(tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	hashedPassword, _ := utils.HashPassword("Password123!")

	testUser := &models.User{
		ID:     "user-123",
		Email:  "test@example.com",
		Status: common.UserStatusActive,
	}

	testAuthProvider := &models.UserAuthProvider{
		ID:           "auth-123",
		UserID:       "user-123",
		Provider:     common.AuthProviderLocal,
		PasswordHash: &hashedPassword,
		IsPrimary:    true,
	}

	tests := []struct {
		name          string
		request       *dto.LoginRequest
		setupMock     func(*MockUserRepository)
		expectedError error
		expectToken   bool
	}{
		{
			name: "successful login",
			request: &dto.LoginRequest{
				Email:    "test@example.com",
				Password: "Password123!",
			},
			setupMock: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByEmail", "test@example.com").Return(testUser, nil)
				mockRepo.On("GetUserAuthProvider", "user-123", common.AuthProviderLocal).Return(testAuthProvider, nil)
				mockRepo.On("UpdateUser", "user-123", mock.Anything).Return(nil)
			},
			expectedError: nil,
			expectToken:   true,
		},
		{
			name: "user not found",
			request: &dto.LoginRequest{
				Email:    "notfound@example.com",
				Password: "Password123!",
			},
			setupMock: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByEmail", "notfound@example.com").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: common.ErrInvalidCredentials,
			expectToken:   false,
		},
		{
			name: "invalid password",
			request: &dto.LoginRequest{
				Email:    "test@example.com",
				Password: "WrongPassword",
			},
			setupMock: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByEmail", "test@example.com").Return(testUser, nil)
				mockRepo.On("GetUserAuthProvider", "user-123", common.AuthProviderLocal).Return(testAuthProvider, nil)
			},
			expectedError: common.ErrInvalidCredentials,
			expectToken:   false,
		},
		{
			name: "inactive user",
			request: &dto.LoginRequest{
				Email:    "inactive@example.com",
				Password: "Password123!",
			},
			setupMock: func(mockRepo *MockUserRepository) {
				inactiveUser := &models.User{
					ID:     "user-inactive",
					Email:  "inactive@example.com",
					Status: common.UserStatusInactive,
				}
				mockRepo.On("GetUserByEmail", "inactive@example.com").Return(inactiveUser, nil)
			},
			expectedError: common.ErrUserInactive,
			expectToken:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			tt.setupMock(mockRepo)

			authService := services.NewAuthService(mockRepo)
			token, err := authService.Login(tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				if tt.expectToken {
					assert.NotEmpty(t, token)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	testUser := &models.User{
		ID:     "user-123",
		Email:  "test@example.com",
		Status: common.UserStatusActive,
	}

	validToken, _ := utils.GenerateToken(testUser.ID)

	tests := []struct {
		name           string
		token          string
		setupMock      func(*MockUserRepository)
		expectedUserID string
		expectedError  error
	}{
		{
			name:  "valid token and active user",
			token: validToken,
			setupMock: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByID", testUser.ID).Return(testUser, nil)
			},
			expectedUserID: testUser.ID,
			expectedError:  nil,
		},
		{
			name:           "invalid token format",
			token:          "invalid.token.format",
			setupMock:      func(mockRepo *MockUserRepository) {},
			expectedUserID: "",
			expectedError:  common.ErrTokenInvalid,
		},
		{
			name:  "user not found",
			token: validToken,
			setupMock: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByID", testUser.ID).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedUserID: "",
			expectedError:  common.ErrUserNotFound,
		},
		{
			name:  "inactive user",
			token: validToken,
			setupMock: func(mockRepo *MockUserRepository) {
				inactiveUser := &models.User{
					ID:     "user-123",
					Email:  "test@example.com",
					Status: common.UserStatusInactive,
				}
				mockRepo.On("GetUserByID", testUser.ID).Return(inactiveUser, nil)
			},
			expectedUserID: "",
			expectedError:  common.ErrUserInactive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			tt.setupMock(mockRepo)

			authService := services.NewAuthService(mockRepo)
			userID, err := authService.ValidateToken(tt.token)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Empty(t, userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUserID, userID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_SignupValidation(t *testing.T) {
	tests := []struct {
		name          string
		request       *dto.SignupRequest
		setupMock     func(*MockUserRepository)
		expectedError string
	}{
		{
			name: "user already exists",
			request: &dto.SignupRequest{
				Email:     "existing@example.com",
				Password:  "Password123!",
				FirstName: "Jane",
				LastName:  "Doe",
			},
			setupMock: func(mockRepo *MockUserRepository) {
				mockRepo.On("ExistsByEmail", "existing@example.com").Return(true, nil)
			},
			expectedError: common.ErrEmailAlreadyExists.Error(),
		},
		{
			name: "weak password",
			request: &dto.SignupRequest{
				Email:     "test@example.com",
				Password:  "123",
				FirstName: "John",
				LastName:  "Doe",
			},
			setupMock:     func(mockRepo *MockUserRepository) {},
			expectedError: common.ErrPasswordTooShort.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			tt.setupMock(mockRepo)

			authService := services.NewAuthService(mockRepo)
			err := authService.Signup(tt.request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)

			mockRepo.AssertExpectations(t)
		})
	}
}
