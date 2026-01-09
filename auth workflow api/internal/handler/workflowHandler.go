package handler

import (
	"auth-workflow/internal/service"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func NewworkflowInstance(db *sql.DB)*Workflowhandler  {
	return &Workflowhandler{
		handler: service.NewserviceInstance(db),
	}
}

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

func (h *Workflowhandler)CheckStatusWorkflow(w http.ResponseWriter,r *http.Request)  {
	//get workflowid from param
	var id string
	vars := mux.Vars(r)

	id=vars["workflowid"]//workflowid same as user id

	workflows,err := h.handler.Processandgetworkflows(id)
	if err!=nil {
		http.Error(w,"User workflow not found",http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(workflows)
}