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
			Extra: map[string]string{
				"README.md": "# README\nUsage details.",
				"CLAUDE.md": "# CLAUDE\nAssistant guidance.",
			},
		},
		DirTree:    "repo/\n├── main.go\n└── go.mod",
		RecentLogs: []string{"abc1234 initial commit"},
	}
}

func TestRenderHandoffNoGit(t *testing.T) {
	s := &collector.Snapshot{
		Timestamp: time.Date(2026, 3, 20, 12, 0, 0, 0, time.UTC),
		Git:       collector.GitInfo{}, // empty = non-git
		Files: collector.ProjectFiles{
			Vision: "# Vision\nBuild something great.",
			Plan:   "# Plan\nStep 1: Do it.",
		},
		DirTree: "project/\n├── main.go\n└── go.mod",
	}
	result, err := RenderHandoff(s, FormatMarkdown)
	if err != nil {
		t.Fatalf("RenderHandoff failed: %v", err)
	}

	if !strings.Contains(result, "Git: Not available") {
		t.Error("should show 'Git: Not available' for empty GitInfo")
	}
	if strings.Contains(result, "Branch:") {
		t.Error("should not show Branch for empty GitInfo")
	}
	for _, check := range []string{"## Vision", "## Plan", "## Directory Structure", "Build something great"} {
		if !strings.Contains(result, check) {
			t.Errorf("output missing %q", check)
		}
	}
}

func TestRenderHandoff(t *testing.T) {
	s := testSnapshot()
	result, err := RenderHandoff(s, FormatMarkdown)
	if err != nil {
		t.Fatalf("RenderHandoff failed: %v", err)
	}

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
		"## Extra: README.md",
		"````\n# README\nUsage details.\n````",
		"## Extra: CLAUDE.md",
		"````\n# CLAUDE\nAssistant guidance.\n````",
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

func TestRenderHandoffXML(t *testing.T) {
	s := testSnapshot()
	result, err := RenderHandoff(s, FormatXML)
	if err != nil {
		t.Fatalf("RenderHandoff failed: %v", err)
	}

	checks := []string{
		"<handoff>",
		"</handoff>",
		"<project>",
		"<branch>main</branch>",
		"<vision>",
		"<plan>",
		`<extra name="CLAUDE.md">`,
		`<extra name="README.md">`,
		"<recent_commits>",
		"<directory_structure>",
	}
	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("xml output missing %q", check)
		}
	}

	// Should NOT contain <instructions>
	if strings.Contains(result, "<instructions>") {
		t.Error("handoff XML should not contain <instructions>")
	}
}

func TestRenderHandoffXMLNoGit(t *testing.T) {
	s := &collector.Snapshot{
		Timestamp: time.Date(2026, 3, 20, 12, 0, 0, 0, time.UTC),
		Git:       collector.GitInfo{},
		Files: collector.ProjectFiles{
			Vision: "# Vision\nBuild something great.",
		},
		DirTree: "project/\n├── main.go\n└── go.mod",
	}
	result, err := RenderHandoff(s, FormatXML)
	if err != nil {
		t.Fatalf("RenderHandoff failed: %v", err)
	}

	if strings.Contains(result, "<project>") {
		t.Error("should not contain <project> section for empty GitInfo")
	}
	for _, check := range []string{"<handoff>", "<vision>", "<directory_structure>"} {
		if !strings.Contains(result, check) {
			t.Errorf("xml output missing %q", check)
		}
	}
}

func TestRenderHandoffExtra(t *testing.T) {
	s := testSnapshot()
	s.Files.Extra = map[string]string{
		"NOTES.md": "# Notes\nSome notes.",
	}

	md, err := RenderHandoff(s, FormatMarkdown)
	if err != nil {
		t.Fatalf("RenderHandoff markdown failed: %v", err)
	}
	if !strings.Contains(md, "## Extra: NOTES.md") {
		t.Error("markdown should contain extra file section with 'Extra: ' prefix")
	}
	if !strings.Contains(md, "Some notes.") {
		t.Error("markdown should contain extra file content")
	}

	xml, err := RenderHandoff(s, FormatXML)
	if err != nil {
		t.Fatalf("RenderHandoff xml failed: %v", err)
	}
	if !strings.Contains(xml, `<extra name="NOTES.md">`) {
		t.Error("xml should contain extra file tag")
	}
}

func TestRenderHandoffInvalidFormat(t *testing.T) {
	s := testSnapshot()
	_, err := RenderHandoff(s, "json")
	if err == nil {
		t.Fatal("RenderHandoff should return error for invalid format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("error should mention unsupported format, got: %v", err)
	}
}

func TestRenderHandoffNoChanges(t *testing.T) {
	s := testSnapshot()
	s.Git.HasChanges = false
	s.Git.DiffSummary = ""

	result, err := RenderHandoff(s, FormatMarkdown)
	if err != nil {
		t.Fatalf("RenderHandoff failed: %v", err)
	}

	if !strings.Contains(result, "Uncommitted changes: no") {
		t.Error("should show 'Uncommitted changes: no'")
	}
	if strings.Contains(result, "### Uncommitted Changes") {
		t.Error("should not show diff section when no changes")
	}
}

func TestRenderHandoffNoRemote(t *testing.T) {
	s := testSnapshot()
	s.Git.RemoteURL = ""

	result, err := RenderHandoff(s, FormatMarkdown)
	if err != nil {
		t.Fatalf("RenderHandoff failed: %v", err)
	}

	if strings.Contains(result, "Repository:") {
		t.Error("should not show Repository when RemoteURL is empty")
	}
	if !strings.Contains(result, "Branch: main") {
		t.Error("should still show Branch")
	}
}

func TestRenderHandoffEmptySnapshot(t *testing.T) {
	s := &collector.Snapshot{
		Timestamp: time.Date(2026, 3, 20, 12, 0, 0, 0, time.UTC),
		Git:       collector.GitInfo{},
		Files:     collector.ProjectFiles{},
		DirTree:   "project/",
	}

	md, err := RenderHandoff(s, FormatMarkdown)
	if err != nil {
		t.Fatalf("RenderHandoff failed: %v", err)
	}
	for _, section := range []string{"## Vision", "## Plan", "## Lessons"} {
		if !strings.Contains(md, section) {
			t.Errorf("missing section %q", section)
		}
	}
	// README/CLAUDE are now Extra — empty files should not produce Extra sections
	for _, section := range []string{"## Extra: README.md", "## Extra: CLAUDE.md"} {
		if strings.Contains(md, section) {
			t.Errorf("empty snapshot should not contain %q", section)
		}
	}
	// Vision, Plan, Lessons show "Not found." (3 sections)
	if strings.Count(md, "Not found.") != 3 {
		t.Errorf("expected 3 'Not found.' markers, got %d", strings.Count(md, "Not found."))
	}

	xml, err := RenderHandoff(s, FormatXML)
	if err != nil {
		t.Fatalf("RenderHandoff XML failed: %v", err)
	}
	// Empty files should not produce XML sections
	for _, tag := range []string{"<vision>", "<plan>", "<lessons>", "<extra"} {
		if strings.Contains(xml, tag) {
			t.Errorf("empty content should not produce %s tag", tag)
		}
	}
}
