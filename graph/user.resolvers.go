package graph

import (
	"context"
	"errors"
	"fmt"

	"tshirt-ecommerce-api/graph/model"
	"tshirt-ecommerce-api/internal/middleware"
	"tshirt-ecommerce-api/internal/models"
	"tshirt-ecommerce-api/internal/repository"
	"tshirt-ecommerce-api/internal/utils"
)

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, input model.RegisterInput) (*model.AuthPayload, error) {
	userRepo := repository.NewUserRepository(r.DB)

	// Check if user exists
	_, err := userRepo.GetUserByEmail(input.Email)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:         input.Email,
		PasswordHash:  hashedPassword,
		Name:          input.Name,
		Role:          "user",
		OAuthProvider: "local",
	}

	if input.Phone != nil {
		user.Phone = *input.Phone
	}
	if input.Address != nil {
		user.Address = *input.Address
	}

	if err := userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &model.AuthPayload{
		Token: token,
		User: &model.User{
			ID:        int(user.ID),
			Email:     user.Email,
			Name:      user.Name,
			Phone:     &user.Phone,
			Address:   &user.Address,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.AuthPayload, error) {
	userRepo := repository.NewUserRepository(r.DB)

	user, err := userRepo.GetUserByEmail(input.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(input.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &model.AuthPayload{
		Token: token,
		User: &model.User{
			ID:        int(user.ID),
			Email:     user.Email,
			Name:      user.Name,
			Phone:     &user.Phone,
			Address:   &user.Address,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// UpdateProfile is the resolver for the updateProfile field.
func (r *mutationResolver) UpdateProfile(ctx context.Context, name *string, phone *string, address *string) (*model.User, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil {
		return nil, errors.New("unauthorized")
	}

	userRepo := repository.NewUserRepository(r.DB)
	user, err := userRepo.GetUserByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	if name != nil {
		user.Name = *name
	}
	if phone != nil {
		user.Phone = *phone
	}
	if address != nil {
		user.Address = *address
	}

	if err := userRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	return &model.User{
		ID:        int(user.ID),
		Email:     user.Email,
		Name:      user.Name,
		Phone:     &user.Phone,
		Address:   &user.Address,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil {
		return nil, errors.New("unauthorized")
	}

	userRepo := repository.NewUserRepository(r.DB)
	user, err := userRepo.GetUserByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:        int(user.ID),
		Email:     user.Email,
		Name:      user.Name,
		Phone:     &user.Phone,
		Address:   &user.Address,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetUser is the resolver for the getUser field.
func (r *queryResolver) GetUser(ctx context.Context, id string) (*model.User, error) {
	claims := middleware.GetUserFromContext(ctx)
	if claims == nil || claims.Role != "admin" {
		return nil, errors.New("unauthorized: admin access required")
	}

	userRepo := repository.NewUserRepository(r.DB)
	userID := uint(0)
	if _, err := fmt.Sscanf(id, "%d", &userID); err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:        int(user.ID),
		Email:     user.Email,
		Name:      user.Name,
		Phone:     &user.Phone,
		Address:   &user.Address,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
