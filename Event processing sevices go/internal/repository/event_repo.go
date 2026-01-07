package repository

import (
	"database/sql"
	"event-processing-service/internal/models"
)

//here to handle with data saving and updation with the database sql
type EventRepository struct {
	DB *sql.DB
}

func (r EventRepository)Save(task *models.Event)error{
	query := `INSERT INTO events (event_id, event_type, payload, status, retry_count)
	VALUES (?, ?, ?, ?, ?)`

	_,err := r.DB.Exec(query)
	return err

}

func (r *EventRepository) UpdateStatus(eventID, status string) error {
	_, err := r.DB.Exec(
		"UPDATE events SET status=? WHERE event_id=?",
		status, eventID,
	)
	return err
}

func (r *EventRepository)IncrementRetry(eventID string)error  {
	query := "UPDATE events SET retry_count = retry_count + 1 WHERE event_id=?"

	_,err := r.DB.Exec(query,eventID)
	return err
}
