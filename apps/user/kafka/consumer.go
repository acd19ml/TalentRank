package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}
type MultiTopicConsumer struct {
	RepoConsumer *KafkaConsumer
	UserConsumer *KafkaConsumer
}

// ConsumeMessage fetches messages from the appropriate topic
func (mc *MultiTopicConsumer) Consume(ctx context.Context, topic string) ([]byte, error) {
	switch topic {
	case "repo_api_tasks":
		return mc.RepoConsumer.Consume(ctx, topic)
	case "user_api_tasks":
		return mc.UserConsumer.Consume(ctx, topic)
	default:
		return nil, fmt.Errorf("unknown topic: %s", topic)
	}
}

// NewKafkaConsumer 初始化 Kafka 消费者
func NewKafkaConsumer(brokers []string, topic string, groupID string) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
	}
}

// Consume 实现 MessageConsumer 接口
func (kc *KafkaConsumer) Consume(ctx context.Context, topic string) ([]byte, error) {
	msg, err := kc.reader.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read message from Kafka: %w", err)
	}

	// log.Printf("Consumed message from topic %s: %s", topic, string(msg.Value))
	return msg.Value, nil
}
