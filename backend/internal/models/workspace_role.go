package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WorkspaceRole represents a role within a workspace with JSON permissions
type WorkspaceRole struct {
	ID          string          `gorm:"type:varchar(36);primaryKey" json:"id"`
	WorkspaceID string          `gorm:"type:varchar(36);not null;index" json:"workspace_id"`
	Name        string          `gorm:"type:varchar(100);not null" json:"name"`
	Description *string         `gorm:"type:text" json:"description,omitempty"`
	Permissions RolePermissions `gorm:"type:json;not null" json:"permissions"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Workspace   Workspace                 `gorm:"constraint:OnDelete:CASCADE" json:"workspace,omitempty"`
	Memberships []UserWorkspaceMembership `gorm:"foreignKey:RoleID" json:"memberships,omitempty"`
}

// RolePermissions represents the JSON structure for permissions
type RolePermissions struct {
	Permissions []string           `json:"permissions"`
	Metadata    PermissionMetadata `json:"metadata"`
}

// PermissionMetadata contains metadata about permissions
type PermissionMetadata struct {
	Version     string    `json:"version"`
	CreatedBy   string    `json:"created_by"`
	UpdatedBy   string    `json:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	Description string    `json:"description,omitempty"`
}

// GORM hooks
func (wr *WorkspaceRole) BeforeCreate(tx *gorm.DB) (err error) {
	if wr.ID == "" {
		wr.ID = uuid.New().String()
	}
	return
}
