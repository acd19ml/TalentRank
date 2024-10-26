package main

import (
	"context"
	"fmt"
	"log"

	"github.com/acd19ml/TalentRank/utils/git"

	"github.com/acd19ml/TalentRank/utils"
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
