package crit

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
	"github.com/spf13/viper"
)

// OptReviewRepo is the configuration/flag string for the review repo path.
const OptReviewRepo = "review-repo"

// ErrNoConfig is the error for no configuration file found for the review repo.
var ErrNoConfig = errors.New("No configuration file for review repos")

// ReviewRepo represents the user's open reviews.
type ReviewRepo struct {
	BaseDir string
	Specs   []ReviewSpec
}

func (rr *ReviewRepo) Save() error {
	jpath := filepath.Join(rr.BaseDir, "review_repo.json")
	f, err := os.Create(jpath)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(rr)
	if err != nil {
		return err
	}
	return f.Sync()
}

// InitReviewRepo creates a new ReviewRepo at path.
func InitReviewRepo(path string) error {
	err := os.MkdirAll(path, 0700)
	if err != nil {
		return err
	}

	rr := &ReviewRepo{
		BaseDir: path,
		Specs:   make([]ReviewSpec, 0),
	}
	return rr.Save()
}

// OpenReviewRepo opens a ReviewRepo at path.
func OpenReviewRepo(path string) (*ReviewRepo, error) {
	jpath := filepath.Join(path, "review_repo.json")
	_, err := os.Stat(jpath)
	if err != nil {
		if err == os.ErrNotExist {
			return nil, ErrNoConfig
		}
		return nil, err
	}
	f, err := os.Open(jpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var rr ReviewRepo
	err = json.NewDecoder(f).Decode(&rr)
	if err != nil {
		return nil, err
	}
	return &rr, nil
}

func GetReviewRepoDir() (out string) {
	out = viper.GetString(OptReviewRepo)
	if out != "" {
		return
	}
	out = configdir.LocalConfig("crit")
	return
}
