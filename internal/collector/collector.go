package collector

import "time"

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
		return nil, err
	}

	files, err := CollectFiles(dir)
	if err != nil {
		return nil, err
	}

	tree, err := BuildDirTree(dir, 3)
	if err != nil {
		return nil, err
	}

	logs, err := RecentCommits(dir, 10)
	if err != nil {
		// non-fatal: repo may have no commits
		logs = nil
	}

	return &Snapshot{
		Timestamp:  time.Now(),
		Git:        *git,
		Files:      *files,
		DirTree:    tree,
		RecentLogs: logs,
	}, nil
}
