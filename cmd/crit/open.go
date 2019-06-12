package main

import (
	"github.com/barakmich/crit"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(openCmd)
}

var openCmd = &cobra.Command{
	Use:   "open NAME",
	Args:  cobra.MinimumNArgs(1),
	Short: "open a review",
	Run: func(cmd *cobra.Command, args []string) {
		path := crit.GetReviewRepoDir()
		rr, err := crit.OpenReviewRepo(path)
		if err != nil {
			fatal(err)
		}
		r, err := rr.OpenReview(args[0])
		if err != nil {
			fatal(err)
		}
		err = crit.ReviewUIMain(r)
		if err != nil {
			fatal(err)
		}
	},
}
