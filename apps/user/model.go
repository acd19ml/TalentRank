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

func (s UserReposSet) Add(developers *UserRepos) {
	s.Developers = append(s.Developers, developers)
}

func (s UserRepos) AddRepos(repo *Repo) {
	s.Repos = append(s.Repos, repo)
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

func (u *User) Validate() error {
	return validate.Struct(u)
}

// char(36)
func (u *User) InjectDefault() {
	if u.Id == "" {
		u.Id = uuid.New().String()
	}
}

func NewUserSet() *UserSet {
	return &UserSet{
		Users: []*User{},
	}
}
func (s *UserSet) Add(user *User) {
	s.Users = append(s.Users, user)
}

type UserSet struct {
	Total int     `json:"total"`
	Users []*User `json:"users"`
}

func NewUser() *User {
	return &User{}
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
	ConfidenceLevel string   `json:"confidence_level"`
	Rankno          string   `json:"rankno"`
}

func (r *Repo) Validate() error {
	return validate.Struct(r)
}

// char(36)
func (r *Repo) InjectDefault() {
	if r.Id == "" {
		r.Id = uuid.New().String()
	}
}

type UserResponceByLLM struct {
	PossibleNation  string `json:"possible_nation"`
	ConfidenceLevel string `json:"confidence_level"`
}

func NewRepo() *Repo {
	return &Repo{}
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

func NewCreateUserReposRequest() *CreateUserReposRequest {
	return &CreateUserReposRequest{}
}

type CreateUserReposRequest struct {
	Username string `json:"username"`
}

func NewQueryUserFromHTTP(r *http.Request) *QueryUserRequest {
	req := NewQueryUserRequest()
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

	req.Location = qs.Get("location")
	return req
}

func NewQueryUserRequest() *QueryUserRequest {
	return &QueryUserRequest{
		PageSize:   20,
		PageNumber: 1,
	}
}

type QueryUserRequest struct {
	PageSize   int    `json:"page_size"`
	PageNumber int    `json:"page_number"`
	Location   string `json:"location"`
}

func (q *QueryUserRequest) OffSet() int64 {
	return int64((q.PageNumber - 1) * q.PageSize)
}

func (q *QueryUserRequest) GetPageSize() uint {
	return uint(q.PageSize)
}

func NewDescribeUserReposRequestFromHTTP(r *http.Request) *DescribeUserReposRequest {
	req := NewDescribeUserReposRequest()
	req.Username = r.URL.Query().Get("username")
	return req
}

func NewDescribeUserReposRequest() *DescribeUserReposRequest {
	return &DescribeUserReposRequest{}
}

type DescribeUserReposRequest struct {
	Username string `json:"username"`
}
