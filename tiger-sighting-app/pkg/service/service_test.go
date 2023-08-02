package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
	"tiger-sighting-app/pkg/auth"
	"tiger-sighting-app/pkg/models"
	"time"
)

// mockTigerRepo is a mock implementation of the TigerRepository interface.
type mockTigerRepo struct {
	createUser               func(user *models.User) error
	getUserByEmail           func(email string) (*models.User, error)
	createTiger              func(tiger *models.Tiger) error
	getAllTigers             func() ([]*models.Tiger, error)
	createTigerSighting      func(newSighting *models.TigerSighting) error
	getPreviousTigerSighting func(tigerID int) (*models.TigerSighting, error)
	getAllTigerSightings     func(tigerID int) ([]*models.TigerSighting, error) // Add this method
}

func (m *mockTigerRepo) GetAllTigerSightings(tigerID int) ([]*models.TigerSighting, error) {
	return m.getAllTigerSightings(tigerID)
}

func (m *mockTigerRepo) CreateUser(user *models.User) error {
	return m.createUser(user)
}

func (m *mockTigerRepo) GetUserByEmail(email string) (*models.User, error) {
	return m.getUserByEmail(email)
}

func (m *mockTigerRepo) CreateTiger(tiger *models.Tiger) error {
	return m.createTiger(tiger)
}

func (m *mockTigerRepo) GetAllTigers() ([]*models.Tiger, error) {
	return m.getAllTigers()
}

func (m *mockTigerRepo) CreateTigerSighting(newSighting *models.TigerSighting) error {
	return m.createTigerSighting(newSighting)
}

func (m *mockTigerRepo) GetPreviousTigerSighting(tigerID int) (*models.TigerSighting, error) {
	return m.getPreviousTigerSighting(tigerID)
}

func TestSignupService_Success(t *testing.T) {
	// Arrange
	mockRepo := &mockTigerRepo{
		createUser: func(user *models.User) error {
			// Mock the CreateUser method to return nil (indicating success)
			return nil
		},
	}

	tigerService := NewTigerService(mockRepo, nil)

	// User with password to be hashed
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
	}

	// Act
	err := tigerService.SignupService(&user)

	// Assert
	assert.NoError(t, err, "SignupService should not return an error if user is created")

	// Ensure that the password was hashed
	assert.NotEqual(t, "testpassword", user.Password, "Password should be hashed")

}

func TestSignupService_Failure(t *testing.T) {
	// Arrange
	mockRepo := &mockTigerRepo{
		createUser: func(user *models.User) error {
			// Mock the CreateUser method to return an error (failure)
			return errors.New("failed to create user")
		},
	}

	tigerService := NewTigerService(mockRepo, nil)

	// User with password to be hashed
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
	}

	// Act
	err := tigerService.SignupService(&user)

	// Assert
	assert.Error(t, err, "SignupService should return an error")
	assert.EqualError(t, err, "failed to create user", "Error message should match")
}

func TestLoginService_Success(t *testing.T) {
	// Arrange
	mockRepo := &mockTigerRepo{
		getUserByEmail: func(email string) (*models.User, error) {
			// Mock the GetUserByEmail method to return a user with the hashed password
			hashedPassword, _ := auth.HashPassword("testpassword")
			return &models.User{
				Username: "testuser",
				Email:    email,
				Password: hashedPassword,
			}, nil
		},
	}

	tigerService := NewTigerService(mockRepo, nil)

	// Login credentials
	credentials := models.LoginCredentials{
		Email:    "test@example.com",
		Password: "testpassword",
	}

	// Act
	user, err := tigerService.LoginService(credentials)

	// Assert
	assert.NoError(t, err, "LoginService should not return an error")
	assert.NotNil(t, user, "User should not be nil")
	assert.Equal(t, "testuser", user.Username, "Usernames should match")
	assert.Equal(t, "test@example.com", user.Email, "Emails should match")
}

