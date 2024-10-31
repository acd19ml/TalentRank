package impl

import (
	"database/sql"

	"github.com/acd19ml/TalentRank/apps"
	"github.com/acd19ml/TalentRank/apps/git"
	"github.com/acd19ml/TalentRank/apps/git/impl"
	"github.com/acd19ml/TalentRank/apps/user"
	"github.com/acd19ml/TalentRank/apps/user/llm"
	_ "github.com/acd19ml/TalentRank/apps/user/llm"
	"github.com/acd19ml/TalentRank/conf"
)

func NewUserServiceImpl() *ServiceImpl {
	svc := impl.NewGitClient()
	llmsvc := llm.NewChatService()
	rspsvc := user.NewUserResponseByLLM()
	return &ServiceImpl{
		db:  conf.C().MySQL.GetDB(),
		svc: svc,
		llm: llmsvc,
		rsp: rspsvc,
	}
}

type ServiceImpl struct {
	db  *sql.DB
	svc git.Service
	llm user.LLMService
	rsp user.LLMResponseService
}

var svcimpl = &ServiceImpl{}

func (s *ServiceImpl) Config() {
	s.db = conf.C().MySQL.GetDB()
}

func (s *ServiceImpl) Name() string {
	return user.AppName
}

// _ import app 自动执行注册逻辑
func init() {
	//  对象注册到ioc层
	apps.RegistryImpl(svcimpl)
}
