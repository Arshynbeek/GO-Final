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
				Admin:    true,
			},
			{
				Name:     "User",
				Username: "user",
				Email:    "user@gmail.com",
				Password: hashPassword("user_password"),
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
			{Name: "MacCoffee", Description: "3 in 1 Original is simply a blend of premium coffee beans, non-dairy creamer and sugar.", Pictures: []string{"/static/images/maccoffee.jpg"}, Quantity: 10, CategoryID: 1},
			{Name: "MacTea", Description: "Instant teas with fruity flavors that everyone enjoys drinking", Pictures: []string{"/static/images/mactea.jpg"}, Quantity: 20, CategoryID: 1},
			{Name: "Pocha (Cheese)", Description: "Tender and satisfying appetizer with cheese", Pictures: []string{"/static/images/pocha.jpg"}, Quantity: 12, CategoryID: 2},
			{Name: "Achma", Description: "Round and chocolate snack, which will be the best combination for tea", Pictures: []string{"/static/images/achma.jpg"}, Quantity: 4, CategoryID: 2},
		}
		if err := tx.Create(&foods).Error; err != nil {
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
