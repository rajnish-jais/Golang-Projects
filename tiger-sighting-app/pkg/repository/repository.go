package repository

import (
	"tiger-sighting-app/pkg/models"
	"tiger-sighting-app/pkg/repository/store"
)

type TigerRepository interface {
	CreateUser(user *models.User) error
	CreateTiger(tiger *models.Tiger) error
	GetAllTigers() ([]*models.Tiger, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateTigerSighting(tigerSighting *models.TigerSighting, resizedImage []byte) error
	GetAllTigerSightings(tigerID int) ([]models.TigerSighting, error)
}

func NewPostgresRepository(connection string) (TigerRepository, error) {
	db, err := store.NewPostgresDB(connection)
	return store.NewPostgresRepository(db), err
}
