package service

import (
	"time"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/models"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/repository"
)

type PromoCodeService struct {
	repo *repository.PromoCodeRepository
}

func NewPromoCodeService(repo *repository.PromoCodeRepository) *PromoCodeService {
	return &PromoCodeService{repo: repo}
}

type ValidationResult struct {
	IsValid        bool
	DiscountAmount float64
	Message        string
}

func (s *PromoCodeService) ValidatePromoCode(code string, orderAmount float64) (*ValidationResult, error) {
	promo, err := s.repo.FindByCode(code)
	if err != nil {
		return &ValidationResult{
			IsValid: false,
			DiscountAmount: 0,
			Message: "Code not found",
		}, nil
	}

	now := time.Now()

	// Check if active
	if !promo.IsActive {
		return &ValidationResult{
			IsValid: false,
			DiscountAmount: 0,
			Message: "Code is inactive",
		}, nil
	}

	// Check validFrom
	if promo.ValidFrom != nil && now.Before(*promo.ValidFrom) {
		return &ValidationResult{
			IsValid: false,
			DiscountAmount: 0,
			Message: "Code not yet valid",
		}, nil
	}

	// Check validUntil
	if now.After(promo.ValidUntil) {
		return &ValidationResult{
			IsValid: false,
			DiscountAmount: 0,
			Message: "Code has expired",
		}, nil
	}

	// Check usage limit
	if promo.UsageLimit != nil && promo.UsageCount >= *promo.UsageLimit {
		return &ValidationResult{
			IsValid: false,
			DiscountAmount: 0,
			Message: "Usage limit reached",
		}, nil
	}

	// Calculate discount
	var discount float64
	if promo.DiscountType == models.DiscountTypePercentage {
		discount = (orderAmount * promo.DiscountValue) / 100.0
	} else {
		discount = promo.DiscountValue
	}

	// Cap discount at order amount
	if discount > orderAmount {
		discount = orderAmount
	}

	return &ValidationResult{
		IsValid: true,
		DiscountAmount: discount,
		Message: "Code applied successfully",
	}, nil
}
