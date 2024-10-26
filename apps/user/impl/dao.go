package impl

import (
	"context"
	"log"

	"github.com/acd19ml/TalentRank/apps/user"
)

func (u *UserServiceImpl) save(ctx context.Context, ins *user.User) error {
	// 开启事务
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// 通过defer处理事务提交
	// 1. 没有报错Commit
	// 2. 有报错Rollback
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("rollback error, %s\n", err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				log.Printf("rollback error, %s\n", err)
			}
		}
	}()

	// 执行插入语句
	_, err = tx.ExecContext(ctx, InsertUserSQL,
		ins.Id, ins.Username, ins.Name, ins.Company, ins.Blog, ins.Location,
		ins.Email, ins.Bio, ins.TotalStar, ins.TotalFork, ins.Followers,
		ins.Dependents, ins.Organizations, ins.Readme, ins.Commits,
	)

	return nil
}
