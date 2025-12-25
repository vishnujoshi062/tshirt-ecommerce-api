package models

import (
	"time"
)

type Order struct {
	ID              uint     `gorm:"primaryKey;autoIncrement"`
	UserID          string   `gorm:"not null;type:varchar(255)"`
	TotalAmount     float64  `gorm:"not null"`
	Discount        float64  `gorm:"default:0"`
	PromoCode       *string  `gorm:"type:varchar(50)"`
	Status          string   `gorm:"not null"`
	ShippingAddress string   `gorm:"not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
	Payment    *Payment    `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID        uint    `gorm:"primaryKey;autoIncrement"`
	OrderID   uint    `gorm:"not null"`
	VariantID uint    `gorm:"not null"`
	Quantity  int     `gorm:"not null"`
	UnitPrice float64 `gorm:"not null"`
	Subtotal  float64 `gorm:"not null"`

	Variant ProductVariant `gorm:"foreignKey:VariantID"`
}

type Payment struct {
	ID            uint    `gorm:"primaryKey;autoIncrement"`
	OrderID       uint    `gorm:"uniqueIndex"`
	Amount        float64 `gorm:"not null"`
	Status        string  `gorm:"not null"`
	PaymentMethod string  `gorm:"not null"`
	TransactionID string
	CreatedAt     time.Time

	Order *Order `gorm:"foreignKey:OrderID"`
}
