package models

import (
	"time"
	"gorm.io/gorm"
)

type Product struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"not null"`
	Description    *string
	DesignImageURL string  `gorm:"not null"`
	BasePrice      float64 `gorm:"not null"`

	Material         string            `gorm:"type:varchar(100)"`
    Neckline         string            `gorm:"type:varchar(50)"`
    SleeveType       string            `gorm:"type:varchar(50)"`
    Fit              string            `gorm:"type:varchar(50)"`
    Brand            string            `gorm:"type:varchar(100)"`
    Category         string            `gorm:"type:varchar(100)"`
    CareInstructions string            `gorm:"type:text"`
    Weight           float64           `gorm:"type:decimal(10,2)"`
    Featured         bool              `gorm:"default:false"`
    
    IsActive         bool              `gorm:"default:true"`
    Variants         []ProductVariant  `gorm:"foreignKey:ProductID"`
    CreatedAt        time.Time
    UpdatedAt        time.Time
    DeletedAt        gorm.DeletedAt    `gorm:"index"`
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
