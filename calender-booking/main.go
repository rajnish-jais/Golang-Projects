/*
Design a book and time suggestion feature of a calendar application

Endpoints:
/time/book
Inputs:

Start time
End time
Output
Message either indicating success or failure
/time/suggest
Inputs

Start time
End time
Meeting duration
Output:
Suggest time slot (start and end time) if possible or failure

Expectations
Clean code
Extensible in nature

Table Booking
id
start_time
end_time

*/

package main

import (
	"log"

	conf "calender-booking/config"
	"calender-booking/pkg/repository"
	"calender-booking/pkg/server"
	"calender-booking/pkg/service"
)

func initializeService(config *conf.Config) (service.BookingService, error) {
	// Initialize the database connection
	dbConnectionString := conf.BuildDBConnectionString(config.Database)

	store, err := repository.NewPostgresRepository(dbConnectionString)
	if err != nil {
		return nil, err
	}

	// Initialize the service
	service := service.NewUserService(store)

	return service, nil
}

func main() {
	// Read the configuration from server.yml
	config, err := conf.ReadConfig("config/local/server.yml")
	if err != nil {
		log.Fatalf("Failed to read configuration: %v", err)
	}

	// Initialize the service
	service, err := initializeService(config)
	if err != nil {
		log.Fatalf("Failed to initialize the service: %v", err)
	}

	// Initialize the server
	srv := server.NewServer()

	// Set up the routes and handler
	srv.SetupRoutes(service)

	// Start the server
	err = srv.Start(config.Server.Port)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
