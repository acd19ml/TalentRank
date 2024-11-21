package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/acd19ml/TalentRank/apps/llm"
	"google.golang.org/grpc"

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

	// 获取用户的所有仓库列表
	repoListResp, err := s.svc.GetRepositories(ctx, &git.GetUsernameRequest{Username: username})
	if err != nil {
		return nil, fmt.Errorf("failed to get repositories: %w", err)
	}
	repoList := repoListResp.Result
	if len(repoList) == 0 {
		return nil, fmt.Errorf("no repositories found for user: %s", username)
	}

	// 初始化结果和错误存储
	repoData := make(map[string]*user.Repo, len(repoList))
	errorCh := make(chan error, len(repoList)*11+1) // 捕获所有错误

	// 推断用户位置
	go func() {
		if err := s.InferUserLocationWithLLM(ctx, userins); err != nil {
			errorCh <- fmt.Errorf("failed to infer user location with LLM: %w", err)
		}
	}()

	// 使用 goroutine 池并发处理每个仓库
	var wg sync.WaitGroup
	repoCh := make(chan string, len(repoList)) // 用于调度仓库

	// 将仓库名发送到 channel 中
	for _, repo := range repoList {
		repoCh <- repo
	}
	close(repoCh)

	// 定义 worker 处理逻辑
	worker := func() {
		for repoName := range repoCh {
			repoInfo := &user.Repo{
				User_id: userins.Id,
				Repo:    repoName,
			}
			repoRequest := &git.RepoRequest{
				Owner:    username,
				RepoName: repoName,
			}

			// 包装所有 RPC 调用函数
			wrappedFuncs := map[string]func(ctx context.Context, req *git.RepoRequest) (*git.IntResponse, error){
				"stars":             wrapRPC(s.svc.GetStarsByRepo),
				"forks":             wrapRPC(s.svc.GetForksByRepo),
				"dependents":        wrapRPC(s.svc.GetDependentRepositoriesByRepo),
				"totalIssues":       wrapRPC(s.svc.GetTotalIssuesByRepo),
				"userIssues":        wrapRPC(s.svc.GetUserSolvedIssuesByRepo),
				"totalPullRequests": wrapRPC(s.svc.GetTotalPullRequestsByRepo),
				"userPullRequests":  wrapRPC(s.svc.GetUserMergedPullRequestsByRepo),
				"totalCodeReviews":  wrapRPC(s.svc.GetTotalCodeReviewsByRepo),
				"userCodeReviews":   wrapRPC(s.svc.GetUserCodeReviewsByRepo),
			}

			// 处理每个调用
			for name, fetchFunc := range wrappedFuncs {
				resp, err := fetchFunc(ctx, repoRequest)
				if err != nil {
					errorCh <- fmt.Errorf("failed to get %s for repo %s: %w", name, repoName, err)
					continue
				}

				switch name {
				case "stars":
					repoInfo.Star = int(resp.Result)
				case "forks":
					repoInfo.Fork = int(resp.Result)
				case "dependents":
					repoInfo.Dependent = int(resp.Result)
				case "totalIssues":
					repoInfo.IssueTotal = int(resp.Result)
				case "userIssues":
					repoInfo.Issue = int(resp.Result)
				case "totalPullRequests":
					repoInfo.PullRequestTotal = int(resp.Result)
				case "userPullRequests":
					repoInfo.PullRequest = int(resp.Result)
				case "totalCodeReviews":
					repoInfo.CodeReviewTotal = int(resp.Result)
				case "userCodeReviews":
					repoInfo.CodeReview = int(resp.Result)
				}
			}

			// 获取增删行数和提交信息
			resp, err := s.svc.GetLineChangesCommitsByRepo(ctx, repoRequest)
			if err == nil {
				repoInfo.LineChange = int(resp.UserChanges)
				repoInfo.LineChangeTotal = int(resp.TotalChanges)
				repoInfo.Commits = int(resp.UserCommits)
				repoInfo.CommitsTotal = int(resp.TotalCommits)
			} else {
				errorCh <- fmt.Errorf("failed to get line changes/commits for repo %s: %w", repoName, err)
			}

			// 注入默认值并存储
			repoInfo.InjectDefault()
			repoData[repoName] = repoInfo
		}
		wg.Done()
	}

	// 启动 worker 池
	numWorkers := 10 // 可根据实际情况调整
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker()
	}

	// 等待所有 worker 完成
	wg.Wait()
	close(errorCh)

	// 检查是否有错误
	for e := range errorCh {
		if e != nil {
			log.Printf("Error occurred during repository construction: %v", e)
			return nil, e
		}
	}

	// 构建Repo列表
	var repos []*user.Repo
	for _, repoInfo := range repoData {
		repos = append(repos, repoInfo)
	}

	userRepos := &user.UserRepos{
		User:  userins,
		Repos: repos,
	}

	// 计算总体评分
	if err = CalculateOverallScore(userRepos); err != nil {
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
		user.Organizations = orgsResp.Result
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

// wrapRPC 包装带 grpc.CallOption 的 RPC 函数，忽略可选参数
func wrapRPC[T any](rpcFunc func(ctx context.Context, req *T, opts ...grpc.CallOption) (*git.IntResponse, error)) func(ctx context.Context, req *T) (*git.IntResponse, error) {
	return func(ctx context.Context, req *T) (*git.IntResponse, error) {
		return rpcFunc(ctx, req) // 忽略 opts
	}
}
