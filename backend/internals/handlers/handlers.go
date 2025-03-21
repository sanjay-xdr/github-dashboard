package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

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

	test := []models.RepoOverview{
		{Name: "Stars", Value: data.Stars},
		{Name: "Forks", Value: data.Forks},
		{Name: "Watchers", Value: data.Watchers},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(test)
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

func FetchPRData(w http.ResponseWriter, r *http.Request) {

	prData, err := github.FetchPRs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	github.InsertPRs(prData)
	data, err := github.QueryPRsByDate(time.Time{}, time.Time{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dashboardData := generatePRDashboard(data)
	sort.Slice(dashboardData, func(i, j int) bool {
		dateI, errI := time.Parse("2006-01-02", dashboardData[i].Date)
		dateJ, errJ := time.Parse("2006-01-02", dashboardData[j].Date)
		if errI != nil || errJ != nil {
			fmt.Println("Error parsing date:", errI, errJ)
			return false
		}
		return dateI.Before(dateJ)
	})
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboardData)

}

func FetchWorkflowData(w http.ResponseWriter, r *http.Request) {

	workflowdata, err := github.FetchWorkflows()
	github.InsertWorkflows(workflowdata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Print(workflowdata)
	// generatePRDashboard(workflowdata)\
	fetchData, err := github.QueryWorkflowsByDate(time.Time{}, time.Time{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	worflowDashboardData := generateWorkflowDashboard(fetchData)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(worflowDashboardData)

}

func generateWorkflowDashboard(workflowdata []models.WorkflowItem) []models.WorkflowSummary {
	dateSummaryMap := make(map[string]*models.WorkflowSummary)

	for _, wf := range workflowdata {
		date := wf.CreatedAt.Format("2006-01-02")
		if _, exists := dateSummaryMap[date]; !exists {
			dateSummaryMap[date] = &models.WorkflowSummary{
				Date:    date,
				Success: 0,
				Failed:  0,
				Pending: 0,
			}
		}

		if wf.Conclusion == "success" {
			dateSummaryMap[date].Success++
		} else if wf.Conclusion == "failure" {
			dateSummaryMap[date].Failed++
		} else if wf.Conclusion == "action_required" {
			dateSummaryMap[date].Pending++
		}
	}

	var summary []models.WorkflowSummary
	for _, v := range dateSummaryMap {
		summary = append(summary, *v)
	}

	return summary
}

func generatePRDashboard(prs []models.PR) []models.PRDashboard {
	prStats := make(map[string]*models.PRDashboard)

	for _, pr := range prs {
		date := pr.CreatedAt.Format("2006-01-02") // Format date as YYYY-MM-DD

		if _, exists := prStats[date]; !exists {
			prStats[date] = &models.PRDashboard{
				Date:     date,
				TotalPR:  0,
				OpenPR:   0,
				MergedPR: 0,
				ClosedPR: 0,
			}
		}

		// Increment total PRs for that date
		prStats[date].TotalPR++

		// Categorize PRs by state
		if pr.State == "open" {
			prStats[date].OpenPR++
		} else if pr.State == "merged" { // Assuming "merged" is a valid PR state
			prStats[date].MergedPR++
		} else if pr.State == "closed" {
			prStats[date].MergedPR++
		}
	}

	// Convert map to slice
	var result []models.PRDashboard
	for _, data := range prStats {
		result = append(result, *data)
	}

	return result
}
