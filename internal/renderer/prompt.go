package renderer

import (
	"strings"

	"github.com/kwrkb/repo-hand-off/internal/collector"
)

const (
	FormatMarkdown = "markdown"
	FormatXML      = "xml"
)

// RenderPrompt generates an AI-ready prompt from a snapshot.
func RenderPrompt(s *collector.Snapshot, format string) string {
	switch format {
	case FormatXML:
		return renderPromptXML(s)
	default:
		return renderPromptMarkdown(s)
	}
}

func renderPromptMarkdown(s *collector.Snapshot) string {
	var b strings.Builder

	b.WriteString("# Project Handoff Context\n\n")
	b.WriteString("You are continuing development on an existing project. ")
	b.WriteString("Below is the current state of the project. ")
	b.WriteString("Read it carefully, then continue development based on the plan and current state.\n\n")
	b.WriteString("---\n\n")
	b.WriteString(RenderHandoff(s, FormatMarkdown))

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
