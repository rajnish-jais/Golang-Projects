package server

import (
	"log"
	"net/http"
	"os"
	"tiger-sighting-app/pkg/auth"
	"tiger-sighting-app/pkg/handlers"
	"tiger-sighting-app/pkg/middleware"
	"tiger-sighting-app/pkg/repository"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	logger *log.Logger
}

func NewServer() *Server {
	return &Server{
		router: mux.NewRouter(),
		logger: log.New(os.Stdout, "[Tiger-Sighting] ", log.LstdFlags),
	}
}

func (s *Server) SetupRoutes(tigerRepo repository.TigerRepository, auth *auth.Auth) {
	handlers := handlers.NewHandlers(tigerRepo, s.logger, auth)

	// Public routes
	s.router.HandleFunc("/signup", handlers.SignupHandler).Methods("POST")
	s.router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	// Protected routes (require authentication)
	s.router.Handle("/tiger/create", middleware.AuthMiddleware(auth, http.HandlerFunc(handlers.CreateTigerHandler))).Methods("POST")
	s.router.HandleFunc("/tigers", handlers.GetAllTigersHandler).Methods("GET")
	s.router.Handle("/tiger-sighting/create", middleware.AuthMiddleware(auth, http.HandlerFunc(handlers.CreateTigerSightingHandler))).Methods("POST")
	s.router.HandleFunc("/tiger/{id}/sightings", handlers.GetAllTigerSightingsHandler).Methods("GET")
}

func (s *Server) Start(port string) error {
	s.logger.Printf("Starting server on port %s...", port)
	return http.ListenAndServe(":"+port, s.router)
}
