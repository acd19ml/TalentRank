package impl

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/acd19ml/TalentRank/apps/user"
)

func (s *ServiceImpl) CreateUserRepos(ctx context.Context, username string) (*user.UserRepos, error) {
	// 使用带有认证的ctx
	ins, err := s.constructUserRepos(s.NewAuthenticatedContext(), username)
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

	// 构建查询语句和参数
	var query string
	var args []interface{}

	// 如果 `location` 有值，则在查询中添加条件
	if req.PossibleNation != "" {
		query = `
			SELECT * FROM (
				SELECT a.id, username, name, company, blog, location, email, bio, 
				       followers, organizations, ROUND(score) AS score, 
				       possible_nation, confidence_level,
				       RANK() OVER (ORDER BY score DESC) AS rankno
				FROM User a
				WHERE 
				(SELECT country_name FROM countries b WHERE a.possible_nation LIKE CONCAT('%', b.country_name, '%') LIMIT 1) = ?
			) AS filtered_users
			LIMIT ? OFFSET ?;
		`
		args = append(args, req.PossibleNation, pageSize, offset)
	} else {
		// `location` 为空时，直接使用基础查询语句
		query = Userquery
		args = append(args, pageSize, offset)
	}

	// 准备查询语句
	stmt, err := s.Db.PrepareContext(ctx, query)
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

	// 查询符合条件的用户总数
	countQuery := "SELECT COUNT(*) FROM User"
	var countArgs []interface{}

	if req.PossibleNation != "" {
		countQuery += " WHERE possible_nation = ?"
		countArgs = append(countArgs, req.PossibleNation)
	}

	err = s.Db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&result.Total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count of users: %w", err)
	}

	return result, nil
}

func (s *ServiceImpl) DescribeUserRepos(ctx context.Context, req *user.DescribeUserReposRequest) (string, error) {
	// 设置 @result 为 NULL
	_, err := s.Db.ExecContext(ctx, "SET @result = NULL;")
	if err != nil {
		return "", fmt.Errorf("failed to set result variable: %w", err)
	}

	// 调用存储过程 GetUserData
	_, err = s.Db.ExecContext(ctx, "CALL GetUserData(?, @result);", req.Username)
	if err != nil {
		return "", fmt.Errorf("failed to execute stored procedure: %w", err)
	}

	// 获取 @result 的值
	var result sql.NullString
	err = s.Db.QueryRowContext(ctx, "SELECT @result;").Scan(&result)
	if err != nil {
		return "", fmt.Errorf("failed to fetch result: %w", err)
	}

	// 检查结果是否为空
	if !result.Valid {
		return "", errors.New("no data found for the specified user")
	}

	return result.String, nil
}

func (s *ServiceImpl) GetLocationCounts(ctx context.Context) ([]*user.GetLocationCountsRequest, error) {
	rows, err := s.Db.QueryContext(ctx, "CALL GetLocation();")
	if err != nil {
		return nil, fmt.Errorf("failed to execute stored procedure: %w", err)
	}
	defer rows.Close()

	var locationCounts []*user.GetLocationCountsRequest

	for rows.Next() {
		locationCount := user.NewGetLocationCountsRequest()
		if err := rows.Scan(&locationCount.CountryName, &locationCount.Count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		locationCounts = append(locationCounts, locationCount)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return locationCounts, nil
}

func (s *ServiceImpl) DeleteUserRepos(ctx context.Context, req *user.DeleteUserReposRequest) (*user.DeleteUserReposResponse, error) {
	var (
		err      error
		ustmt    *sql.Stmt
		rstmt    *sql.Stmt
		qstmt    *sql.Stmt
		username string
	)
	result := user.NewDeleteUserReposResponse()

	// 初始化一个事务
	tx, err := s.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				log.Fatalf("tx rollback error, %s", err)
			}
		} else {
			err := tx.Commit()
			if err != nil {
				log.Fatalf("tx commit error, %s", err)
			}
		}
	}()

	qstmt, err = tx.Prepare("SELECT username FROM user WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query user statement: %w", err)
	}
	defer qstmt.Close()

	err = qstmt.QueryRow(req.Id).Scan(&username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch username for user id %s: %w", req.Id, err)
	}

	ustmt, err = tx.Prepare(DeleteUserSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare delete user statement: %w", err)
	}
	defer ustmt.Close()

	_, err = ustmt.Exec(req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	// 设置删除repos的语句
	rstmt, err = tx.Prepare(DeleteReposSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare delete repos statement: %w", err)
	}
	defer rstmt.Close()

	_, err = rstmt.Exec(req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete repos: %w", err)
	}

	result.Username = username

	return result, nil
}
