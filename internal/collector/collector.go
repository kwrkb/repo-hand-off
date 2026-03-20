package collector

import (
	"errors"
	"time"
)

// Snapshot holds all collected project state.
type Snapshot struct {
	Timestamp  time.Time
	Git        GitInfo
	Files      ProjectFiles
	DirTree    string
	RecentLogs []string
}

// Collect gathers a full project snapshot from the given directory.
func Collect(dir string) (*Snapshot, error) {
	git, err := CollectGit(dir)
	if err != nil {
		if !errors.Is(err, ErrNotGitRepo) {
			return nil, err
		}
		git = &GitInfo{} // non-fatal: not a Git repo
	}

	files, err := CollectFiles(dir)
	if err != nil {
		return nil, err
	}

	tree, err := BuildDirTree(dir, 3)
	if err != nil {
		return nil, err
	}

	var logs []string
	if git.Branch != "" {
		logs, _ = RecentCommits(dir, 10) // non-fatal: repo may have no commits
	}

	return &Snapshot{
		Timestamp:  time.Now(),
		Git:        *git,
		Files:      *files,
		DirTree:    tree,
		RecentLogs: logs,
	}, nil
}
