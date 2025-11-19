package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"tshirt-ecommerce-api/config"
	"tshirt-ecommerce-api/graph"
	"tshirt-ecommerce-api/graph/generated"
	"tshirt-ecommerce-api/internal/database"
	"tshirt-ecommerce-api/internal/middleware"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	database.Connect()
	database.Migrate()

	// Initialize resolver
	resolver := graph.NewResolver()

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
