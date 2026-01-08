package service

import (
	"auth-workflow/internal/models"
	"auth-workflow/internal/repo_workflow"
	"auth-workflow/internal/queue"

	"github.com/google/uuid"
)

type WorkflowService struct {
	service *repo_workflow.WorkflowRepo
}

func (s *WorkflowService)ProcessPayment(workflowid string)(*models.WorkflowInstance,error)  {
	var work *models.WorkflowInstance
	work = &models.WorkflowInstance{
		ID: workflowid,
		Event: "Payment Event",
		Status: "Pending",
	}

	saved,err := s.service.SaveWorkflowInstance(work)
	if err!=nil {
		return nil,err
	}

	return saved,nil
}

func (s *WorkflowService)StepHandler(steps []string,workflowid string)(error)  {
	

	for _,step := range steps{
		//save every step into db

		currstep := &models.WorkflowStep{
			ID: uuid.NewString(),
			WorkflowID: workflowid,
			Name: step,
			Status: "Pending",
			RetryCount: 0,	
		}
			err := s.service.SaveWorkflowStep(currstep)

			//[push stepworflow into workflow queue]

			queue.WorkflowStepQueue <- currstep
			
		if err != nil {
			return err 
		}
	}
	return nil
}

