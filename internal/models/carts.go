package models

import (
    "time"

    "gorm.io/gorm"
)

type Cart struct {
    ID        uint `gorm:"primaryKey"`
    UserID    *string `gorm:"index"` // Make it nullable with *string for guest carts
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`

    // No User relationship - UserID stores Clerk ID directly
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