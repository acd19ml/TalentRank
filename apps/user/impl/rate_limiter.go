package impl

import "time"

type RateLimiter struct {
	tokens    chan struct{}
	interval  time.Duration
	maxTokens int
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(maxTokens int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens:    make(chan struct{}, maxTokens),
		interval:  interval,
		maxTokens: maxTokens,
	}

	// 定时生成令牌
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			select {
			case rl.tokens <- struct{}{}:
			default:
			}
		}
	}()

	return rl
}

// Allow 消耗一个令牌
func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}
