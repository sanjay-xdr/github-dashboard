package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PRStatus struct {
	TotalPR  int16
	OpenPR   int16
	ClosePR  int16
	MergedPR int16
}

type RepoStats struct {
	Stars    int16
	Watchers int16
	Forks    int16
}
type MergedPRsByDate struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type DateRangeRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type RepoOverview struct {
	Name  string `json:"name"`
	Value int16  `json:"value"`
}

type PR struct {
	ID         int                `json:"id"`
	Number     int                `json:"number"`
	Title      string             `json:"title"`
	State      string             `json:"state"`
	Merged     bool               `json:"merged"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	MergedAt   *time.Time         `json:"merged_at,omitempty"`
	Repository string             `json:"repository"`
	InsertedAt primitive.DateTime `bson:"inserted_at"`
}

type PRDashboard struct {
	Date     string `json:"date"`
	TotalPR  int    `json:"totalPR"`
	OpenPR   int    `json:"openPR"`
	MergedPR int    `json:"mergedPR"`
	ClosedPR int    `json:"closedPR"`
}

type WorkflowItem struct {
	ID         int                `json:"id"`
	Name       string             `json:"name"`
	NodeID     string             `json:"node_id"`
	Path       string             `json:"path"`
	Status     string             `json:"status"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	InsertedAt primitive.DateTime `json:"inserted_at"`
	Conclusion string             `json:"conclusion"`
}

type Workflows struct {
	TotalCount int            `json:"total_count"`
	Workflows  []WorkflowItem `json:"workflow_runs"`
}

type WorkflowSummary struct {
	Date    string `json:"date"`
	Success int    `json:"success"`
	Failed  int    `json:"failed"`
	Pending int    `json:"pending"`
}
