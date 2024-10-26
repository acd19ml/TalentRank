package git

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

func (g *Git) GetName(ctx context.Context, username string) (string, error) {
	user, err := g.GetUser(ctx, username)
	if err != nil {
		panic(err)
	}
	return user.GetName(), nil
}

func (g *Git) GetCompany(ctx context.Context, username string) (string, error) {
	user, err := g.GetUser(ctx, username)
	if err != nil {
		panic(err)
	}
	return user.GetCompany(), nil
}

func (g *Git) GetLocation(ctx context.Context, username string) (string, error) {
	user, err := g.GetUser(ctx, username)
	if err != nil {
		panic(err)
	}
	return user.GetLocation(), nil
}

func (g *Git) GetEmail(ctx context.Context, username string) (string, error) {
	user, err := g.GetUser(ctx, username)
	if err != nil {
		panic(err)
	}
	return user.GetEmail(), nil
}

func (g *Git) GetBio(ctx context.Context, username string) (string, error) {
	user, err := g.GetUser(ctx, username)
	if err != nil {
		panic(err)
	}
	return user.GetBio(), nil
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

func (g *Git) GetReadme(ctx context.Context, username string, charLimit int) (string, error) {
	repos, err := g.GetRepositories(ctx, username)
	if err != nil {
		return "", err
	}

	var contents string
	for i, repo := range repos {
		// 只处理前20个仓库
		if i >= 20 {
			break
		}

		readme, _, err := g.client.Repositories.GetReadme(ctx, username, repo, nil)
		if err != nil {
			if githubErr, ok := err.(*github.ErrorResponse); ok && githubErr.Response.StatusCode == 404 {
				continue
			}
			return "", err
		}

		content, err := readme.GetContent()
		if err != nil {
			return "", err
		}

		if len(content) > charLimit {
			content = content[:charLimit] + "..."
		}

		// 确保总内容不会超过最大字符限制
		if len(contents)+len(content) > charLimit*20 {
			contents += content[:charLimit*20-len(contents)] + "..."
			break
		}

		contents += content
	}

	return contents, nil
}

func (g *Git) GetCommits(ctx context.Context, username string, charLimit int) (string, error) {
	repos, err := g.GetRepositories(ctx, username)
	if err != nil {
		return "", err
	}

	var allCommits string

	for i, repo := range repos {
		if i >= 20 {
			break
		}

		opts := &github.CommitsListOptions{
			Author:      username,
			ListOptions: github.ListOptions{PerPage: 10},
		}

		var commits []*github.RepositoryCommit
		var lastErr error

		// 尝试三次请求
		for attempt := 1; attempt <= 3; attempt++ {
			commits, _, err = g.client.Repositories.ListCommits(ctx, username, repo, opts)
			if err == nil {
				lastErr = nil
				break
			}
			lastErr = err
			time.Sleep(2 * time.Second)
		}

		// 如果三次尝试均失败，则跳过该仓库并继续
		if lastErr != nil {
			fmt.Printf("Skipping repo '%s' due to persistent errors: %v\n", repo, lastErr)
			continue
		}

		var repoCommits string
		for _, commit := range commits {
			message := commit.GetCommit().GetMessage()
			if len(repoCommits)+len(message) > charLimit {
				break
			}
			repoCommits += message + "\n"
		}

		if len(repoCommits) > charLimit {
			repoCommits = repoCommits[:charLimit] + "...\n"
		}

		allCommits += fmt.Sprintf("Repo: %s\n%s\n", repo, repoCommits)
		if len(allCommits) > charLimit*20 {
			allCommits = allCommits[:charLimit*20] + "..."
			break
		}
	}

	return allCommits, nil
}
