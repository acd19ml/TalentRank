package impl

const (
	InsertUserSQL = `
	INSERT INTO User (
		Id, Username, Name, Company, Blog, Location, Email, Bio, 
		TotalStar, TotalFork, Followers, Dependents, Organizations, 
		Readme, Commits
	) VALUES (
		?, ?, ?, ?, ?, ?, ?, ?, 
		?, ?, ?, ?, ?, 
		?, ?
	);
	`
)
