package repository

import (
	"strings"
	"gorm.io/gorm"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/models"
)

type PromoCodeRepository struct {
	db *gorm.DB
}

func NewPromoCodeRepository(db *gorm.DB) *PromoCodeRepository {
	return &PromoCodeRepository{db: db}
}

func normalizeCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func (r *PromoCodeRepository) FindAll(isActive *bool) ([]*models.PromoCode, error) {
	var promos []*models.PromoCode
	query := r.db.Order("created_at DESC")
	
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	
	if err := query.Find(&promos).Error; err != nil {
		return nil, err
	}
	return promos, nil
}

func (r *PromoCodeRepository) FindByCode(code string) (*models.PromoCode, error) {
	var promo models.PromoCode
	err := r.db.Where("code = ?", normalizeCode(code)).First(&promo).Error
	if err != nil {
		return nil, err
	}
	return &promo, nil
}

func (r *PromoCodeRepository) Create(promo *models.PromoCode) error {
	promo.Code = normalizeCode(promo.Code)
	return r.db.Create(promo).Error
}

func (r *PromoCodeRepository) Update(id string, updates map[string]interface{}) (*models.PromoCode, error) {
	var promo models.PromoCode
	if err := r.db.Where("id = ?", id).First(&promo).Error; err != nil {
		return nil, err
	}
	
	if code, ok := updates["code"].(string); ok {
		updates["code"] = normalizeCode(code)
	}
	
	if err := r.db.Model(&promo).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &promo, nil
}

func (r *PromoCodeRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.PromoCode{}).Error
}

func (r *PromoCodeRepository) IncrementUsage(code string) error {
	return r.db.Model(&models.PromoCode{}).
		Where("code = ?", normalizeCode(code)).
		Update("usage_count", gorm.Expr("usage_count + ?", 1)).Error
}

func (r *PromoCodeRepository) IncrementUsageWithLimit(tx *gorm.DB, code string, limit int) (int64, error) {
	result := tx.Model(&models.PromoCode{}).
		Where("code = ? AND usage_count < ?", normalizeCode(code), limit).
		Update("usage_count", gorm.Expr("usage_count + ?", 1))
	return result.RowsAffected, result.Error
}
