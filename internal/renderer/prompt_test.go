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
	result, err := RenderPrompt(s, "xml")
	if err != nil {
		t.Fatalf("RenderPrompt failed: %v", err)
	}

	if strings.Contains(result, "<project>") {
		t.Error("should not contain <project> section for empty GitInfo")
	}
	for _, check := range []string{"<handoff>", "<vision>", "<directory_structure>"} {
		if !strings.Contains(result, check) {
			t.Errorf("xml prompt missing %q", check)
		}
	}
}

func TestRenderPromptMarkdown(t *testing.T) {
	s := testSnapshot()
	result, err := RenderPrompt(s, "markdown")
	if err != nil {
		t.Fatalf("RenderPrompt failed: %v", err)
	}

	if !strings.Contains(result, "# Project Handoff Context") {
		t.Error("markdown prompt should contain header")
	}
	if !strings.Contains(result, "# HANDOFF.md") {
		t.Error("markdown prompt should contain handoff content")
	}
}

func TestRenderPromptXML(t *testing.T) {
	s := testSnapshot()
	result, err := RenderPrompt(s, "xml")
	if err != nil {
		t.Fatalf("RenderPrompt failed: %v", err)
	}

	checks := []string{
		"<handoff>",
		"</handoff>",
		"<instructions>",
		"<project>",
		"<vision>",
		"<plan>",
		"<readme>",
		"<claude>",
		"<recent_commits>",
		"<directory_structure>",
	}
	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("xml prompt missing %q", check)
		}
	}
}

func TestRenderPromptInvalidFormat(t *testing.T) {
	s := testSnapshot()
	_, err := RenderPrompt(s, "yaml")
	if err == nil {
		t.Fatal("RenderPrompt should return error for invalid format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("error should mention unsupported format, got: %v", err)
	}
}

func TestIsValidFormat(t *testing.T) {
	tests := []struct {
		format string
		valid  bool
	}{
		{"markdown", true},
		{"xml", true},
		{"json", false},
		{"", false},
		{"MARKDOWN", false},
	}
	for _, tt := range tests {
		if got := IsValidFormat(tt.format); got != tt.valid {
			t.Errorf("IsValidFormat(%q) = %v, want %v", tt.format, got, tt.valid)
		}
	}
}
