package github

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/v60/github"
	"github.com/sanjay-xdr/github-dashboard/backend/internals/models"
)

const (
	owner = "keploy"
	repo  = "website"
)

func FetchAllPullRequestStats() (*models.PRStatus, error) {
	// owner := "keploy"
	// repo := "website"
	client := github.NewClient(nil)

	opt := &github.PullRequestListOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 300},
	}

	var allPRs []*github.PullRequest
	for {
		prs, resp, err := client.PullRequests.List(context.Background(), owner, repo, opt)
		if err != nil {
			log.Fatalf("Error fetching PRs: %v", err)
		}
		allPRs = append(allPRs, prs...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	// prs, _, err := client.PullRequests.List(context.Background(), owner, repo, opt)
	// if err != nil {
	// 	// fmt.Printf("Error fetching pull requests: %v\n", err)
	// 	return nil, fmt.Errorf("error fetching pull requests: %v", err)
	// }

	fmt.Print("============================Printing PRs==========================================\n")

	var openPRs, closedPRs, mergedPRs int16
	for _, pr := range allPRs {
		if pr.GetState() == "open" {
			openPRs++
		} else if pr.GetState() == "closed" {
			closedPRs++
		}
		if !pr.GetMergedAt().IsZero() {
			mergedPRs++
		}

		fmt.Printf("PR #%d: %s (State: %s)\n", pr.GetNumber(), pr.GetTitle(), pr.GetState())

		fmt.Print()

		// break
	}

	totalPRs := int16(len(allPRs))
	return &models.PRStatus{
		TotalPR:  totalPRs,
		ClosePR:  closedPRs,
		MergedPR: mergedPRs,
		OpenPR:   openPRs,
	}, nil
}

func FetchPullRequestStatsUptoDate(date time.Time) (*models.PRStatus, error) {
	owner := "keploy"
	repo := "website"
	client := github.NewClient(nil)

	opt := &github.PullRequestListOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 300},
	}

	var allPRs []*github.PullRequest
	for {
		prs, resp, err := client.PullRequests.List(context.Background(), owner, repo, opt)
		if err != nil {
			log.Fatalf("Error fetching PRs: %v", err)
		}
		for _, pr := range prs {
			if pr.GetCreatedAt().Time.Before(date) || pr.GetCreatedAt().Time.Equal(date) {
				allPRs = append(allPRs, pr)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	fmt.Print("============================Printing PRs==========================================\n")

	var openPRs, closedPRs, mergedPRs int16
	for _, pr := range allPRs {
		if pr.GetState() == "open" {
			openPRs++
		} else if pr.GetState() == "closed" {
			closedPRs++
		}
		if !pr.GetMergedAt().IsZero() {
			mergedPRs++
		}

		fmt.Printf("PR #%d: %s (State: %s)\n", pr.GetNumber(), pr.GetTitle(), pr.GetState())
	}

	totalPRs := int16(len(allPRs))
	return &models.PRStatus{
		TotalPR:  totalPRs,
		ClosePR:  closedPRs,
		MergedPR: mergedPRs,
		OpenPR:   openPRs,
	}, nil
}

func GetRepoStats() (*models.RepoStats, error) {
	// Replace "owner" and "repo" with the repository owner and name
	owner := "keploy"
	repo := "website"
	client := github.NewClient(nil)

	repository, _, err := client.Repositories.Get(context.Background(), owner, repo)
	if err != nil {
		fmt.Printf("Error fetching repository details: %v\n", err)
		return nil, fmt.Errorf("error fetching repository details: %v", err)
	}

	repoStat := &models.RepoStats{

		Stars:    int16(repository.GetStargazersCount()),
		Watchers: int16(repository.GetWatchersCount()),
		Forks:    int16(repository.GetForksCount()),
	}
	return repoStat, nil
}

func GetWorkflowRuns() {
	owner := "keploy"
	repo := "website"
	client := github.NewClient(nil)

	runs, _, err := client.Actions.ListRepositoryWorkflowRuns(context.Background(), owner, repo, nil)

	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Print("Something went wrong")
		return
	}

	for _, run := range runs.WorkflowRuns {
		fmt.Printf("Run ID: %d | Status: %s | Conclusion: %s | Created At: %s\n",
			run.GetID(), run.GetStatus(), run.GetConclusion(), run.GetCreatedAt().Time)
	}
}

func GetMergedPRByDate(startDate, endDate string) ([]models.MergedPRsByDate, error) {

	return getMergedPRByDate(context.Background(), github.NewClient(nil), owner, repo, startDate, endDate)
}

func getMergedPRByDate(ctx context.Context, client *github.Client, owner, repo, startDate, endDate string) ([]models.MergedPRsByDate, error) {
	opt := &github.PullRequestListOptions{
		State:     "closed",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	// Initialize map with 0 counts for all days in the range
	mergedPRsByDate := make(map[string]int)
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		date := d.Format("2006-01-02")
		mergedPRsByDate[date] = 0
	}

	// Fetch PR data
	for {
		prs, resp, err := client.PullRequests.List(ctx, owner, repo, opt)
		if err != nil {
			return nil, err
		}

		for _, pr := range prs {
			if pr.MergedAt != nil && pr.MergedAt.After(start) && pr.MergedAt.Before(end.AddDate(0, 0, 1)) {
				date := pr.MergedAt.Format("2006-01-02")
				mergedPRsByDate[date]++
			} else if pr.MergedAt != nil && pr.MergedAt.Before(start) {
				// Stop processing once past the date range
				break
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	// Convert map to sorted slice
	var result []models.MergedPRsByDate
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		date := d.Format("2006-01-02")
		result = append(result, models.MergedPRsByDate{Date: date, Count: mergedPRsByDate[date]})
	}
	fmt.Println()
	fmt.Print(result)
	fmt.Println()
	return result, nil
}
