package collector

import (
	"fmt"
	"os/exec"
	"strings"
)

// GitInfo holds Git repository state.
type GitInfo struct {
	Branch      string
	ShortHash   string
	RemoteURL   string
	HasChanges  bool
	DiffSummary string
}

// CollectGit gathers Git state from the given directory.
func CollectGit(dir string) (*GitInfo, error) {
	info := &GitInfo{}

	branch, err := gitCmd(dir, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return nil, err
	}
	info.Branch = branch

	hash, err := gitCmd(dir, "rev-parse", "--short", "HEAD")
	if err != nil {
		return nil, err
	}
	info.ShortHash = hash

	remote, _ := gitCmd(dir, "remote", "get-url", "origin")
	info.RemoteURL = remote

	status, _ := gitCmd(dir, "status", "--porcelain")
	info.HasChanges = status != ""

	if info.HasChanges {
		diff, _ := gitCmd(dir, "diff", "--stat")
		staged, _ := gitCmd(dir, "diff", "--cached", "--stat")
		parts := []string{}
		if diff != "" {
			parts = append(parts, diff)
		}
		if staged != "" {
			parts = append(parts, staged)
		}
		info.DiffSummary = strings.Join(parts, "\n")
	}

	return info, nil
}

// RecentCommits returns the last n commit messages.
func RecentCommits(dir string, n int) ([]string, error) {
	out, err := gitCmd(dir, "log", "--oneline", "-n", fmt.Sprintf("%d", n))
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

func gitCmd(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

