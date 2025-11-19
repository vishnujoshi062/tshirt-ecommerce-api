package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"

	razorpay "github.com/razorpay/razorpay-go"
)

type PaymentService struct {
	Client *razorpay.Client
}

func NewPaymentService() *PaymentService {
	client := razorpay.NewClient(
		os.Getenv("RAZORPAY_KEY_ID"),
		os.Getenv("RAZORPAY_KEY_SECRET"),
	)

	return &PaymentService{Client: client}
}

func (s *PaymentService) CreateOrder(amount float64, currency string, receipt string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"amount":   int(amount * 100), // Amount in paise
		"currency": currency,
		"receipt":  receipt,
	}

	body, err := s.Client.Order.Create(data, nil)
	return body, err
}

func (s *PaymentService) VerifySignature(orderID, paymentID, signature string) bool {
	secret := os.Getenv("RAZORPAY_KEY_SECRET")
	data := orderID + "|" + paymentID

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}
