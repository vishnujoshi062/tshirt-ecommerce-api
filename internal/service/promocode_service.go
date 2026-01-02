package service

import (
	"context"
	"time"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/graph/model"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/models"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/repository"
)

type PromoService struct {
	repo *repository.PromoCodeRepository
}

func NewPromoService(repo *repository.PromoCodeRepository) *PromoService {
	return &PromoService{repo: repo}
}

type ValidationResult struct {
	IsValid        bool
	DiscountAmount float64
	Message        string
}

func (s *PromoService) ValidatePromoCode(code string, orderAmount float64) (*ValidationResult, error) {
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

	// Check validUntil (only if set)
	if promo.ValidUntil != nil && now.After(*promo.ValidUntil) {
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


func (s *PromoService) CreatePromoCode(ctx context.Context, input model.PromoCodeInput) (*models.PromoCode, error) {
	promo := &models.PromoCode{
		Code:           input.Code,
		DiscountType:   input.DiscountType,
		DiscountValue:  input.DiscountValue,
		IsActive:       input.IsActive != nil && *input.IsActive,
		UsageLimit:     input.UsageLimit,
		UsageCount:     0,
	}

	if input.ValidFrom != nil {
		validFrom, _ := time.Parse(time.RFC3339, *input.ValidFrom)
		promo.ValidFrom = &validFrom
	}

	if input.ValidUntil != nil {
		validUntil, _ := time.Parse(time.RFC3339, *input.ValidUntil)
		promo.ValidUntil = &validUntil
	}

	if err := s.repo.Create(promo); err != nil {
		return nil, err
	}
	return promo, nil
}

func (s *PromoService) UpdatePromoCode(ctx context.Context, id string, input model.PromoCodeInput) (*models.PromoCode, error) {
	updates := map[string]interface{}{
		"code":            input.Code,
		"discount_type":   input.DiscountType,
		"discount_value":  input.DiscountValue,
		"is_active":       input.IsActive != nil && *input.IsActive,
		"usage_limit":     input.UsageLimit,
	}

	if input.ValidFrom != nil {
		validFrom, _ := time.Parse(time.RFC3339, *input.ValidFrom)
		updates["valid_from"] = validFrom
	}

	if input.ValidUntil != nil {
		validUntil, _ := time.Parse(time.RFC3339, *input.ValidUntil)
		updates["valid_until"] = validUntil
	}

	return s.repo.Update(id, updates)
}

func (s *PromoService) DeletePromoCode(ctx context.Context, id string) (bool, error) {
	err := s.repo.Delete(id)
	return err == nil, err
}

func (s *PromoService) ToggleStatus(ctx context.Context, id string) (*models.PromoCode, error) {
	promo, err := s.repo.FindByCode(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"is_active": !promo.IsActive,
	}

	return s.repo.Update(id, updates)
}

func (s *PromoService) GetAllPromoCodes(ctx context.Context) ([]*models.PromoCode, error) {
	return s.repo.FindAll(nil)
}

func (s *PromoService) GetPromoCodeByID(ctx context.Context, code string) (*models.PromoCode, error) {
	return s.repo.FindByCode(code)
}
