package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"tiger-sighting-app/pkg/auth"
	"tiger-sighting-app/pkg/models"
	"tiger-sighting-app/pkg/service"
	"tiger-sighting-app/pkg/utility"
)

type handlers struct {
	Auth         *auth.Auth
	Logger       *log.Logger
	TigerService service.TigerService
}

func NewHandlers(tigerService service.TigerService, logger *log.Logger, auth *auth.Auth) *handlers {
	return &handlers{
		Auth:         auth,
		Logger:       logger,
		TigerService: tigerService,
	}
}

func (h *handlers) SignupHandler(w http.ResponseWriter, r *http.Request) {
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

	err := h.TigerService.SignupService(&user)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Respond with success status
	utility.RespondWithJSON(w, http.StatusCreated, "success")
}

func (h *handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get login credentials
	var loginCredentials models.LoginCredentials
	if err := json.NewDecoder(r.Body).Decode(&loginCredentials); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}

	user, err := h.TigerService.LoginService(loginCredentials)
	if err != nil {
		utility.RespondWithError(w, http.StatusUnauthorized, err.Error())
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

func (h *handlers) CreateTigerHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get tiger data
	var tiger models.Tiger
	if err := json.NewDecoder(r.Body).Decode(&tiger); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}

	err := h.TigerService.CreateTigerService(tiger)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	// Respond with success status
	utility.RespondWithJSON(w, http.StatusCreated, "success")
}

func (h *handlers) GetAllTigersHandler(w http.ResponseWriter, r *http.Request) {
	tigers, err := h.TigerService.GetAllTigersService()
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	// Respond with the tigers as JSON
	utility.RespondWithJSON(w, http.StatusOK, tigers)
}

func (h *handlers) CreateTigerSightingHandler(w http.ResponseWriter, r *http.Request) {
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

	reporterEmail, ok := auth.GetEmailFromContext(r.Context())
	if !ok {
		utility.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve previous sighting")
		return
	}

	newSighting := models.TigerSighting{TigerID: tigerID, Timestamp: timestamp, Lat: lat, Long: long, ReporterEmail: reporterEmail}

	imageFile, _, err := r.FormFile("image")
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Failed to get image file")
		return
	}
	defer imageFile.Close()

	resizedImage, err := getProcessedImage(imageFile)
	if err != nil {
		h.Logger.Printf("Got Error Resizing Image: %v", err)
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	newSighting.Image = resizedImage
	err = h.TigerService.CreateTigerSightingService(&newSighting)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	// Respond with the newly created tiger sighting
	utility.RespondWithJSON(w, http.StatusCreated, newSighting)
}

func getProcessedImage(imageFile multipart.File) ([]byte, error) {
	// Read the image data into a byte slice
	imageData, err := ioutil.ReadAll(imageFile)
	if err != nil {
		return nil, err
	}

	// Resize the image to 250x200
	resizedImage, err := utility.ResizeImage(imageData, 250, 200)
	if err != nil {
		return nil, err
	}
	return resizedImage, nil
}

func (h *handlers) GetAllTigerSightingsHandler(w http.ResponseWriter, r *http.Request) {
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

	tigerSightings, err := h.TigerService.GetAllTigerSightingsService(tigerIDInt)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	// Respond with the tiger sightings as JSON
	utility.RespondWithJSON(w, http.StatusOK, tigerSightings)
}
