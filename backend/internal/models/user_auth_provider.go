package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserAuthProvider represents authentication provider information
type UserAuthProvider struct {
	ID             string                 `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID         string                 `gorm:"type:varchar(36);not null;index:idx_user_provider" json:"user_id"`
	Provider       string                 `gorm:"type:varchar(50);not null;index:idx_user_provider" json:"provider"`
	ProviderUserID string                 `gorm:"type:varchar(255);not null" json:"provider_user_id"`
	ProviderEmail  *string                `gorm:"type:varchar(255);index" json:"provider_email,omitempty"`
	ProviderData   map[string]interface{} `gorm:"type:json" json:"provider_data,omitempty"`
	PasswordHash   *string                `gorm:"type:varchar(255)" json:"password_hash,omitempty"`
	IsPrimary      bool                   `gorm:"not null;default:false" json:"is_primary"`
	CreatedAt      time.Time              `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time              `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User User `gorm:"constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// GORM hooks
func (uap *UserAuthProvider) BeforeCreate(tx *gorm.DB) (err error) {
	if uap.ID == "" {
		uap.ID = uuid.New().String()
	}
	return
}

// Constants for auth providers
const (
	AuthProviderLocal     = "local"
	AuthProviderGoogle    = "google"
	AuthProviderFacebook  = "facebook"
	AuthProviderGithub    = "github"
	AuthProviderApple     = "apple"
	AuthProviderMicrosoft = "microsoft"
	AuthProviderLinkedin  = "linkedin"
	AuthProviderTwitter   = "twitter"
)
