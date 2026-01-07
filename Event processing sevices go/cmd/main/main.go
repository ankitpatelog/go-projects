package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"event-processing-service/internal/api"
	"event-processing-service/internal/config"
	"event-processing-service/internal/middleware"
	"event-processing-service/internal/queue"
	"event-processing-service/internal/repository"
	"event-processing-service/internal/services"
	"event-processing-service/internal/workers"
)

func main() {

	// init config
	config.ConnectMySQL()
	config.ConnectRedis()

	// init queue
	queue.InitQueue()

	repo := &repository.EventRepository{DB: config.DB}
	service := &services.Eventservice{Repo: repo}
	handler := &api.EventHandler{Service: service}

	// start worker
	workers.StartWorker(repo)

	router := mux.NewRouter()
	router.Use(middleware.RateLimiter)

	router.HandleFunc("/events", handler.HandleEvent).Methods("POST")

	log.Println("🚀 Server running on :8080")
	http.ListenAndServe(":8080", router)
}
