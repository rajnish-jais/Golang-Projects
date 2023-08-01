package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"tiger-sighting-app/pkg/auth"
	"tiger-sighting-app/pkg/messaging"
	"tiger-sighting-app/pkg/models"
	"tiger-sighting-app/pkg/repository"
	"tiger-sighting-app/pkg/utility"
)

type Handlers struct {
	Auth          *auth.Auth
	Logger        *log.Logger
	TigerRepo     repository.TigerRepository
	messageBroker *messaging.MessageBroker
}

func NewHandlers(tigerRepo repository.TigerRepository, messageBroker *messaging.MessageBroker, logger *log.Logger, auth *auth.Auth) *Handlers {
	return &Handlers{
		Auth:          auth,
		Logger:        logger,
		TigerRepo:     tigerRepo,
		messageBroker: messageBroker,
	}
}

func (h *Handlers) SignupHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get user data
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}

	// Validate user data
	if err := auth.ValidateUserData(user); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Hash the user's password before saving to the database
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	user.Password = hashedPassword

	// Create the user in the database
	if err := h.TigerRepo.CreateUser(&user); err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Respond with success status
	utility.RespondWithJSON(w, http.StatusCreated, "success")
}

func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get login credentials
	var loginCredentials models.LoginCredentials
	if err := json.NewDecoder(r.Body).Decode(&loginCredentials); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}

	// Find the user by email in the database
	user, err := h.TigerRepo.GetUserByEmail(loginCredentials.Email)
	if err != nil {
		utility.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Verify the password
	if err := auth.VerifyPassword(user.Password, loginCredentials.Password); err != nil {
		utility.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT token
	token, err := h.Auth.GenerateToken(user.Username, user.Email)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Respond with the token as JSON
	response := map[string]string{"token": token}
	utility.RespondWithJSON(w, http.StatusOK, response)
}

func (h *Handlers) CreateTigerHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get tiger data
	var tiger models.Tiger
	if err := json.NewDecoder(r.Body).Decode(&tiger); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}

	// Create the tiger in the database
	if err := h.TigerRepo.CreateTiger(&tiger); err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, "Failed to create tiger")
		return
	}

	// Respond with success status
	utility.RespondWithJSON(w, http.StatusCreated, "success")
}

func (h *Handlers) GetAllTigersHandler(w http.ResponseWriter, r *http.Request) {
	// Get a list of all tigers from the database
	tigers, err := h.TigerRepo.GetAllTigers()
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch tigers")
		return
	}

	// Sort the tigers by the last seen time (if the last seen time is a time.Time field)
	sort.Slice(tigers, func(i, j int) bool { return tigers[i].LastSeen.After(tigers[j].LastSeen) })

	// Respond with the tigers as JSON
	utility.RespondWithJSON(w, http.StatusOK, tigers)
}

func (h *Handlers) CreateTigerSightingHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the tiger sighting data
	if err := r.ParseMultipartForm(10 << 20); err != nil { // Max memory of 10 MB for file uploads
		utility.RespondWithError(w, http.StatusBadRequest, "Unable to parse form data")
		return
	}

	// Get the form values
	tigerIDStr := r.FormValue("tigerID")
	timestampStr := r.FormValue("timestamp")
	latStr := r.FormValue("lat")
	longStr := r.FormValue("long")

	// Convert the form values to appropriate types
	tigerID, err := strconv.Atoi(tigerIDStr)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid tigerID value")
		return
	}

	timestamp, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid timestamp value")
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid lat value")
		return
	}

	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid long value")
		return
	}

	// Check if the tiger has a previous sighting
	previousSighting, err := h.TigerRepo.GetPreviousTigerSighting(tigerID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve previous sighting")
		return
	}

	// If there is a previous sighting, calculate the distance between the coordinates
	if previousSighting != nil {
		previousCoordinates := models.Coordinates{Lat: previousSighting.Lat, Long: previousSighting.Long}
		currentCoordinates := models.Coordinates{Lat: lat, Long: long}
		distance := utility.CalculateDistance(previousCoordinates, currentCoordinates)

		// If the distance is less than or equal to 5 kilometers, reject the new sighting
		if distance <= 5.0 {
			utility.RespondWithError(w, http.StatusConflict, "A tiger sighting within 5 kilometers already exists")
			return
		}
	}

	reporterEmail, ok := auth.GetEmailFromContext(r.Context())
	if !ok {
		utility.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve previous sighting")
		return
	}

	newSighting := models.TigerSighting{
		TigerID:       tigerID,
		Timestamp:     timestamp,
		Coordinates:   models.Coordinates{Lat: lat, Long: long},
		ReporterEmail: reporterEmail,
	}

	// Check if the required fields are provided
	if newSighting.Lat == 0 || newSighting.Long == 0 || newSighting.Timestamp.IsZero() || newSighting.ReporterEmail == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Latitude, Longitude, Timestamp and ReporterEmail are required")
		return
	}

	imageFile, _, err := r.FormFile("image")
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Failed to get image file")
		return
	}
	defer imageFile.Close()

	resizedImage, err := getProcessedImage(imageFile, w)
	if err != nil {
		h.Logger.Printf("Got Error Resizing Image: %v", err)
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	newSighting.Image = resizedImage
	// Create the tiger sighting in the database
	err = h.TigerRepo.CreateTigerSighting(&newSighting)
	if err != nil {
		h.Logger.Printf("Error creating tiger sighting: %v", err)
		utility.RespondWithError(w, http.StatusInternalServerError, "Failed to create tiger sighting")
		return
	}

	previousSightings, err := h.TigerRepo.GetTigerSightingsByTigerID(tigerID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve previous sightings")
		return
	}

	// Publish a new tiger sighting message
	if err := h.messageBroker.PublishMessage(utility.GetMails(previousSightings)); err != nil {
		log.Printf("failed to publish message: %v", err)
	}

	// Respond with the newly created tiger sighting
	utility.RespondWithJSON(w, http.StatusCreated, newSighting)
}

func getProcessedImage(imageFile multipart.File, w http.ResponseWriter) ([]byte, error) {
	// Read the image data into a byte slice
	imageData, err := ioutil.ReadAll(imageFile)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, "Failed to read image data")
		return nil, err
	}

	// Resize the image to 250x200
	resizedImage, err := utility.ResizeImage(imageData, 250, 200)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, "Error resizing image")
		return nil, err
	}
	return resizedImage, nil
}

func (h *Handlers) GetAllTigerSightingsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tigerID := vars["id"]
	if tigerID == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Missing tiger_id query parameter")
		return
	}

	// Convert the tiger ID to an integer
	tigerIDInt, err := strconv.Atoi(tigerID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid tiger_id query parameter")
		return
	}

	// Get a list of all tiger sightings for the specific tiger from the database
	tigerSightings, err := h.TigerRepo.GetAllTigerSightings(tigerIDInt)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch tiger sightings")
		return
	}

	// Sort the tiger sightings by date (if the timestamp is a time.Time field)
	sort.Slice(tigerSightings, func(i, j int) bool {
		return tigerSightings[i].Timestamp.After(tigerSightings[j].Timestamp)
	})

	// Respond with the tiger sightings as JSON
	utility.RespondWithJSON(w, http.StatusOK, tigerSightings)
}
