package models

import (
	"time"
)

type DiscountType string

const (
	DiscountTypePercentage DiscountType = "percentage"
	DiscountTypeFixed      DiscountType = "fixed"
)

type PromoCode struct {
	ID            string       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Code          string       `gorm:"uniqueIndex;not null" json:"code"`
	DiscountType  DiscountType `gorm:"type:varchar(20);not null" json:"discountType"`
	DiscountValue float64      `gorm:"type:decimal(10,2);not null" json:"discountValue"`
	ValidFrom     *time.Time   `gorm:"type:timestamp;null" json:"validFrom"`
	ValidUntil    *time.Time   `gorm:"type:timestamp;null" json:"validUntil"`
	IsActive      bool         `gorm:"default:true" json:"isActive"`
	UsageLimit    *int         `json:"usageLimit"`
	UsageCount    int          `gorm:"default:0" json:"usageCount"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`
}

func (PromoCode) TableName() string {
	return "promo_codes"
}
