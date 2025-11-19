package repository

import (
	"tshirt-ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type CartRepository struct {
	DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{DB: db}
}

func (r *CartRepository) GetCartByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.DB.Preload("Items").Where("user_id = ?", userID).First(&cart).Error
	return &cart, err
}

func (r *CartRepository) AddItem(item *models.CartItem) error {
	return r.DB.Create(item).Error
}

func (r *CartRepository) UpdateItem(item *models.CartItem) error {
	return r.DB.Save(item).Error
}

func (r *CartRepository) RemoveItem(itemID uint) error {
	return r.DB.Delete(&models.CartItem{}, itemID).Error
}

func (r *CartRepository) ClearCart(cartID uint) error {
	return r.DB.Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}
