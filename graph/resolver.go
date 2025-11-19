package graph

import (
	"context"
	"errors"

	"github.com/vishnujoshi062/tshirt-ecommerce-api/graph/model"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/middleware"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/models"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/repository"
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/utils"
)

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
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Name:         input.Name,
		Phone:        input.Phone,
		Address:      input.Address,
		Role:         "user",
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
			ID:      int(user.ID),
			Email:   user.Email,
			Name:    user.Name,
			Phone:   user.Phone,
			Address: user.Address,
			Role:    user.Role,
		},
	}, nil
}

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
			ID:    int(user.ID),
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		},
	}, nil
}

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
		ID:      int(user.ID),
		Email:   user.Email,
		Name:    user.Name,
		Phone:   user.Phone,
		Address: user.Address,
		Role:    user.Role,
	}, nil
}
