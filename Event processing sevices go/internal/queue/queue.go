package queue

import "event-processing-service/internal/models"

var EventQueue chan *models.Event

func InitQueue()  {
	EventQueue = make(chan *models.Event,100)
}