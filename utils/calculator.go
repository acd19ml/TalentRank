package utils

import (
	"context"
	"log"
	"sync"
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
	wStar := 3.0      // Star 权重
	wFork := 2.0      // Fork 权重
	wDependent := 1.0 // Dependents 权重
	wFollowers := 0.1 // Followers 权重

	var wg sync.WaitGroup
	var totalFollowers int
	var repositories []string
	var totalScore float64

	var followersErr, reposErr error

	// 使用 WaitGroup 和 Goroutines 并发执行每个 Get 方法
	wg.Add(2)

	go func() {
		defer wg.Done()
		totalFollowers, followersErr = osc.service.GetFollowers(ctx, username)
		if followersErr != nil {
			log.Printf("Failed to get followers count: %v", followersErr)
		}
		log.Println("GetFollowers completed")
	}()

	go func() {
		defer wg.Done()
		repositories, reposErr = osc.service.GetRepositories(ctx, username)
		if reposErr != nil {
			log.Printf("Failed to get repositories: %v", reposErr)
		}
		log.Println("GetRepositories completed")
	}()

	// 等待所有数据获取完成
	wg.Wait()

	if followersErr != nil {
		return 0, followersErr
	}
	if reposErr != nil {
		return 0, reposErr
	}

	// 使用 channel 和 WaitGroup 来处理并发项目评分计算
	projectScores := make(chan float64, len(repositories))
	var wgProjects sync.WaitGroup

	// 遍历所有项目，开启协程计算每个项目的技术评分
	for _, repoName := range repositories {
		wgProjects.Add(1)

		go func(repoName string) {
			defer wgProjects.Done()

			// 为项目级别数据获取并发执行增加 WaitGroup
			var wgRepo sync.WaitGroup
			var stars, forks int
			var dependents map[string]int
			var contribution float64

			var err1, err2, err3, err4 error

			wgRepo.Add(4)

			go func() {
				defer wgRepo.Done()
				stars, err1 = osc.service.GetRepoStars(ctx, username, repoName)
				if err1 != nil {
					log.Printf("Failed to get stars for repo %s: %v", repoName, err1)
				}
				log.Printf("GetRepoStars for %s completed\n", repoName)
			}()

			go func() {
				defer wgRepo.Done()
				forks, err2 = osc.service.GetRepoForks(ctx, username, repoName)
				if err2 != nil {
					log.Printf("Failed to get forks for repo %s: %v", repoName, err2)
				}
				log.Printf("GetRepoForks for %s completed\n", repoName)
			}()

			go func() {
				defer wgRepo.Done()
				dependents, err3 = osc.service.GetDependentRepositoriesByRepo(ctx, username)
				if err3 != nil {
					log.Printf("Failed to get dependents for repo %s: %v", repoName, err3)
				}
				log.Printf("GetDependentRepositoriesByRepo for %s completed\n", repoName)
			}()

			go func() {
				defer wgRepo.Done()
				contribution, err4 = osc.CalculateContribution(ctx, username, repoName)
				if err4 != nil {
					log.Printf("Failed to calculate contribution for repo %s: %v", repoName, err4)
				}
				log.Printf("CalculateContribution for %s completed\n", repoName)
			}()

			// 等待项目级别数据获取完成
			wgRepo.Wait()

			// 计算项目的影响力权重
			projectImpact := wStar*float64(stars) + wFork*float64(forks) + wDependent*float64(dependents[repoName])

			// 计算项目的技术评分（贡献度 * 项目影响力）
			projectScore := contribution * projectImpact

			// 将项目分数发送到 channel
			projectScores <- projectScore
		}(repoName)
	}

	// 在一个独立的协程中等待所有项目评分协程完成
	go func() {
		wgProjects.Wait()
		close(projectScores) // 所有项目计算完成后关闭 channel
	}()

	// 汇总所有项目分数
	for score := range projectScores {
		totalScore += score
	}

	// 计算最终技术评分，包括 Followers 的加权影响
	overallScore := totalScore * (1 + wFollowers*float64(totalFollowers))

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

	// 使用 WaitGroup 并发执行每个 Get 方法
	var wg sync.WaitGroup
	var totalIssues, totalPullRequests, totalCodeReviews map[string]int
	var lineChanges map[string][]int
	var userSolvedIssues, userMergedPRs, userCodeReviews map[string]int
	var err1, err2, err3, err4, err5, err6, err7 error

	wg.Add(7)

	go func() {
		defer wg.Done()
		totalIssues, err1 = cc.service.GetTotalIssuesByRepo(ctx, username)
	}()

	go func() {
		defer wg.Done()
		totalPullRequests, err2 = cc.service.GetTotalPullRequestsByRepo(ctx, username)
	}()

	go func() {
		defer wg.Done()
		totalCodeReviews, err3 = cc.service.GetTotalCodeReviewsByRepo(ctx, username)
	}()

	go func() {
		defer wg.Done()
		lineChanges, err4 = cc.service.GetLineChangesByRepo(ctx, username)
	}()

	go func() {
		defer wg.Done()
		userSolvedIssues, err5 = cc.service.GetUserSolvedIssuesByRepo(ctx, username)
	}()

	go func() {
		defer wg.Done()
		userMergedPRs, err6 = cc.service.GetUserMergedPullRequestsByRepo(ctx, username)
	}()

	go func() {
		defer wg.Done()
		userCodeReviews, err7 = cc.service.GetUserCodeReviewsByRepo(ctx, username)
	}()

	// 等待所有并发获取的数据完成
	wg.Wait()

	if err := firstNonNilError(err1, err2, err3, err4, err5, err6, err7); err != nil {
		return 0, err
	}
	// 计算贡献度的每一部分并动态调整总权重
	contribution := 0.0

	// 代码提交贡献
	// []int{totalChanges, userChanges, totalCommits, userCommits}
	if len(lineChanges[repoName]) > 2 {
		contribution += w1 * float64(lineChanges[repoName][3]) / float64(lineChanges[repoName][2])
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
	if len(lineChanges[repoName]) > 0 {
		contribution += w4 * float64(lineChanges[repoName][1]) / float64(lineChanges[repoName][0])
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

// firstNonNilError 返回第一个非空错误
func firstNonNilError(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}
