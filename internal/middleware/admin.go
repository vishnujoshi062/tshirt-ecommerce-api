package middleware

import (
	"context"
	"errors"
)

func RequireAdmin(ctx context.Context) error {
	user := GetUserFromContext(ctx)
	if user == nil {
		return errors.New("not authenticated")
	}
	
	// Check if user has admin role from Clerk metadata
	if user.Role != "admin" {
		return errors.New("access denied: admin role required")
	}
	
	return nil
}
