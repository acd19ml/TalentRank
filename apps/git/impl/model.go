package impl

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/acd19ml/TalentRank/apps/git"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// 接口检查
var _ git.Service = (*Git)(nil)

func NewGitClient() *Git {
	ctx := context.Background()
	// 使用 OAuth2 令牌进行认证
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &Git{
		client: client,
		oauth:  &ts,
	}
}

type Git struct {
	client           *github.Client
	oauth            *oauth2.TokenSource
	reposCache       []string   // 缓存仓库列表
	cacheUsername    string     // 缓存中保存的用户名
	cacheMutex       sync.Mutex // 用于同步缓存访问
	cacheInitialized bool       // 标记缓存是否已初始化
}

// 初始化缓存，当 reposCache 为空或 cacheUsername 不匹配时执行 GetRepositories
func (g *Git) initCache(ctx context.Context, username string) error {
	g.cacheMutex.Lock()
	defer g.cacheMutex.Unlock()

	// 检查缓存是否已经存在且与当前用户名一致
	if g.cacheUsername == username && len(g.reposCache) > 0 {
		return nil
	}

	// 调用 GetRepositories 获取数据并缓存
	reposList, err := g.fetchRepositories(ctx, username)
	if err != nil {
		return err
	}

	// 更新缓存
	g.reposCache = reposList
	g.cacheUsername = username
	return nil
}

func (g *Git) GetUser(ctx context.Context, username string) (*github.User, error) {
	user, _, err := g.client.Users.Get(ctx, username)
	if err != nil {
		panic(err)
	}
	return user, nil
}

// fetchRepositories 真正执行获取指定用户的所有仓库名称
func (g *Git) fetchRepositories(ctx context.Context, username string) ([]string, error) {
	var reposList []string
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	for {
		repos, resp, err := g.client.Repositories.List(ctx, username, opts)
		if err != nil {
			return nil, err
		}

		for _, repo := range repos {
			isContributor, err := g.checkIfUserIsContributor(ctx, username, repo.GetOwner().GetLogin(), repo.GetName())
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

// GetRepositories 获取指定用户的所有仓库名称，依赖缓存
func (g *Git) GetRepositories(ctx context.Context, username string) ([]string, error) {
	// 初始化缓存
	if err := g.initCache(ctx, username); err != nil {
		return nil, err
	}
	return g.reposCache, nil
}

// checkIfUserIsContributor 检查指定用户是否为某个仓库的贡献者
func (g *Git) checkIfUserIsContributor(ctx context.Context, username, owner, repo string) (bool, error) {
	contributors, _, err := g.client.Repositories.ListContributors(ctx, owner, repo, nil)
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

// Commit 结构体用来存储提交信息
type Commit struct {
	Sha    string `json:"sha"`
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
}

// CommitDetail 存储每个提交的代码行变化信息
type CommitDetail struct {
	Stats struct {
		Additions int `json:"additions"`
		Deletions int `json:"deletions"`
	} `json:"stats"`
}
