package main

import (
	"fmt"
	"net/http"

	"auth-workflow/internal/config"
	"auth-workflow/internal/handler"
	"auth-workflow/internal/middleware"
	"auth-workflow/internal/queue"
	"auth-workflow/internal/worker"

	"github.com/gorilla/mux"
)

func main() {
	// DB & infra connections
	DB := config.Connectmysql()
	db := DB
	config.ConnectRedis()

	// init queue & workers
	queue.CreateWorkflow()
	worker.StartWorker()

	router := mux.NewRouter()

	// routes
	router.HandleFunc("/create", handler.NewHnadlerInstance(db).CreateUser).Methods("POST")
	router.HandleFunc("/login", handler.NewHnadlerInstance(db).LoginUser).Methods("POST")
	router.HandleFunc("/profile", handler.NewHnadlerInstance(db).Profile).Methods("GET")
	router.HandleFunc("/payment", handler.NewworkflowInstance(db).Paymenthandler).Methods("POST")
	router.HandleFunc("/ckeckstatus/{workflowid}", handler.NewworkflowInstance(db).CheckStatusWorkflow).Methods("POST")

	// middleware
	router.Use(middleware.RateLimmiter)

	fmt.Println("Server is listening at port: 8000")
	if err := http.ListenAndServe(":8000", router); err != nil {
		panic(err)
	}
}
