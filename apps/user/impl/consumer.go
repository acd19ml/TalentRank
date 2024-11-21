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
			// 消费仓库任务消息
			repoMessage, err := s.Consumer.Consume(ctx, "repo_api_tasks")
			if err == nil {
				processRepoMessage(ctx, s, rateLimiter, repoMessage)
			}

			// 消费用户任务消息
			userMessage, err := s.Consumer.Consume(ctx, "user_api_tasks")
			if err == nil {
				processUserMessage(ctx, s, rateLimiter, userMessage)
			}

			// 短暂休眠以防止过载
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func processRepoMessage(ctx context.Context, s *ServiceImpl, rateLimiter *RateLimiter, message []byte) {
	// 解析消息内容
	var task map[string]interface{}
	if err := json.Unmarshal(message, &task); err != nil {
		log.Printf("Failed to parse repo message: %v", err)
		log.Printf("Raw repo message content: %s", string(message))
		return
	}

	username, ok := task["username"].(string)
	if !ok {
		log.Printf("Invalid or missing 'username' in repo message: %+v", task)
		return
	}
	repoName, ok := task["repo_name"].(string)
	if !ok {
		log.Printf("Invalid or missing 'repo_name' in repo message: %+v", task)
		return
	}
	function, ok := task["function"].(string)
	if !ok {
		log.Printf("Invalid or missing 'function' in repo message: %+v", task)
		return
	}
	parameters, ok := task["parameters"].(map[string]interface{})
	if !ok {
		log.Printf("Invalid or missing 'parameters' in repo message: %+v", task)
		return
	}

	// 控制消费速率
	if !rateLimiter.Allow() {
		log.Println("Rate limiter blocked repo message")
		time.Sleep(1 * time.Second)
		return
	}

	// 调用 gRPC 方法
	var result interface{}
	var err error
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
	}

	// 处理结果
	if err != nil {
		log.Printf("Failed to call repo function %s for %s/%s: %v", function, username, repoName, err)
		return
	}

	// 存储到数据库
	err = s.SaveRepoDataToDB(ctx, username, repoName, function, result)
	if err != nil {
		log.Printf("Failed to save repo data for %s/%s: %v", username, repoName, err)
	} else {
		log.Printf("Successfully updated repo data for %s/%s by function %s", username, repoName, function)
	}
}

func processUserMessage(ctx context.Context, s *ServiceImpl, rateLimiter *RateLimiter, message []byte) {
	// 解析消息内容
	var task map[string]interface{}
	if err := json.Unmarshal(message, &task); err != nil {
		log.Printf("Failed to parse user message: %v", err)
		log.Printf("Raw user message content: %s", string(message))
		return
	}

	username, ok := task["username"].(string)
	if !ok {
		log.Printf("Invalid or missing 'username' in user message: %+v", task)
		return
	}
	function, ok := task["function"].(string)
	if !ok {
		log.Printf("Invalid or missing 'function' in user message: %+v", task)
		return
	}
	parameters, ok := task["parameters"].(map[string]interface{})
	if !ok {
		log.Printf("Invalid or missing 'parameters' in user message: %+v", task)
		return
	}

	// 控制消费速率
	if !rateLimiter.Allow() {
		log.Println("Rate limiter blocked user message")
		time.Sleep(1 * time.Second)
		return
	}

	// 调用 gRPC 方法
	var result interface{}
	var err error
	switch function {
	case "GetName":
		result, err = s.svc.GetName(ctx, &git.GetUsernameRequest{
			Username: parameters["username"].(string),
		})
	case "GetCompany":
		result, err = s.svc.GetCompany(ctx, &git.GetUsernameRequest{
			Username: parameters["username"].(string),
		})
	case "GetLocation":
		result, err = s.svc.GetLocation(ctx, &git.GetUsernameRequest{
			Username: parameters["username"].(string),
		})
	case "GetEmail":
		result, err = s.svc.GetEmail(ctx, &git.GetUsernameRequest{
			Username: parameters["username"].(string),
		})
	case "GetBio":
		result, err = s.svc.GetBio(ctx, &git.GetUsernameRequest{
			Username: parameters["username"].(string),
		})
	case "GetOrganizations":
		result, err = s.svc.GetOrganizations(ctx, &git.GetUsernameRequest{
			Username: parameters["username"].(string),
		})
	case "GetFollowers":
		result, err = s.svc.GetFollowers(ctx, &git.GetUsernameRequest{
			Username: parameters["username"].(string),
		})

	default:
		log.Printf("Unknown function: %s", function)
		return
	}

	// 处理结果
	if err != nil {
		log.Printf("Failed to call user function %s for %s: %v", function, username, err)
		return
	}

	// 存储到数据库
	err = s.SaveUserDataToDB(ctx, username, function, result)
	if err != nil {
		log.Printf("Failed to save user data for %s: %v", username, err)
	} else {
		log.Printf("Successfully updated user data for %s by function %s", username, function)
	}
}
