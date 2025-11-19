package graph

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/vishnujoshi062/tshirt-ecommerce-api/graph/model"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/middleware"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/models"
)

// Products is the resolver for the products field.
func (r *queryResolver) Products(ctx context.Context, isActive *bool) ([]*model.Product, error) {
	var products []models.Product
	query := r.DB.Preload("Variants.Inventory")

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	var result []*model.Product
	for _, p := range products {
		product := &model.Product{
			ID:             strconv.Itoa(int(p.ID)),
			Name:           p.Name,
			Description:    &p.Description,
			DesignImageURL: p.DesignImageURL,
			BasePrice:      p.BasePrice,
			IsActive:       p.IsActive,
			CreatedAt:      p.CreatedAt.Format(time.RFC3339),
		}

		var variants []*model.ProductVariant
		for _, v := range p.Variants {
			availableQty := v.Inventory.StockQuantity - v.Inventory.ReservedQuantity
			variant := &model.ProductVariant{
				ID:            strconv.Itoa(int(v.ID)),
				ProductID:     strconv.Itoa(int(v.ProductID)),
				Size:          v.Size,
				Color:         &v.Color,
				PriceModifier: v.PriceModifier,
				SKU:           v.SKU,
				Price:         p.BasePrice + v.PriceModifier,
				Inventory: &model.Inventory{
					ID:                strconv.Itoa(int(v.Inventory.ID)),
					VariantID:         strconv.Itoa(int(v.ID)),
					StockQuantity:     v.Inventory.StockQuantity,
					ReservedQuantity:  v.Inventory.ReservedQuantity,
					AvailableQuantity: availableQty,
				},
			}
			variants = append(variants, variant)
		}
		product.Variants = variants
		result = append(result, product)
	}

	return result, nil
}

// Product is the resolver for the product field.
func (r *queryResolver) Product(ctx context.Context, id string) (*model.Product, error) {
	productID, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	var product models.Product
	if err := r.DB.Preload("Variants.Inventory").First(&product, productID).Error; err != nil {
		return nil, errors.New("product not found")
	}

	result := &model.Product{
		ID:             strconv.Itoa(int(product.ID)),
		Name:           product.Name,
		Description:    &product.Description,
		DesignImageURL: product.DesignImageURL,
		BasePrice:      product.BasePrice,
		IsActive:       product.IsActive,
		CreatedAt:      product.CreatedAt.Format(time.RFC3339),
	}

	var variants []*model.ProductVariant
	for _, v := range product.Variants {
		availableQty := v.Inventory.StockQuantity - v.Inventory.ReservedQuantity
		variant := &model.ProductVariant{
			ID:            strconv.Itoa(int(v.ID)),
			ProductID:     strconv.Itoa(int(v.ProductID)),
			Size:          v.Size,
			Color:         &v.Color,
			PriceModifier: v.PriceModifier,
			SKU:           v.SKU,
			Price:         product.BasePrice + v.PriceModifier,
			Inventory: &model.Inventory{
				ID:                strconv.Itoa(int(v.Inventory.ID)),
				VariantID:         strconv.Itoa(int(v.ID)),
				StockQuantity:     v.Inventory.StockQuantity,
				ReservedQuantity:  v.Inventory.ReservedQuantity,
				AvailableQuantity: availableQty,
			},
		}
		variants = append(variants, variant)
	}
	result.Variants = variants

	return result, nil
}

// ProductsByCategory is the resolver for the productsByCategory field.
func (r *queryResolver) ProductsByCategory(ctx context.Context, category string) ([]*model.Product, error) {
	// Implement category filtering if you add a category field to Product model
	return nil, errors.New("not implemented yet")
}

// CreateProduct is the resolver for the createProduct field.
func (r *mutationResolver) CreateProduct(ctx context.Context, input model.ProductInput) (*model.Product, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil || claims.Role != "admin" {
		return nil, errors.New("unauthorized: admin access required")
	}

	product := &models.Product{
		Name:           input.Name,
		DesignImageURL: input.DesignImageURL,
		BasePrice:      input.BasePrice,
		IsActive:       true,
	}

	if input.Description != nil {
		product.Description = *input.Description
	}

	if err := r.DB.Create(product).Error; err != nil {
		return nil, err
	}

	return &model.Product{
		ID:             strconv.Itoa(int(product.ID)),
		Name:           product.Name,
		Description:    &product.Description,
		DesignImageURL: product.DesignImageURL,
		BasePrice:      product.BasePrice,
		IsActive:       product.IsActive,
		CreatedAt:      product.CreatedAt.Format(time.RFC3339),
		Variants:       []*model.ProductVariant{},
	}, nil
}

