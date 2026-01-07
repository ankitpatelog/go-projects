package workers

import (
	"log"
	"time"

	"event-processing-service/internal/config"
	"event-processing-service/internal/queue"
	"event-processing-service/internal/repository"
)

func StartWorker(repo *repository.EventRepository) {
	go func() {
		for event := range queue.EventQueue {

			log.Println("Processing event:", event.EventID)

			// mark as processing
			repo.UpdateStatus(event.EventID, "PROCESSING")

			// execute handler
			err := handleEvent(event.EventType)
			if err != nil {

				retryKey := "retry:event:" + event.EventID

				// atomic retry increment
				retryCount, err := config.Redis.Incr(config.Ctx, retryKey).Result()
				if err != nil {
					log.Println("Redis error:", err)
					continue
				}

				// set TTL only on first retry
				if retryCount == 1 {
					config.Redis.Expire(config.Ctx, retryKey, time.Hour)
				}

				// max retry exceeded → FAIL
				if retryCount > 10 {
					repo.UpdateStatus(event.EventID, "FAILED")
					log.Println("Event failed permanently:", event.EventID)
					continue
				}

				// retry allowed → requeue
				log.Println("Retrying event:", event.EventID, "count:", retryCount)
				queue.EventQueue <- event
				continue
			}

			// SUCCESS
			repo.UpdateStatus(event.EventID, "SUCCESS")
			MarkProcessed(event.EventID)
		}
	}()
}


func handleEvent(eventType string) error {
	// simulate success/failure
	return nil
}

func MarkProcessed(eventID string) {
	key := "idem:event:" + eventID
	config.Redis.Set(config.Ctx, key, "done", time.Hour)
}

