package git

import (
	"context"

	"github.com/google/go-github/github"
)

func (g *Git) GetFollowers(ctx context.Context, username string) (int, error) {
	user, err := g.GetUser(ctx, username)
	num_follwers := user.GetFollowers()
	// 打印 followers 总数
	return num_follwers, err

}

// GetTotalStars 获取指定用户的所有仓库 star 总数
func (g *Git) GetTotalStars(ctx context.Context, username string) (int, error) {
	// 初始化星标总数
	totalStars := 0

	// 设置分页参数
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	// 获取所有仓库并累计星标数
	for {
		repos, resp, err := g.client.Repositories.List(ctx, username, opts)
		if err != nil {
			return 0, err
		}

		for _, repo := range repos {
			totalStars += repo.GetStargazersCount()
		}

		// 如果没有下一页，则退出循环
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return totalStars, nil
}

// GetTotalForks 获取指定用户的所有仓库 fork 总数
func (g *Git) GetTotalForks(ctx context.Context, username string) (int, error) {
	// 初始化 fork 总数
	totalForks := 0

	// 设置分页参数
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	// 获取所有仓库并累计 fork 数
	for {
		repos, resp, err := g.client.Repositories.List(ctx, username, opts)
		if err != nil {
			return 0, err
		}

		for _, repo := range repos {
			totalForks += repo.GetForksCount()
		}

		// 如果没有下一页，则退出循环
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return totalForks, nil
}

// GetRepositories 获取指定用户的所有仓库名称
func (g *Git) GetRepositories(ctx context.Context, username string) ([]string, error) {
	var reposList []string

	// 设置分页参数
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	// 获取所有仓库名称
	for {
		repos, resp, err := g.client.Repositories.List(ctx, username, opts)
		if err != nil {
			return nil, err
		}

		for _, repo := range repos {
			reposList = append(reposList, repo.GetName())
		}

		// 如果没有下一页，则退出循环
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return reposList, nil
}

func (g *Git) GetOrganizations(ctx context.Context, username string) ([]string, error) {
	var orgsList []string

	// 设置分页参数
	opts := &github.ListOptions{PerPage: 50}

	// 获取所有组织
	for {
		orgs, resp, err := g.client.Organizations.List(ctx, username, opts)
		if err != nil {
			return nil, err
		}

		for _, org := range orgs {
			orgsList = append(orgsList, org.GetLogin())
		}

		// 如果没有下一页，则退出循环
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return orgsList, nil
}
