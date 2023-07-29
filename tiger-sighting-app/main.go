package main

import (
	"fmt"
	"log"
	conf "tiger-sighting-app/config"
	"tiger-sighting-app/pkg/repository/store"
)

func main() {
	// Read the configuration from server.yml
	config, err := conf.ReadConfig("config/local/server.yml")
	if err != nil {
		log.Fatalf("Failed to read configuration: %v", err)
	}
	fmt.Println(config)

	// Initialize the database connection
	dbConnectionString := conf.BuildDBConnectionString(config.Database)
	fmt.Println(dbConnectionString)
	connection, err := store.NewPostgresDB(dbConnectionString)
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	defer connection.Close()

}
