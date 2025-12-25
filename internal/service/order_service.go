package service

import (
    "github.com/vishnujoshi062/tshirt-ecommerce-api/internal/models"
    "gorm.io/gorm"
)

type OrderService struct {
    DB *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
    return &OrderService{DB: db}
}

func (s *OrderService) MarkOrderPaid(
    razorpayOrderID string,
    razorpayPaymentID string,
) error {
    return s.DB.Model(&models.Order{}).
        Where("razorpay_order_id = ?", razorpayOrderID).
        Updates(map[string]interface{}{
            "status":              "PAID",
            "razorpay_payment_id": razorpayPaymentID,
        }).Error
}
