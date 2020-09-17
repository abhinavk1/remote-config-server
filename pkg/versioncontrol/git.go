package versioncontrol

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"sync"
)

type GitClientConfig struct {
	Url        string
	Path       string
	PrivateKey []byte
}

type GitClient struct {
	repo  *git.Repository
	path  string
	mutex sync.Mutex
}

func NewGitClient(config *GitClientConfig) (*GitClient, error) {

	if config == nil {
		return nil, errors.New("no git configuration provided")
	}

	if err := clearWorkingDir(config.Path); err != nil {
		return nil, err
	}

	options := &git.CloneOptions{
		URL:      config.Url,
		Progress: os.Stdout,
	}

	if config.PrivateKey != nil {
		signer, err := ssh.ParsePrivateKey(config.PrivateKey)
		if err != nil {
			return nil, err
		}

		options.Auth = &gitssh.PublicKeys{
			User:   "git",
			Signer: signer,
		}
	}

	repo, err := git.PlainClone(config.Path, false, options)
	if err != nil {
		return nil, err
	}

	return &GitClient{
		repo:  repo,
		path:  config.Path,
		mutex: sync.Mutex{},
	}, nil
}

func (client *GitClient) Checkout(branch string) error {

	err := client.repo.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		fmt.Println(err)
	}

	workTree, err := client.repo.Worktree()
	if err != nil {
		return err
	}

	client.mutex.Lock()
	err = workTree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		Force:  true,
	})
	client.mutex.Unlock()

	return err
}

func (client *GitClient) Pull() error {

	workTree, err := client.repo.Worktree()
	if err != nil {
		return err
	}

	client.mutex.Lock()
	// Pull the latest changes from the origin remote and merge into the current branch
	err = workTree.Pull(&git.PullOptions{
		RemoteName: "origin",
	})
	if err != nil {
		client.mutex.Unlock()
		return err
	}

	// Print the latest commit that was just pulled
	ref, err := client.repo.Head()
	if err != nil {
		client.mutex.Unlock()
		return err
	}

	commit, err := client.repo.CommitObject(ref.Hash())
	if err != nil {
		client.mutex.Unlock()
		return err
	}

	client.mutex.Unlock()
	log.Printf("latest commit -> %v", commit)
	return nil
}
