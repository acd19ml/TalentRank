package cmd

import (
	"fmt"

	"github.com/acd19ml/TalentRank/version"
	"github.com/spf13/cobra"
)

var vers bool

// RootCmd is the root command for TalentRank
var RootCmd = &cobra.Command{
	Use:   "TalentRank",
	Short: "TalentRank is a demo project",
	Long:  "TalentRank is a demo project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if vers {
			fmt.Println(version.FullVersion())
			return nil
		}
		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&vers, "version", "v", false, "print TalentRank version")
}
