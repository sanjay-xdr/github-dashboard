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

// Fetch PRs from GitHub API for the last 30 days
func FetchPRs() ([]models.PR, error) {
	fmt.Println("Fetching PRs from GitHub...")
	now := time.Now()
	monthAgo := now.AddDate(0, 0, -30).Format("2006-01-02T15:04:05Z")
	fmt.Print("30 days ago: ", monthAgo, "\n")

	url := "https://api.github.com/repos/keploy/website/pulls?state=all&sort=created&direction=desc&per_page=100"
	req, _ := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []models.PR{}, fmt.Errorf("error fetching PRs: %v", err)
	}
	defer resp.Body.Close()

	var prs []models.PR
	err = json.NewDecoder(resp.Body).Decode(&prs)
	if err != nil {
		// log.Fatal("Error decoding response:", err)
		return []models.PR{}, fmt.Errorf("error decoding response: %v", err)
	}

	// Filter PRs within the last 30 days
	var recentPRs []models.PR
	for _, pr := range prs {
		if pr.CreatedAt.After(time.Now().AddDate(0, 0, -30)) {
			pr.Repository = "keploy/website"
			pr.InsertedAt = primitive.NewDateTimeFromTime(time.Now())
			recentPRs = append(recentPRs, pr)
		}
	}

	fmt.Printf("Fetched %d PRs from GitHub\n", len(recentPRs))
	return recentPRs, nil
}

// Insert PRs into MongoDB
func InsertPRs(prs []models.PR) (bool, error) {
	if len(prs) == 0 {
		fmt.Println("No new PRs to insert.")
		return false, nil
	}

	col := database.GetPRCollection()
	ctx := context.TODO()

	for _, pr := range prs {
		// Define the filter to check if PR already exists (by PR number & repository)
		filter := bson.M{"number": pr.Number, "repository": pr.Repository}

		// Define the update operation (set PR data)
		update := bson.M{"$set": pr}

		// Upsert = true ensures:
		// - If PR exists, update it
		// - If PR doesn't exist, insert it
		opts := options.UpdateOne().SetUpsert(true)

		_, err := col.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return false, fmt.Errorf("MongoDB Upsert Error: %v", err)
		}
	}

	fmt.Println("Inserted/Updated PRs successfully!")
	return true, nil
}

func QueryPRsByDate(startDate, endDate time.Time) ([]models.PR, error) {
	col := database.GetPRCollection()
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

	var results []models.PR
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, fmt.Errorf("MongoDB Decoding Error: %v", err)
	}

	fmt.Printf("Found %d PRs from MongoDB\n", len(results))
	return results, nil
}
