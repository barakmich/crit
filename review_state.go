package crit

import (
	"errors"
	"fmt"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const defaultBaseWalkDepth = 500

type reviewState struct {
	review        *Review
	repo          *git.Repository
	baseHash      *plumbing.Hash
	baseHashes    map[plumbing.Hash]bool
	reviewHash    *plumbing.Hash
	rootHash      *plumbing.Hash
	reviewCommits []reviewCommit
}

type reviewCommit struct {
	commit *object.Commit
	base   *plumbing.Hash
}

func (rc reviewCommit) mustGetBaseCommit(rs *reviewState) *object.Commit {
	c, err := rs.repo.CommitObject(*rc.base)
	if err != nil {
		panic(err)
	}
	return c
}

func (rc reviewCommit) String() string {
	return fmt.Sprintf("commit:%s base:%s", rc.commit.ID().String(), rc.base.String())
}

func (r *Review) startReviewState() error {
	rs, err := newReviewState(r)
	r.state = rs
	rs.debugPrint()
	return err

}

func (rs *reviewState) debugPrint() error {
	fmt.Printf("base: %s\n", *rs.baseHash)
	fmt.Printf("review: %s\n", *rs.reviewHash)
	fmt.Printf("alles: %v\n", rs.reviewCommits)
	for _, x := range rs.reviewCommits {
		base := x.mustGetBaseCommit(rs)
		p, err := base.Patch(x.commit)
		if err != nil {
			return err
		}
		for _, fp := range p.FilePatches() {
			for _, chonk := range fp.Chunks() {
				fmt.Println(chonk.Type(), chonk.Content())
			}
		}
	}

	return nil
}

func newReviewState(r *Review) (*reviewState, error) {
	var err error
	rs := &reviewState{
		review: r,
		repo:   r.repo,
	}
	rs.baseHash, err = r.repo.ResolveRevision(plumbing.Revision(fmt.Sprintf("remotes/base/%s", r.Spec.BaseBranch)))
	if err != nil {
		return nil, err
	}
	rs.reviewHash, err = r.repo.ResolveRevision(plumbing.Revision(fmt.Sprintf("remotes/review/%s", r.Spec.ReviewBranch)))
	if err != nil {
		return nil, err
	}
	err = rs.loadReviewCommits()
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (rs *reviewState) loadReviewCommits() error {
	// Walk the base commits back a ways to find the
	// "master" commit history.
	//
	// This may not be enough. Hopefully we can find the root in
	// 500 or so commits. TODO(barakmich): Config option when starting
	// a review for the number of commits.

	base, err := rs.repo.CommitObject(*rs.baseHash)
	if err != nil {
		return err
	}
	rs.baseHashes = make(map[plumbing.Hash]bool)
	queue := make([]*object.Commit, 0)
	queue = append(queue, base)
	for i := 0; i < defaultBaseWalkDepth; i++ {
		if len(queue) == 0 {
			break
		}
		c := queue[0]
		queue = queue[1:]
		rs.baseHashes[c.ID()] = true
		err = c.Parents().ForEach(func(sub *object.Commit) error {
			queue = append(queue, sub)
			return nil
		})
		if err != nil {
			return err
		}
	}

	rev, err := rs.repo.CommitObject(*rs.reviewHash)
	if err != nil {
		return err
	}
	baseHash := &plumbing.Hash{}
	queue = make([]*object.Commit, 0)
	queue = append(queue, rev)
	for len(queue) != 0 {
		c := queue[0]
		queue = queue[1:]
		rc := reviewCommit{
			commit: c,
			base:   baseHash,
		}
		rs.reviewCommits = append(rs.reviewCommits, rc)
		if c.NumParents() == 0 {
			return errors.New("Ran off the edge of the world")
		}
		err = c.Parents().ForEach(func(sub *object.Commit) error {
			if rs.baseHashes[sub.ID()] {
				*baseHash = sub.ID()
				baseHash = &plumbing.Hash{}
			} else {
				queue = append(queue, sub)
			}
			return nil
		})
		if err != nil {
			return err
		}

	}
	return nil
}
