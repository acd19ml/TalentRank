package impl_test

// import (
// 	"context"
// 	"testing"

// 	"github.com/acd19ml/TalentRank/apps/user"
// 	"github.com/acd19ml/TalentRank/apps/user/impl"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// // MockHandler 用于 mock CreateUserByUsername 方法
// type MockHandler struct {
// 	mock.Mock
// }

// // Mock CreateUserByUsername 方法
// func (m *MockHandler) CreateUserByUsername(ctx context.Context, username string) (*user.User, error) {
// 	args := m.Called(ctx, username)
// 	return args.Get(0).(*user.User), args.Error(1)
// }

// func TestCreateUser(t *testing.T) {
// 	// 创建 MockHandler 实例
// 	mockHandler := new(MockHandler)

// 	// 创建测试上下文
// 	ctx := context.Background()
// 	username := "testuser"

// 	// 创建一个模拟的返回结果
// 	mockUser := &user.User{
// 		Username:      username,
// 		Name:          "Test Name",
// 		Company:       "Test Company",
// 		Location:      "Test Location",
// 		Email:         "test@example.com",
// 		Bio:           "Test Bio",
// 		TotalStar:     100,
// 		TotalFork:     50,
// 		Followers:     10,
// 		Dependents:    5,
// 		Organizations: []string{"Org1", "Org2"},
// 		Readme:        "Sample Readme",
// 		Commits:       "Sample Commits",
// 	}

// 	// 设置 MockHandler 的预期返回值
// 	mockHandler.On("CreateUserByUsername", ctx, username).Return(mockUser, nil)

// 	// 创建 UserServiceImpl 实例，并将 mockHandler 作为 Handler 传入
// 	service := &impl.UserServiceImpl{
// 		Handler: mockHandler,
// 	}

// 	// 调用 CreateUser 方法
// 	result, err := service.CreateUser(ctx, username)

// 	// 验证结果
// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.Equal(t, username, result.Username)
// 	assert.Equal(t, "Test Name", result.Name)
// 	assert.Equal(t, "Test Company", result.Company)
// 	assert.Equal(t, "Test Location", result.Location)
// 	assert.Equal(t, "test@example.com", result.Email)
// 	assert.Equal(t, "Test Bio", result.Bio)
// 	assert.Equal(t, 100, result.TotalStar)
// 	assert.Equal(t, 50, result.TotalFork)
// 	assert.Equal(t, 10, result.Followers)
// 	assert.Equal(t, 5, result.Dependents)
// 	assert.Equal(t, []string{"Org1", "Org2"}, result.Organizations)
// 	assert.Equal(t, "Sample Readme", result.Readme)
// 	assert.Equal(t, "Sample Commits", result.Commits)

// 	// 验证 Mock 方法的调用次数
// 	mockHandler.AssertExpectations(t)
// }
