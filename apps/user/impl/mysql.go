package impl

import (
	"database/sql"

	"github.com/acd19ml/TalentRank/conf"
)

func NewUserServiceImpl(db *sql.DB) *UserServiceImpl {
	return &UserServiceImpl{
		db: conf.C().MySQL.GetDB(),
	}
}

type UserServiceImpl struct {
	db *sql.DB
}
