package repo_workflow

import (
	"auth-workflow/internal/models"
	"database/sql"

	"github.com/pelletier/go-toml/query"
	"go.uber.org/mock/mockgen/model"
)

func NewrepoWorkflowInstance(db *sql.DB)*WorkflowRepo  {
	return &WorkflowRepo{
		DB: db,
	}
}

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

func UpdateStatus(db *sql.DB, id string, status string) (string, error) {
	query := `
		UPDATE workflow_step
		SET status = ?
		WHERE id = ?
	`

	_, err := db.Exec(query, status, id)
	if err != nil {
		return "", err
	}

	return status, nil
}

func UpdateRetryCount(db *sql.DB, id string) (int, error) {
	query := `
		SELECT retrycount
		FROM workflow_step
		WHERE id = ?
	`

	rows, _ := db.Query(query, id)

	var prevcounter int
	//increment retry counter by 1
	rows.Scan(&prevcounter)
	prevcounter++

	query2 := `
		UPDATE workflow_step
		SET retrycount =?
		WHERE id = ?
	`

	_,err := db.Query(query2,prevcounter,id)
	if err!=nil {
		return prevcounter,err
	}

	return prevcounter,err
}

func (r *WorkflowRepo) GetworkflowbyID(id string) ([]models.WorkflowStep,error)  {
	query :=`SELECT id,workflowid,name,status,retrycount from WorkflowStep where id = ?`

	
	rows, err := r.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var allworkflows []models.WorkflowStep
	for rows.Next(){
		var workflow models.WorkflowStep

		err := rows.Scan(&workflow.ID,&workflow.WorkflowID,&workflow.Name,&workflow.Status,&workflow.RetryCount)
		if err!=nil {
			return nil,err
		}

		allworkflows = append(allworkflows, workflow)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return allworkflows, nil
}
