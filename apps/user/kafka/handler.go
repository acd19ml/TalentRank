package kafka

import (
	"log"
)

// HandleMessage 处理 Kafka 消息的具体业务逻辑
func HandleMessage(message KafkaMessage) {
	switch message.RequestType {
	case "get_user":
		log.Printf("Handling get_user task for username: %s", message.Username)
		// user := impl.FetchUserData(message.Username)
		// log.Printf("Fetched user data: %+v", user)

	case "get_repo":
		log.Printf("Handling get_repo task for username: %s, repo: %s", message.Username, message.Repo)
		// repo := impl.FetchRepoData(message.Username, message.Repo)
		// log.Printf("Fetched repo data: %+v", repo)

	default:
		log.Printf("Unknown request type: %s", message.RequestType)
	}
}
