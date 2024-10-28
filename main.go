package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// func main() {
// 	service := git.NewGitClient()
// 	names, err := service.GetOrganizations(context.Background(), utils.Username)
// 	if err != nil {
// 		log.Fatalf("Error getting name: %v", err)
// 	}
// 	for _, name := range names {
// 		fmt.Printf("Organization for user %s: %s\n", utils.Username, name)
// 	}
// }

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
	commitsURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", repoOwner, repoName)

	// 发送 GET 请求获取提交列表
	response, err := http.Get(commitsURL)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer response.Body.Close()

	// 检查请求是否成功
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d", response.StatusCode)
	}

	// 解析 JSON 数据
	var commits []Commit
	if err := json.NewDecoder(response.Body).Decode(&commits); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %v", err)
	}

	// 统计贡献者的提交和代码行数变化
	contributorsStats := make(map[string]struct {
		Commits   int
		Additions int
		Deletions int
	})

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

	return contributorsStats, nil
}

func main() {
	repoOwner := "pytorch" // 替换为实际的仓库拥有者
	repoName := "pytorch"  // 替换为实际的仓库名称

	// 调用函数获取贡献者统计数据
	contributorsStats, err := GetContributorsStats(repoOwner, repoName)
	if err != nil {
		log.Fatalf("获取贡献者统计数据失败: %v", err)
	}

	// 输出贡献者信息
	for login, stats := range contributorsStats {
		fmt.Printf("用户名: %s, 提交次数: %d, 增加的行数: %d, 删除的行数: %d\n", login, stats.Commits, stats.Additions, stats.Deletions)
	}
}

// 测试进度条
//package main

//import (
//	"fmt"
//	"net/http"
//	"time"
//
//	"github.com/gin-gonic/gin"
//)
//
//var progressChan = make(chan int)
//
//func main() {
//	router := gin.Default()
//	router.LoadHTMLFiles("index.html")
//
//	router.GET("/", func(c *gin.Context) {
//		c.HTML(http.StatusOK, "index.html", nil)
//	})
//
//	router.POST("/start", func(c *gin.Context) {
//		go func() {
//			for i := 0; i <= 10; i++ {
//				time.Sleep(200 * time.Millisecond) // 增加延迟到200毫秒
//				progressChan <- i                  // 将进度发送到 channel
//			}
//		}()
//		c.Status(http.StatusOK)
//	})
//
//	router.GET("/progress", func(c *gin.Context) {
//		c.Writer.Header().Set("Content-Type", "text/event-stream")
//		c.Writer.Header().Set("Cache-Control", "no-cache")
//		c.Writer.Header().Set("Connection", "keep-alive")
//
//		for progress := range progressChan {
//			fmt.Fprintf(c.Writer, "data: %d%%\n\n", progress) // 发送进度
//			if flusher, ok := c.Writer.(http.Flusher); ok {
//				flusher.Flush() // 刷新数据到客户端
//			}
//		}
//	})
//
//	router.Run(":8080")
//}
