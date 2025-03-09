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
