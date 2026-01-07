package models
import "time"

type Event struct {
	ID         int64
	EventID    string
	EventType  string
	Payload    []byte
	Status     string
	RetryCount int
	CreatedAt  time.Time
}