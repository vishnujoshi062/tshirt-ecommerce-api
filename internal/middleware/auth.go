package middleware

import (
	"context"
	"net/http"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/auth"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/repository"
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

		// No token → proceed WITHOUT user context
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Validate Clerk token
		claims, err := auth.ValidateClerkToken(authHeader)
		if err != nil {
			http.Error(w, "Unauthorized (invalid Clerk token)", http.StatusUnauthorized)
			return
		}

		// Attach user context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthMiddlewareWithSync automatically syncs user to database on authentication
func AuthMiddlewareWithSync(userRepo *repository.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Allow OPTIONS (CORS)
			if r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")

			// No token → proceed WITHOUT user context
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Validate Clerk token
			claims, err := auth.ValidateClerkToken(authHeader)
			if err != nil {
				http.Error(w, "Unauthorized (invalid Clerk token)", http.StatusUnauthorized)
				return
			}

			// Auto-sync user to database in background
			go func() {
				userRepo.UpsertClerkUser(
					claims.UserID,
					claims.Email,
					"", // name not available in JWT
					claims.PhoneNumber,
				)
			}()

			// Attach user context
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) *auth.ClerkClaims {
	user, ok := ctx.Value(UserContextKey).(*auth.ClerkClaims)
	if !ok {
		return nil
	}
	return user
}
