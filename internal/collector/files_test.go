package collector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCollectFiles(t *testing.T) {
	dir := t.TempDir()

	// Create test files
	os.WriteFile(filepath.Join(dir, "VISION.md"), []byte("# Vision"), 0644)
	os.WriteFile(filepath.Join(dir, "PLAN.md"), []byte("# Plan"), 0644)

	files, err := CollectFiles(dir)
	if err != nil {
		t.Fatalf("CollectFiles failed: %v", err)
	}

	if files.Vision != "# Vision" {
		t.Errorf("Vision = %q, want %q", files.Vision, "# Vision")
	}
	if files.Plan != "# Plan" {
		t.Errorf("Plan = %q, want %q", files.Plan, "# Plan")
	}
	if files.Lessons != "" {
		t.Errorf("Lessons = %q, want empty", files.Lessons)
	}
	if files.Readme != "" {
		t.Errorf("Readme = %q, want empty", files.Readme)
	}
}

func TestBuildDirTree(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "src"), 0755)
	os.WriteFile(filepath.Join(dir, "src", "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("# README"), 0644)

	tree, err := BuildDirTree(dir, 2)
	if err != nil {
		t.Fatalf("BuildDirTree failed: %v", err)
	}

	if tree == "" {
		t.Error("tree should not be empty")
	}
}

func TestShouldSkip(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{".git", true},
		{"node_modules", true},
		{"src", false},
		{"main.go", false},
		{".DS_Store", true},
	}
	for _, tt := range tests {
		if got := shouldSkip(tt.name); got != tt.want {
			t.Errorf("shouldSkip(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}
