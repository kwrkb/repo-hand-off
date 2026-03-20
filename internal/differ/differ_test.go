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
		Extra: map[string]string{
			"README.md": "# README\nSame usage.",
			"CLAUDE.md": "# CLAUDE\nOld guidance.",
		},
	}
	current := &collector.ProjectFiles{
		Vision:  "# Vision\nSame content.",
		Plan:    "# Plan\nNew plan.",
		Lessons: "# Lessons\nNew lesson.",
		Extra: map[string]string{
			"README.md": "# README\nSame usage.",
			"CLAUDE.md": "# CLAUDE\nNew guidance.",
		},
	}

	diffs := Compare(parsed, current)

	expected := map[string]string{
		"Vision":    "unchanged",
		"Plan":      "changed",
		"Lessons":   "added",
		"README.md": "unchanged",
		"CLAUDE.md": "changed",
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

func TestCompareExtraFiles(t *testing.T) {
	parsed := &parser.ParsedHandoff{
		Extra: map[string]string{
			"NOTES.md": "old notes",
			"OLD.md":   "only in parsed",
		},
	}
	current := &collector.ProjectFiles{
		Extra: map[string]string{
			"NOTES.md": "new notes",
			"NEW.md":   "only in current",
		},
	}

	diffs := Compare(parsed, current)

	expected := map[string]string{
		"Vision":   "unchanged", // both empty
		"Plan":     "unchanged",
		"Lessons":  "unchanged",
		"NOTES.md": "changed",
		"OLD.md":   "removed",
		"NEW.md":   "added",
	}

	for _, d := range diffs {
		want, ok := expected[d.Name]
		if !ok {
			t.Errorf("unexpected diff entry %q", d.Name)
			continue
		}
		if d.Status != want {
			t.Errorf("%s: Status = %q, want %q", d.Name, d.Status, want)
		}
	}
}

func TestCompareBothEmpty(t *testing.T) {
	parsed := &parser.ParsedHandoff{}
	current := &collector.ProjectFiles{}

	diffs := Compare(parsed, current)
	for _, d := range diffs {
		if d.Status != "unchanged" {
			t.Errorf("%s: Status = %q, want %q", d.Name, d.Status, "unchanged")
		}
	}
}

func TestCompareWhitespaceNormalization(t *testing.T) {
	// TrimSpace strips leading/trailing whitespace from the whole string,
	// so surrounding whitespace differences should be normalized away
	parsed := &parser.ParsedHandoff{
		Vision: "\n  # Vision\nContent.  \n\n",
	}
	current := &collector.ProjectFiles{
		Vision: "# Vision\nContent.",
	}

	diffs := Compare(parsed, current)
	for _, d := range diffs {
		if d.Name == "Vision" && d.Status != "unchanged" {
			t.Errorf("Vision with surrounding whitespace diff: Status = %q, want %q", d.Status, "unchanged")
		}
	}

	// Inner whitespace differences should be detected as "changed"
	parsed2 := &parser.ParsedHandoff{
		Vision: "# Vision\n Content.",
	}
	current2 := &collector.ProjectFiles{
		Vision: "# Vision\nContent.",
	}
	diffs2 := Compare(parsed2, current2)
	for _, d := range diffs2 {
		if d.Name == "Vision" && d.Status != "changed" {
			t.Errorf("Vision with inner whitespace diff: Status = %q, want %q", d.Status, "changed")
		}
	}
}

func TestCompareDeterministicOrder(t *testing.T) {
	parsed := &parser.ParsedHandoff{
		Extra: map[string]string{
			"Z.md": "z",
			"A.md": "a",
			"M.md": "m",
		},
	}
	current := &collector.ProjectFiles{
		Extra: map[string]string{
			"Z.md": "z",
			"A.md": "a",
			"M.md": "m",
		},
	}

	diffs := Compare(parsed, current)
	// First 3 should be the built-in sections (Vision, Plan, Lessons), then extras sorted.
	extraDiffs := diffs[3:]
	for i := 1; i < len(extraDiffs); i++ {
		if extraDiffs[i-1].Name > extraDiffs[i].Name {
			t.Errorf("extra diffs not sorted: %q before %q", extraDiffs[i-1].Name, extraDiffs[i].Name)
		}
	}
}
