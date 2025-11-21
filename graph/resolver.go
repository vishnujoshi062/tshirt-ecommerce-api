package graph

import (
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/repository"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/service"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB                *gorm.DB
	UserRepository    *repository.UserRepository
	ProductRepository *repository.ProductRepository
	CartRepository    *repository.CartRepository
	OrderRepository   *repository.OrderRepository
	PaymentRepository *repository.PaymentRepository
	PaymentService    *service.PaymentService
}
