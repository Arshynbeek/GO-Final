package server

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"project/mod/internal/structs"
	"project/mod/pkg/utils"
)

var DB *gorm.DB
var ERR error

func ConnectDB(connection string) (*gorm.DB, error) {
	var db *gorm.DB
	for i := 0; i < 10; i++ {
		db, ERR = gorm.Open(postgres.Open(connection), &gorm.Config{})
		if ERR == nil {
			break
		}
		log.Printf("Database RUIN failed: %v, retrying...", ERR)
		time.Sleep(5 * time.Second)
	}
	if ERR != nil {
		return nil, ERR
	}

	ERR = db.AutoMigrate(&structs.User{}, &structs.Category{}, &structs.Food{}, &structs.Order{}, &structs.Feedback{})
	if ERR != nil {
		log.Fatalf("Failed to migrate database: %v", ERR)
	}

	utils.Initialize(db)

	fmt.Println("Database RUIN was successful!")
	return db, nil
}
