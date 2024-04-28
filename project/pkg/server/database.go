package server

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"project/mod/internal/structs"
	"project/mod/pkg/utils"
)

var DB *gorm.DB
var ERR error

func ConnectDB(connection string) (*gorm.DB, error) {
	DB, ERR := gorm.Open(postgres.Open(connection), &gorm.Config{})
	if ERR != nil {
		return nil, ERR
	}

	ERR = DB.AutoMigrate(&structs.User{}, &structs.Category{}, &structs.Food{}, &structs.Order{}, &structs.Feedback{})
	if ERR != nil {
		log.Fatalf("Failed to migrate database: %v", ERR)
	}

	utils.Initialize(DB)

	fmt.Println("database RUIN was successful!")
	return DB, nil
}
