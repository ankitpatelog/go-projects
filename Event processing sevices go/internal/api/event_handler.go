package api

import (
	"encoding/json"
	"net/http"

	"event-processing-service/internal/services"
)

type EventRequest struct {
	EventID   string      `json:"event_id"`
	EventType string      `json:"event_type"`
	Payload   interface{} `json:"payload"`
}

type EventHandler struct {
	Service *services.Eventservice
}

func (h *EventHandler) HandleEvent(w http.ResponseWriter, r *http.Request) {

	var req EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.Service.ProcessEvent(req.EventID, req.EventType, req.Payload)
	if err != nil {
		http.Error(w, "Failed to process event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"accepted"}`))
}
