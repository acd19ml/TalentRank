package utils_test

import (
	"context"
	"os"
	"testing"

	"github.com/acd19ml/TalentRank/utils"
	"github.com/acd19ml/TalentRank/utils/git"
)

var (
	client   utils.Service
	ctx      = context.Background()
	username = utils.Username
)

func init() {
	if os.Getenv("GITHUB_TOKEN") == "" {
		panic("GITHUB_TOKEN is not set. Please set it before running tests.")
	}
	client = git.NewGitClient() // 初始化客户端
}

func TestCalculateContribution(t *testing.T) {

	calculator := utils.NewCalculator(client)

	repoName := "alfred-lunar-cal"

	contribution, err := calculator.CalculateContribution(ctx, username, repoName)
	if err != nil {
		t.Fatalf("CalculateContribution failed: %v", err)
	}
	t.Log("Contribution:", contribution)
}

func TestCalculateOverallScore(t *testing.T) {
	calculator := utils.NewCalculator(client)

	overallScore, err := calculator.CalculateOverallScore(ctx, username)
	if err != nil {
		t.Fatalf("CalculateOverallScore failed: %v", err)
	}

	t.Log("Overall score:", overallScore)
}
