package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/config"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/graph"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/graph/generated"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/database"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/middleware"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/repository"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/service"
)

func main() {
	// Try to load .env file from different possible locations
	possiblePaths := []string{
		".env",
		"../.env",
		"../../.env",
		filepath.Join(os.Getenv("PWD"), ".env"),
	}
	envLoaded := false
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err == nil {
				log.Printf("Loaded .env file from: %s", path)
				envLoaded = true
				break
			} else {
				log.Printf("Failed to load .env file from %s: %v", path, err)
			}
		}
	}
	if !envLoaded {
		// Load environment variables from default location
		config.LoadEnv()
	}

	// Initialize Clerk SDK with secret key
	clerkKey := config.GetEnv("CLERK_SECRET_KEY", "")
	if clerkKey == "" {
		log.Fatal("CLERK_SECRET_KEY environment variable is required")
	}
	clerk.SetKey(clerkKey)

	// Connect to database
	database.Connect()
	database.Migrate()

	// Initialize repositories
	userRepo := repository.NewUserRepository(database.DB)
	productRepo := repository.NewProductRepository(database.DB)
	cartRepo := repository.NewCartRepository(database.DB)
	orderRepo := repository.NewOrderRepository(database.DB)
	paymentRepo := repository.NewPaymentRepository(database.DB)
	promoCodeRepo := repository.NewPromoCodeRepository(database.DB)

	// Initialize services
	paymentService := service.NewPaymentService()
	promoCodeService := service.NewPromoCodeService(promoCodeRepo)

	// Initialize resolver
	resolver := &graph.Resolver{
		DB:                database.DB,
		UserRepository:    userRepo,
		ProductRepository: productRepo,
		CartRepository:    cartRepo,
		OrderRepository:   orderRepo,
		PaymentRepository: paymentRepo,
		PaymentService:    paymentService,
		PromoCodeRepo:     promoCodeRepo,
		PromoCodeService:  promoCodeService,
	}

	// Create GraphQL server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	// Setup router
	router := chi.NewRouter()

	// Middleware - CORS must be first
	router.Use(middleware.CorsMiddleware().Handler)
	router.Use(middleware.AuthMiddleware)

	// Routes
	router.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	router.Handle("/query", srv)

	// OAuth routes
	router.Get("/auth/google", handleGoogleLogin)
	router.Get("/auth/google/callback", handleGoogleCallback)

	port := config.GetEnv("PORT", "8080")
	log.Printf("Server starting on http://localhost:%s", port)
	log.Printf("GraphQL endpoint: http://localhost:%s/query", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Implement OAuth login redirect
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Implement OAuth callback handler
}
