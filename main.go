package main

import (
	"acd19ml/TalentRank/utils"
	"acd19ml/TalentRank/utils/git"
	"context"
	"fmt"
	"log"
)

func main() {
	service := git.NewGitClient()
	num, err := service.GetFollowers(context.Background(), utils.Username)
	if err != nil {
		log.Fatalf("Error getting Followers: %v", err)
	}
	fmt.Printf("Followers for user %s: %d\n", utils.Username, num)
}
