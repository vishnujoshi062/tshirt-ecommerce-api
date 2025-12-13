package middleware

import (
	"github.com/go-chi/cors"
)

func CorsMiddleware() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",                             // Local frontend
			"https://ecommerce-frontend-five-nu.vercel.app",    // Vercel frontend
			"https://tshirt-ecommerce-api-production.onrender.com", // Render backend
		},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowedHeaders: []string{
			"*", // Allow all headers including Authorization
		},
		ExposedHeaders: []string{
			"*",
		},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
