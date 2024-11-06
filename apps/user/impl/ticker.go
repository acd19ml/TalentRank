package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/acd19ml/TalentRank/apps/user"
)

func (s *ServiceImpl) StartWeeklyUpdate(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		log.Println("Waiting for next update")
		select {
		case <-ticker.C:
			log.Println("Starting weekly update")

			if err := s.ScheduledUpdateUserRepos(ctx); err != nil {
				log.Printf("Error in weekly update: %v", err)
			}

			// 调用清理函数
			if err := s.DeleteOrphanedRepos(ctx); err != nil {
				log.Printf("Error deleting orphaned repos: %v", err)
			}

			if err := s.RemoveDuplicateUsers(ctx); err != nil {
				log.Printf("Error removing duplicate users: %v", err)
			}

		case <-ctx.Done():
			log.Println("Stopping weekly update")
			return
		}
	}
}

func (s *ServiceImpl) ScheduledUpdateUserRepos(ctx context.Context) error {
	log.Println("Starting getting all usernames.")
	users, err := s.GetAllUsernamesFromDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get usernames: %w", err)
	}

	for _, username := range users {
		newUserRepos, err := s.constructUserRepos(ctx, username)
		if err != nil {
			log.Printf("Failed to construct user repos for %s: %v", username, err)
			continue
		}

		// 比较新数据和数据库中的数据，更新差异
		err = s.CompareAndUpdateUserRepos(ctx, newUserRepos)
		if err != nil {
			log.Printf("Failed to update user repos for %s: %v", username, err)
		} else {
			log.Printf("Successfully updated user repos for: %s", username)
		}
	}
	// 完成所有用户更新后打印日志
	log.Println("Completed updating all user repos.")
	return nil
}

func (s *ServiceImpl) CompareAndUpdateUserRepos(ctx context.Context, newRepos *user.UserRepos) error {
	// 获取数据库中的数据
	existingRepos, err := s.FetchUserReposFromDB(ctx, newRepos.User.Username)
	if err != nil {
		return fmt.Errorf("failed to fetch existing data: %w", err)
	}

	// 比较新数据和数据库中的数据
	if HasDifferences(existingRepos, newRepos) {
		// 删除失效的仓库
		err := s.DeleteObsoleteRepos(ctx, existingRepos, newRepos)
		if err != nil {
			return fmt.Errorf("failed to delete obsolete repos: %w", err)
		}

		// 保存新数据
		err = s.save(ctx, newRepos)
		if err != nil {
			return fmt.Errorf("failed to save updated repos: %w", err)
		}
	}

	return nil
}

func (s *ServiceImpl) DeleteObsoleteRepos(ctx context.Context, oldRepos, newRepos *user.UserRepos) error {
	// 将新数据的仓库名称存入一个集合，以便快速查找
	newRepoMap := make(map[string]struct{})
	for _, newRepo := range newRepos.Repos {
		newRepoMap[newRepo.Repo] = struct{}{}
	}

	// 遍历旧数据中的仓库，删除新数据中不存在的仓库
	for _, oldRepo := range oldRepos.Repos {
		if _, exists := newRepoMap[oldRepo.Repo]; !exists {
			// 如果新数据中不存在该仓库，则从数据库中删除
			_, err := s.Db.ExecContext(ctx, "DELETE FROM repo WHERE id = ?", oldRepo.Id)
			if err != nil {
				log.Printf("Failed to delete obsolete repo %s: %v", oldRepo.Repo, err)
				return fmt.Errorf("failed to delete obsolete repo: %w", err)
			}
			log.Printf("Deleted obsolete repo: %s", oldRepo.Repo)
		}
	}
	return nil
}

func (s *ServiceImpl) GetAllUsernamesFromDB(ctx context.Context) ([]string, error) {
	rows, err := s.Db.QueryContext(ctx, "SELECT username FROM user")
	if err != nil {
		return nil, fmt.Errorf("failed to query usernames: %w", err)
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, fmt.Errorf("failed to scan username: %w", err)
		}
		log.Printf("Found username: %s\n", username)
		usernames = append(usernames, username)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return usernames, nil
}

