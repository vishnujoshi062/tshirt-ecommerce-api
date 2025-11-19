package graph

import (
	"context"
	"errors"
	"strconv"
	"time"

	"tshirt-ecommerce-api/graph/model"
	"tshirt-ecommerce-api/internal/middleware"
	"tshirt-ecommerce-api/internal/models"
)

// MyOrders is the resolver for the myOrders field.
func (r *queryResolver) MyOrders(ctx context.Context) ([]*model.Order, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil {
		return nil, errors.New("unauthorized")
	}

	var orders []models.Order
	if err := r.DB.Preload("OrderItems.Variant.Product").
		Preload("Payment").
		Where("user_id = ?", claims.UserID).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}

	return convertOrdersToModel(orders), nil
}

// Order is the resolver for the order field.
func (r *queryResolver) Order(ctx context.Context, id string) (*model.Order, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil {
		return nil, errors.New("unauthorized")
	}

	orderID, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	var order models.Order
	if err := r.DB.Preload("OrderItems.Variant.Product").
		Preload("Payment").
		First(&order, orderID).Error; err != nil {
		return nil, errors.New("order not found")
	}

	if order.UserID != claims.UserID && claims.Role != "admin" {
		return nil, errors.New("unauthorized")
	}

	return convertOrderToModel(&order), nil
}

// AllOrders is the resolver for the allOrders field.
func (r *queryResolver) AllOrders(ctx context.Context, status *string) ([]*model.Order, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil || claims.Role != "admin" {
		return nil, errors.New("unauthorized: admin access required")
	}

	query := r.DB.Preload("OrderItems.Variant.Product").Preload("Payment")

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	var orders []models.Order
	if err := query.Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}

	return convertOrdersToModel(orders), nil
}

