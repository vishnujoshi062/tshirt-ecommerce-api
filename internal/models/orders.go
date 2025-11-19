package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID              uint    `gorm:"primaryKey"`
	UserID          uint    `gorm:"not null"`
	TotalAmount     float64 `gorm:"not null"`
	Status          string  `gorm:"not null"` // pending, confirmed, shipped, delivered, cancelled
	ShippingAddress string  `gorm:"not null"`
	PaymentID       uint
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	User       User        `gorm:"foreignKey:UserID"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
	Payment    Payment     `gorm:"foreignKey:PaymentID"`
}

type OrderItem struct {
	ID        uint    `gorm:"primaryKey"`
	OrderID   uint    `gorm:"not null"`
	VariantID uint    `gorm:"not null"`
	Quantity  int     `gorm:"not null"`
	UnitPrice float64 `gorm:"not null"`
	Subtotal  float64 `gorm:"not null"`

	Order   Order          `gorm:"foreignKey:OrderID"`
	Variant ProductVariant `gorm:"foreignKey:VariantID"`
}

type Payment struct {
	ID                uint    `gorm:"primaryKey"`
	OrderID           uint    `gorm:"uniqueIndex;not null"`
	Amount            float64 `gorm:"not null"`
	Status            string  `gorm:"not null"` // pending, success, failed
	PaymentMethod     string  `gorm:"not null"` // razorpay, cod
	TransactionID     string
	RazorpayOrderID   string
	RazorpaySignature string
	CreatedAt         time.Time

	Order Order `gorm:"foreignKey:OrderID"`
}
