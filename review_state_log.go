package crit

import (
	"fmt"

	"gopkg.in/src-d/go-git.v4/plumbing"
)

// TODO(barakmich): More expensive calls here can be cached. Initial
// implementations should run through the entire oplist (which is likely N ~ 100)
// but if they add up, we can attack that problem later.

func (rs *reviewState) markRead(file string, hash plumbing.Hash) error {
	if !rs.isRead(file, hash) {
		return fmt.Errorf("File %s@%s is already read", file, hash.String())
	}
	op := ReviewOp{
		Kind:     OpReadFile,
		File:     file,
		Revision: hash,
	}
	rs.review.ReviewOperations = append(rs.review.ReviewOperations, op)
	return rs.review.Save()
}

func (rs *reviewState) isRead(file string, hash plumbing.Hash) bool {
	isRead := false
	for _, x := range rs.review.ReviewOperations {
		if x.Kind == OpReadFile && !isRead {
			if hash == x.Revision && x.File == file {
				isRead = true
			}
		}
		if x.Kind == OpUnreadFile && isRead {
			if hash == x.Revision && x.File == file {
				isRead = false
			}
		}
	}
	return isRead
}
