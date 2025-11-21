package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/config"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/graph"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/graph/generated"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/database"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/middleware"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/repository"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/service"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	database.Connect()
	database.Migrate()

	// Initialize repositories
	userRepo := repository.NewUserRepository(database.DB)
	productRepo := repository.NewProductRepository(database.DB)
	cartRepo := repository.NewCartRepository(database.DB)
	orderRepo := repository.NewOrderRepository(database.DB)
	paymentRepo := repository.NewPaymentRepository(database.DB)

	// Initialize services
	paymentService := service.NewPaymentService()

	// Initialize resolver
	resolver := &graph.Resolver{
		DB:                database.DB,
		UserRepository:    userRepo,
		ProductRepository: productRepo,
		CartRepository:    cartRepo,
		OrderRepository:   orderRepo,
		PaymentRepository: paymentRepo,
		PaymentService:    paymentService,
	}

	// Create GraphQL server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	// Setup router
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.CorsMiddleware().Handler)
	router.Use(middleware.AuthMiddleware)

	// Routes
	router.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	router.Handle("/query", srv)

	// OAuth routes
	router.Get("/auth/google", handleGoogleLogin)
	router.Get("/auth/google/callback", handleGoogleCallback)

	port := os.Getenv("PORT")
	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Implement OAuth login redirect
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Implement OAuth callback handler
}
