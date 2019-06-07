package crit

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4"
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
		Spec: spec,
	}
	err := os.MkdirAll(r.DataDir(), 0700)
	if err != nil {
		return err
	}
	err = os.MkdirAll(r.RepoDir(), 0700)
	if err != nil {
		return err
	}
	// TODO Clone
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
