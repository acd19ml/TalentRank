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
	names, err := service.GetOrganizations(context.Background(), utils.Username)
	if err != nil {
		log.Fatalf("Error getting name: %v", err)
	}
	for _, name := range names {
		fmt.Printf("Organization for user %s: %s\n", utils.Username, name)
	}
}
