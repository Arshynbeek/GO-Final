package main

import (
	"log"

	"project/mod/internal/service"
	"project/mod/pkg/server"
)

func main() {
	db, err := server.ConnectDB("host=localhost user=admin password=root dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to RUIN database: %v", err)
	}
	server.DB = db

	service.Router()
}

// project
// 	cmd
// 		myapp
// 			main.go

// 	deployments
// 		docker
// 			.dockerginore
// 			Dockerfile

// 	internal
// 		api
// 			handler.go

// 		service
// 			service.go

// 		structs
// 			// my structs

// 	pkg
// 		server
// 			database.go

// 		utils
// 			utils.go

// 	frontend
// 		templates
// 			 .html files

// 	docker-compose.yml
// 	go.mod
// 	go.sum
