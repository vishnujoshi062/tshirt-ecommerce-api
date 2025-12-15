package database

import (
	"fmt"
	"log"
	"os"

	"github.com/vishnujoshi062/tshirt-ecommerce-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection established")
}

func Migrate() {
	// Drop existing foreign key constraints that are no longer valid
	// Cart.UserID and Order.UserID store Clerk IDs (strings), not database user IDs (integers)
	// So they cannot have foreign key relationships with User.ID
	// Try to drop common constraint names (IF EXISTS prevents errors if they don't exist)
	dropStatements := []string{
		"ALTER TABLE IF EXISTS carts DROP CONSTRAINT IF EXISTS fk_users_carts",
		"ALTER TABLE IF EXISTS carts DROP CONSTRAINT IF EXISTS users_carts_user_id_fkey",
		"ALTER TABLE IF EXISTS orders DROP CONSTRAINT IF EXISTS fk_users_orders",
		"ALTER TABLE IF EXISTS orders DROP CONSTRAINT IF EXISTS users_orders_user_id_fkey",
	}

	for _, stmt := range dropStatements {
		// IF EXISTS prevents errors if constraints don't exist, so we ignore errors here
		_ = DB.Exec(stmt).Error
	}

	err := DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.ProductVariant{},
		&models.Inventory{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database migration completed")
}
