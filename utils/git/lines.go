package git

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Commit 结构体用来存储提交信息
type Commit struct {
	Sha    string `json:"sha"`
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
}

type CommitDetail struct {
	Stats struct {
		Additions int `json:"additions"`
		Deletions int `json:"deletions"`
	} `json:"stats"`
}

// GetContributorsStats 获取贡献者的提交次数和代码行数变换
func GetContributorsStats(repoOwner, repoName string) (map[string]struct {
	Commits   int
	Additions int
	Deletions int
}, error) {
	contributorsStats := make(map[string]struct {
		Commits   int
		Additions int
		Deletions int
	})

	commitsURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", repoOwner, repoName)
	opts := &http.Client{}

	// 分页获取提交列表
	for page := 1; ; page++ {
		req, err := http.NewRequest("GET", commitsURL, nil)
		if err != nil {
			return nil, fmt.Errorf("创建请求失败: %v", err)
		}
		q := req.URL.Query()
		q.Add("page", fmt.Sprintf("%d", page))
		q.Add("per_page", "100")
		req.URL.RawQuery = q.Encode()

		response, err := opts.Do(req)
		if err != nil {
			return nil, fmt.Errorf("请求失败: %v", err)
		}
		defer response.Body.Close()

		// 检查请求是否成功
		if response.StatusCode != http.StatusOK {
			break
		}

		// 解析 JSON 数据
		var commits []Commit
		if err := json.NewDecoder(response.Body).Decode(&commits); err != nil {
			return nil, fmt.Errorf("解析 JSON 失败: %v", err)
		}

		if len(commits) == 0 {
			break // 如果没有更多提交，退出循环
		}

		// 遍历每个提交
		for _, commit := range commits {
			login := commit.Author.Login

			// 获取提交的详细信息
			detailURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", repoOwner, repoName, commit.Sha)
			detailResponse, err := http.Get(detailURL)
			if err != nil {
				log.Printf("获取详细提交信息失败: %v", err)
				continue
			}
			defer detailResponse.Body.Close()

			// 检查详细请求是否成功
			if detailResponse.StatusCode != http.StatusOK {
				log.Printf("请求详细提交信息失败，状态码: %d", detailResponse.StatusCode)
				continue
			}

			// 解析提交详细信息
			var commitDetail CommitDetail
			if err := json.NewDecoder(detailResponse.Body).Decode(&commitDetail); err != nil {
				log.Printf("解析提交详细信息失败: %v", err)
				continue
			}

			// 更新贡献者统计数据
			stats := contributorsStats[login]
			stats.Commits++
			stats.Additions += commitDetail.Stats.Additions
			stats.Deletions += commitDetail.Stats.Deletions
			contributorsStats[login] = stats
		}
	}

	return contributorsStats, nil
}
