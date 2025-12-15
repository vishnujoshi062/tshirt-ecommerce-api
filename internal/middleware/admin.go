package middleware

import (
	"context"
	"errors"
)

func RequireAdmin(ctx context.Context) error {
	user := ctx.Value("user")
	if user == nil {
		return errors.New("not authenticated")
	}
	
	// Adjust this based on your Clerk user structure
	// Check if user has admin role/metadata
	userMap, ok := user.(map[string]interface{})
	if !ok {
		return errors.New("invalid user context")
	}
	
	role, exists := userMap["role"]
	if !exists || role != "admin" {
		return errors.New("unauthorized: admin access required")
	}
	
	return nil
}
