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

// ✅ Create a new order
func (r *OrderRepository) CreateOrder(order *models.Order) error {
	return r.DB.Create(order).Error
}

// ✅ Fetch single order by ID (admin and user)
func (r *OrderRepository) GetOrderByID(id uint) (*models.Order, error) {
	var order models.Order

	err := r.DB.
		Preload("OrderItems").
		Preload("OrderItems.Variant").
		Preload("OrderItems.Variant.Product").
		Preload("Payment").
		Preload("User").
		First(&order, id).Error

	return &order, err
}

// ✅ Fetch all orders (ADMIN)
func (r *OrderRepository) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order

	err := r.DB.
		Preload("OrderItems").
		Preload("OrderItems.Variant").
		Preload("OrderItems.Variant.Product").
		Preload("Payment").
		Preload("User").
		Find(&orders).Error

	return orders, err
}

// ✅ Fetch orders for a specific user
func (r *OrderRepository) GetOrdersByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order

	err := r.DB.
		Preload("OrderItems").
		Preload("OrderItems.Variant").
		Preload("OrderItems.Variant.Product").
		Where("user_id = ?", userID).
		Find(&orders).Error

	return orders, err
}

// ✅ Update order (status, etc.)
func (r *OrderRepository) UpdateOrder(order *models.Order) error {
	return r.DB.Save(order).Error
}
