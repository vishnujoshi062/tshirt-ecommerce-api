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

	if user.Role != "admin" {
		return errors.New("forbidden")
	}

	return nil
}
