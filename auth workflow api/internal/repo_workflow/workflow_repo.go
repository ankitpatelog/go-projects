package repo_workflow

import (
	"auth-workflow/internal/models"
	"database/sql"
)

type WorkflowRepo struct {
	DB *sql.DB
}

func (r *WorkflowRepo)SaveWorkflowInstance(work *models.WorkflowInstance)(*models.WorkflowInstance,error)  {
	_, err := r.DB.Exec(
		`INSERT INTO workflow_instances (id, event, status)
		 VALUES (?, ?, ?)`,
		work.ID,
		work.Event,
		work.Status,
	)

	if err != nil {
		return nil, err
	}

	return work, nil
}

func (r *WorkflowRepo)SaveWorkflowStep(step *models.WorkflowStep)(error)  {
	_, err := r.DB.Exec(
		`INSERT INTO workflow_step (id, workflowid, name,status,retrycount)
		 VALUES (?, ?, ?, ?, ?)`,
		step.ID,
		step.WorkflowID,
		step.Name,
		step.Status,
		step.RetryCount,
	)

	if err != nil {
		return err
	}

	return nil
}
