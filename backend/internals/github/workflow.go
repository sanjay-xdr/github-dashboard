package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sanjay-xdr/github-dashboard/backend/internals/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FetchWorkflows() ([]models.WorkflowItem, error) {
	fmt.Println("Fetching workflows from GitHub...")
	now := time.Now()
	monthAgo := now.AddDate(0, 0, -30).Format("2006-01-02T15:04:05Z")
	fmt.Print("30 days ago: ", monthAgo, "\n")

	url := "https://api.github.com/repos/keploy/website/actions/workflows/134937554/runs?sort=created&direction=desc"
	req, _ := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	// fmt.Print(resp.Body)
	if err != nil {
		return []models.WorkflowItem{}, fmt.Errorf("error fetching workflows: %v", err)
	}
	defer resp.Body.Close()

	var workflows models.Workflows
	err = json.NewDecoder(resp.Body).Decode(&workflows)
	// fmt.Println()
	// fmt.Print("Printing the workflows")
	// fmt.Print(workflows.Workflows)
	if err != nil {
		// log.Fatal("Error decoding response:", err)
		return []models.WorkflowItem{}, fmt.Errorf("error decoding response: %v", err)
	}

	// Filter workflows within the last 30 days
	var recentWorkflows []models.WorkflowItem
	for _, wf := range workflows.Workflows {
		// fmt.Print("I am here")
		// fmt.Print("Printing the Date \n")
		// fmt.Print(wf.CreatedAt)
		// fmt.Println()

		if wf.CreatedAt.After(time.Now().AddDate(0, 0, -30)) {
			wf.InsertedAt = primitive.NewDateTimeFromTime(time.Now())
			recentWorkflows = append(recentWorkflows, wf)
		}

		// fmt.Print("I am here AGAIN")
	}

	fmt.Printf("Fetched %d workflows from GitHub\n", len(recentWorkflows))
	return recentWorkflows, nil
}
