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
	
	// TODO: Check if user has admin role from Clerk metadata
	// For now, allow all authenticated users
	// In production, check user.PublicMetadata or custom claims for admin role
	
	return nil
}
