package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

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
