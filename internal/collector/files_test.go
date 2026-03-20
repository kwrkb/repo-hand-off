package collector

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCollectFiles(t *testing.T) {
	dir := t.TempDir()

	// Create test files
	os.WriteFile(filepath.Join(dir, "VISION.md"), []byte("# Vision"), 0644)
	os.WriteFile(filepath.Join(dir, "PLAN.md"), []byte("# Plan"), 0644)

	files, err := CollectFiles(dir, nil)
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

func TestCollectFilesExtra(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "NOTES.md"), []byte("# Notes"), 0644)

	files, err := CollectFiles(dir, []string{"NOTES.md", "MISSING.md"})
	if err != nil {
		t.Fatalf("CollectFiles failed: %v", err)
	}

	if files.Extra["NOTES.md"] != "# Notes" {
		t.Errorf("Extra[NOTES.md] = %q, want %q", files.Extra["NOTES.md"], "# Notes")
	}
	if _, ok := files.Extra["MISSING.md"]; ok {
		t.Error("MISSING.md should not be in Extra")
	}
}

func TestBuildDirTree(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "src"), 0755)
	os.WriteFile(filepath.Join(dir, "src", "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("# README"), 0644)

	tree, err := BuildDirTree(dir, 2, nil)
	if err != nil {
		t.Fatalf("BuildDirTree failed: %v", err)
	}

	if tree == "" {
		t.Error("tree should not be empty")
	}
}

func TestBuildDirTreeExclude(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "logs"), 0755)
	os.WriteFile(filepath.Join(dir, "logs", "app.log"), []byte("log"), 0644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0644)

	tree, err := BuildDirTree(dir, 3, []string{"logs"})
	if err != nil {
		t.Fatalf("BuildDirTree failed: %v", err)
	}

	if strings.Contains(tree, "logs") {
		t.Error("tree should not contain excluded directory 'logs'")
	}
	if !strings.Contains(tree, "main.go") {
		t.Error("tree should contain main.go")
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
		if got := shouldSkip(tt.name, nil); got != tt.want {
			t.Errorf("shouldSkip(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestShouldSkipWithExclude(t *testing.T) {
	if !shouldSkip("test.tmp", []string{"*.tmp"}) {
		t.Error("shouldSkip should match *.tmp pattern")
	}
	if shouldSkip("test.go", []string{"*.tmp"}) {
		t.Error("shouldSkip should not match *.tmp for .go file")
	}
}
