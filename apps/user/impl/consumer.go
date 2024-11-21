package impl

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/acd19ml/TalentRank/apps/git"
)

func (s *ServiceImpl) StartConsumer(ctx context.Context, rateLimiter *RateLimiter) {
	log.Println("Kafka Consumer started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer stopped")
			return
		default:
			// 消费一条消息
			message, err := s.Consumer.Consume(ctx, "repo_api_tasks")
			if err != nil {
				log.Printf("Failed to consume message: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// 解析消息内容
			var task map[string]interface{}
			if err := json.Unmarshal(message, &task); err != nil {
				log.Printf("Failed to parse message: %v", err)
				log.Printf("Raw message content: %s", string(message))
				continue
			}

			username, ok := task["username"].(string)
			if !ok {
				log.Printf("Invalid or missing 'username' in message: %+v", task)
				continue
			}
			repoName, ok := task["repo_name"].(string)
			if !ok {
				log.Printf("Invalid or missing 'repo_name' in message: %+v", task)
				continue
			}
			function, ok := task["function"].(string)
			if !ok {
				log.Printf("Invalid or missing 'function' in message: %+v", task)
				continue
			}
			parameters, ok := task["parameters"].(map[string]interface{})
			if !ok {
				log.Printf("Invalid or missing 'parameters' in message: %+v", task)
				continue
			}
			// log.Printf("Message fields - username: %s, repoName: %s, function: %s", username, repoName, function)

			// 使用速率限制器控制消费速率
			if !rateLimiter.Allow() {
				log.Println("Rate limiter blocked this message")
				time.Sleep(1 * time.Second)
				continue
			}

			// 调用对应的远程 gRPC 方法
			var result interface{}
			switch function {
			case "GetStarsByRepo":
				result, err = s.svc.GetStarsByRepo(ctx, &git.RepoRequest{
					Owner:    parameters["owner"].(string),
					RepoName: parameters["repo_name"].(string),
				})
			case "GetForksByRepo":
				result, err = s.svc.GetForksByRepo(ctx, &git.RepoRequest{
					Owner:    parameters["owner"].(string),
					RepoName: parameters["repo_name"].(string),
				})
			case "GetTotalIssuesByRepo":
				result, err = s.svc.GetTotalIssuesByRepo(ctx, &git.RepoRequest{
					Owner:    parameters["owner"].(string),
					RepoName: parameters["repo_name"].(string),
				})
			case "GetUserSolvedIssuesByRepo":
				result, err = s.svc.GetUserSolvedIssuesByRepo(ctx, &git.RepoRequest{
					Owner:    parameters["owner"].(string),
					RepoName: parameters["repo_name"].(string),
				})
			case "GetTotalPullRequestsByRepo":
				result, err = s.svc.GetTotalPullRequestsByRepo(ctx, &git.RepoRequest{
					Owner:    parameters["owner"].(string),
					RepoName: parameters["repo_name"].(string),
				})
			case "GetUserMergedPullRequestsByRepo":
				result, err = s.svc.GetUserMergedPullRequestsByRepo(ctx, &git.RepoRequest{
					Owner:    parameters["owner"].(string),
					RepoName: parameters["repo_name"].(string),
				})
			case "GetTotalCodeReviewsByRepo":
				result, err = s.svc.GetTotalCodeReviewsByRepo(ctx, &git.RepoRequest{
					Owner:    parameters["owner"].(string),
					RepoName: parameters["repo_name"].(string),
				})
			case "GetUserCodeReviewsByRepo":
				result, err = s.svc.GetUserCodeReviewsByRepo(ctx, &git.RepoRequest{
					Owner:    parameters["owner"].(string),
					RepoName: parameters["repo_name"].(string),
				})
			case "GetLineChangesCommitsByRepo":
				result, err = s.svc.GetLineChangesCommitsByRepo(ctx, &git.RepoRequest{
					Owner:    parameters["owner"].(string),
					RepoName: parameters["repo_name"].(string),
				})
			default:
				log.Printf("Unknown function: %s", function)
				continue
			}

			// 错误处理
			if err != nil {
				log.Printf("Failed to call function %s for %s/%s: %v", function, username, repoName, err)
				continue
			}

			// 将结果存储到 MySQL
			err = s.SaveRepoDataToDB(ctx, username, repoName, function, result)
			if err != nil {
				log.Printf("Failed to save data for %s/%s: %v", username, repoName, err)
			} else {
				log.Printf("Ticker: Successfully update data by %s for %s/%s", function, username, repoName)
			}
		}
	}
}
