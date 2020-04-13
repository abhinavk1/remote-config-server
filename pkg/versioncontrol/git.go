package versioncontrol

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
)

func Clone(url, path string, privateKey []byte) (*git.Repository, error) {
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		fmt.Println(err)
	}

	auth := &gitssh.PublicKeys{User: "git", Signer: signer}

	return git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
		Auth:     auth,
	})
}

func HasDiff(repo *git.Repository, fileName string) (*bool, error) {
	commits, err := repo.Log(&git.LogOptions{
		FileName: &fileName,
	})

	if err != nil {
		return nil, err
	}
	defer commits.Close()

	var hasDiff bool
	var retErr error
	var prevCommit *object.Commit
	var prevTree *object.Tree

	for {
		commit, err := commits.Next()
		if err != nil {
			break
		}
		currentTree, err := commit.Tree()
		if err != nil {
			retErr = err
			break
		}

		if prevCommit == nil {
			prevCommit = commit
			prevTree = currentTree
			continue
		}

		changes, err := currentTree.Diff(prevTree)
		if err != nil {
			retErr = err
			break
		}

		for _, c := range changes {
			if c.To.Name == fileName {
				hasDiff = true
				break
			}
		}

		prevCommit = commit
		prevTree = currentTree
	}

	if retErr != nil {
		return nil, retErr
	}

	return &hasDiff, nil
}

func Pull(repo *git.Repository) error {
	// Get the working directory for the repository
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Pull the latest changes from the origin remote and merge into the current branch
	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		return err
	}

	// Print the latest commit that was just pulled
	ref, err := repo.Head()
	if err != nil {
		return err
	}

	commit, err := repo.CommitObject(ref.Hash())
	log.Printf("latest commit -> %v", commit)
	return err
}