func (s *ServiceImpl) FetchUserReposFromDB(ctx context.Context, username string) (*user.UserRepos, error) {
	// 获取用户信息
	row := s.Db.QueryRowContext(ctx, "SELECT * FROM user WHERE username = ?", username)
	var userIns user.User
	var orgsJSON []byte // 用于存储 organizations 字段的 JSON 数据

	if err := row.Scan(&userIns.Id, &userIns.Username, &userIns.Name, &userIns.Company, &userIns.Blog, &userIns.Location,
		&userIns.Email, &userIns.Bio, &userIns.Followers, &orgsJSON, &userIns.Readme, // 使用 orgsJSON 作为临时变量
		&userIns.Commits, &userIns.Score, &userIns.PossibleNation, &userIns.ConfidenceLevel); err != nil {
		return nil, fmt.Errorf("failed to scan user data: %w", err)
	}

	// 将 organizations 从 JSON 字符串解析为字符串数组
	if err := json.Unmarshal(orgsJSON, &userIns.Organizations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal organizations for user %s: %w", userIns.Id, err)
	}

	// 获取用户仓库信息
	rows, err := s.Db.QueryContext(ctx, "SELECT * FROM repo WHERE user_id = ?", userIns.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to query user repositories: %w", err)
	}
	defer rows.Close()

	var repos []*user.Repo
	for rows.Next() {
		var repo user.Repo
		if err := rows.Scan(&repo.Id, &repo.User_id, &repo.Repo, &repo.Star, &repo.Fork, &repo.Dependent, &repo.Commits,
			&repo.CommitsTotal, &repo.Issue, &repo.IssueTotal, &repo.PullRequest, &repo.PullRequestTotal,
			&repo.CodeReview, &repo.CodeReviewTotal, &repo.LineChange, &repo.LineChangeTotal); err != nil {
			return nil, fmt.Errorf("failed to scan repo data: %w", err)
		}
		repos = append(repos, &repo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over repository rows: %w", err)
	}

	return &user.UserRepos{User: &userIns, Repos: repos}, nil
}

func HasDifferences(oldRepos, newRepos *user.UserRepos) bool {
	// 比较用户的基本信息
	if oldRepos.User.Name != newRepos.User.Name ||
		oldRepos.User.Company != newRepos.User.Company ||
		oldRepos.User.Location != newRepos.User.Location ||
		oldRepos.User.Followers != newRepos.User.Followers ||
		oldRepos.User.PossibleNation != newRepos.User.PossibleNation {
		return true
	}

	// 比较每个仓库信息
	oldRepoMap := make(map[string]*user.Repo)
	for _, repo := range oldRepos.Repos {
		oldRepoMap[repo.Repo] = repo
	}

	for _, newRepo := range newRepos.Repos {
		if oldRepo, exists := oldRepoMap[newRepo.Repo]; exists {
			if oldRepo.Star != newRepo.Star || oldRepo.Fork != newRepo.Fork || oldRepo.Dependent != newRepo.Dependent ||
				oldRepo.Commits != newRepo.Commits || oldRepo.Issue != newRepo.Issue || oldRepo.PullRequest != newRepo.PullRequest ||
				oldRepo.CodeReview != newRepo.CodeReview || oldRepo.LineChange != newRepo.LineChange {
				return true
			}
		} else {
			// 如果在旧数据中不存在这个仓库，则表示有差异（新增仓库）
			return true
		}
	}

	return false
}

// DeleteOrphanedRepos 删除user_id在user表中没有匹配记录的repo记录
func (s *ServiceImpl) DeleteOrphanedRepos(ctx context.Context) error {
	_, err := s.Db.ExecContext(ctx, `
		DELETE FROM repo
		WHERE user_id NOT IN (SELECT id FROM user)
	`)
	if err != nil {
		return err
	}
	log.Println("Successfully deleted orphaned repos")
	return nil
}

// RemoveDuplicateUsers 删除重复的username，只保留其中一条记录
func (s *ServiceImpl) RemoveDuplicateUsers(ctx context.Context) error {
	_, err := s.Db.ExecContext(ctx, `
		WITH UserDuplicates AS (
			SELECT id,
			       ROW_NUMBER() OVER (PARTITION BY username ORDER BY id) AS row_num
			FROM user
		)
		DELETE FROM user
		WHERE id IN (
			SELECT id
			FROM UserDuplicates
			WHERE row_num > 1
		)
	`)
	if err != nil {
		return err
	}
	log.Println("Successfully removed duplicate users")
	return nil
}
