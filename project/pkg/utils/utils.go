package utils

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"

	"project/mod/internal/structs"
)

var KEY = []byte("arch")

func GenerateJWT(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(KEY)
	return tokenString, err
}

func Initialize(db *gorm.DB) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var userCount int64
	tx.Model(&structs.User{}).Count(&userCount)
	if userCount == 0 {
		users := []structs.User{
			{
				Name:     "Admin",
				Username: "admin",
				Email:    "admin@admin.admin",
				Password: hashPassword("admin_password"),
				Picture:  "/static/images/icons/camera_off_icon.svg",
				Admin:    true,
			},
			{
				Name:     "User",
				Username: "user",
				Email:    "user@gmail.com",
				Password: hashPassword("user_password"),
				Picture:  "/static/images/icons/camera_off_icon.svg",
			},
		}
		if err := tx.Create(&users).Error; err != nil {
			tx.Rollback()
			return
		}
	}

	var categoryCount int64
	tx.Model(&structs.Category{}).Count(&categoryCount)
	if categoryCount == 0 {
		categories := []structs.Category{
			{Name: "Drinks"},
			{Name: "Snacks"},
		}
		if err := tx.Create(&categories).Error; err != nil {
			tx.Rollback()
			return
		}
	}

	var foodCount int64
	tx.Model(&structs.Food{}).Count(&foodCount)
	if foodCount == 0 {
		foods := []structs.Food{
			{Name: "MacCoffee", Description: "3 in 1 Original is simply a blend of premium coffee beans, non-dairy creamer and sugar.", Pictures: []string{"/static/images/products/maccoffee.jpg", "/static/images/products/table.jpg", "/static/images/products/red.jpg", "/static/images/products/orange.jpg"}, Quantity: 10, Price: 70, CategoryID: 1},
			{Name: "MacTea", Description: "Instant teas with fruity flavors that everyone enjoys drinking", Pictures: []string{"/static/images/products/mactea.jpg", "/static/images/image.jpg", "/static/images/image.jpg", "/static/images/image.jpg"}, Quantity: 20, Price: 60, CategoryID: 1},
			{Name: "Pocha (Cheese)", Description: "Tender and satisfying appetizer with cheese", Pictures: []string{"/static/images/products/pocha.jpg", "/static/images/image.jpg", "/static/images/image.jpg", "/static/images/image.jpg"}, Quantity: 12, Price: 250, CategoryID: 2},
			{Name: "Achma", Description: "Round and chocolate snack, which will be the best combination for tea", Pictures: []string{"/static/images/products/achma.jpg", "/static/images/image.jpg", "/static/images/image.jpg", "/static/images/image.jpg"}, Quantity: 4, Price: 280, CategoryID: 2},
		}
		if err := tx.Create(&foods).Error; err != nil {
			tx.Rollback()
			return
		}
	}

	var orderCount int64
	tx.Model(&structs.Order{}).Count(&orderCount)
	if orderCount == 0 {
		orders := []structs.Order{
			{FoodID: 1, Quantity: 1, UserID: 2},
			{FoodID: 4, Quantity: 1, UserID: 2},
		}
		if err := tx.Create(&orders).Error; err != nil {
			tx.Rollback()
			return
		}
	}

	var feedbackCount int64
	tx.Model(&structs.Feedback{}).Count(&feedbackCount)
	if feedbackCount == 0 {
		feedbacks := []structs.Feedback{
			{UserID: 2, FoodID: 4, Rating: 5, Comment: "Delicious Mediterranean style roll with chocolate. 😋"},
		}
		if err := tx.Create(&feedbacks).Error; err != nil {
			tx.Rollback()
			return
		}
	}

	tx.Commit()
}

func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}
	return string(hashedPassword)
}
