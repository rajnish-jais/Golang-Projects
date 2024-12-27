package store

import (
	"database/sql"
	"fmt"
	"time"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *postgresRepository {
	return &postgresRepository{db: db}
}

func (p *postgresRepository) BookTimeSlot(bookingID int, startTime, endTime time.Time) (int, error) {
	query := `SELECT COUNT(*)
 	FROM bookings
 	WHERE (start_time < $2 and end_time > $1)`

	var count int
	err := p.db.QueryRow(query, startTime, endTime).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error in conflict finding: %v", err)
	}

	if count > 0 {
		return 0, fmt.Errorf("%v conflict found", count)
	}

	query = `INSERT INTO bookings (start_time, end_time)
		VALUES($1, $2)
		RETURNING id`
	err = p.db.QueryRow(query, startTime, endTime).Scan(&bookingID)

	if err != nil {
		return 0, fmt.Errorf("error in inserting booking: %v", err)
	}

	return bookingID, nil
}

func (p *postgresRepository) SuggestTimeSlot(meetingDuration time.Duration, startTime, endTime time.Time) (time.Time, time.Time, error) {
	// Fetch booked slots within range
	query := ` SELECT start_time, end_time FROM bookings WHERE start_time < $2 AND end_time > $1 ORDER BY start_time `
	rows, err := p.db.Query(query, startTime, endTime)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("error fetching bookings %v", err)
	}

	defer rows.Close()

	// Find available slot
	var lastEndTime = startTime
	for rows.Next() {
		var startTime, endTime time.Time
		err := rows.Scan(&startTime, &endTime)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("error scanning bookings %v", err)
		}
		// Check gap between last end time and current start time
		if startTime.Sub(lastEndTime) >= time.Duration(meetingDuration)*time.Minute {
			// Suggest this slot
			suggestedStart := lastEndTime
			suggestedEnd := lastEndTime.Add(time.Duration(meetingDuration) * time.Minute)
			return suggestedStart, suggestedEnd, nil
		}
		lastEndTime = endTime
	}

	// Check if there's space at the end
	if endTime.Sub(lastEndTime) >= time.Duration(meetingDuration)*time.Minute {
		suggestedStart := lastEndTime
		suggestedEnd := lastEndTime.Add(time.Duration(meetingDuration) * time.Minute)

		return suggestedStart, suggestedEnd, nil
	}
	return time.Time{}, time.Time{}, fmt.Errorf("no available slot found")
}
