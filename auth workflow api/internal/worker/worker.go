package worker

import (
	"auth-workflow/internal/config"
	"auth-workflow/internal/models"
	"auth-workflow/internal/queue"
	"auth-workflow/internal/repo_workflow"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func startWorker() {
	go func() {
		// get all the pending workflowsteps
		var db *sql.DB
		steps,_ := Getallpendingworks(db)
		
		for _,step := range steps{
			//added all into channel
			queue.WorkflowStepQueue <- step
		}

		//handle process status
		for step := range queue.WorkflowStepQueue{
			_,err := repo_workflow.UpdateStatus(db,step.ID,"Running")
			
			if err!=nil {
				log.Fatal("Updatiing status problem")
			}
			
			status,err := ExecuteEvent(step)
			if err!=nil {
				log.Fatal("Execution failed problem")
			}
			
			if status=="Success" {
				//update db
				_,err := repo_workflow.UpdateStatus(db,step.ID,"Success")
				if err!=nil {
					log.Fatal("Updatiing status problem")
				}
				return
			}

			if status == "FAILED" {

			//  Update DB status (non-fatal)
			if _,err := repo_workflow.UpdateStatus(db, step.ID, "FAILED"); err != nil {
			log.Println("failed to update status:", err)
				return
			}
		
			// 2️⃣ Increment retry count in DB
			_, err := repo_workflow.UpdateRetryCount(db, step.ID)
			if err != nil {
				log.Println("failed to update retry count:", err)
				return
			}

	// 3️⃣ Redis retry limiter key
	key := fmt.Sprintf("retry:count:%s", step.ID)

	// 4️⃣ Increment retry counter in Redis
	count, err := config.Redis.Incr(config.Ctx, key).Result()
	if err != nil {
		log.Println("redis incr failed:", err)
		return
	}

	// 5️⃣ Set TTL only on first retry
	if count == 1 {
		config.Redis.Expire(config.Ctx, key, time.Hour)
	}

	// 6️⃣ Max retry reached → permanently fail
	if count > 10 {
		log.Println("max retries exceeded for step:", step.ID)

		// mark permanently failed
		_ ,_= repo_workflow.UpdateStatus(db, step.ID, "EVENT_FAILED")

		//send step into channel
		queue.WorkflowStepQueue <- step

		// cleanup redis key
		config.Redis.Del(config.Ctx, key)

		return
	}
  			}			
		}
		

	}()
}

func Getallpendingworks(db *sql.DB) ([]*models.WorkflowStep, error) {
	var steps []*models.WorkflowStep

	query := `
		SELECT id, workflowid, name, status, retrycount
		FROM workflow_step
		WHERE status = 'PENDING'
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		step := &models.WorkflowStep{}

		rows.Scan(&step.ID,
			&step.WorkflowID,
			&step.Name,
			&step.Status,
			&step.RetryCount,
		)

		steps = append(steps, step)
	}
	return steps, nil
}

//simulated a fake event handler for this workflow
func ExecuteEvent(Workflowstep *models.WorkflowStep)(string,error)  {

	switch Workflowstep.Name {
case "VERIFY_PAYMENT":
	status := GenerateRandPossi()
	return status,nil
case "CONFIRM_ORDER":
	status := GenerateRandPossi()
	return status,nil
case "UPDATE_INVENTORY":
	status := GenerateRandPossi()
	return status,nil
case "SEND_NOTIFICATION":
	status := GenerateRandPossi()
	return status,nil
}
return "",nil
}

func GenerateRandPossi() string {
	random:= rand.Intn(10)

	if random>0 &&random<5 {
		return "Failed"
	}
	if random>=5 && random<=9 {
		return "Success"
	}
	return ""
}

