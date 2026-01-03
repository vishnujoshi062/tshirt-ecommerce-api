package auth

import (
	"context"
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

	claims, err := jwt.Verify(context.Background(), &jwt.VerifyParams{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	role := ""
	email := ""

	// Extract email and role from Custom claims (raw JWT claims)
	if raw, ok := claims.Custom.(map[string]interface{}); ok {
		// Extract email
		if v, ok := raw["email"].(string); ok {
			email = v
		}

		// Extract role from public_metadata
		if pm, ok := raw["public_metadata"].(map[string]interface{}); ok {
			if r, ok := pm["role"].(string); ok {
				role = r
			}
		}
	}

	if role == "" {
		log.Printf("CLERK AUTH: No role found for user %s", claims.Subject)
	} else {
		log.Printf("CLERK AUTH: Role '%s' for user %s", role, claims.Subject)
	}

	return &ClerkClaims{
		UserID: claims.Subject,
		Email:  email,
		Role:   role,
	}, nil
}
