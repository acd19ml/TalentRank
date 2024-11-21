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
	ticker := time.NewTicker(24 * time.Second) // 每日触发
	log.Println("The data update cleanup task will start every 24 hours")
	defer ticker.Stop()

	s.DeleteOrphanedRepos(ctx)  // 清理孤立的仓库
	s.RemoveDuplicateUsers(ctx) // 清理重复的用户

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
				// 更新用户相关信息
				userFunctions := []string{
					"GetName",
					"GetCompany",
					"GetLocation",
					"GetEmail",
					"GetBio",
					"GetOrganizations",
					"GetFollowers",
				}

				for _, function := range userFunctions {
					// 控制生产速率
					if !producerRateLimiter.Allow() {
						time.Sleep(1 * time.Second) // 等待限速
					}

					message := map[string]interface{}{
						"username": username,
						"function": function,
						"parameters": map[string]interface{}{
							"username": username,
						},
					}

					// 发送用户信息更新消息到队列
					err := s.Producer.Produce(ctx, "user_api_tasks", message)
					if err != nil {
						log.Printf("Failed to produce message for function %s: %v", function, err)
					}
				}

				// 调用 GetRepositories 获取用户的所有仓库
				repoListResp, err := s.svc.GetRepositories(ctx, &git.GetUsernameRequest{Username: username})
				if err != nil {
					log.Printf("Failed to get repositories for username %s: %v", username, err)
					continue
				}

				// 比较新旧仓库列表，删除无效仓库
				oldeRepos, err := s.FetchReposFromDB(ctx, username)
				if err != nil {
					log.Printf("Failed to get old repositories for username %s: %v", username, err)
					continue
				}
				invalidRepos := HasInvalidRepos(oldeRepos, repoListResp.Result)
				err = s.DeleteInvalidReposFromDB(ctx, invalidRepos)
				if err != nil {
					log.Printf("Failed to delete invalid repositories for username %s: %v", username, err)
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
				// 更新计算用户技术评分
				userrepos, error := s.FetchUserReposFromDB(ctx, username)
				if error != nil {
					log.Printf("Failed to get user repos for username %s: %v", username, error)
					continue
				}
				if err = CalculateOverallScore(userrepos); err != nil {
					log.Printf("Failed to calculate overall score for username %s: %v", username, err)
					continue
				}
				s.UpdateUserScore(ctx, username, userrepos.User.Score)

			}
		}
	}
}
