package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sanjay-xdr/github-dashboard/backend/internals/github"
	"github.com/sanjay-xdr/github-dashboard/backend/internals/models"
)

func GetPullRequestData(w http.ResponseWriter, r *http.Request) {

	data, err := github.FetchAllPullRequestStats()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)

}

func GetRepoData(w http.ResponseWriter, r *http.Request) {

	data, err := github.GetRepoStats()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func GetTestResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	github.GetWorkflowRuns()

	//fetch test result
	json.NewEncoder(w).Encode("Test Result")
}

func GetMergedPRByDate(w http.ResponseWriter, r *http.Request) {

	var dateRange models.DateRangeRequest

	if err := json.NewDecoder(r.Body).Decode(&dateRange); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	startDate := dateRange.StartDate
	endDate := dateRange.EndDate
	if startDate == "" || endDate == "" {
		http.Error(w, "startDate and endDate are required", http.StatusBadRequest)
		return
	}

	mergedPR, err := github.GetMergedPRByDate(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mergedPR)

}
