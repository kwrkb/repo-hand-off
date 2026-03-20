package differ

import (
	"testing"

	"github.com/kwrkb/repo-hand-off/internal/collector"
	"github.com/kwrkb/repo-hand-off/internal/parser"
)

func TestCompareAllStatuses(t *testing.T) {
	parsed := &parser.ParsedHandoff{
		Vision:  "# Vision\nSame content.",
		Plan:    "# Plan\nOld plan.",
		Lessons: "",
	}
	current := &collector.ProjectFiles{
		Vision:  "# Vision\nSame content.",
		Plan:    "# Plan\nNew plan.",
		Lessons: "# Lessons\nNew lesson.",
	}

	diffs := Compare(parsed, current)

	expected := map[string]string{
		"Vision":  "unchanged",
		"Plan":    "changed",
		"Lessons": "added",
	}

	for _, d := range diffs {
		want, ok := expected[d.Name]
		if !ok {
			t.Errorf("unexpected section %q", d.Name)
			continue
		}
		if d.Status != want {
			t.Errorf("%s: Status = %q, want %q", d.Name, d.Status, want)
		}
	}
}

func TestCompareRemoved(t *testing.T) {
	parsed := &parser.ParsedHandoff{
		Vision: "# Vision\nContent.",
	}
	current := &collector.ProjectFiles{}

	diffs := Compare(parsed, current)
	for _, d := range diffs {
		if d.Name == "Vision" && d.Status != "removed" {
			t.Errorf("Vision: Status = %q, want %q", d.Status, "removed")
		}
	}
}
