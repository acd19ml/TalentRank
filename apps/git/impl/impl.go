package impl

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/acd19ml/TalentRank/apps/git"
	"github.com/google/go-github/github"
)

var (
	wg sync.WaitGroup
)

// checkIfUserIsContributor 检查用户是否为仓库的贡献者
func (s *Service) checkIfUserIsContributor(ctx context.Context, username, owner, repo string) (bool, error) {
	client := s.getClientFromContext(ctx)
	contributors, _, err := client.Repositories.ListContributors(ctx, owner, repo, nil)
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

// GetRepositories gRPC 实现：获取用户所有仓库名称
func (s *Service) GetRepositories(ctx context.Context, req *git.GetUsernameRequest) (*git.StringListResponse, error) {

	client := s.getClientFromContext(ctx)

	var reposList []string
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	for {
		repos, resp, err := client.Repositories.List(ctx, req.Username, opts)
		if err != nil {
			return nil, err
		}

		for _, repo := range repos {
			isContributor, err := s.checkIfUserIsContributor(ctx, req.Username, repo.GetOwner().GetLogin(), repo.GetName())
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

	return &git.StringListResponse{Result: reposList}, nil
}

func (s *Service) GetUser(ctx context.Context, req *git.GetUsernameRequest) (*git.UserResponse, error) {
	client := s.getClientFromContext(ctx)
	user, _, err := client.Users.Get(ctx, req.Username)
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

func (s *Service) GetName(ctx context.Context, req *git.GetUsernameRequest) (*git.StringResponse, error) {
	user, err := s.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &git.StringResponse{Result: user.Name}, nil
}

func (s *Service) GetCompany(ctx context.Context, username *git.GetUsernameRequest) (*git.StringResponse, error) {
	user, err := s.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}
	return &git.StringResponse{Result: user.Company}, nil
}

func (s *Service) GetLocation(ctx context.Context, username *git.GetUsernameRequest) (*git.StringResponse, error) {
	user, err := s.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}
	return &git.StringResponse{Result: user.Location}, nil
}

func (s *Service) GetEmail(ctx context.Context, username *git.GetUsernameRequest) (*git.StringResponse, error) {
	user, err := s.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}
	return &git.StringResponse{Result: user.Email}, nil
}

func (s *Service) GetBio(ctx context.Context, username *git.GetUsernameRequest) (*git.StringResponse, error) {
	user, err := s.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}
	// 清理 user.Bio 中的非 UTF-8 字符
	cleanedBio := cleanInvalidUTF8(user.Bio)

	return &git.StringResponse{Result: cleanedBio}, nil
}

func (s *Service) GetOrganizations(ctx context.Context, req *git.GetUsernameRequest) (*git.StringListResponse, error) {
	var orgsList []string
	opts := &github.ListOptions{PerPage: 50} // 设置分页参数

	client := s.getClientFromContext(ctx)
	// 获取所有组织
	for {
		orgs, resp, err := client.Organizations.List(ctx, req.Username, opts)
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

	return &git.StringListResponse{Result: orgsList}, nil
}

func cleanInvalidUTF8(input string) string {
	if utf8.ValidString(input) {
		return input
	}

	// 创建一个新的字符串，过滤掉无效的 UTF-8 字符
	validRunes := make([]rune, 0, len(input))
	for i, r := range input {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(input[i:])
			if size == 1 {
				log.Printf("Skipping invalid UTF-8 character at index %d", i)
				continue // 跳过无效的 UTF-8 字符
			}
		}
		validRunes = append(validRunes, r)
	}
	return string(validRunes)
}

func (s *Service) GetReadme(ctx context.Context, req *git.GetReadmeRequest) (*git.StringResponse, error) {
	client := s.getClientFromContext(ctx)
	repos, err := s.GetRepositories(ctx, &git.GetUsernameRequest{Username: req.Username})
	if err != nil {
		return nil, err
	}

	var contents string
	for i, repo := range repos.Result {
		if i >= int(req.RepoLimit) {
			break
		}

		readme, _, err := client.Repositories.GetReadme(ctx, req.Username, repo, nil)
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

		// 检查和清理非 UTF-8 字符
		content = cleanInvalidUTF8(content)

		// 添加日志，调试每段内容的有效性
		if !utf8.ValidString(content) {
			log.Println("Warning: content still contains invalid UTF-8 characters after cleaning")
			continue
		}

		// 截取内容到指定字符限制
		if len(content) > int(req.CharLimit) {
			content = content[:req.CharLimit] + "..."
		}

		if len(contents)+len(content) > int(req.CharLimit*req.RepoLimit) {
			contents += content[:int(req.CharLimit)*int(req.RepoLimit)-len(contents)] + "..."
			break
		}

		contents += content
	}
	// 再次清理和验证整个内容
	contents = cleanInvalidUTF8(contents)
	if !utf8.ValidString(contents) {
		return nil, fmt.Errorf("resulting content contains invalid UTF-8 characters")
	}

	return &git.StringResponse{Result: contents}, nil
}

