package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/acd19ml/TalentRank/apps/user"
)

func (u *UserServiceImpl) save(ctx context.Context, ins *user.UserRepos) error {

	var (
		err error
	)

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

	ustmt, err := tx.PrepareContext(ctx, InsertUserSQL)
	if err != nil {
		return err
	}
	defer ustmt.Close()

	// 插入的 Organizations 字段为 []string
	organizationsJSON, err := json.Marshal(ins.Organizations)
	if err != nil {
		return fmt.Errorf("failed to marshal options: %v", err)
	}

	// 执行插入语句
	result, err := ustmt.ExecContext(ctx,
		ins.Id, ins.Username, ins.Name, ins.Company, ins.Blog, ins.Location,
		ins.Email, ins.Bio, ins.Followers, string(organizationsJSON), ins.Readme,
		ins.Commits, ins.Score, ins.PossibleNation, ins.ConfidenceLevel,
	)
	if err != nil {
		return err
	} else {
		fmt.Printf("insert user success, %v", result)
	}

	rstmt, err := tx.PrepareContext(ctx, InsertRepoSQL)
	if err != nil {
		return err
	}
	defer rstmt.Close()
	for _, repo := range ins.Repos {
		// 执行插入语句
		result, err := rstmt.ExecContext(ctx,
			repo.Id, ins.Id, repo.Repo, repo.Star, repo.Fork, repo.Dependent,
			repo.Issue, repo.IssueTotal, repo.PullRequest, repo.PullRequestTotal,
			repo.CodeReview, repo.CodeReviewTotal, repo.LineChange, repo.LineChangeTotal,
		)
		if err != nil {
			return err
		} else {
			fmt.Printf("insert repo success, %v", result)
		}
	}

	return nil
}
