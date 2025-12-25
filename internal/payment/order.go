package payment

import "errors"

func (r *RazorpayService) CreateOrder(amount int64, receipt string) (map[string]interface{}, error) {

	data := map[string]interface{}{
		"amount":   amount * 100, // Razorpay uses paise
		"currency": "INR",
		"receipt":  receipt,
	}

	order, err := r.Client.Order.Create(data, nil)
	if err != nil {
		return nil, errors.New("failed to create razorpay order")
	}

	return order, nil
}
