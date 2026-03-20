package collector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCollectGitRepo(t *testing.T) {
	dir := setupGitRepo(t)

	// Add project files
	os.WriteFile(filepath.Join(dir, "VISION.md"), []byte("# Vision\nTest vision."), 0644)
	os.WriteFile(filepath.Join(dir, "PLAN.md"), []byte("# Plan\nTest plan."), 0644)

	snap, err := Collect(dir, CollectOptions{})
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if snap.Git.Branch == "" {
		t.Error("Git.Branch should not be empty in git repo")
	}
	if snap.Files.Vision != "# Vision\nTest vision." {
		t.Errorf("Vision = %q", snap.Files.Vision)
	}
	if snap.Files.Plan != "# Plan\nTest plan." {
		t.Errorf("Plan = %q", snap.Files.Plan)
	}
	if snap.DirTree == "" {
		t.Error("DirTree should not be empty")
	}
	if len(snap.RecentLogs) == 0 {
		t.Error("RecentLogs should not be empty in git repo with commits")
	}
	if snap.Timestamp.IsZero() {
		t.Error("Timestamp should be set")
	}
}

func TestCollectNonGitDir(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "VISION.md"), []byte("# Vision"), 0644)

	snap, err := Collect(dir, CollectOptions{})
	if err != nil {
		t.Fatalf("Collect should succeed in non-git dir, got: %v", err)
	}

	if snap.Git.Branch != "" {
		t.Error("Git.Branch should be empty for non-git dir")
	}
	if snap.Files.Vision != "# Vision" {
		t.Errorf("Vision = %q", snap.Files.Vision)
	}
	if len(snap.RecentLogs) != 0 {
		t.Error("RecentLogs should be empty for non-git dir")
	}
}

func TestCollectWithExtraFiles(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "NOTES.md"), []byte("# Notes"), 0644)

	snap, err := Collect(dir, CollectOptions{
		ExtraFiles: []string{"NOTES.md", "MISSING.md"},
	})
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if snap.Files.Extra["NOTES.md"] != "# Notes" {
		t.Errorf("Extra[NOTES.md] = %q", snap.Files.Extra["NOTES.md"])
	}
	if _, ok := snap.Files.Extra["MISSING.md"]; ok {
		t.Error("MISSING.md should not be in Extra")
	}
}

func TestCollectWithExclude(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "logs"), 0755)
	os.WriteFile(filepath.Join(dir, "logs", "app.log"), []byte("log"), 0644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0644)

	snap, err := Collect(dir, CollectOptions{
		Exclude: []string{"logs"},
	})
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if snap.DirTree == "" {
		t.Error("DirTree should not be empty")
	}
}

func TestCollectWithDepth(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "a", "b", "c", "d"), 0755)
	os.WriteFile(filepath.Join(dir, "a", "b", "c", "d", "deep.txt"), []byte("deep"), 0644)

	snap, err := Collect(dir, CollectOptions{Depth: 1})
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	// Depth 1 should only show the immediate children
	if snap.DirTree == "" {
		t.Error("DirTree should not be empty")
	}
}

func TestCollectEmptyDir(t *testing.T) {
	dir := t.TempDir()

	snap, err := Collect(dir, CollectOptions{})
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if snap.Files.Vision != "" {
		t.Error("Vision should be empty for empty dir")
	}
	if snap.DirTree == "" {
		t.Error("DirTree should show at least the root directory")
	}
}
