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
	// README.md is now auto-extra; should not be in Extra if file doesn't exist
	if _, ok := files.Extra["README.md"]; ok {
		t.Error("README.md should not be in Extra when file doesn't exist")
	}
}

func TestCollectFilesAutoExtra(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("# README"), 0644)
	os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# CLAUDE"), 0644)
	os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte("# Agents"), 0644)

	files, err := CollectFiles(dir, nil)
	if err != nil {
		t.Fatalf("CollectFiles failed: %v", err)
	}

	if files.Extra["README.md"] != "# README" {
		t.Errorf("Extra[README.md] = %q, want %q", files.Extra["README.md"], "# README")
	}
	if files.Extra["CLAUDE.md"] != "# CLAUDE" {
		t.Errorf("Extra[CLAUDE.md] = %q, want %q", files.Extra["CLAUDE.md"], "# CLAUDE")
	}
	if files.Extra["AGENTS.md"] != "# Agents" {
		t.Errorf("Extra[AGENTS.md] = %q, want %q", files.Extra["AGENTS.md"], "# Agents")
	}
	// GEMINI.md doesn't exist, should not be in Extra
	if _, ok := files.Extra["GEMINI.md"]; ok {
		t.Error("GEMINI.md should not be in Extra when file doesn't exist")
	}
}

func TestCollectFilesAutoExtraDedupe(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("# README"), 0644)

	// Pass README.md as explicit extra too — should not duplicate
	files, err := CollectFiles(dir, []string{"README.md"})
	if err != nil {
		t.Fatalf("CollectFiles failed: %v", err)
	}

	if files.Extra["README.md"] != "# README" {
		t.Errorf("Extra[README.md] = %q, want %q", files.Extra["README.md"], "# README")
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

func TestCountTodos(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "main.go"), []byte("// TODO: fix this\n// FIXME: broken\nfmt.Println(\"hello\")\n"), 0644)
	os.WriteFile(filepath.Join(dir, "clean.go"), []byte("package main\n"), 0644)
	os.WriteFile(filepath.Join(dir, "notes.md"), []byte("# TODO list\nFIXME later\n"), 0644)
	// Binary file should be skipped
	os.WriteFile(filepath.Join(dir, "image.png"), []byte("TODO in binary"), 0644)

	count, err := CountTodos(dir, nil)
	if err != nil {
		t.Fatalf("CountTodos failed: %v", err)
	}
	if count != 4 {
		t.Errorf("CountTodos = %d, want 4", count)
	}
}

func TestCountTodosWithExclude(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "vendor"), 0755)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("// TODO: fix\n"), 0644)
	os.WriteFile(filepath.Join(dir, "vendor", "lib.go"), []byte("// TODO: vendor\n"), 0644)

	count, err := CountTodos(dir, nil)
	if err != nil {
		t.Fatalf("CountTodos failed: %v", err)
	}
	// vendor is in skipNames, should be excluded
	if count != 1 {
		t.Errorf("CountTodos = %d, want 1", count)
	}
}

func TestDetectCIFiles(t *testing.T) {
	dir := t.TempDir()

	// No CI files
	files := DetectCIFiles(dir)
	if len(files) != 0 {
		t.Errorf("DetectCIFiles = %v, want empty", files)
	}

	// GitHub Actions
	ghDir := filepath.Join(dir, ".github", "workflows")
	os.MkdirAll(ghDir, 0755)
	os.WriteFile(filepath.Join(ghDir, "ci.yml"), []byte("name: CI"), 0644)
	os.WriteFile(filepath.Join(ghDir, "release.yaml"), []byte("name: Release"), 0644)

	files = DetectCIFiles(dir)
	if len(files) != 2 {
		t.Errorf("DetectCIFiles len = %d, want 2", len(files))
	}

	// GitLab CI
	dir2 := t.TempDir()
	os.WriteFile(filepath.Join(dir2, ".gitlab-ci.yml"), []byte("stages:"), 0644)

	files2 := DetectCIFiles(dir2)
	if len(files2) != 1 {
		t.Errorf("DetectCIFiles len = %d, want 1", len(files2))
	}
	if files2[0] != ".gitlab-ci.yml" {
		t.Errorf("DetectCIFiles[0] = %q, want .gitlab-ci.yml", files2[0])
	}
}
