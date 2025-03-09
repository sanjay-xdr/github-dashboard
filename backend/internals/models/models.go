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
