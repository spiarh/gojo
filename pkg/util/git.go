package util

import (
	"github.com/go-git/go-git/v5"
)

func OpenGitRepo(path string) (*git.Repository, error) {
	opt := &git.PlainOpenOptions{DetectDotGit: true}
	return git.PlainOpenWithOptions(path, opt)
}

func GetGitHeadHash(path string) (string, error) {
	repo, err := OpenGitRepo(path)
	if err != nil {
		return "", err
	}

	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	return head.Hash().String(), nil
}

// GitIsFileClean returns true if the files is in Unmodified status.
func GitIsFileClean(status git.Status, relPath string) bool {
	fStatus := status.File(relPath)
	if fStatus.Worktree != git.Unmodified || fStatus.Staging != git.Unmodified {
		return false
	}
	return true
}

// GitIsFileUnmodifiedWorktree returns true if the files is modified in Worktree.
func GitIsFileUnmodifiedWorktree(status git.Status, relPath string) bool {
	fStatus := status.File(relPath)
	return fStatus.Worktree == git.Unmodified
}

// GitIsFileUnmodifiedStaging returns true if the files is modified in Staging.
func GitIsFileUnmodifiedStaging(status git.Status, relPath string) bool {
	fStatus := status.File(relPath)
	return fStatus.Staging == git.Unmodified
}
