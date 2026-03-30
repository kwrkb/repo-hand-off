package doctor

import "github.com/kwrkb/repo-hand-off/internal/collector"

// Diagnose runs all default rules against the snapshot.
// todoThreshold configures the TODO/FIXME count threshold.
func Diagnose(snapshot *collector.Snapshot, todoThreshold int) []Finding {
	rules := make([]Rule, len(defaultRules))
	copy(rules, defaultRules)
	// Replace TodoFixmeCount with configured threshold
	for i, r := range rules {
		if r.Name() == "todo-fixme-count" {
			rules[i] = &TodoFixmeCount{Threshold: todoThreshold}
			break
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
