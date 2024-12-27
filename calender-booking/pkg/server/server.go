package server

import (
	"log"
	"net/http"
	"os"

	"calender-booking/pkg/handler"
	"calender-booking/pkg/service"

	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
	logger *log.Logger
}

func NewServer() *server {
	return &server{
		router: mux.NewRouter(),
		logger: log.New(os.Stdout, "[Calender Booking APP] ", log.LstdFlags),
	}
}

func (s *server) SetupRoutes(postmanService service.BookingService) {
	handlers := handler.NewHandlers(postmanService, s.logger)

	s.router.HandleFunc("/time/book", handlers.BookTimeSlot).Methods("POST")
	s.router.HandleFunc("/time/suggest", handlers.SuggestTimeSlot).Methods("GET")
}

func (s *server) Start(port string) error {
	s.logger.Printf("Starting server on port %s...", port)
	return http.ListenAndServe(":"+port, s.router)
}
