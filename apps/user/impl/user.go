package impl

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/acd19ml/TalentRank/apps/git"
	"github.com/acd19ml/TalentRank/apps/user"
	"github.com/google/uuid"
)

func (s *ServiceImpl) CreateUserRepos(ctx context.Context, username string) (*user.UserRepos, error) {
	// 使用带有认证的ctx
	ins, err := s.constructUserRepos(s.NewAuthenticatedContext(ctx), username)
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

func (s *ServiceImpl) SaveRepoDataToDB(ctx context.Context, username, repoName, function string, result interface{}) error {
	// 获取 user_id，确保 user_id 存在
	var userID string
	err := s.Db.QueryRowContext(ctx, "SELECT id FROM user WHERE username = ?", username).Scan(&userID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("user not found for username: %s", username)
	} else if err != nil {
		return fmt.Errorf("failed to query user_id for username %s: %w", username, err)
	}

	// 检查是否存在该 repo，若不存在则插入基础数据
	var repoID string
	err = s.Db.QueryRowContext(ctx, "SELECT id FROM Repo WHERE user_id = ? AND repo = ?", userID, repoName).Scan(&repoID)
	if err == sql.ErrNoRows {
		// 如果 repo 不存在，生成新的 repo_id 并插入基础记录
		repoID = uuid.New().String()
		_, err = s.Db.ExecContext(ctx, `
			INSERT INTO Repo (id, user_id, repo) VALUES (?, ?, ?)
		`, repoID, userID, repoName)
		if err != nil {
			return fmt.Errorf("failed to insert new repo: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to query repo_id for repo %s: %w", repoName, err)
	}

	// 根据 function 更新对应的字段
	switch function {
	case "GetStarsByRepo":
		resp := result.(*git.IntResponse)
		_, err = s.Db.ExecContext(ctx, `
			UPDATE Repo SET star = ? WHERE id = ?
		`, resp.Result, repoID)

	case "GetForksByRepo":
		resp := result.(*git.IntResponse)
		_, err = s.Db.ExecContext(ctx, `
			UPDATE Repo SET fork = ? WHERE id = ?
		`, resp.Result, repoID)

	case "GetTotalIssuesByRepo":
		resp := result.(*git.IntResponse)
		_, err = s.Db.ExecContext(ctx, `
			UPDATE Repo SET issue_total = ? WHERE id = ?
		`, resp.Result, repoID)

	case "GetUserSolvedIssuesByRepo":
		resp := result.(*git.IntResponse)
		_, err = s.Db.ExecContext(ctx, `
			UPDATE Repo SET issue = ? WHERE id = ?
		`, resp.Result, repoID)

	case "GetTotalPullRequestsByRepo":
		resp := result.(*git.IntResponse)
		_, err = s.Db.ExecContext(ctx, `
			UPDATE Repo SET pull_request_total = ? WHERE id = ?
		`, resp.Result, repoID)

	case "GetUserMergedPullRequestsByRepo":
		resp := result.(*git.IntResponse)
		_, err = s.Db.ExecContext(ctx, `
			UPDATE Repo SET pull_request = ? WHERE id = ?
		`, resp.Result, repoID)

	case "GetTotalCodeReviewsByRepo":
		resp := result.(*git.IntResponse)
		_, err = s.Db.ExecContext(ctx, `
			UPDATE Repo SET code_review_total = ? WHERE id = ?
		`, resp.Result, repoID)

	case "GetUserCodeReviewsByRepo":
		resp := result.(*git.IntResponse)
		_, err = s.Db.ExecContext(ctx, `
			UPDATE Repo SET code_review = ? WHERE id = ?
		`, resp.Result, repoID)

	case "GetLineChangesCommitsByRepo":
		resp := result.(*git.RepoLineChangesCommitsResponse)
		_, err = s.Db.ExecContext(ctx, `
			UPDATE Repo 
			SET line_change = ?, line_change_total = ?, commits = ?, commits_total = ? 
			WHERE id = ?
		`, resp.UserChanges, resp.TotalChanges, resp.UserCommits, resp.TotalCommits, repoID)

	case "GetDependentRepositoriesByRepo":
		resp := result.(*git.IntResponse)
		_, err = s.Db.ExecContext(ctx, `
			UPDATE Repo SET dependent = ? WHERE id = ?
		`, resp.Result, repoID)

	default:
		return fmt.Errorf("unknown function: %s", function)
	}

	// 检查是否更新成功
	if err != nil {
		return fmt.Errorf("failed to update repo %s for function %s: %w", repoName, function, err)
	}

	return nil
}

func (s *ServiceImpl) GetAllUsernamesFromDB(ctx context.Context) ([]string, error) {
	// 查询所有用户名
	rows, err := s.Db.QueryContext(ctx, "SELECT username FROM user")
	if err != nil {
		return nil, fmt.Errorf("failed to query usernames: %w", err)
	}
	defer rows.Close()

	// 保存用户名的切片
	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, fmt.Errorf("failed to scan username: %w", err)
		}
		usernames = append(usernames, username)
	}

	// 检查是否有错误
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return usernames, nil
}

