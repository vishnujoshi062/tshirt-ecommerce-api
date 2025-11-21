package repository

import (
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) CreateOrder(order *models.Order) error {
	return r.DB.Create(order).Error
}

func (r *OrderRepository) GetOrderByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.DB.Preload("OrderItems").First(&order, id).Error
	return &order, err
}

func (r *OrderRepository) GetOrdersByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.DB.Preload("OrderItems").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) UpdateOrder(order *models.Order) error {
	return r.DB.Save(order).Error
}
