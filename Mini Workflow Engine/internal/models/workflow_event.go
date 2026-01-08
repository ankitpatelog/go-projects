package workflows

import "time"


//this is for new worflow generate first time
type WorkflowInstance struct {
	ID          string    `json:"id"`
	EventType   string    `json:"event_type"`   // signup, order_created
	Status      string    `json:"status"`       // PENDING, RUNNING, FAILED, COMPLETED
	CurrentStep int       `json:"current_step"` // index of step (0,1,2...)
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

//this struct for every no of workflows
type WorkflowStep struct {
	ID           int       `json:"id"`
	WorkflowID   string    `json:"workflow_id"`
	StepName     string    `json:"step_name"`   // validate_user, send_email
	StepOrder    int       `json:"step_order"`  // 0,1,2
	Status       string    `json:"status"`      // PENDING, SUCCESS, FAILED
	RetryCount   int       `json:"retry_count"`
	LastError    string    `json:"last_error"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}