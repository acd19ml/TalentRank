package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer 定义 Kafka 生产者
type Producer struct {
	writer *kafka.Writer
}

// NewProducer 创建 Kafka 生产者实例
func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
			Topic:   topic,
		}),
	}
}

// SendMessage 发送消息到 Kafka
func (p *Producer) SendMessage(ctx context.Context, message KafkaMessage) error {
	// 设置消息的时间戳
	message.Timestamp = time.Now().UTC()

	// 将消息序列化为 JSON
	value, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// 发送消息
	err = p.writer.WriteMessages(ctx, kafka.Message{
		Value: value,
	})
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}

	log.Printf("Message sent: %s", value)
	return nil
}

// Close 关闭生产者
func (p *Producer) Close() error {
	return p.writer.Close()
}
