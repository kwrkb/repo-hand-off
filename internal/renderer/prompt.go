package renderer

import (
	"fmt"
	"strings"

	"github.com/kwrkb/repo-hand-off/internal/collector"
)

const (
	FormatMarkdown = "markdown"
	FormatXML      = "xml"
)

// ValidFormats contains all supported output format names.
var ValidFormats = []string{FormatMarkdown, FormatXML}

// RenderPrompt generates an AI-ready prompt from a snapshot.
func RenderPrompt(s *collector.Snapshot, format string) (string, error) {
	switch format {
	case FormatXML:
		return renderPromptXML(s), nil
	case FormatMarkdown:
		return renderPromptMarkdown(s), nil
	default:
		return "", fmt.Errorf("unsupported format %q (valid: %s)", format, strings.Join(ValidFormats, ", "))
	}
}

func renderPromptMarkdown(s *collector.Snapshot) string {
	var b strings.Builder

	b.WriteString("# Project Handoff Context\n\n")
	b.WriteString("You are continuing development on an existing project. ")
	b.WriteString("Below is the current state of the project. ")
	b.WriteString("Read it carefully, then continue development based on the plan and current state.\n\n")
	b.WriteString("---\n\n")
	handoff, _ := RenderHandoff(s, FormatMarkdown) // format is hardcoded, cannot fail
	b.WriteString(handoff)

	return b.String()
}

func renderPromptXML(s *collector.Snapshot) string {
	var b strings.Builder

	b.WriteString("<handoff>\n")
	b.WriteString("<instructions>\n")
	b.WriteString("You are continuing development on an existing project. ")
	b.WriteString("Below is the current state of the project. ")
	b.WriteString("Read it carefully, then continue development based on the plan and current state.\n")
	b.WriteString("</instructions>\n\n")

	renderHandoffXMLBody(&b, s)

	b.WriteString("</handoff>\n")

	return b.String()
}
