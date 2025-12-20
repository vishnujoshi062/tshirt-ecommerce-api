package repository

import (
	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.DB.Save(user).Error
}

// UpsertClerkUser creates or updates a user from Clerk authentication
func (r *UserRepository) UpsertClerkUser(clerkUserID, email, name, phone string) (*models.User, error) {
	user := &models.User{
		ClerkUserID:   clerkUserID,
		Email:         email,
		Name:          name,
		Phone:         phone,
		PhoneVerified: phone != "", // Mark as verified if phone is provided
		Role:          "user",
		OAuthProvider: "clerk",
		OAuthID:       clerkUserID,
	}

	// Use FirstOrCreate to insert or update
	result := r.DB.Where("clerk_user_id = ?", clerkUserID).
		Assign(user).
		FirstOrCreate(user)

	return user, result.Error
}

// GetUserByClerkID retrieves a user by their Clerk ID
func (r *UserRepository) GetUserByClerkID(clerkUserID string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("clerk_user_id = ?", clerkUserID).First(&user).Error
	return &user, err
}
