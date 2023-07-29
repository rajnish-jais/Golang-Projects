package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/disintegration/imaging"
	"golang.org/x/crypto/bcrypt"
	"image"
	"log"
	"net/http"
	"sort"
	"strconv"
	"tiger-sighting-app/pkg/auth"
	"tiger-sighting-app/pkg/models"
	"tiger-sighting-app/pkg/repository"
	"time"
)

type Handlers struct {
	TigerRepo repository.TigerRepository
	Logger    *log.Logger
	Auth      *auth.Auth
}

func NewHandlers(tigerRepo repository.TigerRepository, logger *log.Logger, auth *auth.Auth) *Handlers {
	return &Handlers{
		TigerRepo: tigerRepo,
		Logger:    logger,
		Auth:      auth,
	}
}

func (h *Handlers) SignupHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get user data
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Validate user data (e.g., check for required fields, email format, password strength, etc.)
	if err := validateUserData(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash the user's password before saving to the database
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	// Create the user in the database
	if err := h.TigerRepo.CreateUser(&user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusCreated)
}

func validateUserData(user models.User) error {
	if user.Username == "" {
		return errors.New("username is required")
	} else if user.Email == "" {
		return errors.New("email is required")
	} else if user.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func hashPassword(password string) (string, error) {
	// Generate the hash of the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Convert the hashed password to a string and return it
	return string(hashedPassword), nil
}

func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get login credentials
	var loginCredentials models.LoginCredentials
	if err := json.NewDecoder(r.Body).Decode(&loginCredentials); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Find the user by email in the database
	user, err := h.TigerRepo.GetUserByEmail(loginCredentials.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Verify the password
	if err := verifyPassword(user.Password, loginCredentials.Password); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := h.Auth.GenerateToken(user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Respond with the token as JSON
	response := map[string]string{
		"token": token,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		return
	}
}

func verifyPassword(hashedPassword, plainPassword string) error {
	// Compare the hashed password with the plain-text password using bcrypt.CompareHashAndPassword
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		// If the comparison fails, return an error indicating an invalid password
		return err
	}

	// If the comparison succeeds, return nil (no error)
	return nil
}

func (h *Handlers) CreateTigerHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get tiger data
	var tiger models.Tiger
	if err := json.NewDecoder(r.Body).Decode(&tiger); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Create the tiger in the database
	if err := h.TigerRepo.CreateTiger(&tiger); err != nil {
		http.Error(w, "Failed to create tiger", http.StatusInternalServerError)
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) GetAllTigersHandler(w http.ResponseWriter, r *http.Request) {
	// Get a list of all tigers from the database
	tigers, err := h.TigerRepo.GetAllTigers()
	if err != nil {
		http.Error(w, "Failed to fetch tigers", http.StatusInternalServerError)
		return
	}

	// Sort the tigers by the last seen time (if the last seen time is a time.Time field)
	sort.Slice(tigers, func(i, j int) bool {
		return tigers[i].LastSeen.After(tigers[j].LastSeen)
	})

	// Respond with the tigers as JSON
	if err := json.NewEncoder(w).Encode(tigers); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) CreateTigerSightingHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the tiger sighting data
	var newSighting models.TigerSighting
	err := json.NewDecoder(r.Body).Decode(&newSighting)
	if err != nil {
		h.Logger.Printf("Error decoding request body: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Check if the required fields are provided
	if newSighting.Latitude == 0 || newSighting.Longitude == 0 || newSighting.Timestamp.IsZero() {
		respondWithError(w, http.StatusBadRequest, "Latitude, Longitude, and Timestamp are required")
		return
	}

	// Resize the image to 250x200 (not implemented here, you can use external libraries like imaging)
	resizedImage, err := resizeImage(newSighting.Image, 250, 200)
	if err != nil {
		h.Logger.Printf("Error resizing image: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error resizing image")
		return
	}

	// Create the tiger sighting in the database
	err = h.TigerRepo.CreateTigerSighting(&newSighting, resizedImage)
	if err != nil {
		h.Logger.Printf("Error creating tiger sighting: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create tiger sighting")
		return
	}

	// Respond with the newly created tiger sighting
	respondWithJSON(w, http.StatusCreated, newSighting)
}

func resizeImage(imageBytes []byte, width, height int) ([]byte, error) {
	// Decode the imageBytes into an image.Image
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}

	// Resize the image using the Lanczos filter
	resizedImg := imaging.Resize(img, width, height, imaging.Lanczos)

	// Encode the resized image back to bytes
	var buf bytes.Buffer
	err = imaging.Encode(&buf, resizedImg, imaging.JPEG)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (h *Handlers) GetAllTigerSightingsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the tiger ID from the URL query parameter
	tigerID := r.URL.Query().Get("tiger_id")
	if tigerID == "" {
		http.Error(w, "Missing tiger_id query parameter", http.StatusBadRequest)
		return
	}

	// Convert the tiger ID to an integer
	tigerIDInt, err := strconv.Atoi(tigerID)
	if err != nil {
		http.Error(w, "Invalid tiger_id query parameter", http.StatusBadRequest)
		return
	}

	// Get a list of all tiger sightings for the specific tiger from the database
	tigerSightings, err := h.TigerRepo.GetAllTigerSightings(tigerIDInt)
	if err != nil {
		http.Error(w, "Failed to fetch tiger sightings", http.StatusInternalServerError)
		return
	}

	// Sort the tiger sightings by date (if the timestamp is a time.Time field)
	sort.Slice(tigerSightings, func(i, j int) bool {
		return tigerSightings[i].Timestamp.After(tigerSightings[j].Timestamp)
	})

	// Respond with the tiger sightings as JSON
	if err := json.NewEncoder(w).Encode(tigerSightings); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		return
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	response := map[string]string{"error": message}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResponse)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	jsonResponse, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResponse)
}

func (h *Handlers) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		start := time.Now()
		next.ServeHTTP(w, r)
		h.Logger.Printf("Request processed in %v", time.Since(start))
	})
}
