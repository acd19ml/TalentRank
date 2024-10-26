package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var (
	validate = validator.New()
)

func NewUser() *User {
	return &User{}
}

type User struct {
	Id            string   `json:"id"`
	Username      string   `json:"username"`
	Name          string   `json:"name"`
	Company       string   `json:"company"`
	Blog          string   `json:"blog"`
	Location      string   `json:"location"`
	Email         string   `json:"email"`
	Bio           string   `json:"bio"`
	TotalStar     int      `json:"total_star"`
	TotalFork     int      `json:"total_fork"`
	Followers     int      `json:"followers"`
	Dependents    int      `json:"dependents"`
	Organizations []string `json:"organizations"`
	Readme        string   `json:"readme"`
	Commits       string   `json:"commits"`
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
