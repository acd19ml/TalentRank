package http

import (
	"encoding/json"
	"fmt"
	"github.com/acd19ml/TalentRank/apps/user"
	"github.com/acd19ml/TalentRank/myredis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"log"
	"net/http"
	"os"
	"time"
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

	// 清空 Redis 数据库
	err = myredis.Rdb.FlushDB(c).Err() // 清空当前数据库
	// 如果需要清空所有数据库，则使用 FlushAll：
	// err = myredis.Rdb.FlushAll(c).Err()

	if err != nil {
		log.Printf("Failed to clear Redis database: %v", err)
		c.JSON(500, gin.H{"error": "Failed to clear Redis database"})
		return
	}

	log.Println("Redis database cleared successfully")

	// 返回结果
	c.JSON(200, userRepos)
}

func (h *Handler) QueryUsers(c *gin.Context) {
	// 从请求中获取参数
	req := user.NewQueryUserFromHTTP(c.Request)
	cacheKey := fmt.Sprintf("users:possible_nation:%s:offset:%d:pagesize:%d", req.PossibleNation, req.OffSet(), req.GetPageSize())

	// 尝试从 Redis 获取缓存数据
	cachedData, err := myredis.Rdb.Get(c, cacheKey).Result()
	if err == nil {
		// 缓存命中，将 JSON 字符串反序列化为结构体
		var userSet *user.UserSet
		if jsonErr := json.Unmarshal([]byte(cachedData), &userSet); jsonErr == nil {
			c.JSON(200, userSet) // 直接返回缓存数据
			log.Printf("缓存命中")
			return
		} else {
			// 如果反序列化失败，记录日志但继续处理
			log.Printf("failed to unmarshal cached data: %v", jsonErr)
		}
	} else if err != redis.Nil {
		// 如果是非缓存未命中的错误，直接返回 500
		c.JSON(500, gin.H{"error": "Failed to get data from cache"})
		return
	}

	// 调用服务层查询数据
	set, err := h.svc.QueryUsers(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 缓存查询结果
	cacheData, jsonErr := json.Marshal(set)
	if jsonErr != nil {
		log.Printf("failed to marshal user set to cache: %v", jsonErr)
	} else {
		// 设置缓存过期时间（例如 5 分钟）
		myredis.Rdb.Set(c, cacheKey, cacheData, 5*time.Minute)
		log.Printf("缓存存储成功")
	}

	// 返回查询结果
	c.JSON(200, set)
}

func (h *Handler) DescribeUserRepos(c *gin.Context) {
	// 从请求中获取参数
	req := user.NewDescribeUserReposRequestFromHTTP(c.Request)

	// 构造缓存键，假设根据请求的参数生成唯一的缓存键
	cacheKey := fmt.Sprintf("userRepos:id:%s", req.Username) // 假设缓存是按用户ID进行区分

	// 尝试从 Redis 获取缓存数据
	cachedData, err := myredis.Rdb.Get(c, cacheKey).Result()
	if err == nil {
		// 缓存命中，返回缓存数据
		var jsonResponse map[string]interface{}
		if err := json.Unmarshal([]byte(cachedData), &jsonResponse); err != nil {
			c.JSON(500, gin.H{"error": "failed to parse cached JSON response"})
			return
		}
		c.JSON(200, jsonResponse)
		log.Println("缓存命中")
		return
	} else if err != redis.Nil {
		// 如果 Redis 发生其他错误，返回服务器错误
		c.JSON(500, gin.H{"error": "Failed to get data from cache"})
		return
	}

	// 缓存未命中，调用服务获取数据
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

	// 缓存数据到 Redis，设置缓存有效期（比如 1 小时）
	cacheData, _ := json.Marshal(jsonResponse)
	err = myredis.Rdb.Set(c, cacheKey, cacheData, time.Hour).Err()
	if err != nil {
		log.Printf("Failed to cache data for key %s: %v", cacheKey, err)
	}

	// 返回格式化后的 JSON
	c.JSON(200, jsonResponse)
	log.Println("缓存未命中！")
}

func (h *Handler) GetLocationCounts(c *gin.Context) {
	ctx := c.Request.Context()

	// 缓存键
	cacheKey := "locationCounts"

	// 尝试从 Redis 获取缓存数据
	cachedData, err := myredis.Rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中，将 JSON 字符串反序列化为结构体
		var locationCounts []*user.GetLocationCountsRequest
		if jsonErr := json.Unmarshal([]byte(cachedData), &locationCounts); jsonErr == nil {
			c.JSON(200, locationCounts) // 直接返回缓存数据
			log.Printf("缓存命中")
			return
		} else {
			// 如果反序列化失败，记录日志但继续处理
			log.Printf("failed to unmarshal cached data: %v", jsonErr)
		}
	} else if err != redis.Nil {
		// 如果是非缓存未命中的错误，直接返回 500
		c.JSON(500, gin.H{"error": "Failed to get data from cache"})
		return
	}

	// 如果缓存未命中，调用底层服务获取数据
	locationCounts, err := h.svc.GetLocationCounts(ctx)
	log.Printf("缓存未命中！")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get location counts"})
		return
	}

	// 序列化数据为 JSON 字符串
	data, err := json.Marshal(locationCounts)
	if err != nil {
		log.Printf("failed to marshal locationCounts: %v", err)
		c.JSON(500, gin.H{"error": "Failed to serialize data"})
		return
	}

	// 将结果缓存到 Redis，设置 1 小时过期时间
	err = myredis.Rdb.Set(ctx, cacheKey, data, time.Hour).Err()
	if err != nil {
		log.Printf("failed to cache data: %v", err)
		// 缓存失败不影响主流程，继续返回数据
	}

	// 返回获取到的数据
	c.JSON(200, locationCounts)
}

func (h *Handler) DeleteUserRepos(c *gin.Context) {
	req := &user.DeleteUserReposRequest{
		Id: c.Param("id"),
	}

	// 调用服务删除用户仓库
	result, err := h.svc.DeleteUserRepos(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 清空 Redis 数据库
	err = myredis.Rdb.FlushDB(c).Err() // 清空当前数据库
	// 如果需要清空所有数据库，则使用 FlushAll：
	// err = myredis.Rdb.FlushAll(c).Err()

	if err != nil {
		log.Printf("Failed to clear Redis database: %v", err)
		c.JSON(500, gin.H{"error": "Failed to clear Redis database"})
		return
	}

	log.Println("Redis database cleared successfully")

	// 返回删除结果
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

	// 将 Token 写入 Cookie
	c.SetCookie("githubToken", body.Token, int(7*24*time.Hour.Seconds()), "/", "", false, true)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Token set successfully"})
}
