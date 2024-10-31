package impl

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/acd19ml/TalentRank/apps/user"
)

func (s *ServiceImpl) CreateUserRepos(ctx context.Context, username string) (*user.UserRepos, error) {
	ins, err := s.constructUserRepos(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to construct user repos: %w", err)
	}

	if err = s.save(ctx, ins); err != nil {
		return nil, fmt.Errorf("failed to save user repos: %w", err)
	}
	return ins, nil
}

func (s *ServiceImpl) QueryUsers(ctx context.Context, req *user.QueryUserRequest) (*user.UserSet, error) {
	// 初始化返回结果
	result := user.NewUserSet()
	offset := req.OffSet()
	pageSize := req.GetPageSize()

	// 动态构建查询语句
	query := Userquery
	var args []interface{}

	// 如果 `location` 有值，则在查询中添加条件
	if req.Location != "" {
		query = `
			SELECT id, username, name, company, blog, location, email, bio, 
				   followers, organizations, score, possible_nation, confidence_level
			FROM User
			WHERE location = ?
			LIMIT ? OFFSET ?;
		`
		args = append(args, req.Location, pageSize, offset)
	} else {
		// `location` 为空时，直接使用基础查询语句
		query = Userquery
		args = append(args, pageSize, offset)
	}

	// 准备查询语句
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query for users: %w", err)
	}
	defer stmt.Close()

	// 执行查询
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute user query: %w", err)
	}
	defer rows.Close()

	// 遍历查询结果
	for rows.Next() {
		user := user.NewUser()
		var orgs string // 临时存储 organizations 字段的 JSON 数据

		// 扫描用户数据
		if err := rows.Scan(&user.Id, &user.Username, &user.Name, &user.Company, &user.Blog,
			&user.Location, &user.Email, &user.Bio, &user.Followers, &orgs, &user.Score,
			&user.PossibleNation, &user.ConfidenceLevel); err != nil {
			return nil, fmt.Errorf("failed to scan user data: %w", err)
		}

		// 将 organizations 从 JSON 转为字符串数组
		if err := json.Unmarshal([]byte(orgs), &user.Organizations); err != nil {
			return nil, fmt.Errorf("failed to unmarshal organizations for user %s: %w", user.Id, err)
		}

		// 将 user 添加到结果集
		result.Add(user)
	}

	// 检查遍历错误
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during user rows iteration: %w", err)
	}

	// 查询用户总数（考虑 location 条件）
	countQuery := "SELECT COUNT(*) FROM User"
	var countArgs []interface{}

	if req.Location != "" {
		countQuery += " WHERE location = ?"
		countArgs = append(countArgs, req.Location)
	}

	err = s.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&result.Total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count of users: %w", err)
	}

	return result, nil
}

func (s *ServiceImpl) DescribeUserRepos(ctx context.Context, req *user.DescribeUserReposRequest) (*user.UserRepos, error) {
	// 初始化返回结果
	userRepos := user.NewUserRepos()

	// 准备并执行用户查询
	err := s.db.QueryRowContext(ctx, QueryUser, req.Username).Scan(
		&userRepos.User.Id,
		&userRepos.User.Username,
		&userRepos.User.Name,
		&userRepos.User.Company,
		&userRepos.User.Blog,
		&userRepos.User.Location,
		&userRepos.User.Email,
		&userRepos.User.Bio,
		&userRepos.User.Followers,
		new(string), // 临时变量存储 organizations
		&userRepos.User.Score,
		&userRepos.User.PossibleNation,
		&userRepos.User.ConfidenceLevel,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with username %s not found", req.Username)
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	// 解析 organizations 字段
	var orgs string
	if err := json.Unmarshal([]byte(orgs), &userRepos.User.Organizations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal organizations for user %s: %w", userRepos.User.Id, err)
	}

	// 准备 Repo 查询语句
	rows, err := s.db.QueryContext(ctx, QueryRepos, userRepos.User.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to query repos for user %s: %w", userRepos.User.Id, err)
	}
	defer rows.Close()

	// 遍历查询结果，填充 Repo 信息
	for rows.Next() {
		repo := user.NewRepo()
		if err := rows.Scan(
			&repo.Id,
			&repo.User_id,
			&repo.Repo,
			&repo.Star,
			&repo.Fork,
			&repo.Dependent,
			&repo.Commits,
			&repo.CommitsTotal,
			&repo.Issue,
			&repo.IssueTotal,
			&repo.PullRequest,
			&repo.PullRequestTotal,
			&repo.CodeReview,
			&repo.CodeReviewTotal,
			&repo.LineChange,
			&repo.LineChangeTotal,
		); err != nil {
			return nil, fmt.Errorf("failed to scan repo data for user %s: %w", userRepos.User.Id, err)
		}

		// 添加 Repo 到用户的 Repos 列表
		userRepos.AddRepos(repo)
	}

	// 检查遍历错误
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during repo rows iteration for user %s: %w", userRepos.User.Id, err)
	}

	return userRepos, nil
}
