package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	ClerkUserID    string          `gorm:"uniqueIndex;not null" json:"clerk_user_id"`
	Email          string          `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash   string          `json:"-"`
	Name           string          `json:"name"`
	Phone          string          `json:"phone_number"`
	PhoneVerified  bool            `gorm:"default:false" json:"phone_verified"`
	PhoneUpdatedAt sql.NullTime    `json:"phone_updated_at,omitempty"`
	Address        string          `json:"address"`
	Role           string          `gorm:"default:'user'" json:"role"`
	OAuthProvider  string          `gorm:"default:'clerk'" json:"oauth_provider"`
	OAuthID        string          `json:"oauth_id"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      sql.NullTime    `gorm:"index" json:"-"`
}