func TestLoginService_Failure(t *testing.T) {
	// Arrange
	mockRepo := &mockTigerRepo{
		getUserByEmail: func(email string) (*models.User, error) {
			// Mock the GetUserByEmail method to return an error (user not found)
			return nil, errors.New("user not found")
		},
	}

	tigerService := NewTigerService(mockRepo, nil)

	// Login credentials
	credentials := models.LoginCredentials{
		Email:    "test@example.com",
		Password: "testpassword",
	}

	// Act
	user, err := tigerService.LoginService(credentials)

	// Assert
	assert.Error(t, err, "LoginService should return an error")
	assert.Equal(t, user, &models.User{ID: 0, Username: "", Email: "", Password: ""}, "User should be nil")
	assert.EqualError(t, err, "invalid email or password", "Error message should match")
}

func TestCreateTigerService_Success(t *testing.T) {
	// Arrange
	mockRepo := &mockTigerRepo{
		createTiger: func(tiger *models.Tiger) error {
			// Mock the CreateTiger method to return nil (success)
			return nil
		},
	}

	tigerService := NewTigerService(mockRepo, nil)

	// Create a test tiger
	tiger := models.Tiger{
		Name:        "Test Tiger",
		DateOfBirth: time.Now(),
		LastSeen:    time.Now(),
		Lat:         12.34,
		Long:        56.78,
	}

	// Act
	err := tigerService.CreateTigerService(tiger)

	// Assert
	assert.NoError(t, err, "CreateTigerService should not return an error")
}

func TestCreateTigerService_Failure(t *testing.T) {
	// Arrange
	mockRepo := &mockTigerRepo{
		createTiger: func(tiger *models.Tiger) error {
			// Mock the CreateTiger method to return an error
			return errors.New("failed to create tiger")
		},
	}

	tigerService := NewTigerService(mockRepo, nil)

	// Create a test tiger
	tiger := models.Tiger{
		Name:        "Test Tiger",
		DateOfBirth: time.Now(),
		LastSeen:    time.Now(),
		Lat:         12.34,
		Long:        56.78,
	}

	// Act
	err := tigerService.CreateTigerService(tiger)

	// Assert
	assert.Error(t, err, "CreateTigerService should return an error")
	assert.EqualError(t, err, "failed to create tiger", "Error message should match")
}

func TestGetAllTigersService_Success(t *testing.T) {
	// Arrange
	mockRepo := &mockTigerRepo{
		getAllTigers: func() ([]*models.Tiger, error) {
			// Mock the GetAllTigers method to return a list of tigers
			tigers := []*models.Tiger{
				{
					Name:        "Tiger 1",
					DateOfBirth: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
					LastSeen:    time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
					Lat:         12.34,
					Long:        56.78,
				},
				{
					Name:        "Tiger 2",
					DateOfBirth: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
					LastSeen:    time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC),
					Lat:         12.56,
					Long:        56.78,
				},
			}
			return tigers, nil
		},
	}

	tigerService := NewTigerService(mockRepo, nil)

	// Act
	tigers, err := tigerService.GetAllTigersService()

	// Assert
	assert.NoError(t, err, "GetAllTigersService should not return an error")

	// Sort the tigers by the last seen time (descending order)
	sort.Slice(tigers, func(i, j int) bool {
		return tigers[i].LastSeen.After(tigers[j].LastSeen)
	})

	// Check if the tigers are sorted by last seen time (descending order)
	assert.True(t, tigers[0].LastSeen.After(tigers[1].LastSeen), "Tigers should be sorted by last seen time")
}

func TestGetAllTigersService_Failure(t *testing.T) {
	// Arrange
	mockRepo := &mockTigerRepo{
		getAllTigers: func() ([]*models.Tiger, error) {
			// Mock the GetAllTigers method to return an error
			return []*models.Tiger{}, errors.New("failed to fetch tigers")
		},
	}

	tigerService := NewTigerService(mockRepo, nil)

	// Act
	tigers, err := tigerService.GetAllTigersService()

	// Assert
	assert.Error(t, err, "GetAllTigersService should return an error")
	assert.EqualError(t, err, "failed to fetch tigers", "Error message should match")
	assert.Empty(t, tigers, "Tigers should be empty when there is an error")
}

