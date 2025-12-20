package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint   `gorm:"primaryKey"`
	ClerkUserID   string `gorm:"uniqueIndex;not null"` // Clerk user ID
	Email         string `gorm:"uniqueIndex;not null"`
	PasswordHash  string
	Name          string
	Phone         string
	PhoneVerified bool   `gorm:"default:false"`
	Address       string
	Role          string `gorm:"default:'user'"` // user, admin
	OAuthProvider string `gorm:"default:'clerk'"`
	OAuthID       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	// Note: Cart.UserID and Order.UserID store Clerk IDs (strings), not database user IDs
	// Therefore, we cannot use foreign key relationships here
}
