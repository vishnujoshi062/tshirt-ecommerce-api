package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type contextKey string

const UserContextKey = contextKey("user")

// ClerkClaims represents the user information from Clerk
type ClerkClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		// ✅ Allow public requests if no token
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		// TODO: Implement Clerk token validation here
		// For now, we'll create a placeholder that assumes the token is valid
		// and extract user info from custom headers that Clerk might provide

		// In a real implementation, you would validate the Clerk JWT token here
		// using the Clerk SDK or by calling their API

		// For demonstration purposes, we'll extract user info from headers
		// (Clerk typically adds user info to headers)
		userIDStr := r.Header.Get("Clerk-User-ID")
		email := r.Header.Get("Clerk-User-Email")

		// Convert string user ID to uint
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil || userIDStr == "" || email == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid or missing authentication",
			})
			return
		}

		// Create claims object with user info
		claims := &ClerkClaims{
			UserID: uint(userID),
			Email:  email,
			Role:   "user", // Default role, could be customized based on Clerk user data
		}

		// ✅ Attach user context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) *ClerkClaims {
	user, ok := ctx.Value(UserContextKey).(*ClerkClaims)
	if !ok {
		return nil
	}
	return user
}
