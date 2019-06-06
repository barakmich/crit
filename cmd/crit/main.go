package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/barakmich/crit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "crit",
	Short: "crit is a code review tool for git repositories",
	Run: func(cmd *cobra.Command, args []string) {
		crit.OpenReviewRepo()
	},
}

func main() {
	viper.SetEnvPrefix("crit")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
