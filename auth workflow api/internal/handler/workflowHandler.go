package handler

import (
	"auth-workflow/internal/service"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Workflowhandler struct {
	handler *service.WorkflowService
}

func (h *Workflowhandler)Paymenthandler(w http.ResponseWriter,r *http.Request)  {
	//create workflow into db
	workflowID := uuid.NewString()
	workflow ,err := h.handler.ProcessPayment(workflowID)
	if err!=nil {
		http.Error(w,"Workflow saving error",http.StatusInternalServerError)
		return
	}
	//make workflow steps for payment
	steps := []string{
		"VERIFY_PAYMENT",
		"CONFIRM_ORDER",
		"UPDATE_INVENTORY",
		"SEND_NOTIFICATION",
	}

	//genearte each step into database
	err = h.handler.StepHandler(steps,workflowID)
	if err!=nil {
		http.Error(w,"error in saving workflow_step",http.StatusInternalServerError)
	}
	
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(workflow)

}