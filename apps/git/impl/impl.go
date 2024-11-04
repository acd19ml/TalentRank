package impl

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/acd19ml/TalentRank/apps/git"
	"github.com/google/go-github/github"
)

// GetRepositories gRPC 实现：获取用户所有仓库名称，使用缓存
func (g *Service) GetRepositories(ctx context.Context, req *git.GetUsernameRequest) (*git.RepoResponse, error) {

	if err := g.initCache(ctx, req.Username); err != nil {
		return nil, err
	}

	return &git.RepoResponse{Repos: g.reposCache}, nil
}

func (s *Service) initCache(ctx context.Context, username string) error {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	// 如果缓存已经存在并且与当前用户名匹配，跳过初始化
	if s.cacheUsername == username && len(s.reposCache) > 0 {
		return nil
	}

	// 调用 fetchRepositories 并缓存结果
	reposList, err := s.fetchRepositories(ctx, username)
	if err != nil {
		return err
	}
	s.reposCache = reposList
	s.cacheUsername = username
	return nil
}

// fetchRepositories 获取用户的所有仓库名称
func (s *Service) fetchRepositories(ctx context.Context, username string) ([]string, error) {
	var reposList []string
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	for {
		repos, resp, err := s.client.Repositories.List(ctx, username, opts)
		if err != nil {
			return nil, err
		}

		for _, repo := range repos {
			isContributor, err := s.checkIfUserIsContributor(ctx, username, repo.GetOwner().GetLogin(), repo.GetName())
			if err != nil {
				return nil, fmt.Errorf("failed to check if user is a contributor: %v", err)
			}
			if isContributor {
				reposList = append(reposList, repo.GetName())
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return reposList, nil
}

// checkIfUserIsContributor 检查用户是否为仓库的贡献者
func (s *Service) checkIfUserIsContributor(ctx context.Context, username, owner, repo string) (bool, error) {
	contributors, _, err := s.client.Repositories.ListContributors(ctx, owner, repo, nil)
	if err != nil {
		return false, fmt.Errorf("failed to list contributors: %v", err)
	}

	for _, contributor := range contributors {
		if contributor.GetLogin() == username {
			return true, nil
		}
	}

	return false, nil
}

func (g *Service) GetUser(ctx context.Context, req *git.GetUsernameRequest) (*git.UserResponse, error) {
	user, _, err := g.client.Users.Get(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	// 将获取到的用户信息映射到 UserResponse 中
	return &git.UserResponse{
		Name:      user.GetName(),
		Company:   user.GetCompany(),
		Location:  user.GetLocation(),
		Email:     user.GetEmail(),
		Bio:       user.GetBio(),
		Followers: int32(user.GetFollowers()),
	}, nil
}

func (g *Service) GetName(ctx context.Context, req *git.GetUsernameRequest) (*git.StringResponse, error) {
	user, err := g.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &git.StringResponse{Result: user.Name}, nil
}

func (g *Service) GetCompany(ctx context.Context, username *git.GetUsernameRequest) (*git.StringResponse, error) {
	user, err := g.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}
	return &git.StringResponse{Result: user.Company}, nil
}

func (g *Service) GetLocation(ctx context.Context, username *git.GetUsernameRequest) (*git.StringResponse, error) {
	user, err := g.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}
	location := user.GetLocation()
	if location == "" {
		log.Printf("warning: location for user %s is empty", username)
	}
	return &git.StringResponse{Result: user.Location}, nil
}

func (g *Service) GetEmail(ctx context.Context, username *git.GetUsernameRequest) (*git.StringResponse, error) {
	user, err := g.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}
	return &git.StringResponse{Result: user.Email}, nil
}

func (g *Service) GetBio(ctx context.Context, username *git.GetUsernameRequest) (*git.StringResponse, error) {
	user, err := g.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}
	return &git.StringResponse{Result: user.Bio}, nil
}

