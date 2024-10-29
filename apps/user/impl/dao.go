package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/acd19ml/TalentRank/apps/user"
)

func (u *ServiceImpl) save(ctx context.Context, ins *user.UserRepos) error {

	var (
		err error
	)

	// 开启事务
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// 通过defer处理事务提交
	// 1. 没有报错Commit
	// 2. 有报错Rollback
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("rollback error, %s\n", err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				log.Printf("rollback error, %s\n", err)
			}
		}
	}()

	ustmt, err := tx.PrepareContext(ctx, InsertUserSQL)
	if err != nil {
		return err
	}
	defer ustmt.Close()

	// 插入的 Organizations 字段为 []string
	organizationsJSON, err := json.Marshal(ins.Organizations)
	if err != nil {
		return fmt.Errorf("failed to marshal options: %v", err)
	}

	// 执行插入语句
	result, err := ustmt.ExecContext(ctx,
		ins.Id, ins.Username, ins.Name, ins.Company, ins.Blog, ins.Location,
		ins.Email, ins.Bio, ins.Followers, string(organizationsJSON), ins.Readme,
		ins.Commits, ins.Score, ins.PossibleNation, ins.ConfidenceLevel,
	)
	if err != nil {
		return err
	} else {
		fmt.Printf("insert user success, %v", result)
	}

	rstmt, err := tx.PrepareContext(ctx, InsertRepoSQL)
	if err != nil {
		return err
	}
	defer rstmt.Close()
	for _, repo := range ins.Repos {
		// 执行插入语句
		result, err := rstmt.ExecContext(ctx,
			repo.Id, ins.Id, repo.Repo, repo.Star, repo.Fork, repo.Dependent, repo.Commits,
			repo.CommitsTotal, repo.Issue, repo.IssueTotal, repo.PullRequest, repo.PullRequestTotal,
			repo.CodeReview, repo.CodeReviewTotal, repo.LineChange, repo.LineChangeTotal,
		)
		if err != nil {
			return err
		} else {
			fmt.Printf("insert repo success, %v", result)
		}
	}

	return nil
}

