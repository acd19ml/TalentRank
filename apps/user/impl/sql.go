package impl

const (
	InsertUserSQL = `
		INSERT INTO user (id, username, name, company, blog, location, email, bio, followers, organizations, readme, commits, score, possible_nation, confidence_level)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			name = VALUES(name),
			company = VALUES(company),
			blog = VALUES(blog),
			location = VALUES(location),
			email = VALUES(email),
			bio = VALUES(bio),
			followers = VALUES(followers),
			organizations = VALUES(organizations),
			readme = VALUES(readme),
			commits = VALUES(commits),
			score = VALUES(score),
			possible_nation = VALUES(possible_nation),
			confidence_level = VALUES(confidence_level);

	`

	InsertRepoSQL = `
		INSERT INTO repo (id, user_id, repo, star, fork, dependent, commits, commits_total, issue, issue_total, pull_request, pull_request_total, code_review, code_review_total, line_change, line_change_total)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			star = VALUES(star),
			fork = VALUES(fork),
			dependent = VALUES(dependent),
			commits = VALUES(commits),
			commits_total = VALUES(commits_total),
			issue = VALUES(issue),
			issue_total = VALUES(issue_total),
			pull_request = VALUES(pull_request),
			pull_request_total = VALUES(pull_request_total),
			code_review = VALUES(code_review),
			code_review_total = VALUES(code_review_total),
			line_change = VALUES(line_change),
			line_change_total = VALUES(line_change_total);

	`

	Userquery = `
		SELECT a.id, username, name, company, blog,
       location,
       email, bio, 
       followers, organizations, round(score) score, 
       possible_nation, confidence_level,
       rank() OVER (ORDER BY score DESC) AS rankno
FROM User a 
		LIMIT ? OFFSET ?;
	

	`

	QueryUser = `
		SELECT id, username, name, company, blog, location, email, bio, 
			   followers, organizations, score, possible_nation, confidence_level
		FROM User
		WHERE username = ?;
	`

	QueryRepos = `
		SELECT id, user_id, repo, star, fork, dependent, commits, commits_total, 
			   issue, issue_total, pull_request, pull_request_total, 
			   code_review, code_review_total, line_change, line_change_total
		FROM Repo
		WHERE user_id = ?;
	`

	DeleteReposSQL = `
		DELETE FROM repo WHERE user_id = ?;
	`

	DeleteUserSQL = `
		DELETE FROM user WHERE id = ?;
	`

	UpdateUserScoreSQL = `
		INSERT INTO User (username, score)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE
		score = VALUES(score)
	`
)

// CREATE TABLE User (
//     id CHAR(36) PRIMARY KEY,
//     username VARCHAR(255) NOT NULL,
//     name VARCHAR(255) DEFAULT '',
//     company VARCHAR(255) DEFAULT '',
//     blog VARCHAR(255) DEFAULT '',
//     location VARCHAR(255) DEFAULT '',
//     email VARCHAR(255) DEFAULT '',
//     bio TEXT DEFAULT '',
//     followers INT DEFAULT 0,
//     organizations JSON DEFAULT '[]',
//     readme TEXT DEFAULT '',
//     commits TEXT DEFAULT '',
//     score DOUBLE DEFAULT 0,
//     possible_nation VARCHAR(255) DEFAULT '',
//     confidence_level TINYINT DEFAULT 0
// );

// CREATE TABLE Repo (
//     id CHAR(36) PRIMARY KEY,
//     user_id CHAR(36) NOT NULL,
//     repo VARCHAR(255) NOT NULL,
//     star INT DEFAULT 0,
//     fork INT DEFAULT 0,
//     dependent INT DEFAULT 0,
// 	   commits INT DEFAULT 0,
// 	   commits_total INT DEFAULT 0,
//     issue INT DEFAULT 0,
//     issue_total INT DEFAULT 0,
//     pull_request INT DEFAULT 0,
//     pull_request_total INT DEFAULT 0,
//     code_review INT DEFAULT 0,
//     code_review_total INT DEFAULT 0,
//     line_change INT DEFAULT 0,
//     line_change_total INT DEFAULT 0
// );