// FetchUserFromDB 从数据库中获取用户信息
func (s *ServiceImpl) FetchUserFromDB(ctx context.Context, username string) (*user.User, error) {
	row := s.Db.QueryRowContext(ctx, "SELECT * FROM user WHERE username = ?", username)
	var userIns user.User
	var orgsJSON []byte // 用于存储 organizations 字段的 JSON 数据

	if err := row.Scan(&userIns.Id, &userIns.Username, &userIns.Name, &userIns.Company, &userIns.Blog, &userIns.Location,
		&userIns.Email, &userIns.Bio, &userIns.Followers, &orgsJSON, &userIns.Readme,
		&userIns.Commits, &userIns.Score, &userIns.PossibleNation, &userIns.ConfidenceLevel); err != nil {
		return nil, fmt.Errorf("failed to scan user data: %w", err)
	}

	// 将 organizations 从 JSON 字符串解析为字符串数组
	if err := json.Unmarshal(orgsJSON, &userIns.Organizations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal organizations for user %s: %w", userIns.Id, err)
	}

	return &userIns, nil
}

// FetchReposFromDB 从数据库中获取用户的仓库名称
// FetchReposFromDB 从数据库中获取用户的仓库名称
func (s *ServiceImpl) FetchReposFromDB(ctx context.Context, username string) ([]string, error) {
	// 从 user 表中获取 user_id
	var userID string
	err := s.Db.QueryRowContext(ctx, "SELECT id FROM user WHERE username = ?", username).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user ID for username %s: %w", username, err)
	}

	// 使用 user_id 查询 repo 表中的 repo 字段
	rows, err := s.Db.QueryContext(ctx, "SELECT repo FROM repo WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query repositories for user ID %d: %w", userID, err)
	}
	defer rows.Close()

	// 收集所有 repo 名称
	var repoNames []string
	for rows.Next() {
		var repoName string
		if err := rows.Scan(&repoName); err != nil {
			return nil, fmt.Errorf("failed to scan repo name: %w", err)
		}
		repoNames = append(repoNames, repoName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over repository rows: %w", err)
	}

	return repoNames, nil
}

// HasDifferencesUser 比较两个用户的基本信息是否存在差异
func HasDifferencesUser(oldUser, newUser *user.User) bool {
	return oldUser.Name != newUser.Name ||
		oldUser.Company != newUser.Company ||
		oldUser.Location != newUser.Location ||
		oldUser.Followers != newUser.Followers ||
		oldUser.PossibleNation != newUser.PossibleNation
}

func HasInvalidRepos(oldRepos, newRepos []string) []string {
	// 创建一个 map 来存储 newRepos 中的所有仓库
	newRepoSet := make(map[string]struct{})
	for _, repo := range newRepos {
		newRepoSet[repo] = struct{}{}
	}

	// 遍历 oldRepos，查找不在 newRepos 中的仓库
	var invalidRepos []string
	for _, repo := range oldRepos {
		if _, exists := newRepoSet[repo]; !exists {
			invalidRepos = append(invalidRepos, repo)
		}
	}

	return invalidRepos
}

func (s *ServiceImpl) DeleteInvalidReposFromDB(ctx context.Context, repos []string) error {
	if len(repos) == 0 {
		return nil // 如果没有要删除的仓库，直接返回
	}

	// 构建动态查询语句
	placeholders := make([]string, len(repos))
	args := make([]interface{}, len(repos))
	for i, repo := range repos {
		placeholders[i] = "?"
		args[i] = repo
	}
	query := fmt.Sprintf("DELETE FROM repo WHERE repo IN (%s)", strings.Join(placeholders, ","))

	// 执行删除操作
	_, err := s.Db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete invalid repos: %w", err)
	}

	return nil
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
	// log.Println("Successfully deleted orphaned repos")
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
	// log.Println("Successfully removed duplicate users")
	return nil
}

func (s *ServiceImpl) SaveUserDataToDB(ctx context.Context, username, function string, result interface{}) error {
	var query string
	var value interface{}

	switch function {
	case "GetName":
		query = "UPDATE User SET name = ? WHERE username = ?"
		value = result.(*git.StringResponse).Result
	case "GetCompany":
		query = "UPDATE User SET company = ? WHERE username = ?"
		value = result.(*git.StringResponse).Result
	case "GetLocation":
		query = "UPDATE User SET location = ? WHERE username = ?"
		value = result.(*git.StringResponse).Result
	case "GetEmail":
		query = "UPDATE User SET email = ? WHERE username = ?"
		value = result.(*git.StringResponse).Result
	case "GetBio":
		query = "UPDATE User SET bio = ? WHERE username = ?"
		value = result.(*git.StringResponse).Result
	case "GetOrganizations":
		query = "UPDATE User SET organizations = ? WHERE username = ?"
		valueBytes, _ := json.Marshal(result.(*git.StringListResponse).Result)
		value = string(valueBytes) // 转为 JSON 字符串存储
	case "GetFollowers":
		query = "UPDATE User SET followers = ? WHERE username = ?"
		value = result.(*git.IntResponse).Result
	default:
		return fmt.Errorf("unknown function: %s", function)
	}

	_, err := s.Db.ExecContext(ctx, query, value, username)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}
	return nil
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

func (s *ServiceImpl) UpdateUserScore(ctx context.Context, username string, score float64) error {

	// Prepare the statement
	stmt, err := s.Db.PrepareContext(ctx, UpdateUserScoreSQL)
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return err
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.ExecContext(ctx, username, score)
	if err != nil {
		log.Printf("Failed to execute statement for user %s: %v", username, err)
		return err
	}

	log.Printf("Successfully updated score for user %s to %.2f", username, score)
	return nil
}
