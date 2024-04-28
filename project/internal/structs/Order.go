package structs

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	FoodID uint `gorm:"not null"`
	UserID uint `gorm:"not null"`
	Food   Food `gorm:"foreignKey:FoodID"`
	User   User `gorm:"foreignKey:UserID"`
}
