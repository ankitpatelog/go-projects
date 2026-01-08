package models

type WorkflowInstance struct {
	ID     string
	Event  string
	Status string
}

type WorkflowStep struct {
	ID         string
	WorkflowID string
	Name       string
	Status     string
	RetryCount int
}
