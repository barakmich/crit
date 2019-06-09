package main

import (
	"github.com/barakmich/crit"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a review repository",
	Run: func(cmd *cobra.Command, args []string) {

		path := crit.GetReviewRepoDir()
		err := crit.InitReviewRepo(path)
		if err != nil {
			fatal(err)
		}
	},
}
