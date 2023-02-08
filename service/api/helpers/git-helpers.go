package helpers

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/notherealmarco/WhaleDeployer/service/structures"
)

func CloneOrPullSSH(p *structures.Project, key transport.AuthMethod, progress *os.File) error {

	_, err := git.PlainClone(p.Path, false, &git.CloneOptions{
		Auth:          key,
		URL:           p.GitURL,
		Progress:      progress,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", p.GitBranch)),
	})

	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		// do a git pull

		progress.Write([]byte("\nGit: repository already exists, pulling...\n"))

		r, err := git.PlainOpen(p.Path)

		if err != nil {
			return fmt.Errorf("Git error opening repository: %w", err)
		}

		wt, err := r.Worktree()

		if err != nil {
			return fmt.Errorf("Git error setting worktree: %w", err)
		}

		err = wt.Pull(&git.PullOptions{
			RemoteName:    "origin",
			Auth:          key,
			Progress:      progress,
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", p.GitBranch)),
		})

		if err == git.NoErrAlreadyUpToDate {
			progress.Write([]byte("Git: " + err.Error() + "\n"))
			return nil
		}
	}

	if err != nil {
		progress.Write([]byte("Git error: " + err.Error() + "\n"))
	}

	return err
}

func CloneOrPull(p *structures.Project, progress *os.File) error {

	_, err := git.PlainClone(p.Path, false, &git.CloneOptions{
		URL:           p.GitURL,
		Progress:      progress,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", p.GitBranch)),
	})

	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		// do a git pull

		progress.Write([]byte("\nGit: repository already exists, pulling...\n"))

		r, err := git.PlainOpen(p.Path)

		if err != nil {
			return fmt.Errorf("Git error opening repository: %w", err)
		}

		wt, err := r.Worktree()

		if err != nil {
			return fmt.Errorf("Git error setting worktree: %w", err)
		}

		err = wt.Pull(&git.PullOptions{
			RemoteName:    "origin",
			Progress:      progress,
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", p.GitBranch)),
		})

		if err == git.NoErrAlreadyUpToDate {
			progress.Write([]byte("Git: " + err.Error() + "\n"))
			return nil
		}
	}

	if err != nil {
		progress.Write([]byte("Git error: " + err.Error() + "\n"))
	}

	return err
}
