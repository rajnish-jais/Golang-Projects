package store

import (
	"database/sql"
	"fmt"
	"tiger-sighting-app/pkg/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (p *PostgresRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
	`
	_, err := p.db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresRepository) CreateTiger(tiger *models.Tiger) error {
	query := `
		INSERT INTO tigers (name, date_of_birth, last_seen, lat, long)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := p.db.Exec(query, tiger.Name, tiger.DateOfBirth, tiger.LastSeen, tiger.Lat, tiger.Long)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresRepository) GetAllTigers() ([]*models.Tiger, error) {
	query := `
		SELECT id, name, date_of_birth, last_seen, lat, long
		FROM tigers
		ORDER BY last_seen DESC
	`
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tigers := []*models.Tiger{}
	for rows.Next() {
		tiger := &models.Tiger{}
		err := rows.Scan(&tiger.ID, &tiger.Name, &tiger.DateOfBirth, &tiger.LastSeen, &tiger.Lat, &tiger.Long)
		if err != nil {
			return nil, err
		}
		tigers = append(tigers, tiger)
	}

	return tigers, nil
}

func (p *PostgresRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
        SELECT id, username, email, password
        FROM users
        WHERE email = $1
    `

	user := &models.User{}
	err := p.db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (p *PostgresRepository) CreateTigerSighting(tigerSighting *models.TigerSighting, resizedImage []byte) error {

	query := `
        INSERT INTO tiger_sightings (tiger_id, timestamp, latitude, longitude, image_url)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	err := p.db.QueryRow(query, tigerSighting.TigerID, tigerSighting.Timestamp, tigerSighting.Latitude, tigerSighting.Longitude, tigerSighting.Image).Scan(&tigerSighting.ID)
	if err != nil {
		return fmt.Errorf("failed to create tiger sighting: %v", err)
	}

	return nil
}

func (p *PostgresRepository) GetAllTigerSightings(tigerID int) ([]models.TigerSighting, error) {
	query := `
        SELECT id, tiger_id, timestamp, latitude, longitude, image_url
        FROM tiger_sightings
        WHERE tiger_id = $1
        ORDER BY timestamp DESC
    `

	rows, err := p.db.Query(query, tigerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tiger sightings: %v", err)
	}
	defer rows.Close()

	sightings := []models.TigerSighting{}
	for rows.Next() {
		var sighting models.TigerSighting
		err := rows.Scan(&sighting.ID, &sighting.TigerID, &sighting.Timestamp, &sighting.Latitude, &sighting.Longitude, &sighting.Image)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tiger sighting: %v", err)
		}
		sightings = append(sightings, sighting)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error processing tiger sightings rows: %v", err)
	}

	return sightings, nil
}
