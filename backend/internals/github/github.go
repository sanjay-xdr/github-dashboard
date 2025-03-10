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

func Run() {

	GetMergeFrequency(context.Background(), github.NewClient(nil), owner, repo)
}

// make these functions for Daily Weekly and Monthly use

// func GetTotalMerges(ctx context.Context, client *github.Client, owner, repo string) (int, error) {
// 	now := time.Now()
// 	daily := 0
// 	weekly := 0
// 	monthly := 0
// 	mergedPRs, _, err := client.Search.Issues(ctx, fmt.Sprintf("repo:%s/%s is:pr is:merged", owner, repo), nil)
// 	if err != nil {
// 		return 0, err
// 	}

// 	fmt.Print("Total Merged PRs")
// 	fmt.Print(mergedPRs.GetTotal())
// 	return mergedPRs.GetTotal(), nil
// }

func GetMergeFrequency(ctx context.Context, client *github.Client, owner, repo string) (int, int, int, error) {
	now := time.Now()
	daily := 0
	weekly := 0
	monthly := 0

	options := &github.SearchOptions{
		Sort:  "updated",
		Order: "desc",
	}

	mergedPRs, _, err := client.Search.Issues(ctx, fmt.Sprintf("repo:%s/%s is:pr is:merged", owner, repo), options)
	if err != nil {
		return 0, 0, 0, err
	}

	for _, pr := range mergedPRs.Issues {
		if pr.GetClosedAt().After(now.AddDate(0, 0, -1)) {
			daily++
		}
		if pr.GetClosedAt().After(now.AddDate(0, 0, -7)) {
			weekly++
		}
		if pr.GetClosedAt().After(now.AddDate(0, -1, 0)) {
			monthly++
		}
	}

	fmt.Print("Daily ", daily)
	fmt.Print("Weekly ", weekly)
	fmt.Print("Monthly ", monthly)

	return daily, weekly, monthly, nil
}

func GetPendingMerges(ctx context.Context, client *github.Client, owner, repo string) (int, int, int, error) {
	now := time.Now()
	dailyStartTime := now.AddDate(0, 0, -1)
	weeklyStartTime := now.AddDate(0, 0, -7)
	monthlyStartTime := now.AddDate(0, -1, 0)

	// Construct queries
	dailyQuery := fmt.Sprintf("repo:%s/%s is:pr is:open created:>=%s", owner, repo, dailyStartTime.Format(time.RFC3339))
	weeklyQuery := fmt.Sprintf("repo:%s/%s is:pr is:open created:>=%s", owner, repo, weeklyStartTime.Format(time.RFC3339))
	monthlyQuery := fmt.Sprintf("repo:%s/%s is:pr is:open created:>=%s", owner, repo, monthlyStartTime.Format(time.RFC3339))

	// Fetch daily pending merges
	dailyPRs, _, err := client.Search.Issues(ctx, dailyQuery, nil)
	if err != nil {
		return 0, 0, 0, err
	}

	// Fetch weekly pending merges
	weeklyPRs, _, err := client.Search.Issues(ctx, weeklyQuery, nil)
	if err != nil {
		return 0, 0, 0, err
	}

	// Fetch monthly pending merges
	monthlyPRs, _, err := client.Search.Issues(ctx, monthlyQuery, nil)
	if err != nil {
		return 0, 0, 0, err
	}

	return dailyPRs.GetTotal(), weeklyPRs.GetTotal(), monthlyPRs.GetTotal(), nil
}

func GetMergeSuccessRate(ctx context.Context, client *github.Client, owner, repo string) (float64, error) {
	mergedPRs, _, err := client.Search.Issues(ctx, fmt.Sprintf("repo:%s/%s is:pr is:merged", owner, repo), nil)
	if err != nil {
		return 0, err
	}

	closedPRs, _, err := client.Search.Issues(ctx, fmt.Sprintf("repo:%s/%s is:pr is:closed", owner, repo), nil)
	if err != nil {
		return 0, err
	}

	if closedPRs.GetTotal() == 0 {
		return 0, nil
	}

	successRate := float64(mergedPRs.GetTotal()) / float64(closedPRs.GetTotal())
	return successRate, nil
}
type MergedPRsByDate struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

func getMergedPRs(ctx context.Context, client *github.Client, owner, repo string) ([]MergedPRsByDate, error) {
	opt := &github.PullRequestListOptions{
		State:     "closed",
		Sort:      "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	now := time.Now()
	sevenDaysAgo := now.AddDate(0, 0, -7)
	mergedPRsByDate := make(map[string]int)

	for {
		prs, resp, err := client.PullRequests.List(ctx, owner, repo, opt)
		if err != nil {
			return nil, err
		}

		for _, pr := range prs {
			if pr.MergedAt != nil && pr.MergedAt.After(sevenDaysAgo) {
				date := pr.MergedAt.Format("2006-01-02")
				mergedPRsByDate[date]++
			} else if pr.MergedAt != nil && pr.MergedAt.Before(sevenDaysAgo) {
				// If we encounter a PR merged before the 7-day window, we can stop
				break
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var result []MergedPRsByDate
	for date, count := range mergedPRsByDate {
		result = append(result, MergedPRsByDate{Date: date, Count: count})
	}

	return result, nil
}