package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProviderData map[string]interface{}

func (pd ProviderData) Value() (driver.Value, error) {
	if pd == nil {
		return nil, nil
	}
	return json.Marshal(pd)
}

func (pd *ProviderData) Scan(value interface{}) error {
	if value == nil {
		*pd = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into ProviderData", value)
	}

	return json.Unmarshal(bytes, pd)
}

type UserAuthProvider struct {
	ID             string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID         string         `gorm:"type:varchar(36);not null;index:idx_user_provider" json:"user_id"`
	Provider       string         `gorm:"type:varchar(50);not null;index:idx_user_provider" json:"provider"`
	ProviderUserID string         `gorm:"type:varchar(255);not null" json:"provider_user_id"`
	ProviderEmail  *string        `gorm:"type:varchar(255);index" json:"provider_email,omitempty"`
	ProviderData   ProviderData   `gorm:"type:json" json:"provider_data,omitempty"`
	PasswordHash   *string        `gorm:"type:varchar(255)" json:"password_hash,omitempty"`
	IsPrimary      bool           `gorm:"not null;default:false" json:"is_primary"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
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

// Import constants from common package
// Auth provider constants are now in internal/common/auth_constants.go
