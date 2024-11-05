package impl

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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
			SELECT a.id, username, name, company, blog,
       COALESCE(a.location, '') AS Location,  -- 使用 COALESCE 替代 NULL 值
       email, bio, 
       followers, organizations, round(score) AS score, 
       possible_nation, confidence_level,
       rank() OVER (ORDER BY score DESC) AS rankno
FROM User a
JOIN countries c
    ON a.location LIKE CONCAT('%', c.country_name, '%')
WHERE c.country_name = ?
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
			&user.PossibleNation, &user.ConfidenceLevel, &user.Rankno); err != nil {
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

func (s *ServiceImpl) DescribeUserRepos(ctx context.Context, req *user.DescribeUserReposRequest) (string, error) {
	// 设置 @result 为 NULL
	_, err := s.db.ExecContext(ctx, "SET @result = NULL;")
	if err != nil {
		return "", fmt.Errorf("failed to set result variable: %w", err)
	}

	// 调用存储过程 GetUserData
	_, err = s.db.ExecContext(ctx, "CALL GetUserData(?, @result);", req.Username)
	if err != nil {
		return "", fmt.Errorf("failed to execute stored procedure: %w", err)
	}

	// 获取 @result 的值
	var result sql.NullString
	err = s.db.QueryRowContext(ctx, "SELECT @result;").Scan(&result)
	if err != nil {
		return "", fmt.Errorf("failed to fetch result: %w", err)
	}

	// 检查结果是否为空
	if !result.Valid {
		return "", errors.New("no data found for the specified user")
	}

	return result.String, nil
}
