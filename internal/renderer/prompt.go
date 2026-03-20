package renderer

import (
	"fmt"
	"strings"

	"github.com/kwrkb/repo-hand-off/internal/collector"
)

// RenderPrompt generates an AI-ready prompt from a snapshot.
func RenderPrompt(s *collector.Snapshot, format string) string {
	switch format {
	case "xml":
		return renderXML(s)
	default:
		return renderMarkdown(s)
	}
}

func renderMarkdown(s *collector.Snapshot) string {
	var b strings.Builder

	b.WriteString("# Project Handoff Context\n\n")
	b.WriteString("You are continuing development on an existing project. ")
	b.WriteString("Below is the current state of the project. ")
	b.WriteString("Read it carefully, then continue development based on the plan and current state.\n\n")
	b.WriteString("---\n\n")
	b.WriteString(RenderHandoff(s))

	return b.String()
}

func renderXML(s *collector.Snapshot) string {
	var b strings.Builder

	b.WriteString("<handoff>\n")
	b.WriteString("<instructions>\n")
	b.WriteString("You are continuing development on an existing project. ")
	b.WriteString("Below is the current state of the project. ")
	b.WriteString("Read it carefully, then continue development based on the plan and current state.\n")
	b.WriteString("</instructions>\n\n")

	b.WriteString("<project>\n")
	if s.Git.RemoteURL != "" {
		b.WriteString(fmt.Sprintf("  <repository>%s</repository>\n", s.Git.RemoteURL))
	}
	b.WriteString(fmt.Sprintf("  <branch>%s</branch>\n", s.Git.Branch))
	b.WriteString(fmt.Sprintf("  <commit>%s</commit>\n", s.Git.ShortHash))
	b.WriteString(fmt.Sprintf("  <uncommitted_changes>%t</uncommitted_changes>\n", s.Git.HasChanges))
	b.WriteString("</project>\n\n")

	writeXMLSection(&b, "vision", s.Files.Vision)
	writeXMLSection(&b, "plan", s.Files.Plan)
	writeXMLSection(&b, "lessons", s.Files.Lessons)

	if len(s.RecentLogs) > 0 {
		b.WriteString("<recent_commits>\n")
		for _, log := range s.RecentLogs {
			b.WriteString(fmt.Sprintf("  <commit>%s</commit>\n", log))
		}
		b.WriteString("</recent_commits>\n\n")
	}

	if s.Git.DiffSummary != "" {
		b.WriteString("<uncommitted_changes>\n")
		b.WriteString(s.Git.DiffSummary)
		b.WriteString("\n</uncommitted_changes>\n\n")
	}

	b.WriteString("<directory_structure>\n")
	b.WriteString(s.DirTree)
	b.WriteString("\n</directory_structure>\n")

	b.WriteString("</handoff>\n")

	return b.String()
}

func writeXMLSection(b *strings.Builder, tag, content string) {
	if content == "" {
		return
	}
	b.WriteString(fmt.Sprintf("<%s>\n", tag))
	b.WriteString(strings.TrimSpace(content))
	b.WriteString(fmt.Sprintf("\n</%s>\n\n", tag))
}
