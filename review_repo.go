package crit

import (
	"fmt"

	"github.com/spf13/viper"
)

type ReviewRepo struct {
	BaseDir string
	Specs   []ReviewSpec
}

func InitReviewRepo(path string) {

}

func OpenReviewRepo() {
	rr := viper.GetString("review-repo")
	fmt.Println(rr)
}
