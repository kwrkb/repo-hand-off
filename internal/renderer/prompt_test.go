package renderer

import (
	"strings"
	"testing"
	"time"

	"github.com/kwrkb/repo-hand-off/internal/collector"
)

func TestRenderPromptXMLNoGit(t *testing.T) {
	s := &collector.Snapshot{
		Timestamp: time.Date(2026, 3, 20, 12, 0, 0, 0, time.UTC),
		Git:       collector.GitInfo{}, // empty = non-git
		Files: collector.ProjectFiles{
			Vision: "# Vision\nBuild something great.",
		},
		DirTree: "project/\n├── main.go\n└── go.mod",
	}
	result := RenderPrompt(s, "xml")

	if strings.Contains(result, "<project>") {
		t.Error("should not contain <project> section for empty GitInfo")
	}
	// Other sections should still render
	for _, check := range []string{"<handoff>", "<vision>", "<directory_structure>"} {
		if !strings.Contains(result, check) {
			t.Errorf("xml prompt missing %q", check)
		}
	}
}

func TestRenderPromptMarkdown(t *testing.T) {
	s := testSnapshot()
	result := RenderPrompt(s, "markdown")

	if !strings.Contains(result, "# Project Handoff Context") {
		t.Error("markdown prompt should contain header")
	}
	if !strings.Contains(result, "# HANDOFF.md") {
		t.Error("markdown prompt should contain handoff content")
	}
}

func TestRenderPromptXML(t *testing.T) {
	s := testSnapshot()
	result := RenderPrompt(s, "xml")

	checks := []string{
		"<handoff>",
		"</handoff>",
		"<instructions>",
		"<project>",
		"<vision>",
		"<plan>",
		"<recent_commits>",
		"<directory_structure>",
	}
	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("xml prompt missing %q", check)
		}
	}
}
