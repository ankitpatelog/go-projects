package repository

import (
	"database/sql"
	"net/http"

	"workflow-engine/internal/models"
)

type WorkflowRepository struct {
	DB *sql.DB
}

func (r WorkflowRepository) CreateWorkflow(workflow *models.WorkflowInstance)error  {
	query := `
		INSERT INTO workflow_instances
		(id, event_type, status, current_step, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		`

		query2 := `
		INSERT INTO workflow_step
		(id, event_type, status, current_step, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		`

		_, err := r.DB.Exec(
		query,
		workflow.ID,
		workflow.EventType,
		workflow.Status,
		workflow.CurrentStep,
		workflow.CreatedAt,
		workflow.UpdatedAt,
	    )
	
	if err != nil {
		return err
	}

	return nil
	
}

func (r WorkflowRepository) CreateWorkflow(workflow *models.WorkflowInstance,steps []models.WorkflowStep,) error {

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	// 1️⃣ Insert workflow instance
	workflowQuery := `
	INSERT INTO workflow_instances
	(id, event_type, status, current_step, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = tx.Exec(
		workflowQuery,
		workflow.ID,
		workflow.EventType,
		workflow.Status,
		workflow.CurrentStep,
		workflow.CreatedAt,
		workflow.UpdatedAt,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	//  Insert workflow steps
	stepQuery := `
	INSERT INTO workflow_steps
	(workflow_id, step_name, step_order, status, retry_count, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	for _, step := range steps {
		_, err := tx.Exec(
			stepQuery,
			workflow.ID,
			step.StepName,
			step.StepOrder,
			step.Status,
			step.RetryCount,
			step.CreatedAt,
			step.UpdatedAt,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	//  Commit transaction
	return tx.Commit()
}
