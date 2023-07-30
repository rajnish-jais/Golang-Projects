package models

import "time"

type TigerSighting struct {
	ID          int       `json:"id"`
	TigerID     int       `json:"tigerID"`
	Timestamp   time.Time `json:"timestamp"`
	Coordinates `json:"coordinates omitempty"`
	Image       []byte `json:"image, omitempty"`
}
