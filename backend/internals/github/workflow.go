package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sanjay-xdr/github-dashboard/backend/internals/database"
	"github.com/sanjay-xdr/github-dashboard/backend/internals/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

// Insert PRs into MongoDB
func InsertWorkflows(workflow []models.WorkflowItem) (bool, error) {
	if len(workflow) == 0 {
		fmt.Println("No new PRs to insert.")
		return false, nil
	}

	col := database.GetWorkflowCollection()
	ctx := context.TODO()

	for _, item := range workflow {
		// Define the filter to check if PR already exists (by PR number & repository)
		filter := bson.M{"Id": item.ID, "Name": item.Name}

		// Define the update operation (set PR data)
		update := bson.M{"$set": item}

		// Upsert = true ensures:
		// - If PR exists, update it
		// - If PR doesn't exist, insert it
		opts := options.UpdateOne().SetUpsert(true)

		_, err := col.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return false, fmt.Errorf("MongoDB Upsert Error: %v", err)
		}
	}

	fmt.Println("Inserted/Updated Workflows successfully!")
	return true, nil
}

func QueryWorkflowsByDate(startDate, endDate time.Time) ([]models.WorkflowItem, error) {
	col := database.GetWorkflowCollection()
	filter := bson.M{}
	if !startDate.IsZero() && !endDate.IsZero() {
		filter["created_at"] = bson.M{
			"$gte": startDate.UTC(),
			"$lte": endDate.UTC(),
		}
	}

	cursor, err := col.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("MongoDB Query Error: %v", err)
	}
	defer cursor.Close(context.TODO())

	var results []models.WorkflowItem
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, fmt.Errorf("MongoDB Decoding Error: %v", err)
	}

	fmt.Printf("Found %d PRs from MongoDB\n", len(results))
	return results, nil
}