func (g *Service) GetOrganizations(ctx context.Context, req *git.GetUsernameRequest) (*git.OrgListResponse, error) {
	var orgsList []string
	opts := &github.ListOptions{PerPage: 50} // 设置分页参数

	// 获取所有组织
	for {
		orgs, resp, err := g.client.Organizations.List(ctx, req.Username, opts)
		if err != nil {
			return nil, err
		}

		for _, org := range orgs {
			orgsList = append(orgsList, org.GetLogin())
		}

		// 检查是否有下一页
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return &git.OrgListResponse{Organizations: orgsList}, nil
}

func (g *Service) GetReadme(ctx context.Context, req *git.GetReadmeRequest) (*git.StringResponse, error) {
	repos, err := g.GetRepositories(ctx, &git.GetUsernameRequest{Username: req.Username})
	if err != nil {
		return nil, err
	}

	var contents string
	for i, repo := range repos.Repos {
		if i >= int(req.RepoLimit) {
			break
		}

		readme, _, err := g.client.Repositories.GetReadme(ctx, req.Username, repo, nil)
		if err != nil {
			if githubErr, ok := err.(*github.ErrorResponse); ok && githubErr.Response.StatusCode == 404 {
				continue
			}
			return nil, err
		}

		content, err := readme.GetContent()
		if err != nil {
			return nil, err
		}

		if len(content) > int(req.CharLimit) {
			content = content[:req.CharLimit] + "..."
		}

		if len(contents)+len(content) > int(req.CharLimit*req.RepoLimit) {
			contents += content[:int(req.CharLimit)*int(req.RepoLimit)-len(contents)] + "..."
			break
		}

		contents += content
	}

	return &git.StringResponse{Result: contents}, nil
}

func (g *Service) GetCommits(ctx context.Context, req *git.GetCommitsRequest) (*git.StringResponse, error) {
	repos, err := g.GetRepositories(ctx, &git.GetUsernameRequest{Username: req.Username})
	if err != nil {
		return nil, err
	}

	var allCommits string

	for i, repo := range repos.Repos {
		if i >= int(req.RepoLimit) {
			break
		}

		opts := &github.CommitsListOptions{
			Author:      req.Username,
			ListOptions: github.ListOptions{PerPage: 10},
		}

		var commits []*github.RepositoryCommit
		var lastErr error

		// 尝试三次请求
		for attempt := 1; attempt <= 3; attempt++ {
			commits, _, err = g.client.Repositories.ListCommits(ctx, req.Username, repo, opts)
			if err == nil {
				lastErr = nil
				break
			}
			lastErr = err
			time.Sleep(2 * time.Second)
		}

		if lastErr != nil {
			fmt.Printf("Skipping repo '%s' due to persistent errors: %v\n", repo, lastErr)
			continue
		}

		var repoCommits string
		for _, commit := range commits {
			message := commit.GetCommit().GetMessage()
			if len(repoCommits)+len(message) > int(req.CharLimit) {
				break
			}
			repoCommits += message + "\n"
		}

		if len(repoCommits) > int(req.CharLimit) {
			repoCommits = repoCommits[:req.CharLimit] + "...\n"
		}

		allCommits += fmt.Sprintf("Repo: %s\n%s\n", repo, repoCommits)
		if len(allCommits) > int(req.CharLimit*req.RepoLimit) {
			allCommits = allCommits[:req.CharLimit*req.RepoLimit] + "..."
			break
		}
	}

	return &git.StringResponse{Result: allCommits}, nil
}

func (g *Service) GetFollowers(ctx context.Context, req *git.GetUsernameRequest) (*git.IntResponse, error) {
	user, err := g.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &git.IntResponse{Result: user.Followers}, err

}

func (g *Service) GetRepoStars(ctx context.Context, req *git.RepoRequest) (*git.IntResponse, error) {
	// 获取指定仓库的信息
	repo, _, err := g.client.Repositories.Get(ctx, req.Owner, req.RepoName)
	if err != nil {
		return nil, err
	}

	// 返回包含 star 数量的 IntResponse
	return &git.IntResponse{Result: int32(repo.GetStargazersCount())}, nil
}

// 返回单个fork数量
func (g *Service) GetRepoForks(ctx context.Context, req *git.RepoRequest) (*git.IntResponse, error) {
	// 获取指定仓库的信息
	repo, _, err := g.client.Repositories.Get(ctx, req.Owner, req.RepoName)
	if err != nil {
		return nil, err
	}

	// 返回包含 fork 数量的 IntResponse
	return &git.IntResponse{Result: int32(repo.GetForksCount())}, nil
}

func (g *Service) GetStarsByRepo(ctx context.Context, req *git.GetUsernameRequest) (*git.RepoIntMapResponse, error) {
	// 初始化一个 map，用于存储仓库名和 stars 数量
	repoStarsMap := make(map[string]int32)

	// 获取该用户的所有仓库名称
	repos, err := g.GetRepositories(ctx, req)
	if err != nil {
		return nil, err
	}

	// 遍历每个仓库，获取 stars 数量
	for _, repoName := range repos.Repos {
		// 调用 GetRepoStars 获取当前仓库的 star 数量
		starsResp, err := g.GetRepoStars(ctx, &git.RepoRequest{Owner: req.Username, RepoName: repoName})
		if err != nil {
			continue // 跳过出错的仓库
		}
		repoStarsMap[repoName] = starsResp.Result
	}

	return &git.RepoIntMapResponse{RepoMap: repoStarsMap}, nil
}

// GetForksByRepo 获取指定用户的所有仓库 fork 数量
func (g *Service) GetForksByRepo(ctx context.Context, req *git.GetUsernameRequest) (*git.RepoIntMapResponse, error) {
	// 初始化一个 map，用于存储仓库名和 forks 数量
	repoForksMap := make(map[string]int32)

	// 获取该用户的所有仓库名称
	repos, err := g.GetRepositories(ctx, req)
	if err != nil {
		return nil, err
	}

	// 遍历每个仓库，获取 forks 数量
	for _, repoName := range repos.Repos {
		// 调用 GetRepoForks 获取当前仓库的 forks 数量
		forksResp, err := g.GetRepoForks(ctx, &git.RepoRequest{Owner: req.Username, RepoName: repoName})
		if err != nil {
			continue // 跳过出错的仓库
		}
		repoForksMap[repoName] = forksResp.Result
	}

	return &git.RepoIntMapResponse{RepoMap: repoForksMap}, nil
}

func (g *Service) GetTotalCommitsByRepo(ctx context.Context, req *git.GetUsernameRequest) (*git.RepoIntMapResponse, error) {
	repoCommitsCount := make(map[string]int32)

	// 获取用户的所有仓库
	repos, err := g.GetRepositories(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get repositories for user %s: %w", req.Username, err)
	}

	// 遍历每个仓库，抓取并解析提交总数
	for _, repo := range repos.Repos {
		url := fmt.Sprintf("https://github.com/%s/%s", req.Username, repo)

		// 发起 HTTP GET 请求获取 HTML
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch repo page for %s: %w", repo, err)
		}
		defer resp.Body.Close()

		// 使用 goquery 加载 HTML 文档
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to parse HTML for repo %s: %w", repo, err)
		}

		// 查找提交总数元素
		var commitCount int32
		doc.Find("span[data-component='text'] .fgColor-default").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			parts := strings.Fields(text)
			if len(parts) > 0 {
				countStr := strings.ReplaceAll(parts[0], ",", "")
				count, err := strconv.Atoi(countStr)
				if err == nil {
					log.Printf("failed to parse commit count for repo %s: %v", repo, err)
				}
				commitCount = int32(count)
			}
		})

		// 将结果存入 map
		repoCommitsCount[repo] = commitCount
	}

	return &git.RepoIntMapResponse{RepoMap: repoCommitsCount}, nil
}

