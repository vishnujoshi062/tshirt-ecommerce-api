package graph

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/vishnujoshi062/tshirt-ecommerce-api/graph/model"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/middleware"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/models"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/service"
)

// CreateRazorpayOrder is the resolver for the createRazorpayOrder field.
func (r *mutationResolver) CreateRazorpayOrder(ctx context.Context, orderID string) (*model.RazorpayOrder, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil {
		return nil, errors.New("unauthorized")
	}

	oID, err := strconv.Atoi(orderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	// Get order
	var order models.Order
	if err := r.DB.First(&order, oID).Error; err != nil {
		return nil, errors.New("order not found")
	}

	if order.UserID != claims.UserID {
		return nil, errors.New("unauthorized")
	}

	if order.Status != "pending" {
		return nil, errors.New("order is not pending")
	}

	// Create Razorpay order
	paymentService := service.NewPaymentService()
	receipt := fmt.Sprintf("order_%d_%d", order.ID, time.Now().Unix())

	rzpOrder, err := paymentService.CreateOrder(order.TotalAmount, "INR", receipt)
	if err != nil {
		return nil, err
	}

	// Create payment record
	payment := &models.Payment{
		OrderID:         order.ID,
		Amount:          order.TotalAmount,
		Status:          "pending",
		PaymentMethod:   "razorpay",
		RazorpayOrderID: rzpOrder["id"].(string),
	}

	if err := r.DB.Create(payment).Error; err != nil {
		return nil, err
	}

	// Update order with payment ID
	order.PaymentID = payment.ID
	r.DB.Save(&order)

	return &model.RazorpayOrder{
		ID:       rzpOrder["id"].(string),
		Amount:   order.TotalAmount,
		Currency: "INR",
		Receipt:  &receipt,
	}, nil
}

// VerifyPayment is the resolver for the verifyPayment field.
func (r *mutationResolver) VerifyPayment(ctx context.Context, input model.VerifyPaymentInput) (*model.Payment, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil {
		return nil, errors.New("unauthorized")
	}

	oID, err := strconv.Atoi(input.OrderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	// Get payment
	var payment models.Payment
	if err := r.DB.Where("order_id = ?", oID).First(&payment).Error; err != nil {
		return nil, errors.New("payment not found")
	}

	// Verify signature
	paymentService := service.NewPaymentService()
	isValid := paymentService.VerifySignature(
		input.RazorpayOrderID,
		input.RazorpayPaymentID,
		input.RazorpaySignature,
	)

	if !isValid {
		payment.Status = "failed"
		r.DB.Save(&payment)
		return nil, errors.New("payment verification failed")
	}

	// Update payment
	payment.TransactionID = input.RazorpayPaymentID
	payment.RazorpaySignature = input.RazorpaySignature
	payment.Status = "success"

	if err := r.DB.Save(&payment).Error; err != nil {
		return nil, err
	}

	// Update order status
	var order models.Order
	if err := r.DB.First(&order, oID).Error; err == nil {
		order.Status = "confirmed"
		r.DB.Save(&order)
	}

	return &model.Payment{
		ID:            strconv.Itoa(int(payment.ID)),
		OrderID:       strconv.Itoa(int(payment.OrderID)),
		Amount:        payment.Amount,
		Status:        payment.Status,
		PaymentMethod: payment.PaymentMethod,
		TransactionID: &payment.TransactionID,
		CreatedAt:     payment.CreatedAt.Format(time.RFC3339),
	}, nil
}
