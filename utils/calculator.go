package utils

import (
	"context"
	"fmt"
	"log"
)

type Calculator struct {
	service Service
}

func NewCalculator(service Service) *Calculator {
	return &Calculator{service: service}
}

// CalculateOverallScore 计算开发者的最终技术评分
func (osc *Calculator) CalculateOverallScore(ctx context.Context, username string) (float64, error) {
	// 权重设置
	wStar := 3.0          // Star 权重
	wFork := 2.0          // Fork 权重
	wDependent := 1.0     // Dependents 权重
	wFollowers := 0.05    // Followers 权重
	wAchievements := 0.05 // Achievements 权重

	// 获取用户的影响力数据
	totalFollowers, err := osc.service.GetFollowers(ctx, username)
	if err != nil {
		return 0, fmt.Errorf("failed to get followers count: %w", err)
	}
	log.Println("GetFollowers completed")

	// 获取用户的项目数据
	repositories, err := osc.service.GetRepositories(ctx, username)
	if err != nil {
		return 0, fmt.Errorf("failed to get repositories: %w", err)
	}
	log.Println("GetRepositories completed")

	// 初始化总分
	totalScore := 0.0

	// 遍历所有项目，计算每个项目的贡献度和技术评分
	for _, repoName := range repositories {
		// 计算项目贡献度
		contribution, err := osc.CalculateContribution(ctx, username, repoName)
		if err != nil {
			return 0, fmt.Errorf("failed to calculate contribution for repo %s: %w", repoName, err)
		}
		log.Printf("CalculateContribution for %s completed\n", repoName)

		// 获取项目的 Star、Fork、Dependents 数据
		stars, err := osc.service.GetRepoStars(ctx, username, repoName)
		if err != nil {
			return 0, fmt.Errorf("failed to get stars for repo %s: %w", repoName, err)
		}
		log.Printf("GetRepoStars for %s completed\n", repoName)

		forks, err := osc.service.GetRepoForks(ctx, username, repoName)
		if err != nil {
			return 0, fmt.Errorf("failed to get forks for repo %s: %w", repoName, err)
		}
		log.Printf("GetRepoForks for %s completed\n", repoName)

		dependents, err := osc.service.GetDependentRepositoriesByRepo(ctx, username)
		if err != nil {
			return 0, fmt.Errorf("failed to get dependents for repo %s: %w", repoName, err)
		}
		log.Printf("GetDependentRepositoriesByRepo for %s completed\n", repoName)

		// 计算项目的影响力权重
		projectImpact := wStar*float64(stars) + wFork*float64(forks) + wDependent*float64(dependents[repoName])

		// 计算项目的技术评分（贡献度 * 项目影响力）
		projectScore := contribution * projectImpact

		// 累加到总分
		totalScore += projectScore
	}

	// 获取 Achievements 数量
	achievements, err := osc.service.GetDependentRepositories(ctx, username)
	if err != nil {
		return 0, fmt.Errorf("failed to get achievements count: %w", err)
	}
	log.Println("GetDependentRepositories (Achievements count) completed")

	// 计算最终技术评分，包括 Followers 和 Achievements 的加权影响
	overallScore := totalScore * (1 + wFollowers*float64(totalFollowers) + wAchievements*float64(achievements))

	return overallScore, nil
}

