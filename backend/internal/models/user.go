package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the main user entity
type User struct {
	ID              string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	Email           string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Phone           *string    `gorm:"type:varchar(20);uniqueIndex" json:"phone,omitempty"`
	EmailVerifiedAt *time.Time `gorm:"type:timestamp" json:"email_verified_at,omitempty"`
	PhoneVerifiedAt *time.Time `gorm:"type:timestamp" json:"phone_verified_at,omitempty"`
	GlobalRole      string     `gorm:"type:varchar(50);not null;default:'customer';index" json:"global_role"`
	Status          string     `gorm:"type:varchar(50);not null;default:'pending';index" json:"status"`
	LastLoginAt     *time.Time `gorm:"type:timestamp" json:"last_login_at,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Profile              *UserProfile              `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"profile,omitempty"`
	AuthProviders        []UserAuthProvider        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"auth_providers,omitempty"`
	WorkspaceMemberships []UserWorkspaceMembership `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"workspace_memberships,omitempty"`
	OwnedWorkspaces      []Workspace               `gorm:"foreignKey:OwnerID" json:"owned_workspaces,omitempty"`
	InvitedMemberships   []UserWorkspaceMembership `gorm:"foreignKey:InvitedBy" json:"invited_memberships,omitempty"`
}

// UserProfile represents extended user information
type UserProfile struct {
	UserID       string     `gorm:"type:varchar(36);primaryKey" json:"user_id"`
	FirstName    string     `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName     string     `gorm:"type:varchar(100);not null" json:"last_name"`
	DisplayName  *string    `gorm:"type:varchar(200);index" json:"display_name,omitempty"`
	AvatarURL    *string    `gorm:"type:varchar(500)" json:"avatar_url,omitempty"`
	DateOfBirth  *time.Time `gorm:"type:date" json:"date_of_birth,omitempty"`
	AddressLine1 *string    `gorm:"type:varchar(255)" json:"address_line1,omitempty"`
	AddressLine2 *string    `gorm:"type:varchar(255)" json:"address_line2,omitempty"`
	City         *string    `gorm:"type:varchar(100)" json:"city,omitempty"`
	State        *string    `gorm:"type:varchar(100)" json:"state,omitempty"`
	PostalCode   *string    `gorm:"type:varchar(20)" json:"postal_code,omitempty"`
	Country      *string    `gorm:"type:varchar(100);index" json:"country,omitempty"`
	Timezone     string     `gorm:"type:varchar(50);not null;default:'UTC'" json:"timezone"`
	Locale       string     `gorm:"type:varchar(10);not null;default:'en'" json:"locale"`
	Bio          *string    `gorm:"type:text" json:"bio,omitempty"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User User `gorm:"constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// GORM hooks
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}
