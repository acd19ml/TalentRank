package kafka

import "time"

// KafkaMessage 定义 Kafka 消息的结构
type KafkaMessage struct {
	RequestType string            `json:"request_type"`           // 任务类型，如 "get_user"
	Username    string            `json:"username"`               // GitHub 用户名
	Repo        string            `json:"repo,omitempty"`         // 仓库名，仓库任务需要
	ExtraParams map[string]string `json:"extra_params,omitempty"` // 额外参数
	Timestamp   time.Time         `json:"timestamp"`              // 消息生成时间
}
