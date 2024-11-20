package http

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/acd19ml/TalentRank/apps/user"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateUserRepos(c *gin.Context) {

	// 从 Cookie 中获取 Token
	token, err := c.Cookie("githubToken")
	if err != nil {
		token = os.Getenv("GITHUB_TOKEN") // 如果没有 Cookie，使用默认的环境变量
	}

	// 将 client 保存到上下文中
	c.Set("githubToken", token)

	ins := user.NewCreateUserReposRequest()
	log.Println("Received POST request to create UserRepos")
	// 用户传递过来的参数进行解析, 实现了一个json 的unmarshal
	if err := c.Bind(ins); err != nil {
		// 参数绑定失败
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 调用服务
	userRepos, err := h.svc.CreateUserRepos(c, ins.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 返回结果
	c.JSON(200, userRepos)
}

func (h *Handler) QueryUsers(c *gin.Context) {
	// 从请求中获取参数
	req := user.NewQueryUserFromHTTP(c.Request)

	set, err := h.svc.QueryUsers(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 返回结果
	c.JSON(200, set)
}

func (h *Handler) DescribeUserRepos(c *gin.Context) {
	// 从请求中获取参数
	req := user.NewDescribeUserReposRequestFromHTTP(c.Request)

	userRepos, err := h.svc.DescribeUserRepos(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 将 JSON 字符串解析为 Go 的 map 类型
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal([]byte(userRepos), &jsonResponse); err != nil {
		c.JSON(500, gin.H{"error": "failed to parse JSON response"})
		return
	}

	// 返回格式化后的 JSON
	c.JSON(200, jsonResponse)
}

func (h *Handler) GetLocationCounts(c *gin.Context) {
	// 调用服务
	req, err := h.svc.GetLocationCounts(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 返回结果
	c.JSON(200, req)
}

func (h *Handler) DeleteUserRepos(c *gin.Context) {
	req := &user.DeleteUserReposRequest{
		Id: c.Param("id"),
	}

	// 调用服务
	result, err := h.svc.DeleteUserRepos(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 返回结果
	c.JSON(200, result)
}

func (h *Handler) setTokenHandler(c *gin.Context) {
	var body struct {
		Token string `json:"token"`
	}
	// 绑定请求的 JSON 数据到结构体
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	// 检查 Token 是否为空
	if body.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token cannot be empty"})
		return
	}

	// 验证 Token 是否有效
	if !user.VerifyToken(body.Token) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid GitHub Token"})
		return
	}

	// 调用 Service 的 Config 方法，配置 OAuth2 客户端
	//Todo：更新环境中的 token

	// 将 Token 写入 Cookie
	c.SetCookie("githubToken", body.Token, int(7*24*time.Hour.Seconds()), "/", "", false, true)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Token set successfully"})
}
