package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserWorkspaceMembership represents user membership in a workspace
type UserWorkspaceMembership struct {
	ID          string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID      string     `gorm:"type:varchar(36);not null;index" json:"user_id"`
	WorkspaceID string     `gorm:"type:varchar(36);not null;index" json:"workspace_id"`
	RoleID      string     `gorm:"type:varchar(36);not null" json:"role_id"`
	Status      string     `gorm:"type:varchar(50);not null;default:'pending';index" json:"status"`
	InvitedBy   *string    `gorm:"type:varchar(36)" json:"invited_by,omitempty"`
	JoinedAt    *time.Time `gorm:"type:timestamp" json:"joined_at,omitempty"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User      User          `gorm:"constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Workspace Workspace     `gorm:"constraint:OnDelete:CASCADE" json:"workspace,omitempty"`
	Role      WorkspaceRole `gorm:"constraint:OnDelete:RESTRICT" json:"role,omitempty"`
	Inviter   *User         `gorm:"foreignKey:InvitedBy" json:"inviter,omitempty"`
}

// GORM hooks
func (uwm *UserWorkspaceMembership) BeforeCreate(tx *gorm.DB) (err error) {
	if uwm.ID == "" {
		uwm.ID = uuid.New().String()
	}
	return
}

// Constants for membership status
const (
	MembershipStatusActive    = "active"
	MembershipStatusInactive  = "inactive"
	MembershipStatusPending   = "pending"
	MembershipStatusRejected  = "rejected"
	MembershipStatusSuspended = "suspended"
)