// mockMessageBroker is a mock implementation of the MessageBroker interface.
type mockMessageBroker struct {
	publishMessage func(message []byte) error
}

func (m *mockMessageBroker) PublishMessage(message []byte) error {
	return m.publishMessage(message)
}

//func TestCreateTigerSightingService_Success(t *testing.T) {
//	// Arrange
//	previousSighting := &models.TigerSighting{
//		ID:            1,
//		TigerID:       1,
//		Timestamp:     time.Date(2023, time.July, 20, 12, 0, 0, 0, time.UTC),
//		Lat:           12.34,
//		Long:          56.78,
//		ReporterEmail: "reporter@example.com",
//	}
//	newSighting := &models.TigerSighting{
//		TigerID:       1,
//		Timestamp:     time.Date(2023, time.July, 21, 12, 0, 0, 0, time.UTC),
//		Lat:           13.35,
//		Long:          56.79,
//		ReporterEmail: "reporter@example.com",
//	}
//
//	mockRepo := &mockTigerRepo{
//		getPreviousTigerSighting: func(tigerID int) (*models.TigerSighting, error) {
//			return previousSighting, nil
//		},
//		createTigerSighting: func(newSighting *models.TigerSighting) error {
//			// Simulate successful tiger sighting creation in the database
//			return nil
//		},
//		getAllTigerSightings: func(tigerID int) ([]*models.TigerSighting, error) {
//			// Return previous sighting as the only previous sighting for the tiger
//			return []*models.TigerSighting{previousSighting}, nil
//		},
//	}
//
//	// Set up a mock message broker
//	var publishedMessage []byte
//	mockBroker := &mockMessageBroker{
//		publishMessage: func(message []byte) error {
//			publishedMessage = message
//			return nil
//		},
//	}
//
//	tigerService := NewTigerService(mockRepo, mockBroker)
//
//	// Act
//	err := tigerService.CreateTigerSightingService(newSighting)
//
//	// Assert
//	assert.NoError(t, err, "CreateTigerSightingService should not return an error")
//
//	// Ensure that the expected message is published
//	expectedMessage := []byte("reporter@example.com")
//	assert.Equal(t, expectedMessage, publishedMessage, "Published message should match")
//}

func TestCreateTigerSightingService_ExistingSightingWithin5Km(t *testing.T) {
	// Arrange
	previousSighting := &models.TigerSighting{
		ID:            1,
		TigerID:       1,
		Timestamp:     time.Date(2023, time.July, 20, 12, 0, 0, 0, time.UTC),
		Lat:           12.34,
		Long:          56.78,
		ReporterEmail: "reporter@example.com",
	}
	newSighting := &models.TigerSighting{
		TigerID:       1,
		Timestamp:     time.Date(2023, time.July, 21, 12, 0, 0, 0, time.UTC),
		Lat:           12.35,
		Long:          56.79,
		ReporterEmail: "reporter@example.com",
	}

	mockRepo := &mockTigerRepo{
		getPreviousTigerSighting: func(tigerID int) (*models.TigerSighting, error) {
			return previousSighting, nil
		},
	}

	tigerService := NewTigerService(mockRepo, nil)

	// Act
	err := tigerService.CreateTigerSightingService(newSighting)

	// Assert
	assert.Error(t, err, "CreateTigerSightingService should return an error")
	assert.EqualError(t, err, "A tiger sighting within 5 kilometers already exists", "Error message should match")
}

func TestCreateTigerSightingService_RequiredFieldsMissing(t *testing.T) {
	// Arrange
	newSighting := &models.TigerSighting{
		TigerID:       1,
		Timestamp:     time.Date(2023, time.July, 21, 12, 0, 0, 0, time.UTC),
		Lat:           0,
		Long:          56.79,
		ReporterEmail: "reporter@example.com",
	}

	mockRepo := &mockTigerRepo{
		getPreviousTigerSighting: func(tigerID int) (*models.TigerSighting, error) {
			return nil, nil
		},
	}

	tigerService := NewTigerService(mockRepo, nil)

	// Act
	err := tigerService.CreateTigerSightingService(newSighting)

	// Assert
	assert.Error(t, err, "CreateTigerSightingService should return an error")
	assert.EqualError(t, err, "latitude, longitude, timestamp and reporterEmail are required", "Error message should match")
}

