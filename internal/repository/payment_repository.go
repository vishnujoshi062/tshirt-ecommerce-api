package repository

import (
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	DB *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{DB: db}
}

func (r *PaymentRepository) CreatePayment(payment *models.Payment) error {
	return r.DB.Create(payment).Error
}

func (r *PaymentRepository) GetPaymentByOrderID(orderID uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.DB.Where("order_id = ?", orderID).First(&payment).Error
	return &payment, err
}

func (r *PaymentRepository) UpdatePayment(payment *models.Payment) error {
	return r.DB.Save(payment).Error
}
