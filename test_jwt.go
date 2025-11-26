package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/utils"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("No .env file found")
	}

	// Check if JWT_SECRET is loaded
	jwtSecret := os.Getenv("JWT_SECRET")
	fmt.Printf("JWT_SECRET: %s\n", jwtSecret)

	// Test token generation
	token, err := utils.GenerateToken(1, "test@example.com", "user")
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		return
	}

	fmt.Printf("Generated token: %s\n", token)

	// Test token validation
	claims, err := utils.ValidateToken(token)
	if err != nil {
		fmt.Printf("Error validating token: %v\n", err)
		return
	}

	fmt.Printf("Token validated successfully. UserID: %d, Email: %s, Role: %s\n", 
		claims.UserID, claims.Email, claims.Role)
}