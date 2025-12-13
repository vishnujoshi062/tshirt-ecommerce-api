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
            if origin == "https://ecommerce-frontend-bigshow.vercel.app/" {
                return true
            }
            
            // Allow any localhost or 127.0.0.1 for development
            if strings.HasPrefix(origin, "http://localhost:") || 
               strings.HasPrefix(origin, "http://127.0.0.1:") {
                return true
            }
            
            return false
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
