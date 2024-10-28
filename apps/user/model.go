package user

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var (
	validate = validator.New()
)

func NewUserReposSet() *UserReposSet {
	return &UserReposSet{
		Developers: []*UserRepos{},
	}
}

type UserReposSet struct {
	Total      int          `json:"total"`
	Developers []*UserRepos `json:"developers"`
}

func NewUserRepos() *UserRepos {
	return &UserRepos{
		User:  &User{},
		Repos: []*Repo{},
	}
}

type UserRepos struct {
	*User
	Repos []*Repo
}

func (ur *UserRepos) Validate() error {
	return validate.Struct(ur)
}

// char(36)
func (ur *UserRepos) InjectDefault() {
	if ur.Id == "" {
		ur.Id = uuid.New().String()
	}
}

type User struct {
	Id              string   `json:"id"`
	Username        string   `json:"username"`
	Name            string   `json:"name"`
	Company         string   `json:"company"`
	Blog            string   `json:"blog"`
	Location        string   `json:"location"`
	Email           string   `json:"email"`
	Bio             string   `json:"bio"`
	Followers       int      `json:"followers"`
	Organizations   []string `json:"organizations"`
	Readme          string   `json:"readme"`
	Commits         string   `json:"commits"`
	Score           float64  `json:"score"`
	PossibleNation  string   `json:"possible_nation"`
	ConfidenceLevel int      `json:"confidence_level"`
}

type Repo struct {
	Id               string `json:"id"`
	User_id          string `json:"user_id"`
	Repo             string `json:"repo"`
	Star             int    `json:"star"`
	Fork             int    `json:"fork"`
	Dependent        int    `json:"dependent"`
	Commits          int    `json:"commits"`
	CommitsTotal     int    `json:"commits_total"`
	Issue            int    `json:"issue"`
	IssueTotal       int    `json:"issue_total"`
	PullRequest      int    `json:"pull_request"`
	PullRequestTotal int    `json:"pull_request_total"`
	CodeReview       int    `json:"code_review"`
	CodeReviewTotal  int    `json:"code_review_total"`
	LineChange       int    `json:"line_change"`
	LineChangeTotal  int    `json:"line_change_total"`
}

func NewQueryUserReposFromHTTP(r *http.Request) *QueryUserReposRequest {
	req := NewQueryUserReposRequest()
	// query string
	qs := r.URL.Query()
	pss := qs.Get("page_size")
	if pss != "" {
		req.PageSize, _ = strconv.Atoi(pss)
	}

	pns := qs.Get("page_number")
	if pns != "" {
		req.PageNumber, _ = strconv.Atoi(pns)
	}

	req.Keywords = qs.Get("kws")
	return req
}

func NewQueryUserReposRequest() *QueryUserReposRequest {
	return &QueryUserReposRequest{
		PageSize:   20,
		PageNumber: 1,
	}
}

type QueryUserReposRequest struct {
	PageSize   int    `json:"page_size"`
	PageNumber int    `json:"page_number"`
	Keywords   string `json:"kws"`
}

func (q *QueryUserReposRequest) OffSet() int64 {
	return int64((q.PageNumber - 1) * q.PageSize)
}

func (q *QueryUserReposRequest) GetPageSize() uint {
	return uint(q.PageSize)
}
