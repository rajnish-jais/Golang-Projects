package main

import (
	"log"
	conf "tiger-sighting-app/config"
	"tiger-sighting-app/pkg/auth"
	"tiger-sighting-app/pkg/repository"
	"tiger-sighting-app/pkg/server"
)

func main() {
	// Read the configuration from server.yml
	config, err := conf.ReadConfig("config/local/server.yml")
	if err != nil {
		log.Fatalf("Failed to read configuration: %v", err)
	}

	// Initialize the database connection
	dbConnectionString := conf.BuildDBConnectionString(config.Database)

	store, err := repository.NewPostgresRepository(dbConnectionString)
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}

	// Initialize the JWT authentication
	auth := auth.NewAuth(config.JWT.SecretKey)

	// Initialize the server
	srv := server.NewServer()

	// Set up the routes and handlers
	srv.SetupRoutes(store, auth)

	// Start the server
	port := "8080" // You can change the port as needed
	err = srv.Start(port)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
