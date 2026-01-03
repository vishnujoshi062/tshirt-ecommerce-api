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

// extractRoleFromMetadata attempts to extract the role from various possible
// locations in the JWT claims where Clerk might store public metadata
func extractRoleFromMetadata(rawClaims map[string]interface{}) string {
	// Check multiple possible locations for public metadata
	possiblePaths := []string{
		"org_metadata.role",
		"public_metadata.role",
		"metadata.role",
		"org_role",
		"role",
	}

	for _, path := range possiblePaths {
		if role := getNestedValue(rawClaims, path); role != "" {
			return role
		}
	}

	return ""
}

// getNestedValue extracts a value from a nested map using dot notation
func getNestedValue(data map[string]interface{}, path string) string {
	parts := strings.Split(path, ".")
	current := data

	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part - extract the value
			if val, exists := current[part]; exists {
				if str, ok := val.(string); ok {
					return str
				}
			}
			return ""
		}

		// Navigate deeper
		if val, exists := current[part]; exists {
			if nextMap, ok := val.(map[string]interface{}); ok {
				current = nextMap
			} else {
				return ""
			}
		} else {
			return ""
		}
	}

	return ""
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

	// Extract email and role from the raw JWT claims
	email := ""
	role := ""
	
	if rawClaims, ok := claims.Custom.(map[string]interface{}); ok {
		// Extract email
		if emailVal, exists := rawClaims["email"]; exists {
			if emailStr, ok := emailVal.(string); ok {
				email = emailStr
			}
		}
		
		// Extract role from public metadata
		// Clerk stores public metadata in the JWT claims
		// The exact location may vary, so we check multiple possible paths
		role = extractRoleFromMetadata(rawClaims)
		
		// Log for debugging (can be removed in production)
		if role == "" {
			log.Printf("CLERK AUTH: No role found in JWT claims for user %s", claims.Subject)
		} else {
			log.Printf("CLERK AUTH: Found role '%s' for user %s", role, claims.Subject)
		}
	}

	return &ClerkClaims{
		UserID: claims.Subject,
		Email:  email,
		Role:   role,
	}, nil
}