// CreateOrder is the resolver for the createOrder field.
func (r *mutationResolver) CreateOrder(ctx context.Context, input model.CreateOrderInput) (*model.Order, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil {
		return nil, errors.New("unauthorized")
	}

	// Get user's cart
	var cart models.Cart
	if err := r.DB.Preload("CartItems.Variant.Product").
		Preload("CartItems.Variant.Inventory").
		Where("user_id = ?", claims.UserID).
		First(&cart).Error; err != nil {
		return nil, errors.New("cart not found")
	}

	if len(cart.CartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Start transaction
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Calculate total and reserve inventory
	var totalAmount float64
	for _, item := range cart.CartItems {
		inventory := item.Variant.Inventory
		availableQty := inventory.StockQuantity - inventory.ReservedQuantity

		if availableQty < item.Quantity {
			tx.Rollback()
			return nil, errors.New("insufficient stock for item: " + item.Variant.SKU)
		}

		// Reserve inventory
		inventory.ReservedQuantity += item.Quantity
		if err := tx.Save(&inventory).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		itemPrice := item.Variant.Product.BasePrice + item.Variant.PriceModifier
		totalAmount += itemPrice * float64(item.Quantity)
	}

	// Create order
	order := &models.Order{
		UserID:          claims.UserID,
		TotalAmount:     totalAmount,
		Status:          "pending",
		ShippingAddress: input.ShippingAddress,
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create order items
	for _, cartItem := range cart.CartItems {
		itemPrice := cartItem.Variant.Product.BasePrice + cartItem.Variant.PriceModifier
		orderItem := &models.OrderItem{
			OrderID:   order.ID,
			VariantID: cartItem.VariantID,
			Quantity:  cartItem.Quantity,
			UnitPrice: itemPrice,
			Subtotal:  itemPrice * float64(cartItem.Quantity),
		}
		if err := tx.Create(orderItem).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Clear cart
	if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Reload order with associations
	if err := r.DB.Preload("OrderItems.Variant.Product").First(order, order.ID).Error; err != nil {
		return nil, err
	}

	return convertOrderToModel(order), nil
}

// UpdateOrderStatus is the resolver for the updateOrderStatus field.
func (r *mutationResolver) UpdateOrderStatus(ctx context.Context, orderID string, status string) (*model.Order, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil || claims.Role != "admin" {
		return nil, errors.New("unauthorized: admin access required")
	}

	oID, err := strconv.Atoi(orderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	validStatuses := []string{"pending", "confirmed", "shipped", "delivered", "cancelled"}
	isValid := false
	for _, s := range validStatuses {
		if s == status {
			isValid = true
			break
		}
	}

	if !isValid {
		return nil, errors.New("invalid order status")
	}

	var order models.Order
	if err := r.DB.First(&order, oID).Error; err != nil {
		return nil, errors.New("order not found")
	}

	order.Status = status
	if err := r.DB.Save(&order).Error; err != nil {
		return nil, err
	}

	// If cancelled, release reserved inventory
	if status == "cancelled" {
		var orderItems []models.OrderItem
		r.DB.Preload("Variant.Inventory").Where("order_id = ?", order.ID).Find(&orderItems)

		for _, item := range orderItems {
			inventory := item.Variant.Inventory
			inventory.ReservedQuantity -= item.Quantity
			if inventory.ReservedQuantity < 0 {
				inventory.ReservedQuantity = 0
			}
			r.DB.Save(&inventory)
		}
	}

	// If delivered, deduct from stock
	if status == "delivered" {
		var orderItems []models.OrderItem
		r.DB.Preload("Variant.Inventory").Where("order_id = ?", order.ID).Find(&orderItems)

		for _, item := range orderItems {
			inventory := item.Variant.Inventory
			inventory.StockQuantity -= item.Quantity
			inventory.ReservedQuantity -= item.Quantity
			if inventory.StockQuantity < 0 {
				inventory.StockQuantity = 0
			}
			if inventory.ReservedQuantity < 0 {
				inventory.ReservedQuantity = 0
			}
			r.DB.Save(&inventory)
		}
	}

	if err := r.DB.Preload("OrderItems.Variant.Product").Preload("Payment").First(&order, oID).Error; err != nil {
		return nil, err
	}

	return convertOrderToModel(&order), nil
}

// CancelOrder is the resolver for the cancelOrder field.
func (r *mutationResolver) CancelOrder(ctx context.Context, orderID string) (*model.Order, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil {
		return nil, errors.New("unauthorized")
	}

	oID, err := strconv.Atoi(orderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	var order models.Order
	if err := r.DB.First(&order, oID).Error; err != nil {
		return nil, errors.New("order not found")
	}

	if order.UserID != claims.UserID {
		return nil, errors.New("unauthorized")
	}

	if order.Status == "shipped" || order.Status == "delivered" {
		return nil, errors.New("cannot cancel order in current status")
	}

	return r.UpdateOrderStatus(ctx, orderID, "cancelled")
}

// Helper functions
func convertOrderToModel(order *models.Order) *model.Order {
	var items []*model.OrderItem
	for _, item := range order.OrderItems {
		orderItem := &model.OrderItem{
			ID:        strconv.Itoa(int(item.ID)),
			OrderID:   strconv.Itoa(int(item.OrderID)),
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Subtotal:  item.Subtotal,
		}
		items = append(items, orderItem)
	}

	modelOrder := &model.Order{
		ID:              strconv.Itoa(int(order.ID)),
		UserID:          strconv.Itoa(int(order.UserID)),
		Items:           items,
		TotalAmount:     order.TotalAmount,
		Status:          order.Status,
		ShippingAddress: order.ShippingAddress,
		CreatedAt:       order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       order.UpdatedAt.Format(time.RFC3339),
	}

	if order.Payment.ID != 0 {
		modelOrder.Payment = &model.Payment{
			ID:            strconv.Itoa(int(order.Payment.ID)),
			OrderID:       strconv.Itoa(int(order.Payment.OrderID)),
			Amount:        order.Payment.Amount,
			Status:        order.Payment.Status,
			PaymentMethod: order.Payment.PaymentMethod,
			TransactionID: &order.Payment.TransactionID,
			CreatedAt:     order.Payment.CreatedAt.Format(time.RFC3339),
		}
	}

	return modelOrder
}

func convertOrdersToModel(orders []models.Order) []*model.Order {
	var result []*model.Order
	for _, order := range orders {
		result = append(result, convertOrderToModel(&order))
	}
	return result
}
