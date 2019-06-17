package crit

import (
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type ReviewOp struct {
	Kind     ReviewOpKind
	File     string
	Revision plumbing.Hash
}

type ReviewOpKind string

const (
	OpReadFile   ReviewOpKind = "read"
	OpUnreadFile              = "unread"
)
