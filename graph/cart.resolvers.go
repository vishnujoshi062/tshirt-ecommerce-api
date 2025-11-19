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

// MyCart is the resolver for the myCart field.
func (r *queryResolver) MyCart(ctx context.Context) (*model.Cart, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil {
		return nil, errors.New("unauthorized")
	}

	var cart models.Cart
	err := r.DB.Preload("CartItems.Variant.Product").
		Preload("CartItems.Variant.Inventory").
		Where("user_id = ?", claims.UserID).
		First(&cart).Error

	if err != nil {
		// Create new cart if doesn't exist
		cart = models.Cart{
			UserID: claims.UserID,
		}
		if err := r.DB.Create(&cart).Error; err != nil {
			return nil, err
		}
	}

	var items []*model.CartItem
	var subtotal float64

	for _, item := range cart.CartItems {
		itemPrice := item.Variant.Product.BasePrice + item.Variant.PriceModifier
		itemTotal := itemPrice * float64(item.Quantity)
		subtotal += itemTotal

		cartItem := &model.CartItem{
			ID:        strconv.Itoa(int(item.ID)),
			CartID:    strconv.Itoa(int(item.CartID)),
			Quantity:  item.Quantity,
			ItemTotal: itemTotal,
			AddedAt:   item.AddedAt.Format(time.RFC3339),
		}
		items = append(items, cartItem)
	}

	return &model.Cart{
		ID:        strconv.Itoa(int(cart.ID)),
		UserID:    strconv.Itoa(int(cart.UserID)),
		Items:     items,
		Subtotal:  subtotal,
		CreatedAt: cart.CreatedAt.Format(time.RFC3339),
	}, nil
}

// AddToCart is the resolver for the addToCart field.
func (r *mutationResolver) AddToCart(ctx context.Context, input model.AddToCartInput) (*model.Cart, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil {
		return nil, errors.New("unauthorized")
	}

	variantID, err := strconv.Atoi(input.VariantID)
	if err != nil {
		return nil, errors.New("invalid variant ID")
	}

	// Check inventory availability
	var inventory models.Inventory
	if err := r.DB.Where("variant_id = ?", variantID).First(&inventory).Error; err != nil {
		return nil, errors.New("product variant not found")
	}

	availableQty := inventory.StockQuantity - inventory.ReservedQuantity
	if availableQty < input.Quantity {
		return nil, errors.New("insufficient stock available")
	}

	// Get or create cart
	var cart models.Cart
	err = r.DB.Where("user_id = ?", claims.UserID).First(&cart).Error
	if err != nil {
		cart = models.Cart{UserID: claims.UserID}
		if err := r.DB.Create(&cart).Error; err != nil {
			return nil, err
		}
	}

	// Check if item already in cart
	var existingItem models.CartItem
	err = r.DB.Where("cart_id = ? AND variant_id = ?", cart.ID, variantID).First(&existingItem).Error

	if err == nil {
		// Update quantity if item exists
		newQuantity := existingItem.Quantity + input.Quantity
		if availableQty < newQuantity {
			return nil, errors.New("insufficient stock for updated quantity")
		}
		existingItem.Quantity = newQuantity
		if err := r.DB.Save(&existingItem).Error; err != nil {
			return nil, err
		}
	} else {
		// Create new cart item
		cartItem := &models.CartItem{
			CartID:    cart.ID,
			VariantID: uint(variantID),
			Quantity:  input.Quantity,
			AddedAt:   time.Now(),
		}
		if err := r.DB.Create(cartItem).Error; err != nil {
			return nil, err
		}
	}

	// Return updated cart
	return r.MyCart(ctx)
}

// RemoveFromCart is the resolver for the removeFromCart field.
func (r *mutationResolver) RemoveFromCart(ctx context.Context, cartItemID string) (*model.Cart, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil {
		return nil, errors.New("unauthorized")
	}

	itemID, err := strconv.Atoi(cartItemID)
	if err != nil {
		return nil, errors.New("invalid cart item ID")
	}

	// Verify ownership
	var cartItem models.CartItem
	if err := r.DB.Preload("Cart").First(&cartItem, itemID).Error; err != nil {
		return nil, errors.New("cart item not found")
	}

	if cartItem.Cart.UserID != claims.UserID {
		return nil, errors.New("unauthorized")
	}

	if err := r.DB.Delete(&cartItem).Error; err != nil {
		return nil, err
	}

	return r.MyCart(ctx)
}

// UpdateCartItemQuantity is the resolver for the updateCartItemQuantity field.
func (r *mutationResolver) UpdateCartItemQuantity(ctx context.Context, cartItemID string, quantity int) (*model.Cart, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil {
		return nil, errors.New("unauthorized")
	}

	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	itemID, err := strconv.Atoi(cartItemID)
	if err != nil {
		return nil, errors.New("invalid cart item ID")
	}

	var cartItem models.CartItem
	if err := r.DB.Preload("Cart").Preload("Variant.Inventory").First(&cartItem, itemID).Error; err != nil {
		return nil, errors.New("cart item not found")
	}

	if cartItem.Cart.UserID != claims.UserID {
		return nil, errors.New("unauthorized")
	}

	// Check inventory
	availableQty := cartItem.Variant.Inventory.StockQuantity - cartItem.Variant.Inventory.ReservedQuantity
	if availableQty < quantity {
		return nil, errors.New("insufficient stock available")
	}

	cartItem.Quantity = quantity
	if err := r.DB.Save(&cartItem).Error; err != nil {
		return nil, err
	}

	return r.MyCart(ctx)
}

// ClearCart is the resolver for the clearCart field.
func (r *mutationResolver) ClearCart(ctx context.Context) (bool, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil {
		return false, errors.New("unauthorized")
	}

	var cart models.Cart
	if err := r.DB.Where("user_id = ?", claims.UserID).First(&cart).Error; err != nil {
		return false, err
	}

	if err := r.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
		return false, err
	}

	return true, nil
}
