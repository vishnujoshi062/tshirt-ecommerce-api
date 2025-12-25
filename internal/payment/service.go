package payment

import (
	"os"
	"github.com/razorpay/razorpay-go"
)

type RazorpayService struct {
	Client *razorpay.Client
}

func NewRazorpayService() *RazorpayService {
	client := razorpay.NewClient(
		os.Getenv("RAZORPAY_KEY_ID"),
		os.Getenv("RAZORPAY_KEY_SECRET"),
	)

	return &RazorpayService{
		Client: client,
	}
}