func (g *Service) GetUserCommitsByRepo(ctx context.Context, req *git.GetUsernameRequest) (*git.RepoIntMapResponse, error) {
	userCommitsCount := make(map[string]int32)

	// 获取用户的所有仓库
	repos, err := g.GetRepositories(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get repositories for user %s: %w", req.Username, err)
	}

	for _, repo := range repos.Repos {
		opts := &github.CommitsListOptions{
			Author:      req.Username,
			ListOptions: github.ListOptions{PerPage: 100},
		}

		var userCommits int32

		for {
			commits, resp, err := g.client.Repositories.ListCommits(ctx, req.Username, repo, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to get commits for repo %s by user %s: %w", repo, req.Username, err)
			}

			userCommits += int32(len(commits))

			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}

		userCommitsCount[repo] = userCommits
	}

	return &git.RepoIntMapResponse{RepoMap: userCommitsCount}, nil
}

func (g *Service) GetTotalIssuesByRepo(ctx context.Context, req *git.GetUsernameRequest) (*git.RepoIntMapResponse, error) {
	issuesCount := make(map[string]int32)

	// 获取用户的所有仓库
	repos, err := g.GetRepositories(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get repositories for user %s: %w", req.Username, err)
	}

	for _, repo := range repos.Repos {
		opts := &github.IssueListByRepoOptions{
			State:       "closed",
			ListOptions: github.ListOptions{PerPage: 100},
		}

		var totalIssues int32

		for {
			issues, resp, err := g.client.Issues.ListByRepo(ctx, req.Username, repo, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list issues for repo %s: %w", repo, err)
			}
			totalIssues += int32(len(issues))

			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}

		issuesCount[repo] = totalIssues
	}

	return &git.RepoIntMapResponse{RepoMap: issuesCount}, nil
}

