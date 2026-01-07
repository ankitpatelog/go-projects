package services

import (
	"encoding/json"
	"event-processing-service/internal/config"
	"event-processing-service/internal/models"
	"event-processing-service/internal/repository"
	"event-processing-service/internal/queue"
	
)

type Eventservice struct {
	Repo *repository.EventRepository
}

func (s *Eventservice)ProcessEvent(eventID,eventType string,payload interface{})error  {
	
	// redis mai event jo success ho chuk hai uska data store hai
	key := "idem:event" + eventID

	exists,_ := config.Redis.Exists(config.Ctx,key).Result()
	if exists==1 {
		//event alredy processed qith success
		return nil
	}

	body, _ := json.Marshal(payload)

	event := &models.Event{
		EventID:    eventID,
		EventType:  eventType,
		Payload:    body,
		Status:     "PENDING",
		RetryCount: 0,
	}
	
	//save to mysql
	err := s.Repo.Save(event)
	if err !=nil {
		return err
	}

	//push to queue for workers to process
	queue.EventQueue <- event
	return nil
}

