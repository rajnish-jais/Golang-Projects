package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"calender-booking/pkg/models"
	"calender-booking/pkg/service"
	"calender-booking/pkg/utils"
	"strconv"
)

type handlers struct {
	Logger      *log.Logger
	BookService service.BookingService
}

func NewHandlers(bookService service.BookingService, logger *log.Logger) *handlers {
	return &handlers{
		Logger:      logger,
		BookService: bookService,
	}
}

func (h *handlers) BookTimeSlot(w http.ResponseWriter, r *http.Request) {
	var booking models.Booking
	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Input")
		return
	}
	if booking.StartTime.After(booking.EndTime) {
		utils.RespondWithError(w, http.StatusBadRequest, "start time must be before end time")
		return
	}

	booking.ID, err = h.BookService.BookTimeSlot(booking.ID, booking.StartTime, booking.EndTime)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Booking Failed bcs:  %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Booking Success",
		"booking": booking,
	})
}

func (h *handlers) SuggestTimeSlot(w http.ResponseWriter, r *http.Request) {
	startTime := r.URL.Query().Get("start_time")
	endTime := r.URL.Query().Get("end_time")
	meetingDuration := r.URL.Query().Get("meeting_duration")

	// Parse and validate inputs
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		http.Error(w, "Invalid start_time format", http.StatusBadRequest)
		return
	}
	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		http.Error(w, "Invalid end_time format", http.StatusBadRequest)
		return
	}

	duration, err := strconv.Atoi(meetingDuration)
	if err != nil || duration <= 0 {
		http.Error(w, "Invalid meeting_duration", http.StatusBadRequest)
		return
	}

	suggestedStart, suggestedEnd, err := h.BookService.SuggestTimeSlot(time.Duration(duration), start, end)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "No available time slot found")
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"suggested_start_time": suggestedStart,
		"suggested_end_time":   suggestedEnd,
	})
}
