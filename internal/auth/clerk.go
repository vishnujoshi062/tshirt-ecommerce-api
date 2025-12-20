package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2/jwt"
)

type ClerkClaims struct {
	UserID      string
	Email       string
	PhoneNumber string
}

func ValidateClerkToken(authHeader string) (*ClerkClaims, error) {
	if authHeader == "" {
		return nil, errors.New("missing Authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := jwt.Verify(context.Background(), &jwt.VerifyParams{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	// Extract email and phone from the raw JWT claims
	email := ""
	phoneNumber := ""
	if rawClaims, ok := claims.Custom.(map[string]interface{}); ok {
		if emailVal, exists := rawClaims["email"]; exists {
			if emailStr, ok := emailVal.(string); ok {
				email = emailStr
			}
		}
		if phoneVal, exists := rawClaims["phone_number"]; exists {
			if phoneStr, ok := phoneVal.(string); ok {
				phoneNumber = phoneStr
			}
		}
	}

	return &ClerkClaims{
		UserID:      claims.Subject,
		Email:       email,
		PhoneNumber: phoneNumber,
	}, nil
}
