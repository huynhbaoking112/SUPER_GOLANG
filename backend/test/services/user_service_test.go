package services

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"go-backend-v2/internal/common"
	"go-backend-v2/internal/models"
	"go-backend-v2/internal/services"
)

func TestNewUserService(t *testing.T) {
	mockRepo := &MockUserRepository{}

	userService := services.NewUserService(mockRepo)

	assert.NotNil(t, userService)
	assert.Implements(t, (*services.UserServiceInterface)(nil), userService)
}

func TestUserService_GetUserWithWorkspaces(t *testing.T) {
	testUser := &models.User{
		ID:        "user-123",
		Email:     "test@example.com",
		Status:    common.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Profile: &models.UserProfile{
			UserID:    "user-123",
			FirstName: "John",
			LastName:  "Doe",
		},
		WorkspaceMemberships: []models.UserWorkspaceMembership{
			{
				ID:          "membership-1",
				UserID:      "user-123",
				WorkspaceID: "workspace-1",
				RoleID:      "role-1",
			},
		},
	}

	tests := []struct {
		name          string
		userID        string
		setupMock     func(*MockUserRepository)
		expectedError error
		expectUser    bool
	}{
		{
			name:   "successful retrieval with workspaces",
			userID: "user-123",
			setupMock: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserWithWorkspaces", "user-123").Return(testUser, nil)
			},
			expectedError: nil,
			expectUser:    true,
		},
		{
			name:   "user not found",
			userID: "nonexistent-user",
			setupMock: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserWithWorkspaces", "nonexistent-user").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
			expectUser:    false,
		},
		{
			name:   "database error",
			userID: "user-123",
			setupMock: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserWithWorkspaces", "user-123").Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			expectUser:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			tt.setupMock(mockRepo)

			userService := services.NewUserService(mockRepo)
			user, err := userService.GetUserWithWorkspaces(tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, testUser.ID, user.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserProfile(t *testing.T) {
	testUser := &models.User{
		ID:        "user-123",
		Email:     "test@example.com",
		Status:    common.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Profile: &models.UserProfile{
			UserID:    "user-123",
			FirstName: "John",
			LastName:  "Doe",
		},
	}

	tests := []struct {
		name          string
		userID        string
		setupMock     func(*MockUserRepository)
		expectedError error
		expectUser    bool
	}{
		{
			name:   "successful profile retrieval",
			userID: "user-123",
			setupMock: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserWithProfile", "user-123").Return(testUser, nil)
			},
			expectedError: nil,
			expectUser:    true,
		},
		{
			name:   "user not found",
			userID: "nonexistent-user",
			setupMock: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserWithProfile", "nonexistent-user").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
			expectUser:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			tt.setupMock(mockRepo)

			userService := services.NewUserService(mockRepo)
			user, err := userService.GetUserProfile(tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, testUser.ID, user.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Integration(t *testing.T) {
	t.Run("compare profile vs workspaces retrieval", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		userService := services.NewUserService(mockRepo)

		userID := "integration-user"

		// User data for profile retrieval
		profileUser := &models.User{
			ID:     userID,
			Email:  "integration@example.com",
			Status: common.UserStatusActive,
			Profile: &models.UserProfile{
				UserID:    userID,
				FirstName: "Integration",
				LastName:  "Test",
			},
		}

		// User data for workspace retrieval (same user, different associations)
		workspaceUser := &models.User{
			ID:     userID,
			Email:  "integration@example.com",
			Status: common.UserStatusActive,
			WorkspaceMemberships: []models.UserWorkspaceMembership{
				{
					ID:          "membership-1",
					UserID:      userID,
					WorkspaceID: "workspace-1",
				},
			},
		}

		mockRepo.On("GetUserWithProfile", userID).Return(profileUser, nil)
		mockRepo.On("GetUserWithWorkspaces", userID).Return(workspaceUser, nil)

		// Test profile retrieval
		profileResult, err := userService.GetUserProfile(userID)
		assert.NoError(t, err)
		assert.Equal(t, "Integration", profileResult.Profile.FirstName)
		assert.Empty(t, profileResult.WorkspaceMemberships)

		// Test workspace retrieval
		workspaceResult, err := userService.GetUserWithWorkspaces(userID)
		assert.NoError(t, err)
		assert.Len(t, workspaceResult.WorkspaceMemberships, 1)
		// Profile might not be loaded in workspace query
		assert.Equal(t, userID, workspaceResult.ID)

		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_EdgeCases(t *testing.T) {
	t.Run("user with special characters in ID", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		userService := services.NewUserService(mockRepo)

		specialUserID := "user-123!@#$%^&*()"
		user := &models.User{
			ID:    specialUserID,
			Email: "special@example.com",
		}

		mockRepo.On("GetUserWithProfile", specialUserID).Return(user, nil)

		result, err := userService.GetUserProfile(specialUserID)
		assert.NoError(t, err)
		assert.Equal(t, specialUserID, result.ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("user with many workspace memberships", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		userService := services.NewUserService(mockRepo)

		userID := "user-many-workspaces"

		// Create user with 100 workspace memberships
		memberships := make([]models.UserWorkspaceMembership, 100)
		for i := 0; i < 100; i++ {
			memberships[i] = models.UserWorkspaceMembership{
				ID:          fmt.Sprintf("membership-%d", i),
				UserID:      userID,
				WorkspaceID: fmt.Sprintf("workspace-%d", i),
			}
		}

		user := &models.User{
			ID:                   userID,
			Email:                "manyworkspaces@example.com",
			WorkspaceMemberships: memberships,
		}

		mockRepo.On("GetUserWithWorkspaces", userID).Return(user, nil)

		result, err := userService.GetUserWithWorkspaces(userID)
		assert.NoError(t, err)
		assert.Len(t, result.WorkspaceMemberships, 100)

		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Performance(t *testing.T) {
	t.Run("concurrent access to same user", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		userService := services.NewUserService(mockRepo)

		userID := "concurrent-user"
		user := &models.User{
			ID:    userID,
			Email: "concurrent@example.com",
		}

		// Setup mock to handle multiple concurrent calls
		mockRepo.On("GetUserWithProfile", userID).Return(user, nil).Times(10)

		// Run 10 concurrent requests
		errChan := make(chan error, 10)
		for i := 0; i < 10; i++ {
			go func() {
				_, err := userService.GetUserProfile(userID)
				errChan <- err
			}()
		}

		// Collect results
		for i := 0; i < 10; i++ {
			err := <-errChan
			assert.NoError(t, err)
		}

		mockRepo.AssertExpectations(t)
	})
}