// constructUserRepos 创建用户及其仓库信息，开启多个 goroutine 并发获取数据
func (s *ServiceImpl) constructUserRepos(ctx context.Context, username string) (*user.UserRepos, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	// 创建User实例
	userins, err := s.constructUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	var wg sync.WaitGroup
	mu := sync.Mutex{}
	errCh := make(chan error, 10) // 用于捕获错误

	// 定义结果 map
	starsByRepo, forksByRepo, dependentsByRepo := map[string]int{}, map[string]int{}, map[string]int{}
	totalCommitsByRepo, totalIssuesByRepo := map[string][]int{}, map[string]int{}
	userIssuesByRepo, totalPRsByRepo, userPRsByRepo := map[string]int{}, map[string]int{}, map[string]int{}
	totalReviewsByRepo, userReviewsByRepo := map[string]int{}, map[string]int{}

	// 启动 goroutines 获取每个仓库信息
	wg.Add(10)

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetStarsByRepo(ctx, username); err == nil {
			mu.Lock()
			starsByRepo = result
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get stars by repo: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetForksByRepo(ctx, username); err == nil {
			mu.Lock()
			forksByRepo = result
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get forks by repo: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetDependentRepositoriesByRepo(ctx, username); err == nil {
			mu.Lock()
			dependentsByRepo = result
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get dependents by repo: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetLineChangesByRepo(ctx, username); err == nil {
			mu.Lock()
			totalCommitsByRepo = result
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get line changes by repo: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetTotalIssuesByRepo(ctx, username); err == nil {
			mu.Lock()
			totalIssuesByRepo = result
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get total issues by repo: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetUserSolvedIssuesByRepo(ctx, username); err == nil {
			mu.Lock()
			userIssuesByRepo = result
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get user solved issues by repo: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetTotalPullRequestsByRepo(ctx, username); err == nil {
			mu.Lock()
			totalPRsByRepo = result
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get total pull requests by repo: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetUserMergedPullRequestsByRepo(ctx, username); err == nil {
			mu.Lock()
			userPRsByRepo = result
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get user merged pull requests by repo: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetTotalCodeReviewsByRepo(ctx, username); err == nil {
			mu.Lock()
			totalReviewsByRepo = result
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get total code reviews by repo: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetUserCodeReviewsByRepo(ctx, username); err == nil {
			mu.Lock()
			userReviewsByRepo = result
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get user code reviews by repo: %w", err)
		}
	}()

	// 等待所有 goroutine 完成
	wg.Wait()
	close(errCh)

	// 检查是否有错误
	for e := range errCh {
		if e != nil {
			return nil, e
		}
	}

	// 构建Repo列表
	var repos []*user.Repo
	for repoName := range starsByRepo {
		repo := &user.Repo{
			User_id:          userins.Id,
			Repo:             repoName,
			Star:             starsByRepo[repoName],
			Fork:             forksByRepo[repoName],
			Dependent:        dependentsByRepo[repoName],
			Commits:          totalCommitsByRepo[repoName][3], // 用户提交
			CommitsTotal:     totalCommitsByRepo[repoName][2], // 总提交
			Issue:            userIssuesByRepo[repoName],
			IssueTotal:       totalIssuesByRepo[repoName],
			PullRequest:      userPRsByRepo[repoName],
			PullRequestTotal: totalPRsByRepo[repoName],
			CodeReview:       userReviewsByRepo[repoName],
			CodeReviewTotal:  totalReviewsByRepo[repoName],
			LineChange:       totalCommitsByRepo[repoName][1], // 用户变更行数
			LineChangeTotal:  totalCommitsByRepo[repoName][0], // 总变更行数
		}
		repo.InjectDefault()
		repos = append(repos, repo)
	}

	// 创建并返回 UserRepos 结构体
	userRepos := &user.UserRepos{
		User:  userins,
		Repos: repos,
	}

	// 计算用户的最终技术评分，并插入结构体
	if err = calculateOverallScore(ctx, userRepos); err != nil {
		return nil, fmt.Errorf("failed to calculate overall score: %w", err)
	}

	return userRepos, nil
}

func (s *ServiceImpl) constructUser(ctx context.Context, username string) (*user.User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	// 获取用户信息
	name, err := s.svc.GetName(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user name: %w", err)
	}

	company, err := s.svc.GetCompany(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user company: %w", err)
	}

	location, err := s.svc.GetLocation(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user location: %w", err)
	}

	email, err := s.svc.GetEmail(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user email: %w", err)
	}

	bio, err := s.svc.GetBio(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user bio: %w", err)
	}

	organizations, err := s.svc.GetOrganizations(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user organizations: %w", err)
	}

	followers, err := s.svc.GetFollowers(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user followers: %w", err)
	}

	readme, err := s.svc.GetReadme(ctx, username, 500)
	if err != nil {
		return nil, fmt.Errorf("failed to get user readme: %w", err)
	}

	commits, err := s.svc.GetCommits(ctx, username, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to get user commits: %w", err)
	}

	// 创建并返回User结构体
	user := &user.User{
		Id:            "", // 后续通过 InjectDefault() 填充
		Username:      username,
		Name:          name,
		Company:       company,
		Location:      location,
		Email:         email,
		Bio:           bio,
		Followers:     followers,
		Organizations: organizations,
		Readme:        readme,
		Commits:       commits,
		// Score、PossibleNation 和 ConfidenceLevel 需要根据特定逻辑或数据计算
		Score:           0,
		PossibleNation:  "",
		ConfidenceLevel: 0,
	}

	// 为新用户生成默认 ID
	user.InjectDefault()

	return user, nil
}

// CalculateOverallScore 计算开发者的最终技术评分
func calculateOverallScore(ctx context.Context, userRepos *user.UserRepos) error {
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
			contribution, err := calculateContribution(ctx, repo)
			if err != nil {
				log.Printf("Failed to calculate contribution for repo %s: %v", repo.Repo, err)
				projectScores <- 0
				return
			}

			projectScore := contribution * projectImpact
			log.Printf("Received %s score: %f\n", repo.Repo, projectScore)
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

	return nil
}

// CalculateContribution 计算给定用户在给定项目中的贡献度
func calculateContribution(ctx context.Context, repo *user.Repo) (float64, error) {
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
