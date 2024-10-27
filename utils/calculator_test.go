package utils_test

import (
	"context"
	"os"
	"testing"

	"github.com/acd19ml/TalentRank/utils"
	"github.com/acd19ml/TalentRank/utils/git"
)

var (
	client utils.Service
	ctx    = context.Background()
)

func init() {
	if os.Getenv("GITHUB_TOKEN") == "" {
		panic("GITHUB_TOKEN is not set. Please set it before running tests.")
	}
	client = git.NewGitClient() // 初始化客户端
}

// func TestCalculateContribution(t *testing.T) {

// 	calculator := utils.NewCalculator(client)

// 	ctx := context.Background()
// 	username := utils.Username
// 	repoName := "testrepo"

// 	contribution, err := calculator.CalculateContribution(ctx, username, repoName)
// 	if err != nil {
// 		t.Fatalf("CalculateContribution failed: %v", err)
// 	}

// 	// 根据你在 Service 实现中的测试数据，验证计算结果是否符合预期
// 	expectedContribution := 0.5 // 假设的期望贡献度
// 	if contribution != expectedContribution {
// 		t.Errorf("expected contribution %v, got %v", expectedContribution, contribution)
// 	}
// }

func TestCalculateOverallScore(t *testing.T) {
	calculator := utils.NewCalculator(client)

	ctx := context.Background()
	username := utils.Username

	overallScore, err := calculator.CalculateOverallScore(ctx, username)
	if err != nil {
		t.Fatalf("CalculateOverallScore failed: %v", err)
	}

	t.Log("Overall score:", overallScore)
}
