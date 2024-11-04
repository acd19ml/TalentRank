package impl

// Commit 结构体用来存储提交信息
type Commit struct {
	Sha    string `json:"sha"`
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
}

// CommitDetail 存储每个提交的代码行变化信息
type CommitDetail struct {
	Stats struct {
		Additions int `json:"additions"`
		Deletions int `json:"deletions"`
	} `json:"stats"`
}
