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
	WorkDir    string   // absolute path to the project directory
	TodoCount  int      // number of TODO/FIXME comments in source files
	CIFiles    []string // detected CI config file paths (relative)
}

// CollectOptions configures the Collect function.
type CollectOptions struct {
	ExtraFiles []string
	Exclude    []string
	Depth      int // 0 → default 3
}

// Collect gathers a full project snapshot from the given directory.
func Collect(dir string, opts CollectOptions) (*Snapshot, error) {
	git, err := CollectGit(dir)
	if err != nil {
		if !errors.Is(err, ErrNotGitRepo) {
			return nil, err
		}
		git = &GitInfo{} // non-fatal: not a Git repo
	}

	files, err := CollectFiles(dir, opts.ExtraFiles)
	if err != nil {
		return nil, err
	}

	tree, err := BuildDirTree(dir, opts.Depth, opts.Exclude)
	if err != nil {
		return nil, err
	}

	var logs []string
	if git.Branch != "" {
		logs, _ = RecentCommits(dir, 10) // non-fatal: repo may have no commits
	}

	todoCount, err := CountTodos(dir, opts.Exclude)
	if err != nil {
		return nil, err
	}

	ciFiles := DetectCIFiles(dir)

	return &Snapshot{
		Timestamp:  time.Now(),
		Git:        *git,
		Files:      *files,
		DirTree:    tree,
		RecentLogs: logs,
		WorkDir:    dir,
		TodoCount:  todoCount,
		CIFiles:    ciFiles,
	}, nil
}
