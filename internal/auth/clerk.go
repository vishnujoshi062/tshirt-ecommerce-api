package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2/jwt"
)

type ClerkClaims struct {
	UserID string
	Email  string
	Role   string
}

func ValidateClerkToken(authHeader string) (*ClerkClaims, error) {
	if authHeader == "" {
		return nil, errors.New("missing Authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// 1️⃣ Verify signature + expiry with Clerk
	verified, err := jwt.Verify(context.Background(), &jwt.VerifyParams{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	// 2️⃣ Decode JWT payload manually
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWT format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	// 3️⃣ Extract fields
	email, _ := claims["email"].(string)

	role := ""
	if pm, ok := claims["public_metadata"].(map[string]interface{}); ok {
		if r, ok := pm["role"].(string); ok {
			role = r
		}
	}

	log.Printf("CLERK AUTH DEBUG: user=%s role=%s", verified.Subject, role)

	return &ClerkClaims{
		UserID: verified.Subject,
		Email:  email,
		Role:   role,
	}, nil
}
