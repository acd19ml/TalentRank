package http

import (
	"encoding/json"
	"github.com/acd19ml/TalentRank/apps/user"
	"github.com/gin-gonic/gin"
	"log"
)

func (h *Handler) CreateUserRepos(c *gin.Context) {
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

//func (h *Handler) GetLocationCounts(c *gin.Context) {
//	rows, err := h.svc.GetLocationCounts() // svc 是一个 Service 类型的字段
//	if err != nil {
//		fmt.Println("Error calling stored procedure:", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
//		return
//	}
//	defer rows.Close() // 确保 rows 是 *sql.Rows 类型
//
//	var locationCounts []LocationCount
//	for rows.Next() {
//		var locationCount LocationCount
//		if err := rows.Scan(&locationCount.Location, &locationCount.Count); err != nil {
//			fmt.Println("Error scanning row:", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
//			return
//		}
//		locationCounts = append(locationCounts, locationCount)
//	}
//
//	c.JSON(http.StatusOK, locationCounts)
//}
//
//type LocationCount struct {
//	Location string `json:"location"`
//	Count    int    `json:"count"`
//}