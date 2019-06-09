package main

import (
	"fmt"

	"github.com/barakmich/crit"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list open reviews",
	Run: func(cmd *cobra.Command, args []string) {
		//TODO(barakmich): Better tabulate the results
		path := crit.GetReviewRepoDir()
		rr, err := crit.OpenReviewRepo(path)
		if err != nil {
			fatal(err)
		}
		fmt.Println("Name\tBase URL\tBase Branch\tReview Branch")
		fmt.Println("-----------")
		for _, x := range rr.Specs {
			fmt.Printf("%s\t%s\t%s\t%s\n", x.Name, x.BaseURL, x.BaseBranch, x.ReviewBranch)
		}
	},
}
