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
	if claims.PublicMetadata != nil {
		if r, ok := claims.PublicMetadata["role"].(string); ok {
			role = r
		}
	}

	if role == "" {
		log.Printf("CLERK AUTH: No admin role for user %s", claims.Subject)
	}

	return &ClerkClaims{
		UserID: claims.Subject,
		Email:  claims.Email,
		Role:   role,
	}, nil
}
