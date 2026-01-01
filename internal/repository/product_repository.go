package repository

import (
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) CreateProduct(product *models.Product) error {
	return r.DB.Create(product).Error
}

func (r *ProductRepository) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.DB.Preload("Variants", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Inventory")
	}).First(&product, id).Error
	return &product, err
}

func (r *ProductRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	err := r.DB.Preload("Variants", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Inventory")
	}).Find(&products).Error
	return products, err
}

func (r *ProductRepository) UpdateProduct(product *models.Product) error {
	return r.DB.Save(product).Error
}

func (r *ProductRepository) DeleteProduct(id uint) error {
	return r.DB.Delete(&models.Product{}, id).Error
}
