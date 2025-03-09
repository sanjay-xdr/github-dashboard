package github

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/github"
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
			if pr.GetCreatedAt().Before(date) || pr.GetCreatedAt().Equal(date) {
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
