package impl

import (
	"database/sql"

	"github.com/acd19ml/TalentRank/apps/git"
	"github.com/acd19ml/TalentRank/apps/git/impl"
	"github.com/acd19ml/TalentRank/conf"
)

func NewUserServiceImpl() *ServiceImpl {
	svc := impl.NewGitClient()
	return &ServiceImpl{
		db:  conf.C().MySQL.GetDB(),
		svc: svc,
	}
}

type ServiceImpl struct {
	db  *sql.DB
	svc git.Service
}
