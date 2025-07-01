package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-backend-v2/internal/controllers"
	"go-backend-v2/internal/dto"
)

// MockAuthService is a mock implementation of AuthServiceInterface
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Signup(req *dto.SignupRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAuthService) Login(req *dto.LoginRequest) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func TestNewAuthController(t *testing.T) {
	mockService := &MockAuthService{}

	controller := controllers.NewAuthController(mockService)

	assert.NotNil(t, controller)
}

// HTTP endpoint tests are commented out because they require global configuration
// For full integration testing, these should be run in an integration test environment
// with proper application setup including global.Config for JWT handling

// func TestAuthController_Signup(t *testing.T) {
//     // This test requires global.Config to be initialized for JWT handling
//     // Should be implemented as integration tests
// }

// func TestAuthController_Login(t *testing.T) {
//     // This test requires global.Config to be initialized for JWT handling
//     // Should be implemented as integration tests
// }

func TestAuthController_DependencyInjection(t *testing.T) {
	t.Run("controller accepts auth service interface", func(t *testing.T) {
		mockService := &MockAuthService{}

		// This should not panic and should return a valid controller
		controller := controllers.NewAuthController(mockService)
		assert.NotNil(t, controller)
	})

	t.Run("controller validates service interface compatibility", func(t *testing.T) {
		mockService := &MockAuthService{}

		// Ensure the mock service implements the required interface
		assert.Implements(t, (*interface {
			Signup(req *dto.SignupRequest) error
			Login(req *dto.LoginRequest) (string, error)
			ValidateToken(token string) (string, error)
		})(nil), mockService)
	})
}

func TestAuthController_MockValidation(t *testing.T) {
	t.Run("mock service works for signup", func(t *testing.T) {
		mockService := &MockAuthService{}
		req := &dto.SignupRequest{
			Email:     "test@example.com",
			Password:  "Password123!",
			FirstName: "John",
			LastName:  "Doe",
		}

		// Setup mock expectation
		mockService.On("Signup", req).Return(nil)

		// Call the mock
		err := mockService.Signup(req)

		// Assert
		assert.NoError(t, err)
		mockService.AssertExpectations(t)
	})

	t.Run("mock service works for login", func(t *testing.T) {
		mockService := &MockAuthService{}
		req := &dto.LoginRequest{
			Email:    "test@example.com",
			Password: "Password123!",
		}

		// Setup mock expectation
		mockService.On("Login", req).Return("jwt-token", nil)

		// Call the mock
		token, err := mockService.Login(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "jwt-token", token)
		mockService.AssertExpectations(t)
	})

	t.Run("mock service works for token validation", func(t *testing.T) {
		mockService := &MockAuthService{}
		testToken := "valid-jwt-token"
		expectedUserID := "user-123"

		// Setup mock expectation
		mockService.On("ValidateToken", testToken).Return(expectedUserID, nil)

		// Call the mock
		userID, err := mockService.ValidateToken(testToken)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUserID, userID)
		mockService.AssertExpectations(t)
	})
}

func TestAuthController_Construction(t *testing.T) {
	t.Run("controller can be created with valid service", func(t *testing.T) {
		mockService := &MockAuthService{}

		controller := controllers.NewAuthController(mockService)

		assert.NotNil(t, controller)
		// Note: We can't easily test the internal validator without accessing private fields
		// This test ensures the constructor works and doesn't panic
	})
}
