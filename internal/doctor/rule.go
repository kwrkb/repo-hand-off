package doctor

import "github.com/kwrkb/repo-hand-off/internal/collector"

// Severity represents the importance level of a finding.
type Severity int

const (
	Info    Severity = iota // informational
	Warning                 // improvement recommended
	Error                   // blocks handoff readiness
)

// String returns the lowercase name of the severity.
func (s Severity) String() string {
	switch s {
	case Info:
		return "info"
	case Warning:
		return "warning"
	case Error:
		return "error"
	default:
		return "unknown"
	}
}

// Finding represents a single diagnostic result from a rule.
type Finding struct {
	Rule     string   // rule name (e.g. "vision-exists")
	Severity Severity // severity level
	Message  string   // human-readable description
	Action   string   // suggested next step (empty if none)
}

// Rule defines a single diagnostic check.
type Rule interface {
	Name() string
	Run(snapshot *collector.Snapshot) []Finding
}
