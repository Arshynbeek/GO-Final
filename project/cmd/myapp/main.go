package main

import (
	"log"

	"project/mod/internal/service"
	"project/mod/pkg/server"
)

func main() {
	db, err := server.ConnectDB("host=final_db user=admin password=root dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to RUIN database: %v", err)
	}
	server.DB = db

	service.Router()
}
