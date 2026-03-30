package doctor

import "github.com/kwrkb/repo-hand-off/internal/collector"

// DiagnoseOptions configures the Diagnose function.
type DiagnoseOptions struct {
	TodoThreshold int
	OutputPath    string // handoff output file path (for handoff-exists rule)
}

// Diagnose runs all default rules against the snapshot.
func Diagnose(snapshot *collector.Snapshot, opts DiagnoseOptions) []Finding {
	rules := make([]Rule, len(defaultRules))
	copy(rules, defaultRules)
	for i, r := range rules {
		switch r.Name() {
		case "todo-fixme-count":
			rules[i] = &TodoFixmeCount{Threshold: opts.TodoThreshold}
		case "handoff-exists":
			rules[i] = &HandoffExists{OutputPath: opts.OutputPath}
		}
	}
	return DiagnoseWith(snapshot, rules)
}

// DiagnoseWith runs the given rules against the snapshot.
func DiagnoseWith(snapshot *collector.Snapshot, rules []Rule) []Finding {
	var findings []Finding
	for _, r := range rules {
		findings = append(findings, r.Run(snapshot)...)
	}
	return findings
}
