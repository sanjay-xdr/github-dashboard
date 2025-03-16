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
