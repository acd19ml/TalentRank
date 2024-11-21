package impl

import (
	"log"
	"sync"

	"github.com/acd19ml/TalentRank/apps/user"
)

// CalculateOverallScore 计算开发者的最终技术评分
func CalculateOverallScore(userRepos *user.UserRepos) error {
	// 权重设置
	wStar := 3.0      // Star 权重
	wFork := 2.0      // Fork 权重
	wDependent := 1.0 // Dependents 权重
	wFollowers := 0.1 // Followers 权重

	// 从 UserRepos 中获取总粉丝数和仓库列表
	totalFollowers := userRepos.User.Followers
	totalScore := 0.0

	// 使用 channel 和 WaitGroup 来处理并发项目评分计算
	projectScores := make(chan float64, len(userRepos.Repos))
	var wgProjects sync.WaitGroup

	// 遍历所有项目，开启协程计算每个项目的技术评分
	for _, repo := range userRepos.Repos {
		wgProjects.Add(1)

		go func(repo *user.Repo) {
			defer wgProjects.Done()

			// 计算项目的影响力权重
			projectImpact := wStar*float64(repo.Star) + wFork*float64(repo.Fork) + wDependent*float64(repo.Dependent)

			// 计算项目的技术评分（贡献度 * 项目影响力）
			contribution, err := calculateContribution(repo)
			if err != nil {
				log.Printf("Failed to calculate contribution for repo %s: %v", repo.Repo, err)
				projectScores <- 0
				return
			}

			projectScore := contribution * projectImpact
			// 将项目分数发送到 channel
			projectScores <- projectScore
		}(repo)
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
	userRepos.User.Score = overallScore // 将最终评分存储在 UserRepos 的 Score 字段中
	// log.Printf("Overall score for %s: %f\n", userRepos.User.Username, overallScore)
	return nil
}

// CalculateContribution 计算给定用户在给定项目中的贡献度
func calculateContribution(repo *user.Repo) (float64, error) {
	// 权重设置
	w1 := 0.25 // 代码提交数的权重
	w2 := 0.3  // 解决的 Issue 和合并的 PR 数的权重
	w3 := 0.2  // 代码评审数的权重
	w4 := 0.25 // 代码行变更数的权重

	// 初始总权重
	totalWeight := 0.0
	contribution := 0.0

	// 代码提交贡献
	if repo.CommitsTotal > 0 {
		contribution += w1 * float64(repo.Commits) / float64(repo.CommitsTotal)
		totalWeight += w1
	}

	// 解决的 Issue 和合并的 PR 数贡献
	totalIssuesAndPRs := repo.IssueTotal + repo.PullRequestTotal
	if totalIssuesAndPRs > 0 {
		userIssuesAndPRs := repo.Issue + repo.PullRequest
		contribution += w2 * float64(userIssuesAndPRs) / float64(totalIssuesAndPRs)
		totalWeight += w2
	}

	// 代码评审贡献
	if repo.CodeReviewTotal > 0 {
		contribution += w3 * float64(repo.CodeReview) / float64(repo.CodeReviewTotal)
		totalWeight += w3
	}

	// 代码行变更贡献
	if repo.LineChangeTotal > 0 {
		contribution += w4 * float64(repo.LineChange) / float64(repo.LineChangeTotal)
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
