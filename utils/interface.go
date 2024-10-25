package utils

type Services interface {
	GetFollowers(username string) (int, error)
	GetTotalStars(username string) (int, error)
	GetTotalForks(username string) (int, error)
}