func (g *Service) GetUserSolvedIssuesByRepo(ctx context.Context, req *git.GetUsernameRequest) (*git.RepoIntMapResponse, error) {
	userIssuesCount := make(map[string]int32)
	repos, err := g.GetRepositories(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get repositories: %w", err)
	}

	for _, repo := range repos.Repos {
		opts := &github.IssueListByRepoOptions{
			State:       "closed",
			ListOptions: github.ListOptions{PerPage: 100},
		}
		userIssues := int32(0)

		for {
			issues, resp, err := g.client.Issues.ListByRepo(ctx, req.Username, repo, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list issues for repo %s: %w", repo, err)
			}
			for _, issue := range issues {
				if issue.Assignee != nil && issue.Assignee.GetLogin() == req.Username {
					userIssues++
				}
			}
			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}
		userIssuesCount[repo] = userIssues
	}

	return &git.RepoIntMapResponse{RepoMap: userIssuesCount}, nil
}

func (g *Service) GetTotalPullRequestsByRepo(ctx context.Context, req *git.GetUsernameRequest) (*git.RepoIntMapResponse, error) {
	prCount := make(map[string]int32)
	repos, err := g.GetRepositories(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get repositories: %w", err)
	}

	for _, repo := range repos.Repos {
		opts := &github.PullRequestListOptions{
			State:       "closed",
			ListOptions: github.ListOptions{PerPage: 100},
		}
		totalPRs := int32(0)

		for {
			prs, resp, err := g.client.PullRequests.List(ctx, req.Username, repo, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list pull requests for repo %s: %w", repo, err)
			}
			totalPRs += int32(len(prs))
			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}
		prCount[repo] = totalPRs
	}

	return &git.RepoIntMapResponse{RepoMap: prCount}, nil
}

func (g *Service) GetUserMergedPullRequestsByRepo(ctx context.Context, req *git.GetUsernameRequest) (*git.RepoIntMapResponse, error) {
	userPRCount := make(map[string]int32)
	repos, err := g.GetRepositories(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get repositories: %w", err)
	}

	for _, repo := range repos.Repos {
		opts := &github.PullRequestListOptions{
			State:       "closed",
			ListOptions: github.ListOptions{PerPage: 100},
		}
		userPRs := int32(0)

		for {
			prs, resp, err := g.client.PullRequests.List(ctx, req.Username, repo, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list pull requests for repo %s: %w", repo, err)
			}
			for _, pr := range prs {
				if pr.User != nil && pr.User.GetLogin() == req.Username {
					userPRs++
				}
			}
			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}
		userPRCount[repo] = userPRs
	}

	return &git.RepoIntMapResponse{RepoMap: userPRCount}, nil
}

func (g *Service) GetTotalCodeReviewsByRepo(ctx context.Context, req *git.GetUsernameRequest) (*git.RepoIntMapResponse, error) {
	repos, err := g.GetRepositories(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get repositories for user %s: %w", req.Username, err)
	}

	reviewCount := make(map[string]int32)

	for _, repo := range repos.Repos {
		opts := &github.PullRequestListOptions{
			ListOptions: github.ListOptions{PerPage: 100},
		}
		totalReviews := int32(0)

		for {
			pullRequests, resp, err := g.client.PullRequests.List(ctx, req.Username, repo, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list pull requests for repo %s: %w", repo, err)
			}

			for _, pr := range pullRequests {
				reviews, _, err := g.client.PullRequests.ListReviews(ctx, req.Username, repo, pr.GetNumber(), nil)
				if err != nil {
					return nil, fmt.Errorf("failed to list reviews for repo %s PR #%d: %w", repo, pr.GetNumber(), err)
				}
				totalReviews += int32(len(reviews))
			}

			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}

		reviewCount[repo] = totalReviews
	}

	return &git.RepoIntMapResponse{RepoMap: reviewCount}, nil
}

