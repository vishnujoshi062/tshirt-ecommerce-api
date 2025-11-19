package models

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User      User       `gorm:"foreignKey:UserID"`
	CartItems []CartItem `gorm:"foreignKey:CartID"`
}

type CartItem struct {
	ID        uint `gorm:"primaryKey"`
	CartID    uint `gorm:"not null"`
	VariantID uint `gorm:"not null"`
	Quantity  int  `gorm:"not null"`
	AddedAt   time.Time

	Cart    Cart           `gorm:"foreignKey:CartID"`
	Variant ProductVariant `gorm:"foreignKey:VariantID"`
}
