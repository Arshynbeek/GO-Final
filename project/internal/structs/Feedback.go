package structs

import (
	"gorm.io/gorm"
)

type Feedback struct {
	gorm.Model
	UserID  uint   `gorm:"not null"`
	FoodID  uint   `gorm:"not null"`
	Rating  int    `gorm:"check:rating >= 0 AND rating <= 5"`
	Comment string `gorm:"size:255;not null;"`
	User    User   `gorm:"foreignKey:UserID"`
	Food    Food   `gorm:"foreignKey:FoodID"`
}