// CalculateContribution 计算给定用户在给定项目中的贡献度
func (cc *Calculator) CalculateContribution(ctx context.Context, username, repoName string) (float64, error) {
	// 权重设置
	w1 := 0.25 // 代码提交数的权重
	w2 := 0.3  // 解决的 Issue 和合并的 PR 数的权重
	w3 := 0.2  // 代码评审数的权重
	w4 := 0.25 // 代码行变更数的权重

	// 初始总权重
	totalWeight := 0.0

	// 获取项目总数据
	totalCommits, err := cc.service.GetTotalCommitsByRepo(ctx, username)
	if err != nil {
		return 0, fmt.Errorf("failed to get total commits: %w", err)
	}
	log.Println("GetTotalCommitsByRepo completed")

	totalIssues, err := cc.service.GetTotalIssuesByRepo(ctx, username)
	if err != nil {
		return 0, fmt.Errorf("failed to get total issues: %w", err)
	}
	log.Println("GetTotalIssuesByRepo completed")

	totalPullRequests, err := cc.service.GetTotalPullRequestsByRepo(ctx, username)
	if err != nil {
		return 0, fmt.Errorf("failed to get total pull requests: %w", err)
	}
	log.Println("GetTotalPullRequestsByRepo completed")

	totalCodeReviews, err := cc.service.GetTotalCodeReviewsByRepo(ctx, username)
	if err != nil {
		return 0, fmt.Errorf("failed to get total code reviews: %w", err)
	}
	log.Println("GetTotalCodeReviewsByRepo completed")

	totalLineChanges, err := cc.service.GetTotalLineChanges(ctx, username)
	if err != nil {
		return 0, fmt.Errorf("failed to get total line changes: %w", err)
	}
	log.Println("GetTotalLineChanges completed")

	// 获取用户在项目中的贡献数据
	userCommits, err := cc.service.GetUserCommitsByRepo(ctx, username)
	if err != nil {
		return 0, fmt.Errorf("failed to get user commits: %w", err)
	}
	log.Println("GetUserCommitsByRepo completed")

	userSolvedIssues, err := cc.service.GetUserSolvedIssuesByRepo(ctx, username)
	if err != nil {
		return 0, fmt.Errorf("failed to get user solved issues: %w", err)
	}
	log.Println("GetUserSolvedIssuesByRepo completed")

	userMergedPRs, err := cc.service.GetUserMergedPullRequestsByRepo(ctx, username)
	if err != nil {
		return 0, fmt.Errorf("failed to get user merged pull requests: %w", err)
	}
	log.Println("GetUserMergedPullRequestsByRepo completed")

	userCodeReviews, err := cc.service.GetUserCodeReviewsByRepo(ctx, username)
	if err != nil {
		return 0, fmt.Errorf("failed to get user code reviews: %w", err)
	}
	log.Println("GetUserCodeReviewsByRepo completed")

	userLineChanges, err := cc.service.GetLineChanges(ctx, username, repoName)
	if err != nil {
		return 0, fmt.Errorf("failed to get user line changes: %w", err)
	}
	log.Println("GetLineChanges completed")

	// 计算贡献度的每一部分并动态调整总权重
	contribution := 0.0

	// 代码提交贡献
	if totalCommits[repoName] > 0 {
		contribution += w1 * float64(userCommits[repoName]) / float64(totalCommits[repoName])
		totalWeight += w1
	}

	// 解决的 Issue 和合并的 PR 数贡献
	totalIssuesAndPRs := totalIssues[repoName] + totalPullRequests[repoName]
	if totalIssuesAndPRs > 0 {
		userIssuesAndPRs := userSolvedIssues[repoName] + userMergedPRs[repoName]
		contribution += w2 * float64(userIssuesAndPRs) / float64(totalIssuesAndPRs)
		totalWeight += w2
	}

	// 代码评审贡献
	if totalCodeReviews[repoName] > 0 {
		contribution += w3 * float64(userCodeReviews[repoName]) / float64(totalCodeReviews[repoName])
		totalWeight += w3
	}

	// 代码行变更贡献
	if totalLineChanges > 0 {
		contribution += w4 * float64(userLineChanges) / float64(totalLineChanges)
		totalWeight += w4
	}

	// 如果总权重为 0，说明项目没有有效数据，返回贡献度为 0
	if totalWeight == 0 {
		return 0, nil
	}

	// 计算最终贡献度
	finalContribution := contribution / totalWeight

	return finalContribution, nil
}
