package models

import (
	"time"
)

type TigerSighting struct {
	ID        int       `json:"id"`
	TigerID   int       `json:"tiger_id"`
	Timestamp time.Time `json:"timestamp"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Image     []byte    `json:"image_url"`
}
