package kafka

import (
	"context"
	"log"
)

func RunKafkaServices(brokers []string, producerTopic, consumerTopic, groupID string) {
	ctx := context.Background()

	// 初始化生产者
	producer := NewProducer(brokers, producerTopic)
	defer producer.Close()

	// 初始化消费者
	consumer := NewConsumer(brokers, groupID, consumerTopic)
	defer consumer.Close()

	// 启动消费者
	go func() {
		log.Printf("Starting consumer for topic: %s", consumerTopic)
		err := consumer.ListenAndConsume(ctx, HandleMessage)
		if err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()

	// 示例：发送测试消息
	producer.SendMessage(ctx, KafkaMessage{
		RequestType: "get_user",
		Username:    "exampleUser",
	})
}
