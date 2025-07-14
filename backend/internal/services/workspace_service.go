package services

import (
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

type WorkspaceService struct {
	workspaceRepo repo.WorkspaceRepositoryInterface
	userRepo      repo.UserRepositoryInterface
}

func NewWorkspaceService(workspaceRepo repo.WorkspaceRepositoryInterface, userRepo repo.UserRepositoryInterface) WorkspaceServiceInterface {
	return &WorkspaceService{
		workspaceRepo: workspaceRepo,
		userRepo:      userRepo,
	}
}

func (s *WorkspaceService) CreateWorkspace(userID string, req *dto.CreateWorkspaceRequest) (*dto.WorkspaceResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, common.ErrUserNotFound
	}
	if user.Status != common.UserStatusActive {
		return nil, common.ErrUserInactive
	}
	if user.GlobalRole != common.GlobalRoleSuperAdmin {
		return nil, common.ErrWorkspaceCreateForbidden
	}

	uniqueSlug, err := utils.GenerateUniqueSlug(req.Name, s.workspaceRepo.ExistsBySlug)
	if err != nil {
		return nil, common.ErrWorkspaceSlugGeneration
	}

	var result *dto.WorkspaceResponse

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		workspace := &models.Workspace{
			Name:        req.Name,
			Slug:        uniqueSlug,
			Description: &req.Description,
			OwnerID:     userID,
			Status:      common.ActiveStatus,
		}

		if err := s.workspaceRepo.CreateWorkspace(tx, workspace); err != nil {
			return fmt.Errorf("failed to create workspace: %w", err)
		}

		adminRole := &models.WorkspaceRole{
			WorkspaceID: workspace.ID,
			Name:        common.WorkspaceRoleAdmin,
			Description: stringPtr("Full administrative access to workspace"),
			Permissions: models.RolePermissions{
				Permissions: []string{"all"},
				Metadata: models.PermissionMetadata{
					Version:     "1.0",
					CreatedBy:   "system",
					UpdatedBy:   "system",
					UpdatedAt:   time.Now(),
					Description: "Default admin role with full permissions",
				},
			},
			Status: common.ActiveStatus,
		}

		if err := s.workspaceRepo.CreateWorkspaceRole(tx, adminRole); err != nil {
			return fmt.Errorf("failed to create admin role: %w", err)
		}

		membership := &models.UserWorkspaceMembership{
			UserID:      userID,
			WorkspaceID: workspace.ID,
			RoleID:      adminRole.ID,
			Status:      common.ActiveStatus,
			JoinedAt:    &time.Time{},
		}
		now := time.Now()
		membership.JoinedAt = &now

		if err := s.workspaceRepo.CreateMembership(tx, membership); err != nil {
			return fmt.Errorf("failed to create membership: %w", err)
		}

		result = &dto.WorkspaceResponse{
			ID:          workspace.ID,
			Name:        workspace.Name,
			Slug:        workspace.Slug,
			Description: getStringValue(workspace.Description),
			OwnerID:     workspace.OwnerID,
			Status:      workspace.Status,
			CreatedAt:   workspace.CreatedAt,
			UpdatedAt:   workspace.UpdatedAt,
		}

		return nil
	})

	if err != nil {
		return nil, common.ErrWorkspaceCreateFailed
	}

	if global.EventTopicPublisher != nil {

		payload := &dto.WorkspaceCreatedPayload{
			WorkspaceID:   result.ID,
			WorkspaceName: result.Name,
			WorkspaceSlug: result.Slug,
			CreatedByID:   userID,
		}

		go func() {
			if err := global.EventTopicPublisher.Publish(common.WorkspaceCreatedLog, payload); err != nil {
				fmt.Printf("Error publishing workspace creation event: %v\n", err)
			}
		}()
	}

	return result, nil
}

func stringPtr(s string) *string {
	return &s
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
