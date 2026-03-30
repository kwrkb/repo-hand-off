package renderer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kwrkb/repo-hand-off/internal/doctor"
)

const (
	FormatText = "text"
	FormatJSON = "json"
)

// RenderDoctorText renders findings as human-readable text output.
func RenderDoctorText(findings []doctor.Finding, repoName string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "handoff doctor — %s\n", repoName)

	if len(findings) == 0 {
		b.WriteString("\n  All checks passed.\n")
		return b.String()
	}

	b.WriteString("\n")
	for _, f := range findings {
		label := strings.ToUpper(f.Severity.String())
		fmt.Fprintf(&b, "  [%-7s] %s: %s\n", label, f.Rule, f.Message)
		if f.Action != "" {
			fmt.Fprintf(&b, "            → %s\n", f.Action)
		}
	}

	summary := countSeverities(findings)
	fmt.Fprintf(&b, "\nSummary: %d error, %d warning, %d info\n", summary.Errors, summary.Warnings, summary.Infos)
	return b.String()
}

// RenderDoctorJSON renders findings as JSON output.
func RenderDoctorJSON(findings []doctor.Finding) (string, error) {
	type jsonFinding struct {
		Rule     string `json:"rule"`
		Severity string `json:"severity"`
		Message  string `json:"message"`
		Action   string `json:"action,omitempty"`
	}

	type jsonSummary struct {
		Error   int `json:"error"`
		Warning int `json:"warning"`
		Info    int `json:"info"`
	}

	type jsonOutput struct {
		Findings []jsonFinding `json:"findings"`
		Summary  jsonSummary   `json:"summary"`
	}

	jf := make([]jsonFinding, len(findings))
	for i, f := range findings {
		jf[i] = jsonFinding{
			Rule:     f.Rule,
			Severity: f.Severity.String(),
			Message:  f.Message,
			Action:   f.Action,
		}
	}

	summary := countSeverities(findings)
	out := jsonOutput{
		Findings: jf,
		Summary: jsonSummary{
			Error:   summary.Errors,
			Warning: summary.Warnings,
			Info:    summary.Infos,
		},
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal doctor output: %w", err)
	}
	return string(data) + "\n", nil
}

type severityCounts struct {
	Errors   int
	Warnings int
	Infos    int
}

func countSeverities(findings []doctor.Finding) severityCounts {
	var c severityCounts
	for _, f := range findings {
		switch f.Severity {
		case doctor.Error:
			c.Errors++
		case doctor.Warning:
			c.Warnings++
		case doctor.Info:
			c.Infos++
		}
	}
	return c
}

// CountErrors returns the number of Error-severity findings.
func CountErrors(findings []doctor.Finding) int {
	return countSeverities(findings).Errors
}
