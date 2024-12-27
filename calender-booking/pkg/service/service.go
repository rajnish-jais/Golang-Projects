package service

import (
	"calender-booking/pkg/repository"
	"time"
)

type service struct {
	BookRepo repository.BookingRepository
}

func NewUserService(bookRepo repository.BookingRepository) BookingService {
	return service{
		BookRepo: bookRepo,
	}
}

type BookingService interface {
	BookTimeSlot(bookingID int, startTime, endTime time.Time) (int, error)
	SuggestTimeSlot(meetingDuration time.Duration, startTime, endTime time.Time) (time.Time, time.Time, error)
}

func (s service) BookTimeSlot(bookingID int, startTime, endTime time.Time) (int, error) {
	bookingId, err := s.BookRepo.BookTimeSlot(bookingID, startTime, endTime)
	if err != nil {
		return 0, err
	}
	return bookingId, nil
}

func (s service) SuggestTimeSlot(meetingDuration time.Duration, startTime, endTime time.Time) (time.Time, time.Time, error) {
	return s.BookRepo.SuggestTimeSlot(meetingDuration, startTime, endTime)
}