func (g *Service) GetUserCodeReviewsByRepo(ctx context.Context, req *git.GetUsernameRequest) (*git.RepoIntMapResponse, error) {
	repos, err := g.GetRepositories(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get repositories for user %s: %w", req.Username, err)
	}

	userReviewCount := make(map[string]int32)

	for _, repo := range repos.Repos {
		opts := &github.PullRequestListOptions{
			ListOptions: github.ListOptions{PerPage: 100},
		}
		userReviews := int32(0)

		for {
			pullRequests, resp, err := g.client.PullRequests.List(ctx, req.Username, repo, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list pull requests for repo %s: %w", repo, err)
			}

			for _, pr := range pullRequests {
				reviews, _, err := g.client.PullRequests.ListReviews(ctx, req.Username, repo, pr.GetNumber(), nil)
				if err != nil {
					return nil, fmt.Errorf("failed to list reviews for repo %s PR #%d: %w", repo, pr.GetNumber(), err)
				}

				for _, review := range reviews {
					if review.GetUser().GetLogin() == req.Username {
						userReviews++
					}
				}
			}

			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}

		userReviewCount[repo] = userReviews
	}

	return &git.RepoIntMapResponse{RepoMap: userReviewCount}, nil
}

// GetDependentRepositorie 获取仓库被依赖数量
func GetDependentRepositorie(url string) (*git.IntResponse, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch pages: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	repoCount := 0
	re := regexp.MustCompile(`(\d{1,3}(?:,\d{3})*)\s+Repositories`)

	doc.Find("a.btn-link.selected").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			countStr := strings.ReplaceAll(matches[1], ",", "")
			fmt.Sscanf(countStr, "%d", &repoCount)
		}
	})

	return &git.IntResponse{Result: int32(repoCount)}, nil
}

// GetDependentRepositoriesByRepo 获取每个仓库的依赖数量
func (s *Service) GetDependentRepositoriesByRepo(ctx context.Context, req *git.GetUsernameRequest) (*git.RepoIntMapResponse, error) {
	repos, err := s.GetRepositories(ctx, req)
	if err != nil {
		return nil, err
	}

	repoDependents := make(map[string]int32)
	for _, repo := range repos.Repos {
		countResp, err := GetDependentRepositorie(fmt.Sprintf("https://github.com/%s/%s/network/dependents", req.Username, repo))
		if err != nil {
			continue
		}
		repoDependents[repo] = countResp.Result
	}

	return &git.RepoIntMapResponse{RepoMap: repoDependents}, nil
}

// getLineChanges 获取仓库的总增删行数和指定用户的增删行数，并统计提交次数
func (g *Service) getLineChanges(ctx context.Context, repoOwner, repoName, username string) (int32, int32, int32, int32, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", repoOwner, repoName)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to create request: %v", err)
	}

	var commits []Commit
	_, err = g.client.Do(ctx, req, &commits)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("request failed: %v", err)
	}

	totalAdditions, totalDeletions := int32(0), int32(0)
	userAdditions, userDeletions := int32(0), int32(0)
	totalCommits, userCommits := int32(0), int32(0)

	for _, commit := range commits {
		totalCommits++

		detailURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", repoOwner, repoName, commit.Sha)
		detailReq, err := http.NewRequestWithContext(ctx, "GET", detailURL, nil)
		if err != nil {
			log.Printf("failed to create request for commit detail: %v", err)
			continue
		}

		var commitDetail CommitDetail
		_, err = g.client.Do(ctx, detailReq, &commitDetail)
		if err != nil {
			log.Printf("failed to fetch commit details: %v", err)
			continue
		}

		totalAdditions += int32(commitDetail.Stats.Additions)
		totalDeletions += int32(commitDetail.Stats.Deletions)

		if commit.Author.Login == username {
			userCommits++
			userAdditions += int32(commitDetail.Stats.Additions)
			userDeletions += int32(commitDetail.Stats.Deletions)
		}
	}

	return totalAdditions + totalDeletions, userAdditions + userDeletions, totalCommits, userCommits, nil
}

// GetLineChangesByRepo 获取用户所有仓库的增删行数信息，包含总提交和用户提交
func (g *Service) GetLineChangesByRepo(ctx context.Context, req *git.GetUsernameRequest) (*git.RepoLineChangesResponse, error) {
	repos, err := g.GetRepositories(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get repositories: %v", err)
	}

	lineChanges := make(map[string]*git.LineChangeStats)
	for _, repo := range repos.Repos {
		totalChanges, userChanges, totalCommits, userCommits, err := g.getLineChanges(ctx, req.Username, repo, req.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to get line changes for repo %s: %v", repo, err)
		}

		// 将增删行数和提交信息存储到 LineChangeStats 中
		lineChanges[repo] = &git.LineChangeStats{
			TotalChanges: int32(totalChanges),
			UserChanges:  int32(userChanges),
			TotalCommits: int32(totalCommits),
			UserCommits:  int32(userCommits),
		}
	}

	return &git.RepoLineChangesResponse{RepoChanges: lineChanges}, nil
}
