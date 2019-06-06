package crit

import (
	"gopkg.in/src-d/go-git.v4"
)

type Review struct {
	Name string
	repo *git.Repository
}

type ReviewSpec struct {
	Name         string
	BaseBranch   string
	ReviewBranch string
	OriginURL    string
}

func CreateReview(spec ReviewSpec) (*Review, error) {
	panic("unimplemented")
}

func OpenReview(name string) (*Review, error) {
	panic("unimplemented")
}
