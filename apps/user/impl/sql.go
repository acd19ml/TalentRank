package impl

const (
	InsertUserSQL = `
	INSERT INTO User (
		id, username, name, company, blog, location, 
		email, bio, followers, organizations, readme, 
		commits, score, possible_nation, confidence_level
	) VALUES (
		?, ?, ?, ?, ?, ?, ?, ?, 
		?, ?, ?, ?, ?, ?, ?
	);

	`

	InsertRepoSQL = `
	INSERT INTO Repo (
		id, user_id, repo, star, fork, dependent, commits, commits_total,
		issue, issue_total, pull_request, pull_request_total, 
		code_review, code_review_total, line_change, line_change_total
	) VALUES (
	 	?, ?, ?, ?, ?, ?, ?, 
		?, ?, ?, ?, ?, ?, ?
	);

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
