package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
)

func VerifySignature(orderID, paymentID, signature string) bool {

	secret := os.Getenv("RAZORPAY_KEY_SECRET")

	data := orderID + "|" + paymentID

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))

	expected := hex.EncodeToString(h.Sum(nil))

	return expected == signature
}
