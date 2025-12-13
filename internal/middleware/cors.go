package middleware

import (
    "net/http"
    "strings"
    
    "github.com/go-chi/cors"
)

func CorsMiddleware() *cors.Cors {
	return cors.New(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			// Allow production frontend
			if strings.Contains(origin, "vercel.app") {
				return true
			}

			// Allow localhost in development
			if strings.HasPrefix(origin, "http://localhost:") ||
				strings.HasPrefix(origin, "http://127.0.0.1:") {
				return true
			}

			return true // Allow all for now to debug
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

