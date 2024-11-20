package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/acd19ml/TalentRank/apps/llm"
	"log"
	"sync"

	"github.com/acd19ml/TalentRank/apps"
	"github.com/acd19ml/TalentRank/apps/git"
	"github.com/acd19ml/TalentRank/apps/user"
)

func (u *ServiceImpl) save(ctx context.Context, ins *user.UserRepos) error {

	var (
		err error
	)

	// 开启事务
	tx, err := u.Db.BeginTx(ctx, nil)
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
	_, err = ustmt.ExecContext(ctx,
		ins.Id, ins.Username, ins.Name, ins.Company, ins.Blog, ins.Location,
		ins.Email, ins.Bio, ins.Followers, string(organizationsJSON), ins.Readme,
		ins.Commits, ins.Score, ins.PossibleNation, ins.ConfidenceLevel,
	)
	if err != nil {
		return err
	} else {
		log.Printf("insert user %s success\n", ins.Username)
	}

	rstmt, err := tx.PrepareContext(ctx, InsertRepoSQL)
	if err != nil {
		return err
	}
	defer rstmt.Close()
	for _, repo := range ins.Repos {
		// 执行插入语句
		_, err = rstmt.ExecContext(ctx,
			repo.Id, ins.Id, repo.Repo, repo.Star, repo.Fork, repo.Dependent, repo.Commits,
			repo.CommitsTotal, repo.Issue, repo.IssueTotal, repo.PullRequest, repo.PullRequestTotal,
			repo.CodeReview, repo.CodeReviewTotal, repo.LineChange, repo.LineChangeTotal,
		)
		if err != nil {
			return err
		}
	}
	log.Printf("insert repo %s success\n", ins.Username)
	return nil
}

