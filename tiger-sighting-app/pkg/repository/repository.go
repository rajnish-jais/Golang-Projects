package repository

import (
	"database/sql"
	"tiger-sighting-app/pkg/repository/store"
)

type TigerRepository interface {
	//CreateUser(user *models.User) error
	//CreateTiger(tiger *models.Tiger) error
	//GetAllTigers() ([]*models.Tiger, error)
}

func NewDBConnection(db *sql.DB) TigerRepository {
	return store.NewPostgresRepository(db)
}