// UpdateProduct is the resolver for the updateProduct field.
func (r *mutationResolver) UpdateProduct(ctx context.Context, id string, input model.ProductInput) (*model.Product, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil || claims.Role != "admin" {
		return nil, errors.New("unauthorized: admin access required")
	}

	productID, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	var product models.Product
	if err := r.DB.First(&product, productID).Error; err != nil {
		return nil, errors.New("product not found")
	}

	product.Name = input.Name
	if input.Description != nil {
		product.Description = *input.Description
	}
	product.DesignImageURL = input.DesignImageURL
	product.BasePrice = input.BasePrice

	if err := r.DB.Save(&product).Error; err != nil {
		return nil, err
	}

	return &model.Product{
		ID:             strconv.Itoa(int(product.ID)),
		Name:           product.Name,
		Description:    &product.Description,
		DesignImageURL: product.DesignImageURL,
		BasePrice:      product.BasePrice,
		IsActive:       product.IsActive,
		CreatedAt:      product.CreatedAt.Format(time.RFC3339),
	}, nil
}

// DeleteProduct is the resolver for the deleteProduct field.
func (r *mutationResolver) DeleteProduct(ctx context.Context, id string) (bool, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil || claims.Role != "admin" {
		return false, errors.New("unauthorized: admin access required")
	}

	productID, err := strconv.Atoi(id)
	if err != nil {
		return false, errors.New("invalid product ID")
	}

	if err := r.DB.Delete(&models.Product{}, productID).Error; err != nil {
		return false, err
	}

	return true, nil
}

// CreateProductVariant is the resolver for the createProductVariant field.
func (r *mutationResolver) CreateProductVariant(ctx context.Context, input model.ProductVariantInput) (*model.ProductVariant, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil || claims.Role != "admin" {
		return nil, errors.New("unauthorized: admin access required")
	}

	productID, err := strconv.Atoi(input.ProductID)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	variant := &models.ProductVariant{
		ProductID:     uint(productID),
		Size:          input.Size,
		PriceModifier: input.PriceModifier,
		SKU:           input.SKU,
	}

	if input.Color != nil {
		variant.Color = *input.Color
	}

	if err := r.DB.Create(variant).Error; err != nil {
		return nil, err
	}

	// Create inventory record
	inventory := &models.Inventory{
		VariantID:        variant.ID,
		StockQuantity:    input.StockQuantity,
		ReservedQuantity: 0,
	}

	if err := r.DB.Create(inventory).Error; err != nil {
		return nil, err
	}

	return &model.ProductVariant{
		ID:            strconv.Itoa(int(variant.ID)),
		ProductID:     input.ProductID,
		Size:          variant.Size,
		Color:         &variant.Color,
		PriceModifier: variant.PriceModifier,
		SKU:           variant.SKU,
		Inventory: &model.Inventory{
			ID:                strconv.Itoa(int(inventory.ID)),
			VariantID:         strconv.Itoa(int(variant.ID)),
			StockQuantity:     inventory.StockQuantity,
			ReservedQuantity:  0,
			AvailableQuantity: inventory.StockQuantity,
		},
	}, nil
}

// UpdateInventory is the resolver for the updateInventory field.
func (r *mutationResolver) UpdateInventory(ctx context.Context, variantID string, quantity int) (*model.Inventory, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil || claims.Role != "admin" {
		return nil, errors.New("unauthorized: admin access required")
	}

	vID, err := strconv.Atoi(variantID)
	if err != nil {
		return nil, errors.New("invalid variant ID")
	}

	var inventory models.Inventory
	if err := r.DB.Where("variant_id = ?", vID).First(&inventory).Error; err != nil {
		return nil, errors.New("inventory not found")
	}

	inventory.StockQuantity = quantity

	if err := r.DB.Save(&inventory).Error; err != nil {
		return nil, err
	}

	availableQty := inventory.StockQuantity - inventory.ReservedQuantity

	return &model.Inventory{
		ID:                strconv.Itoa(int(inventory.ID)),
		VariantID:         variantID,
		StockQuantity:     inventory.StockQuantity,
		ReservedQuantity:  inventory.ReservedQuantity,
		AvailableQuantity: availableQty,
	}, nil
}
