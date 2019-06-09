package crit

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
)

var ErrReviewNotExist = errors.New("review doesn't exist")

type Review struct {
	Spec       ReviewSpec
	repo       *git.Repository
	reviewRepo *ReviewRepo
}

type ReviewSpec struct {
	Name         string
	BaseBranch   string
	ReviewBranch string
	BaseURL      string
	ReviewURL    string
	CloneURL     string
}

func (r *Review) Save() error {
	jpath := filepath.Join(r.DataDir(), "review.json")
	f, err := os.Create(jpath)
	if err != nil {
		return err
	}

	defer f.Close()
	err = json.NewEncoder(f).Encode(r)
	if err != nil {
		return err
	}
	return f.Sync()
}

func (r *Review) RepoDir() string {
	return filepath.Join(r.reviewRepo.BaseDir, "repos", r.Spec.Name)
}

func (r *Review) DataDir() string {
	return reviewDataDir(r.reviewRepo, r.Spec.Name)
}

func reviewDataDir(rr *ReviewRepo, name string) string {
	return filepath.Join(rr.BaseDir, "data", name)
}

func (rr *ReviewRepo) CreateReview(spec ReviewSpec) error {
	if spec.CloneURL == "" {
		spec.CloneURL = spec.BaseURL
	}
	if spec.Name == "" {
		return errors.New("No name provided in ReviewSpec")
	}
	r := &Review{
		Spec:       spec,
		reviewRepo: rr,
	}
	err := os.MkdirAll(r.DataDir(), 0700)
	if err != nil {
		return err
	}
	err = os.MkdirAll(r.RepoDir(), 0700)
	if err != nil {
		return err
	}
	err = r.initClone()
	if err != nil {
		return err
	}
	err = r.fetchAndUpdate()
	if err != nil {
		return err
	}
	err = r.Save()
	if err != nil {
		return err
	}
	rr.Specs = append(rr.Specs, spec)
	err = rr.Save()
	if err != nil {
		return err
	}
	return nil
}

func (r *Review) fetchRemote(remotename string) error {
	remote, err := r.repo.Remote(remotename)
	if err != nil {
		return err
	}
	err = remote.Fetch(&git.FetchOptions{
		// TODO(barakmich): fetch only the branch needed
		RefSpecs: []config.RefSpec{
			config.RefSpec(fmt.Sprintf("+refs/heads/*:refs/remotes/%s/*", remotename)),
		},
		Force: true,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}
	return nil
}

func (r *Review) fetchAndUpdate() error {
	err := r.fetchRemote("base")
	if err != nil {
		return err
	}
	err = r.fetchRemote("review")
	if err != nil {
		return err
	}
	return nil
}

func (r *Review) initClone() error {
	repo, err := git.PlainClone(r.RepoDir(), true, &git.CloneOptions{
		URL:        r.Spec.CloneURL,
		NoCheckout: true,
		Tags:       git.NoTags,
	})
	if err != nil {
		return err
	}
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "base",
		URLs: []string{r.Spec.BaseURL},
	})
	if err != nil {
		return err
	}
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "review",
		URLs: []string{r.Spec.ReviewURL},
	})
	if err != nil {
		return err
	}
	r.repo = repo
	return nil
}

func (rr *ReviewRepo) OpenReview(name string) (*Review, error) {
	for _, spec := range rr.Specs {
		if name == spec.Name {
			return rr.openSpec(spec)
		}
	}
	return nil, ErrReviewNotExist
}

func (rr *ReviewRepo) openSpec(spec ReviewSpec) (*Review, error) {
	f, err := os.Open(filepath.Join(reviewDataDir(rr, spec.Name), "review.json"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var r Review
	err = json.NewDecoder(f).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
