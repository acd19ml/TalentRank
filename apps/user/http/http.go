package http

import (
	"github.com/acd19ml/TalentRank/apps"
	"github.com/acd19ml/TalentRank/apps/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc user.Service
}

var handler = &Handler{}

func (h *Handler) Registry(r gin.IRouter) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},                   // 允许的源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // 允许的方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 允许的头
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/user", h.QueryUsers)
	r.GET("/userRepos", h.DescribeUserRepos)
	r.GET("/api/locations", h.GetLocationCounts)
	r.DELETE("/userRepos:id", h.DeleteUserRepos)
	r.POST("/setToken", h.setTokenHandler)
	r.POST("/userRepos", h.CreateUserRepos)
}

func (h *Handler) Config() {
	if apps.UserService == nil {
		panic("UserService required")
	}

	// 从ioc获取UserService的实例，代替原来从构造函数传入
	h.svc = apps.GetImpl(user.AppName).(user.Service)
}

func (h *Handler) Name() string {
	return user.AppName
}

func init() {
	apps.RegistryGin(handler)
}
