package main

import (
	"github.com/acd19ml/TalentRank/apps/user/kafka"
)

func main() {
	brokers := []string{"localhost:9092"}
	producerTopic := "user_tasks"
	consumerTopic := "user_tasks"
	groupID := "user_service_group"

	// 启动 Kafka 服务
	kafka.RunKafkaServices(brokers, producerTopic, consumerTopic, groupID)
}
