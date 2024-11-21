package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer 初始化 Kafka 生产者
func NewKafkaProducer(brokers []string) *KafkaProducer {
	return &KafkaProducer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
		}),
	}
}

// Produce 实现 MessageProducer 接口
func (kp *KafkaProducer) Produce(ctx context.Context, topic string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = kp.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: data,
	})
	if err != nil {
		return fmt.Errorf("failed to write message to Kafka: %w", err)
	}

	log.Printf("Produced message to topic %s: %s", topic, string(data))
	return nil
}
