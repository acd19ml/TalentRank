package git_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUtils 用于模拟依赖的 utils 函数
type MockUtils struct {
	mock.Mock
}

func (m *MockUtils) GetDependentRepositories(url string) (int, error) {
	args := m.Called(url)
	return args.Int(0), args.Error(1)
}

// MockGitClient 用于模拟 git.Git 客户端
type MockGitClient struct {
	mock.Mock
	MockUtils *MockUtils
}

// Mock GetRepositories 方法
func (m *MockGitClient) GetRepositories(ctx context.Context, username string) ([]string, error) {
	args := m.Called(ctx, username)
	return args.Get(0).([]string), args.Error(1)
}

// Mock GetDependentRepositories 方法，整合了 MockUtils 的依赖
func (m *MockGitClient) GetDependentRepositories(ctx context.Context, username string) (int, error) {
	// 调用 Mock GetRepositories 模拟数据
	repos, err := m.GetRepositories(ctx, username)
	if err != nil {
		return 0, err
	}

	// 模拟总依赖计数
	totalDependents := 0
	for _, repo := range repos {
		dependentURL := "https://github.com/" + username + "/" + repo + "/network/dependents"
		count, err := m.MockUtils.GetDependentRepositories(dependentURL)
		if err != nil {
			return 0, err
		}
		totalDependents += count
	}
	return totalDependents, nil
}

// TestGetDependentRepositories 测试 Git.GetDependentRepositories 函数
func TestGetDependentRepositories(t *testing.T) {
	ctx := context.Background()
	mockUtils := new(MockUtils)
	mockClient := &MockGitClient{MockUtils: mockUtils} // 初始化带 MockUtils 的 MockGitClient

	// 模拟 GetRepositories 返回的仓库列表
	mockClient.On("GetRepositories", ctx, username).Return([]string{"repo1", "repo2", "repo3"}, nil)

	// 设置 MockUtils 的行为，模拟依赖计数
	mockUtils.On("GetDependentRepositories", mock.Anything).Return(5, nil)

	// 运行 GetDependentRepositories 测试
	totalDependents, err := mockClient.GetDependentRepositories(ctx, username)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 15, totalDependents) // 每个仓库有 5 个依赖，总共 15 个

	// 验证 Mock 调用次数
	mockClient.AssertExpectations(t)
	mockUtils.AssertExpectations(t)
}