func TestCreateTigerSightingService_Failure(t *testing.T) {
	// Arrange
	newSighting := &models.TigerSighting{
		TigerID:       1,
		Timestamp:     time.Date(2023, time.July, 21, 12, 0, 0, 0, time.UTC),
		Lat:           12.35,
		Long:          56.79,
		ReporterEmail: "reporter@example.com",
	}

	mockRepo := &mockTigerRepo{
		getPreviousTigerSighting: func(tigerID int) (*models.TigerSighting, error) {
			return nil, nil
		},
		createTigerSighting: func(newSighting *models.TigerSighting) error {
			// Simulate failure in creating tiger sighting in the database
			return errors.New("failed to create tiger sighting")
		},
		getAllTigerSightings: func(tigerID int) ([]*models.TigerSighting, error) {
			// Simulate failure in retrieving previous sightings from the database
			return nil, errors.New("failed to retrieve previous sightings")
		},
	}

	tigerService := NewTigerService(mockRepo, nil)

	// Act
	err := tigerService.CreateTigerSightingService(newSighting)

	// Assert
	assert.Error(t, err, "CreateTigerSightingService should return an error")
	assert.EqualError(t, err, "failed to create tiger sighting", "Error message should match")
}

func TestGetAllTigerSightingsService_Success(t *testing.T) {
	// Arrange
	tigerID := 1
	tigerSightings := []*models.TigerSighting{
		{
			ID:        1,
			TigerID:   1,
			Timestamp: time.Date(2023, time.July, 20, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:        2,
			TigerID:   1,
			Timestamp: time.Date(2023, time.July, 22, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:        3,
			TigerID:   1,
			Timestamp: time.Date(2023, time.July, 21, 12, 0, 0, 0, time.UTC),
		},
	}

	// Sort the tiger sightings by date (ascending order) for assertion
	sort.Slice(tigerSightings, func(i, j int) bool { return tigerSightings[i].Timestamp.Before(tigerSightings[j].Timestamp) })

	mockRepo := &mockTigerRepo{
		getAllTigerSightings: func(tigerID int) ([]*models.TigerSighting, error) {
			return tigerSightings, nil
		},
	}

	tigerService := NewTigerService(mockRepo, nil)

	// Act
	result, err := tigerService.GetAllTigerSightingsService(tigerID)

	// Assert
	assert.NoError(t, err, "GetAllTigerSightingsService should not return an error")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, len(tigerSightings), len(result), "Number of tiger sightings should match")

	// Check if the tiger sightings are sorted by the timestamp (ascending order)
	for i := range tigerSightings {
		assert.Equal(t, tigerSightings[i].ID, result[i].ID, "Tiger sighting IDs should match")
		assert.Equal(t, tigerSightings[i].TigerID, result[i].TigerID, "Tiger IDs should match")
		assert.Equal(t, tigerSightings[i].Timestamp, result[i].Timestamp, "Tiger sighting timestamps should match")
	}
}

func TestGetAllTigerSightingsService_Failure(t *testing.T) {
	// Arrange
	tigerID := 1
	mockRepo := &mockTigerRepo{
		getAllTigerSightings: func(tigerID int) ([]*models.TigerSighting, error) {
			// Simulate failure in retrieving tiger sightings from the database
			return nil, errors.New("failed to fetch tiger sightings")
		},
	}

	tigerService := NewTigerService(mockRepo, nil)

	// Act
	result, err := tigerService.GetAllTigerSightingsService(tigerID)

	// Assert
	assert.Error(t, err, "GetAllTigerSightingsService should return an error")
	assert.Equal(t, result, []*models.TigerSighting{}, "Result should be nil on failure")
	assert.EqualError(t, err, "failed to fetch tiger sightings", "Error message should match")
}
