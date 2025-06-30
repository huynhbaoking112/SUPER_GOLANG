package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Resource represents a system resource that can have permissions
type Resource struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	DisplayName string    `gorm:"type:varchar(150);not null" json:"display_name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	Category    *string   `gorm:"type:varchar(50);index" json:"category,omitempty"`
	Icon        *string   `gorm:"type:varchar(100)" json:"icon,omitempty"`
	Status      string    `gorm:"type:varchar(50);not null;default:'active';index" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// GORM hooks
func (r *Resource) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return
}

// Constants for resource status
const (
	ResourceStatusActive     = "active"
	ResourceStatusInactive   = "inactive"
	ResourceStatusDeprecated = "deprecated"
)

// Constants for resource categories
const (
	ResourceCategoryCore        = "core"
	ResourceCategoryFeature     = "feature"
	ResourceCategoryAdmin       = "admin"
	ResourceCategoryIntegration = "integration"
)

// GetValidActions returns valid actions for any resource
func GetValidActions() []string {
	return []string{"read", "create", "update", "delete", "all"}
}
