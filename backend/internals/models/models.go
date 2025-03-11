package models

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
