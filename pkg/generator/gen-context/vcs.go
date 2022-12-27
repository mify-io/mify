package gencontext

import (
	"fmt"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/index"
)

type VcsIntegration struct {
	git       *git.Repository
	index     *index.Index
	workTree  *git.Worktree
	status    *git.Status
	initMutex sync.Mutex
}

func initVcsIntegration(path string) (*VcsIntegration, error) {
	rep, err := git.PlainOpen(path)
	if err != nil && err != git.ErrRepositoryNotExists {
		return nil, fmt.Errorf("can't load git repository: %w", err)
	}

	if rep == nil {
		return &VcsIntegration{
			git:      nil,
			workTree: nil,
		}, nil
	}

	workTree, err := rep.Worktree()
	if err != nil {
		return nil, err
	}

	return &VcsIntegration{
		git:      rep,
		workTree: workTree,
	}, nil
}

func (vcs *VcsIntegration) FileHasUncommitedChanges(relPath string) (bool, error) {
	if vcs.status == nil {
		vcs.initMutex.Lock()
		defer vcs.initMutex.Unlock()

		index, err := vcs.git.Storer.Index()
		if err != nil {
			return false, err
		}
		vcs.index = index

		status, err := vcs.workTree.Status()
		if err != nil {
			return false, err
		}
		vcs.status = &status
	}

	_, err := vcs.index.Entry(relPath)
	if err == index.ErrEntryNotFound {
		return true, nil
	}
	if err != nil {
		return false, err
	}

	worktreeStatus := vcs.status.File(relPath).Worktree
	if worktreeStatus != git.Unmodified && worktreeStatus != git.Untracked {
		return true, nil
	}

	stagingStatus := vcs.status.File(relPath).Staging
	if stagingStatus != git.Unmodified && stagingStatus != git.Untracked {
		return true, nil
	}

	return false, nil
}
