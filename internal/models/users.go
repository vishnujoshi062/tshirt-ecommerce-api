package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint   `gorm:"primaryKey"`
	Email         string `gorm:"uniqueIndex;not null"`
	PasswordHash  string `gorm:"not null"`
	Name          string `gorm:"not null"`
	Phone         string
	Address       string
	Role          string `gorm:"default:'user'"`  // user, admin
	OAuthProvider string `gorm:"default:'local'"` // local, google
	OAuthID       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	// Note: Cart.UserID and Order.UserID store Clerk IDs (strings), not database user IDs
	// Therefore, we cannot use foreign key relationships here
}
