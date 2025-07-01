package services

import (
	"fmt"
	"go-backend-v2/internal/models"
	"go-backend-v2/internal/repo"
)

type UserServiceInterface interface {
	GetUserWithWorkspaces(userID string) (*models.User, error)
	GetUserProfile(userID string) (*models.User, error)
}

type UserService struct {
	userRepo repo.UserRepositoryInterface
}

func NewUserService(userRepo repo.UserRepositoryInterface) UserServiceInterface {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserWithWorkspaces(userID string) (*models.User, error) {
	user, err := s.userRepo.GetUserWithWorkspaces(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user with workspaces: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *UserService) GetUserProfile(userID string) (*models.User, error) {
	user, err := s.userRepo.GetUserWithProfile(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}
