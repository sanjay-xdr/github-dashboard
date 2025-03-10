package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sanjay-xdr/github-dashboard/backend/internals/github"
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

func Test(W http.ResponseWriter, r *http.Request) {

	github.Run()
}
