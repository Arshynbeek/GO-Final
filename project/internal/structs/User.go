package structs

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"size:255;not null"`
	Username string `gorm:"size:255;not null;unique"`
	Email    string `gorm:"size:255;not null;unique"`
	Password string `gorm:"size:255;not null"`
	Picture  string `gorm:"size:255"`
	Admin    bool   `gorm:"default:false"`
}
