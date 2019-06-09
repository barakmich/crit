package main

import (
	"github.com/barakmich/crit"
	"github.com/spf13/cobra"
)

var startBaseURL string
var startBaseBranch string
var startReviewURL string
var startReviewBranch string

func init() {
	startCmd.Flags().StringVar(&startBaseURL, "base-url", "", "Base URL")
	startCmd.Flags().StringVar(&startReviewURL, "review-url", "", "Review URL")
	startCmd.Flags().StringVar(&startBaseBranch, "base-branch", "master", "Base Branch")
	startCmd.Flags().StringVar(&startReviewBranch, "review-branch", "master", "Review Branch")
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start NAME",
	Args:  cobra.MinimumNArgs(1),
	Short: "start a new review",
	Run: func(cmd *cobra.Command, args []string) {
		path := crit.GetReviewRepoDir()
		rr, err := crit.OpenReviewRepo(path)
		if err != nil {
			fatal(err)
		}
		if startReviewURL == "" {
			startReviewURL = startBaseURL
		}
		spec := crit.ReviewSpec{
			Name:         args[0],
			BaseBranch:   startBaseBranch,
			BaseURL:      startBaseURL,
			ReviewBranch: startReviewBranch,
			ReviewURL:    startReviewURL,
		}
		err = rr.CreateReview(spec)
		if err != nil {
			fatal(err)
		}
	},
}
