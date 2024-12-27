package repository

import (
	"calender-booking/pkg/repository/store"
	"time"
)

func NewPostgresRepository(connection string) (BookingRepository, error) { // paas the model interface
	db, err := store.NewPostgresDB(connection)
	return store.NewPostgresRepository(db), err
}

type BookingRepository interface {
	BookTimeSlot(bookingID int, startTime, endTime time.Time) (int, error)
	SuggestTimeSlot(meetingDuration time.Duration, startTime, endTime time.Time) (time.Time, time.Time, error)
}
