package queue

import "auth-workflow/internal/models"

var WorkflowStepQueue  chan *models.WorkflowStep

func CreateWorkflow()  {
	WorkflowStepQueue = make(chan *models.WorkflowStep, 100)
}