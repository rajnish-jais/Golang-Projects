package models

import "time"

type TigerSighting struct {
	ID            int       `json:"id"`
	TigerID       int       `json:"tigerID"`
	Timestamp     time.Time `json:"timestamp"`
	Coordinates   `json:"coordinates"`
	Image         []byte `json:"image"`
	ReporterEmail string `json:"reporterEmail"`
}
