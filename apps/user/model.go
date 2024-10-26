package user

type User struct {
	Username      string
	Name          string
	Company       string
	Blog          string
	Location      string
	Email         string
	Bio           string
	TotalStar     int
	TotalFork     int
	Followers     int
	Dependents    int
	Organizations []string
	Readme        string
	Commits       string
}
