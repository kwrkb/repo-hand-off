package parser

import (
	"strings"
)

// ParsedHandoff represents a parsed HANDOFF.md file.
type ParsedHandoff struct {
	Vision  string
	Plan    string
	Lessons string
	Readme  string
	Claude  string
	Extra   map[string]string
}

// handoffSections are the known top-level HANDOFF.md section names.
var handoffSections = map[string]bool{
	"Project":             true,
	"Vision":              true,
	"Plan":                true,
	"Lessons":             true,
	"README":              true,
	"CLAUDE":              true,
	"Current State":       true,
	"Directory Structure": true,
}

// skipSections are headers that should not be imported.
var skipSections = map[string]bool{
	"Project":             true,
	"Current State":       true,
	"Directory Structure": true,
	"HANDOFF.md":          true,
}

const extraPrefix = "Extra: "

// isHandoffSection returns true if the header is a known HANDOFF.md section
// or an extra file section (prefixed with "Extra: ").
func isHandoffSection(name string) bool {
	if handoffSections[name] {
		return true
	}
	return strings.HasPrefix(name, extraPrefix)
}

// ParseHandoffMarkdown parses a HANDOFF.md file into structured sections.
func ParseHandoffMarkdown(content string) (*ParsedHandoff, error) {
	result := &ParsedHandoff{
		Extra: make(map[string]string),
	}

	sections := splitSections(content)
	for name, body := range sections {
		body = normalizeBody(body)
		switch name {
		case "Vision":
			result.Vision = body
		case "Plan":
			result.Plan = body
		case "Lessons":
			result.Lessons = body
		case "README":
			result.Readme = stripCodeFence(body)
		case "CLAUDE":
			result.Claude = stripCodeFence(body)
		default:
			if skipSections[name] {
				continue
			}
			if strings.HasPrefix(name, extraPrefix) {
				result.Extra[strings.TrimPrefix(name, extraPrefix)] = body
			}
		}
	}

	return result, nil
}

// splitSections splits markdown content by ## headers that are known
// HANDOFF.md sections or "Extra: " prefixed sections. Nested ## headers
// within content are preserved as part of the section body. Code fences
// (```` or ```) are tracked so that ## headers inside them are not treated
// as section boundaries.
func splitSections(content string) map[string]string {
	sections := make(map[string]string)
	var currentName string
	var currentBody strings.Builder
	inFence := false

	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "````") || (!inFence && strings.HasPrefix(trimmed, "```")) {
			inFence = !inFence
		}

		if !inFence && strings.HasPrefix(line, "## ") {
			name := strings.TrimPrefix(line, "## ")
			if isHandoffSection(name) {
				// Save previous section
				if currentName != "" {
					sections[currentName] = currentBody.String()
				}
				currentName = name
				currentBody.Reset()
				continue
			}
		}
		// Append to current section body
		if currentName != "" {
			if currentBody.Len() > 0 || line != "" {
				if currentBody.Len() > 0 {
					currentBody.WriteString("\n")
				}
				currentBody.WriteString(line)
			}
		}
	}
	// Save last section
	if currentName != "" {
		sections[currentName] = currentBody.String()
	}

	return sections
}

// stripCodeFence removes a wrapping code fence (````markdown ... ````) if present.
func stripCodeFence(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	lines := strings.Split(s, "\n")
	if len(lines) >= 2 && strings.HasPrefix(lines[0], "````") && strings.HasPrefix(lines[len(lines)-1], "````") {
		inner := strings.Join(lines[1:len(lines)-1], "\n")
		return strings.TrimSpace(inner)
	}
	return s
}

// normalizeBody trims whitespace and treats "Not found." as empty.
func normalizeBody(body string) string {
	body = strings.TrimSpace(body)
	if body == "Not found." {
		return ""
	}
	return body
}
