package git_test

import (
	"github.com/acd19ml/TalentRank/utils/git"
	"log"
	"testing"
)

// TestGetContributorsStats 测试 GetContributorsStats 函数 ZJK
func TestGetContributorsStats(t *testing.T) {
	repoOwner := "THUDM"      // 替换为实际的仓库拥有者
	repoName := "GLM-4-Voice" // 替换为实际的仓库名称

	// 调用函数获取贡献者统计数据
	contributorsStats, err := git.GetContributorsStats(repoOwner, repoName)
	if err != nil {
		t.Fatalf("获取贡献者统计数据失败: %v", err)
	}

	// 检查是否返回了有效的数据
	if len(contributorsStats) == 0 {
		t.Error("未返回任何贡献者统计数据")
	}

	// 输出贡献者信息（可选）
	for login, stats := range contributorsStats {
		log.Printf("用户名: %s, 提交次数: %d, 增加的行数: %d, 删除的行数: %d\n", login, stats.Commits, stats.Additions, stats.Deletions)
	}
}