// constructUserRepos 创建用户及其仓库信息，开启多个 goroutine 并发获取数据
func (s *ServiceImpl) constructUserRepos(ctx context.Context, username string) (*user.UserRepos, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in constructUserRepos: %v", r)
		}
	}()

	// 创建User实例
	userins, err := s.ConstructUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	if userins == nil {
		return nil, errors.New("constructUser returned nil for userins")
	}
	if userins.Username == "" {
		return nil, errors.New("username is empty in userins")
	}

	var wg sync.WaitGroup
	mu := sync.Mutex{}
	errCh := make(chan error, 11) // 用于捕获错误

	// 定义结果 map
	starsByRepo, forksByRepo, dependentsByRepo := map[string]int32{}, map[string]int32{}, map[string]int32{}
	totalCommitsByRepo, totalIssuesByRepo := map[string][]int32{}, map[string]int32{}
	userIssuesByRepo, totalPRsByRepo, userPRsByRepo := map[string]int32{}, map[string]int32{}, map[string]int32{}
	totalReviewsByRepo, userReviewsByRepo := map[string]int32{}, map[string]int32{}

	req := &git.GetUsernameRequest{Username: username}

	// 启动 goroutines 获取每个仓库信息
	wg.Add(11)

	go func() {
		defer wg.Done()

		if err := s.InferUserLocationWithLLM(ctx, userins); err != nil {
			errCh <- fmt.Errorf("failed to infer user location: %w", err)
		}

	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetStarsByRepo(ctx, req); err == nil {
			mu.Lock()
			starsByRepo = result.RepoMap
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get stars by repo: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetForksByRepo(ctx, req); err == nil {
			mu.Lock()
			forksByRepo = result.RepoMap
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get forks by repo: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetDependentRepositoriesByRepo(ctx, req); err == nil {
			mu.Lock()
			dependentsByRepo = result.RepoMap
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get dependents by repo: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetLineChangesByRepo(ctx, req); err == nil {
			mu.Lock()
			for repoName, changes := range result.RepoChanges {
				totalCommitsByRepo[repoName] = []int32{changes.TotalChanges, changes.UserChanges, changes.TotalCommits, changes.UserCommits}
			}
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get line changes by repo: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetTotalIssuesByRepo(ctx, req); err == nil {
			mu.Lock()
			totalIssuesByRepo = result.RepoMap
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get total issues by repo: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetUserSolvedIssuesByRepo(ctx, req); err == nil {
			mu.Lock()
			userIssuesByRepo = result.RepoMap
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get user solved issues by repo: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetTotalPullRequestsByRepo(ctx, req); err == nil {
			mu.Lock()
			totalPRsByRepo = result.RepoMap
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get total pull requests by repo: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetUserMergedPullRequestsByRepo(ctx, req); err == nil {
			mu.Lock()
			userPRsByRepo = result.RepoMap
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get user merged pull requests by repo: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetTotalCodeReviewsByRepo(ctx, req); err == nil {
			mu.Lock()
			totalReviewsByRepo = result.RepoMap
			mu.Unlock()
		} else {
			errCh <- fmt.Errorf("failed to get total code reviews by repo: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if result, err := s.svc.GetUserCodeReviewsByRepo(ctx, req); err == nil {
			mu.Lock()
			userReviewsByRepo = result.RepoMap
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
			log.Printf("Error occurred during repository construction: %v", e)
			return nil, e
		}
	}

	// 构建Repo列表
	var repos []*user.Repo
	for repoName := range starsByRepo {
		if len(totalCommitsByRepo[repoName]) < 4 {
			log.Printf("Warning: Skipping repo %s due to insufficient data in totalCommitsByRepo", repoName)
			totalCommitsByRepo[repoName] = []int32{0, 0, 0, 0}
		}

		repo := &user.Repo{
			User_id:          userins.Id,
			Repo:             repoName,
			Star:             int(starsByRepo[repoName]),
			Fork:             int(forksByRepo[repoName]),
			Dependent:        int(dependentsByRepo[repoName]),
			Commits:          int(totalCommitsByRepo[repoName][3]),
			CommitsTotal:     int(totalCommitsByRepo[repoName][2]),
			Issue:            int(userIssuesByRepo[repoName]),
			IssueTotal:       int(totalIssuesByRepo[repoName]),
			PullRequest:      int(userPRsByRepo[repoName]),
			PullRequestTotal: int(totalPRsByRepo[repoName]),
			CodeReview:       int(userReviewsByRepo[repoName]),
			CodeReviewTotal:  int(totalReviewsByRepo[repoName]),
			LineChange:       int(totalCommitsByRepo[repoName][1]),
			LineChangeTotal:  int(totalCommitsByRepo[repoName][0]),
		}
		repo.InjectDefault()
		repos = append(repos, repo)
	}

	userRepos := &user.UserRepos{
		User:  userins,
		Repos: repos,
	}

	if err = calculateOverallScore(userRepos); err != nil {
		return nil, fmt.Errorf("failed to calculate overall score: %w", err)
	}

	return userRepos, nil
}

func (s *ServiceImpl) ConstructUser(ctx context.Context, username string) (*user.User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	req := &git.GetUsernameRequest{Username: username}

	user := &user.User{
		Username: username,
	}

	var wg sync.WaitGroup
	mu := sync.Mutex{}
	errCh := make(chan error, 10) // 缓存通道，用于收集错误

	// 使用 goroutine 发起并发请求
	wg.Add(9)

	go func() {
		defer wg.Done()
		nameResp, err := s.svc.GetName(ctx, req)
		if err != nil {
			errCh <- fmt.Errorf("failed to get user name: %w", err)
			return
		}
		mu.Lock()
		user.Name = nameResp.Result
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		companyResp, err := s.svc.GetCompany(ctx, req)
		if err != nil {
			errCh <- fmt.Errorf("failed to get user company: %w", err)
			return
		}
		mu.Lock()
		user.Company = companyResp.Result
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		locationResp, err := s.svc.GetLocation(ctx, req)
		if err != nil {
			errCh <- fmt.Errorf("failed to get user location: %w", err)
			return
		}
		mu.Lock()
		user.Location = locationResp.Result
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		emailResp, err := s.svc.GetEmail(ctx, req)
		if err != nil {
			errCh <- fmt.Errorf("failed to get user email: %w", err)
			return
		}
		mu.Lock()
		user.Email = emailResp.Result
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		bioResp, err := s.svc.GetBio(ctx, req)
		if err != nil {
			errCh <- fmt.Errorf("failed to get user bio: %w", err)
			return
		}
		mu.Lock()
		user.Bio = bioResp.Result
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		orgsResp, err := s.svc.GetOrganizations(ctx, req)
		if err != nil {
			errCh <- fmt.Errorf("failed to get user organizations: %w", err)
			return
		}
		mu.Lock()
		user.Organizations = orgsResp.Organizations
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		followersResp, err := s.svc.GetFollowers(ctx, req)
		if err != nil {
			errCh <- fmt.Errorf("failed to get user followers: %w", err)
			return
		}
		mu.Lock()
		user.Followers = int(followersResp.Result)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		readmeReq := &git.GetReadmeRequest{
			Username:  username,
			CharLimit: apps.CharLimit,
			RepoLimit: apps.RepoLimit,
		}
		readmeResp, err := s.svc.GetReadme(ctx, readmeReq)
		if err != nil {
			errCh <- fmt.Errorf("failed to get user readme: %w", err)
			return
		}
		mu.Lock()
		user.Readme = readmeResp.Result
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		commitsReq := &git.GetCommitsRequest{
			Username:  username,
			CharLimit: apps.CharLimit,
			RepoLimit: apps.RepoLimit,
		}
		commitsResp, err := s.svc.GetCommits(ctx, commitsReq)
		if err != nil {
			errCh <- fmt.Errorf("failed to get user commits: %w", err)
			return
		}
		mu.Lock()
		user.Commits = commitsResp.Result
		mu.Unlock()
	}()

	// 等待所有 goroutines 完成
	wg.Wait()
	close(errCh) // 关闭错误通道

	// 检查是否有错误
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	// 填充 ID
	user.InjectDefault()

	return user, nil
}

// CalculateOverallScore 计算开发者的最终技术评分
func calculateOverallScore(userRepos *user.UserRepos) error {
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

func (s *ServiceImpl) InferUserLocationWithLLM(ctx context.Context, userins *user.User) error {
	mu := sync.Mutex{}
	inputJSON, err := user.GetUserReposJSONWithRequestDoubao(ctx, userins)
	if err != nil {
		return fmt.Errorf("failed to create JSON request: %w", err)
	}

	//调用 gRPC 服务
	req := &llm.ChatRequest{InputJson: string(inputJSON)}
	resp, err := s.llm.ProcessChatCompletion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to call LLM service: %w", err)
	}

	// 将返回的 JSON 反序列化为 UserResponseByLLM 结构体
	var llmResponse user.UserResponceByLLM
	if err := json.Unmarshal([]byte(resp.OutputJson), &llmResponse); err != nil {
		return fmt.Errorf("failed to unmarshal LLM response: %w", err)
	}

	if llmResponse.PossibleNation == "" || llmResponse.ConfidenceLevel == "" {
		return fmt.Errorf("LLM returned incomplete response: %v", llmResponse)
	}

	if llmResponse.PossibleNation == "N/A" || llmResponse.ConfidenceLevel == "0" {
		log.Printf("Doubao returned no possible nation for user %s, falling back to GPT-4o", userins.Username)

		json, err := user.GetUserReposJSONWithRequestGPT(ctx, userins)
		if err != nil {
			return fmt.Errorf("failed to create JSON request: %w", err)
		}

		resp1, err := user.PostAnalyze(json)
		if err != nil {
			return fmt.Errorf("failed to call GPT service: %w", err)
		}
		nation, level := user.ExtractFields(resp1)
		log.Printf("GPT returned possible nation %s with confidence level %s for user %s", nation, level, userins.Username)
		mu.Lock()
		userins.PossibleNation = nation
		userins.ConfidenceLevel = level
		mu.Unlock()
	} else {
		log.Printf("Doubao returned possible nation %s with confidence level %s for user %s", llmResponse.PossibleNation, llmResponse.ConfidenceLevel, userins.Username)
		mu.Lock()
		userins.PossibleNation = llmResponse.PossibleNation
		userins.ConfidenceLevel = llmResponse.ConfidenceLevel
		mu.Unlock()
	}

	return nil
}
