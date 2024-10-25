package main

import (
	"acd19ml/TalentRank/utils/git"
	"fmt"
	"log"
)

const (
	username = "alibaba"
)

func main() {
	service := git.NewGitClient()
	num, err := service.GetFollowers(username)
	if err != nil {
		log.Fatalf("Error getting Followers: %v", err)
	}
	fmt.Printf("Followers for user %s: %d\n", username, num)

	totalStars, err := service.GetTotalStars(username)
	if err != nil {
		log.Fatalf("Error getting total stars: %v", err)
	}

	fmt.Printf("Total stars for user %s: %d\n", username, totalStars)

	totalForks, err := service.GetTotalForks(username)
	if err != nil {
		log.Fatalf("Error getting total forks: %v", err)
	}
	fmt.Printf("Total forks for user %s: %d\n", username, totalForks)
}
