package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Workspace represents a multi-tenant workspace
type Workspace struct {
	ID          string                 `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name        string                 `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string                 `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	Description *string                `gorm:"type:text" json:"description,omitempty"`
	AvatarURL   *string                `gorm:"type:varchar(500)" json:"avatar_url,omitempty"`
	OwnerID     string                 `gorm:"type:varchar(36);not null;index" json:"owner_id"`
	Settings    map[string]interface{} `gorm:"type:json" json:"settings,omitempty"`
	Status      string                 `gorm:"type:varchar(50);not null;default:'active';index" json:"status"`
	CreatedAt   time.Time              `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time              `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Owner       User                      `gorm:"constraint:OnDelete:RESTRICT" json:"owner,omitempty"`
	Roles       []WorkspaceRole           `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"roles,omitempty"`
	Memberships []UserWorkspaceMembership `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"memberships,omitempty"`
}

// GORM hooks
func (w *Workspace) BeforeCreate(tx *gorm.DB) (err error) {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	return
}

// Constants for workspace status
const (
	WorkspaceStatusActive    = "active"
	WorkspaceStatusInactive  = "inactive"
	WorkspaceStatusSuspended = "suspended"
	WorkspaceStatusDeleted   = "deleted"
)
