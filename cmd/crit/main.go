package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "crit",
	Short: "crit is a code review tool for git repositories",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Woo")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
