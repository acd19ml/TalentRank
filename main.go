package main

import (
	"fmt"

	"github.com/acd19ml/TalentRank/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
