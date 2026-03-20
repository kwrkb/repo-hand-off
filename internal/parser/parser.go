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
			result.Readme = body
		case "CLAUDE":
			result.Claude = body
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
// within content are preserved as part of the section body.
func splitSections(content string) map[string]string {
	sections := make(map[string]string)
	var currentName string
	var currentBody strings.Builder

	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "## ") {
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

// normalizeBody trims whitespace and treats "Not found." as empty.
func normalizeBody(body string) string {
	body = strings.TrimSpace(body)
	if body == "Not found." {
		return ""
	}
	return body
}
