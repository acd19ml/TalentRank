package impl

import (
	"context"
	"log"
	"time"

	"github.com/acd19ml/TalentRank/apps/git"
	"golang.org/x/time/rate"
)

var producerRateLimiter = rate.NewLimiter(rate.Every(5*time.Second), 1) // 每 5 秒生成一条消息

func (s *ServiceImpl) StartProducer(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour) // 每日触发
	log.Println("The data update cleanup task will start every 24 hours")
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Producer stopped")
			return
		case <-ticker.C:
			log.Println("Starting user task production")

			// 从数据库中获取所有用户名
			usernames, err := s.GetAllUsernamesFromDB(ctx)
			if err != nil {
				log.Printf("Failed to get usernames: %v", err)
				continue
			}

			for _, username := range usernames {
				// 调用 GetRepositories 获取用户的所有仓库
				repoListResp, err := s.svc.GetRepositories(ctx, &git.GetUsernameRequest{Username: username})
				if err != nil {
					log.Printf("Failed to get repositories for username %s: %v", username, err)
					continue
				}

				// 为每个仓库生成任务消息
				for _, repo := range repoListResp.Result {
					functions := []string{
						"GetStarsByRepo",
						"GetForksByRepo",
						"GetTotalIssuesByRepo",
						"GetUserSolvedIssuesByRepo",
						"GetTotalPullRequestsByRepo",
						"GetUserMergedPullRequestsByRepo",
						"GetTotalCodeReviewsByRepo",
						"GetUserCodeReviewsByRepo",
						"GetLineChangesCommitsByRepo",
					}

					for _, function := range functions {
						// 控制生产速率
						if !producerRateLimiter.Allow() {
							time.Sleep(1 * time.Second) // 等待限速
						}

						message := map[string]interface{}{
							"username":  username,
							"repo_name": repo,
							"function":  function,
							"parameters": map[string]interface{}{
								"owner":     username,
								"repo_name": repo,
							},
						}

						// 发送消息到队列
						err := s.Producer.Produce(ctx, "repo_api_tasks", message)
						if err != nil {
							log.Printf("Failed to produce message for function %s: %v", function, err)
						}
					}
				}
			}
		}
	}
}
