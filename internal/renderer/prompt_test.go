package renderer

import (
	"strings"
	"testing"
)

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