func (s *Service) GetCommits(ctx context.Context, req *git.GetCommitsRequest) (*git.StringResponse, error) {
	client := s.getClientFromContext(ctx)
	repos, err := s.GetRepositories(ctx, &git.GetUsernameRequest{Username: req.Username})
	if err != nil {
		return nil, err
	}

	var allCommits string

	for i, repo := range repos.Result {
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
			commits, _, err = client.Repositories.ListCommits(ctx, req.Username, repo, opts)
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

func (s *Service) GetFollowers(ctx context.Context, req *git.GetUsernameRequest) (*git.IntResponse, error) {
	user, err := s.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &git.IntResponse{Result: user.Followers}, err

}

func (s *Service) GetStarsByRepo(ctx context.Context, req *git.RepoRequest) (*git.IntResponse, error) {
	client := s.getClientFromContext(ctx)
	// 获取指定仓库的信息
	repo, _, err := client.Repositories.Get(ctx, req.Owner, req.RepoName)
	if err != nil {
		return nil, err
	}

	// 返回包含 star 数量的 IntResponse
	return &git.IntResponse{Result: int32(repo.GetStargazersCount())}, nil
}

// 返回单个fork数量
func (s *Service) GetForksByRepo(ctx context.Context, req *git.RepoRequest) (*git.IntResponse, error) {
	client := s.getClientFromContext(ctx)
	// 获取指定仓库的信息
	repo, _, err := client.Repositories.Get(ctx, req.Owner, req.RepoName)
	if err != nil {
		return nil, err
	}

	// 返回包含 fork 数量的 IntResponse
	return &git.IntResponse{Result: int32(repo.GetForksCount())}, nil
}

// GetTotalCommitsByRepo 获取指定用户的所有仓库的提交总数
func (s *Service) GetTotalCommitsByRepo(ctx context.Context, req *git.RepoRequest) (*git.IntResponse, error) {

	url := fmt.Sprintf("https://github.com/%s/%s", req.Owner, req.RepoName)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("failed to fetch repo page for %s: %v", req.RepoName, err)
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("failed to parse HTML for repo %s: %v", req.RepoName, err)
		return nil, err
	}

	var commitCount int32
	doc.Find("span[data-component='text'] .fgColor-default").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		parts := strings.Fields(text)
		if len(parts) > 0 {
			countStr := strings.ReplaceAll(parts[0], ",", "")
			count, err := strconv.Atoi(countStr)
			if err == nil {
				commitCount = int32(count)
			}
		}
	})

	return &git.IntResponse{Result: commitCount}, nil
}

