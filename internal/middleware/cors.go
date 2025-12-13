package middleware

import "github.com/go-chi/cors"

func CorsMiddleware() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://127.0.0.1:3000",
			"https://ecommerce-frontend-five-nu.vercel.app",
			// Add wildcard for development if needed (be careful in production)
			// "http://localhost:*",
		},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
			"PATCH",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-Requested-With",
			"X-CSRF-Token",
			"Content-Length",
		},
		ExposedHeaders: []string{
			"Link",
		},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
