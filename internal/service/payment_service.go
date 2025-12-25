package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/razorpay/razorpay-go"
)

type PaymentService struct {
	Client *razorpay.Client
}

func NewPaymentService() *PaymentService {
	client := razorpay.NewClient(
		os.Getenv("RAZORPAY_KEY_ID"),
		os.Getenv("RAZORPAY_KEY_SECRET"),
	)

	return &PaymentService{
		Client: client,
	}
}

func (ps *PaymentService) CreateOrder(amount float64, currency string, receipt string) (map[string]interface{}, error) {
	if ps.Client == nil {
		return nil, fmt.Errorf("payment client not initialized")
	}

	// Convert amount to paise (multiply by 100)
	amountPaise := int64(amount * 100)

	params := map[string]interface{}{
		"amount":   amountPaise,
		"currency": currency,
		"receipt":  receipt,
	}

	order, err := ps.Client.Order.Create(params, nil)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (ps *PaymentService) VerifySignature(razorpayOrderID, razorpayPaymentID, razorpaySignature string) bool {
	secret := os.Getenv("RAZORPAY_KEY_SECRET")
	if secret == "" {
		return false
	}

	data := razorpayOrderID + "|" + razorpayPaymentID
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return expectedSignature == razorpaySignature
}
