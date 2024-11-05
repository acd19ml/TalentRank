package impl_test

import (
	"context"
	"log"
	"testing"
	"unicode/utf8"

	"github.com/acd19ml/TalentRank/apps/git"
	"github.com/acd19ml/TalentRank/apps/git/impl"
	"github.com/stretchr/testify/assert"
)

func TestGetReadme(t *testing.T) {
	// 创建一个服务实例
	g := &impl.Service{}
	g.Config()

	// 设置测试请求
	req := &git.GetReadmeRequest{
		Username:  "lunarianss",
		RepoLimit: 10,
		CharLimit: 200,
	}

	// 执行测试函数
	response, err := g.GetReadme(context.Background(), req)

	// 确认函数没有报错
	assert.NoError(t, err, "Expected no error from GetReadme")

	// 确认返回结果不为空
	assert.NotNil(t, response, "Expected a non-nil response from GetReadme")
	assert.NotEmpty(t, response.Result, "Expected non-empty result in response")

	// 确认结果内容为有效的 UTF-8 字符串
	assert.True(t, utf8.ValidString(response.Result), "Expected result to contain valid UTF-8 characters")

	// 打印输出以便调试
	log.Printf("Readme Content: %s", response.Result)
}
