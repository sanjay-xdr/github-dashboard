package github

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/v60/github"
	"github.com/sanjay-xdr/github-dashboard/backend/internals/models"
)

func FetchAllPullRequestStats() (*models.PRStatus, error) {
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
