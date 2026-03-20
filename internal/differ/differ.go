package differ

import (
	"sort"
	"strings"

	"github.com/kwrkb/repo-hand-off/internal/collector"
	"github.com/kwrkb/repo-hand-off/internal/parser"
)

// SectionDiff describes the status of a single section comparison.
type SectionDiff struct {
	Name   string
	Status string // "unchanged", "changed", "added", "removed"
}

// Compare compares parsed HANDOFF.md sections against current project files.
func Compare(parsed *parser.ParsedHandoff, current *collector.ProjectFiles) []SectionDiff {
	type pair struct {
		name    string
		parsed  string
		current string
	}

	pairs := []pair{
		{"Vision", parsed.Vision, current.Vision},
		{"Plan", parsed.Plan, current.Plan},
		{"Lessons", parsed.Lessons, current.Lessons},
		{"README", parsed.Readme, current.Readme},
		{"CLAUDE", parsed.Claude, current.Claude},
	}

	var diffs []SectionDiff
	for _, p := range pairs {
		diffs = append(diffs, SectionDiff{
			Name:   p.name,
			Status: compareStrings(p.parsed, p.current),
		})
	}

	// Compare Extra files
	allExtra := make(map[string]bool)
	for k := range parsed.Extra {
		allExtra[k] = true
	}
	for k := range current.Extra {
		allExtra[k] = true
	}
	keys := make([]string, 0, len(allExtra))
	for k := range allExtra {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		diffs = append(diffs, SectionDiff{
			Name:   name,
			Status: compareStrings(parsed.Extra[name], current.Extra[name]),
		})
	}

	return diffs
}

func compareStrings(a, b string) string {
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)

	switch {
	case a == "" && b == "":
		return "unchanged"
	case a != "" && b == "":
		return "removed"
	case a == "" && b != "":
		return "added"
	case a == b:
		return "unchanged"
	default:
		return "changed"
	}
}
