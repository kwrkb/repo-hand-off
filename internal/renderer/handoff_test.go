package renderer

import (
	"strings"
	"testing"
	"time"

	"github.com/kwrkb/repo-hand-off/internal/collector"
)

func testSnapshot() *collector.Snapshot {
	return &collector.Snapshot{
		Timestamp: time.Date(2026, 3, 20, 12, 0, 0, 0, time.UTC),
		Git: collector.GitInfo{
			Branch:      "main",
			ShortHash:   "abc1234",
			RemoteURL:   "https://github.com/test/repo.git",
			HasChanges:  true,
			DiffSummary: " file.go | 3 +++",
		},
		Files: collector.ProjectFiles{
			Vision: "# Vision\nBuild something great.",
			Plan:   "# Plan\nStep 1: Do it.",
		},
		DirTree:    "repo/\n├── main.go\n└── go.mod",
		RecentLogs: []string{"abc1234 initial commit"},
	}
}

func TestRenderHandoff(t *testing.T) {
	s := testSnapshot()
	result := RenderHandoff(s)

	checks := []string{
		"# HANDOFF.md",
		"2026-03-20 12:00:00",
		"Branch: main @ abc1234",
		"Uncommitted changes: yes",
		"## Vision",
		"Build something great",
		"## Plan",
		"Step 1: Do it",
		"## Lessons",
		"Not found",
		"## Current State",
		"abc1234 initial commit",
		"## Directory Structure",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("output missing %q", check)
		}
	}
}
