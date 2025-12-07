package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/auth"
)

type contextKey string

const UserContextKey = contextKey("user")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// ✅ Allow public OPTIONS requests (CORS)
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// ✅ Allow GraphQL Playground without token
		if r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}

		// ✅ Read auth token
		authHeader := r.Header.Get("Authorization")

		// ✅ Allow unauthenticated operations
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		// ✅ Validate Clerk token
		claims, err := auth.ValidateClerkToken(authHeader)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Unauthorized (invalid Clerk token)",
			})
			return
		}

		// ✅ Attach user context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) *auth.ClerkClaims {
	user, ok := ctx.Value(UserContextKey).(*auth.ClerkClaims)
	if !ok {
		return nil
	}
	return user
}
