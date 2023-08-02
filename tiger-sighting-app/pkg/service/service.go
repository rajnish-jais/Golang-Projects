package service

import (
	"errors"
	"log"
	"sort"

	"tiger-sighting-app/pkg/auth"
	"tiger-sighting-app/pkg/messaging"
	"tiger-sighting-app/pkg/models"
	"tiger-sighting-app/pkg/repository"
	"tiger-sighting-app/pkg/utility"
)

type service struct {
	TigerRepo     repository.TigerRepository
	messageBroker *messaging.MessageBroker
}

func NewTigerService(tigerRepository repository.TigerRepository, broker *messaging.MessageBroker) TigerService {
	return service{
		TigerRepo:     tigerRepository,
		messageBroker: broker,
	}
}

type TigerService interface {
	SignupService(*models.User) error
	LoginService(models.LoginCredentials) (*models.User, error)
	CreateTigerService(tiger models.Tiger) error
	GetAllTigersService() ([]*models.Tiger, error)
	CreateTigerSightingService(*models.TigerSighting) error
	GetAllTigerSightingsService(int) ([]*models.TigerSighting, error)
}

func (s service) SignupService(user *models.User) error {
	// Hash the user's password before saving to the database
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = hashedPassword

	// Create the user in the database
	if err := s.TigerRepo.CreateUser(user); err != nil {
		return errors.New("failed to create user")
	}
	return err
}

func (s service) LoginService(credentials models.LoginCredentials) (*models.User, error) {
	// Find the user by email in the database
	user, err := s.TigerRepo.GetUserByEmail(credentials.Email)
	if err != nil {
		return &models.User{}, errors.New("invalid email or password")
	}

	// Verify the password
	if err := auth.VerifyPassword(user.Password, credentials.Password); err != nil {
		return &models.User{}, errors.New("invalid email or password")
	}
	return user, nil
}

func (s service) CreateTigerService(tiger models.Tiger) error {
	// Create the tiger in the database
	if err := s.TigerRepo.CreateTiger(&tiger); err != nil {
		return errors.New("failed to create tiger")
	}
	return nil
}

func (s service) GetAllTigersService() ([]*models.Tiger, error) {
	// Get a list of all tigers from the database
	tigers, err := s.TigerRepo.GetAllTigers()
	if err != nil {
		return []*models.Tiger{}, errors.New("failed to fetch tigers")
	}

	// Sort the tigers by the last seen time (if the last seen time is a time.Time field)
	sort.Slice(tigers, func(i, j int) bool { return tigers[i].LastSeen.After(tigers[j].LastSeen) })
	return tigers, nil
}

func (s service) CreateTigerSightingService(newSighting *models.TigerSighting) error {
	// Check if the tiger has a previous sighting
	previousSighting, err := s.TigerRepo.GetPreviousTigerSighting(newSighting.TigerID)
	if err != nil {
		return errors.New("failed to retrieve previous sighting")
	}

	// If there is a previous sighting, calculate the distance between the coordinates
	if previousSighting != nil {
		previousCoordinates := models.Coordinates{Lat: previousSighting.Lat, Long: previousSighting.Long}
		currentCoordinates := models.Coordinates{Lat: newSighting.Lat, Long: newSighting.Long}
		distance := utility.CalculateDistance(previousCoordinates, currentCoordinates)

		// If the distance is less than or equal to 5 kilometers, reject the new sighting
		if distance <= 5.0 {
			return errors.New("A tiger sighting within 5 kilometers already exists")
		}
	}

	// Check if the required fields are provided
	if newSighting.Lat == 0 || newSighting.Long == 0 || newSighting.Timestamp.IsZero() || newSighting.ReporterEmail == "" {
		return errors.New("latitude, longitude, timestamp and reporterEmail are required")
	}

	// Create the tiger sighting in the database
	err = s.TigerRepo.CreateTigerSighting(newSighting)
	if err != nil {
		return errors.New("failed to create tiger sighting")
	}

	previousSightings, err := s.TigerRepo.GetAllTigerSightings(newSighting.TigerID)
	if err != nil {
		return errors.New("failed to retrieve previous sightings")
	}

	// Publish a new tiger sighting message
	if err := s.messageBroker.PublishMessage(utility.GetMails(previousSightings)); err != nil {
		log.Printf("failed to publish message: %v", err)
	}

	return nil
}

func (s service) GetAllTigerSightingsService(tigerID int) ([]*models.TigerSighting, error) {
	// Get a list of all tiger sightings for the specific tiger from the database
	tigerSightings, err := s.TigerRepo.GetAllTigerSightings(tigerID)
	if err != nil {

		return []*models.TigerSighting{}, errors.New("failed to fetch tiger sightings")
	}

	// Sort the tiger sightings by date (if the timestamp is a time.Time field)
	sort.Slice(tigerSightings, func(i, j int) bool {
		return tigerSightings[i].Timestamp.After(tigerSightings[j].Timestamp)
	})
	return tigerSightings, nil
}
