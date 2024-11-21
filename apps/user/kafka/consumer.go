package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

// Consumer 定义 Kafka 消费者
type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer 创建 Kafka 消费者实例
func NewConsumer(brokers []string, groupID, topic string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
	}
}

// ListenAndConsume 监听并消费 Kafka 消息
func (c *Consumer) ListenAndConsume(ctx context.Context, handleMessage func(KafkaMessage)) error {
	for {
		// 读取 Kafka 消息
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return err
		}

		// 解析消息
		var message KafkaMessage
		err = json.Unmarshal(msg.Value, &message)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		// 调用消息处理函数
		handleMessage(message)
	}
}

// Close 关闭消费者
func (c *Consumer) Close() error {
	return c.reader.Close()
}
