package middleware

import (
	"context"
	"log"
	"net/http"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/auth"
)

type contextKey string

const UserContextKey = contextKey("user")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Allow OPTIONS (CORS)
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		log.Printf("AUTH MIDDLEWARE: method=%s path=%s auth_header=%s", r.Method, r.URL.Path, authHeader)

		// If token provided, try to validate it
		if authHeader != "" {
			claims, err := auth.ValidateClerkToken(authHeader)
			if err == nil {
				// Valid token → attach user context
				log.Printf("AUTH MIDDLEWARE: Valid token for user %s", claims.UserID)
				ctx := context.WithValue(r.Context(), UserContextKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			// Invalid token → log but continue (dev mode)
			log.Printf("AUTH MIDDLEWARE: Invalid token, error: %v, continuing without auth", err)
		}

		// No token or invalid token → proceed WITHOUT user context
		// Resolvers will handle auth checks as needed
		log.Printf("AUTH MIDDLEWARE: No auth, proceeding without user context")
		next.ServeHTTP(w, r)
	})
}

func GetUserFromContext(ctx context.Context) *auth.ClerkClaims {
	user, ok := ctx.Value(UserContextKey).(*auth.ClerkClaims)
	if !ok {
		return nil
	}
	return user
}
