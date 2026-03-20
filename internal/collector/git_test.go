package collector

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func setupGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("setup %v failed: %v\n%s", args, err, out)
		}
	}

	// Create and commit a file
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Test"), 0644)
	for _, args := range [][]string{
		{"git", "add", "."},
		{"git", "commit", "-m", "initial"},
	} {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("setup %v failed: %v\n%s", args, err, out)
		}
	}

	return dir
}

func TestCollectGit(t *testing.T) {
	dir := setupGitRepo(t)

	info, err := CollectGit(dir)
	if err != nil {
		t.Fatalf("CollectGit failed: %v", err)
	}

	if info.Branch == "" {
		t.Error("Branch should not be empty")
	}
	if info.ShortHash == "" {
		t.Error("ShortHash should not be empty")
	}
	if info.HasChanges {
		t.Error("HasChanges should be false for clean repo")
	}
}

func TestCollectGitWithChanges(t *testing.T) {
	dir := setupGitRepo(t)

	// Make a change
	os.WriteFile(filepath.Join(dir, "new.txt"), []byte("new file"), 0644)

	info, err := CollectGit(dir)
	if err != nil {
		t.Fatalf("CollectGit failed: %v", err)
	}

	if !info.HasChanges {
		t.Error("HasChanges should be true after adding a file")
	}
}

func TestCollectGitNonGitDir(t *testing.T) {
	dir := t.TempDir() // plain directory, not a git repo

	_, err := CollectGit(dir)
	if err == nil {
		t.Fatal("CollectGit should return error for non-git directory")
	}
	if !errors.Is(err, ErrNotGitRepo) {
		t.Errorf("expected ErrNotGitRepo, got: %v", err)
	}
}

func TestRecentCommits(t *testing.T) {
	dir := setupGitRepo(t)

	logs, err := RecentCommits(dir, 5)
	if err != nil {
		t.Fatalf("RecentCommits failed: %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("expected 1 commit, got %d", len(logs))
	}
}
