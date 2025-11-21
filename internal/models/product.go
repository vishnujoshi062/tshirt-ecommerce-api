package models

import (
	"time"
)

type Product struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"not null"`
	Description    *string
	DesignImageURL string  `gorm:"not null"`
	BasePrice      float64 `gorm:"not null"`
	IsActive       bool    `gorm:"default:true"`
	CreatedAt      time.Time

	Variants []ProductVariant `gorm:"foreignKey:ProductID"`
}

type ProductVariant struct {
	ID            uint   `gorm:"primaryKey"`
	ProductID     uint   `gorm:"not null"`
	Size          string `gorm:"not null"`
	Color         *string
	PriceModifier float64 `gorm:"not null"`
	SKU           string  `gorm:"uniqueIndex;not null"`
	CreatedAt     time.Time

	Product   *Product   `gorm:"foreignKey:ProductID"`
	Inventory *Inventory `gorm:"foreignKey:VariantID"`
}

type Inventory struct {
	ID               uint `gorm:"primaryKey"`
	VariantID        uint `gorm:"uniqueIndex;not null"`
	StockQuantity    int  `gorm:"not null;default:0"`
	ReservedQuantity int  `gorm:"not null;default:0"`

	UpdatedAt time.Time

	Variant *ProductVariant `gorm:"foreignKey:VariantID"`
}
