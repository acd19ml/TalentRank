package http

import (
	"github.com/acd19ml/TalentRank/apps"
	"github.com/acd19ml/TalentRank/apps/user"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc user.Service
}

var handler = &Handler{}

func (h *Handler) Registry(r gin.IRouter) {
	r.POST("/userRepos", h.CreateUserRepos)
	r.GET("/user", h.QueryUsers)
	r.GET("/userRepos", h.DescribeUserRepos)
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