// GetUserCommitsByRepo 获取指定用户指定repo的提交数量
func (s *Service) GetUserCommitsByRepo(ctx context.Context, req *git.RepoRequest) (*git.IntResponse, error) {
	client := s.getClientFromContext(ctx)

	opts := &github.CommitsListOptions{
		Author:      req.Owner,
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var userCommits int32
	for {
		commits, resp, err := client.Repositories.ListCommits(ctx, req.Owner, req.RepoName, opts)
		if err != nil {
			log.Printf("failed to get commits for repo %s by user %s: %v", req.RepoName, req.Owner, err)
			return nil, err
		}
		userCommits += int32(len(commits))
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return &git.IntResponse{Result: userCommits}, nil
}

// GetTotalIssuesByRepo 获取指定用户的所有仓库的总 issues 数量
func (s *Service) GetTotalIssuesByRepo(ctx context.Context, req *git.RepoRequest) (*git.IntResponse, error) {
	client := s.getClientFromContext(ctx)

	opts := &github.IssueListByRepoOptions{
		State:       "closed",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var totalIssues int32
	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, req.Owner, req.RepoName, opts)
		if err != nil {
			return nil, err
		}
		totalIssues += int32(len(issues))
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return &git.IntResponse{Result: totalIssues}, nil
}

// GetUserSolvedIssuesByRepo 获取指定用户解决的所有仓库的 issues 数量
func (s *Service) GetUserSolvedIssuesByRepo(ctx context.Context, req *git.RepoRequest) (*git.IntResponse, error) {
	client := s.getClientFromContext(ctx)

	opts := &github.IssueListByRepoOptions{
		State:       "closed",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	userIssues := int32(0)
	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, req.Owner, req.RepoName, opts)
		if err != nil {
			return nil, err
		}
		for _, issue := range issues {
			if issue.Assignee != nil && issue.Assignee.GetLogin() == req.Owner {
				userIssues++
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	wg.Wait()
	return &git.IntResponse{Result: userIssues}, nil
}

// GetTotalPullRequestsByRepo 获取指定用户所有仓库的总 PR 数量
func (s *Service) GetTotalPullRequestsByRepo(ctx context.Context, req *git.RepoRequest) (*git.IntResponse, error) {
	client := s.getClientFromContext(ctx)

	opts := &github.PullRequestListOptions{
		State:       "closed",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	totalPRs := int32(0)
	for {
		prs, resp, err := client.PullRequests.List(ctx, req.Owner, req.RepoName, opts)
		if err != nil {
			return nil, err
		}
		totalPRs += int32(len(prs))
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return &git.IntResponse{Result: totalPRs}, nil
}

// GetUserMergedPullRequestsByRepo 获取指定用户合并的所有仓库的 PR 数量
func (s *Service) GetUserMergedPullRequestsByRepo(ctx context.Context, req *git.RepoRequest) (*git.IntResponse, error) {
	client := s.getClientFromContext(ctx)

	opts := &github.PullRequestListOptions{
		State:       "closed",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	userPRs := int32(0)
	for {
		prs, resp, err := client.PullRequests.List(ctx, req.Owner, req.RepoName, opts)
		if err != nil {
			return nil, err
		}
		for _, pr := range prs {
			if pr.User != nil && pr.User.GetLogin() == req.Owner {
				userPRs++
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return &git.IntResponse{Result: userPRs}, nil
}

// GetTotalCodeReviewsByRepo 获取每个仓库的代码审查总数
func (s *Service) GetTotalCodeReviewsByRepo(ctx context.Context, req *git.RepoRequest) (*git.IntResponse, error) {
	client := s.getClientFromContext(ctx)

	opts := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	totalReviews := int32(0)

	for {
		pullRequests, resp, err := client.PullRequests.List(ctx, req.Owner, req.RepoName, opts)
		if err != nil {
			return nil, err
		}

		for _, pr := range pullRequests {
			reviews, _, err := client.PullRequests.ListReviews(ctx, req.Owner, req.RepoName, pr.GetNumber(), nil)
			if err != nil {
				return nil, err
			}
			totalReviews += int32(len(reviews))
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return &git.IntResponse{Result: totalReviews}, nil
}

// GetUserCodeReviewsByRepo 获取用户在每个仓库中的代码审查数量
func (s *Service) GetUserCodeReviewsByRepo(ctx context.Context, req *git.RepoRequest) (*git.IntResponse, error) {
	client := s.getClientFromContext(ctx)

	opts := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	userReviews := int32(0)

	for {
		pullRequests, resp, err := client.PullRequests.List(ctx, req.Owner, req.RepoName, opts)
		if err != nil {
			return nil, err
		}

		for _, pr := range pullRequests {
			reviews, _, err := client.PullRequests.ListReviews(ctx, req.Owner, req.RepoName, pr.GetNumber(), nil)
			if err != nil {
				return nil, err
			}

			for _, review := range reviews {
				if review.GetUser().GetLogin() == req.Owner {
					userReviews++
				}
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return &git.IntResponse{Result: userReviews}, nil
}

// GetDependentRepositoriesByRepo 获取仓库被依赖数量
func getDependentRepositorie(url string) (*git.IntResponse, error) {
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
func (s *Service) GetDependentRepositoriesByRepo(ctx context.Context, req *git.RepoRequest) (*git.IntResponse, error) {

	countResp, err := getDependentRepositorie(fmt.Sprintf("https://github.com/%s/%s/network/dependents", req.Owner, req.RepoName))
	if err != nil {
		return nil, err
	}

	return countResp, nil
}

// getLineChanges 获取仓库的总增删行数和指定用户的增删行数，并统计提交次数
func (s *Service) getLineChanges(ctx context.Context, repoOwner, repoName, username string) (int32, int32, int32, int32, error) {
	client := s.getClientFromContext(ctx)
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", repoOwner, repoName)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to create request: %v", err)
	}

	var commits []Commit
	_, err = client.Do(ctx, req, &commits)
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
		_, err = client.Do(ctx, detailReq, &commitDetail)
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
func (s *Service) GetLineChangesCommitsByRepo(ctx context.Context, req *git.RepoRequest) (*git.RepoLineChangesCommitsResponse, error) {

	totalChanges, userChanges, totalCommits, userCommits, err := s.getLineChanges(ctx, req.Owner, req.RepoName, req.Owner)
	if err != nil {
		return nil, err
	}

	return &git.RepoLineChangesCommitsResponse{
		TotalChanges: int32(totalChanges),
		UserChanges:  int32(userChanges),
		TotalCommits: int32(totalCommits),
		UserCommits:  int32(userCommits)}, nil
}
