package doctor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kwrkb/repo-hand-off/internal/collector"
)

func TestDiagnoseWith(t *testing.T) {
	s := &collector.Snapshot{}
	// A custom rule that always returns one finding
	rules := []Rule{&VisionExists{}}
	findings := DiagnoseWith(s, rules)
	if len(findings) != 1 {
		t.Errorf("DiagnoseWith findings = %d, want 1", len(findings))
	}
	if findings[0].Rule != "vision-exists" {
		t.Errorf("finding rule = %q, want vision-exists", findings[0].Rule)
	}
}

func TestDiagnoseWithNoFindings(t *testing.T) {
	s := &collector.Snapshot{
		Files: collector.ProjectFiles{
			Vision: "# Vision\nThis is a project that aims to solve a very important problem in the world.",
		},
	}
	rules := []Rule{&VisionExists{}, &VisionNotEmpty{}}
	findings := DiagnoseWith(s, rules)
	if len(findings) != 0 {
		t.Errorf("DiagnoseWith findings = %d, want 0", len(findings))
	}
}

func TestDiagnose(t *testing.T) {
	dir := t.TempDir()
	// Create a minimal valid project
	os.WriteFile(filepath.Join(dir, "VISION.md"), []byte("# Vision\nThis is a project that aims to solve a very important problem in the world."), 0644)
	os.WriteFile(filepath.Join(dir, "PLAN.md"), []byte("# Plan\nPhase 1: Set up the project and implement all the core features needed."), 0644)
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("# README\nSome content"), 0644)
	os.WriteFile(filepath.Join(dir, "LICENSE"), []byte("MIT"), 0644)
	os.WriteFile(filepath.Join(dir, "LESSONS.md"), []byte("# Lessons"), 0644)
	os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("*.o\n"), 0644)
	os.WriteFile(filepath.Join(dir, "HANDOFF.md"), []byte("# HANDOFF"), 0644)
	ghDir := filepath.Join(dir, ".github", "workflows")
	os.MkdirAll(ghDir, 0755)
	os.WriteFile(filepath.Join(ghDir, "ci.yml"), []byte("name: CI"), 0644)

	s := &collector.Snapshot{
		WorkDir: dir,
		Files: collector.ProjectFiles{
			Vision:  "# Vision\nThis is a project that aims to solve a very important problem in the world.",
			Plan:    "# Plan\nPhase 1: Set up the project and implement all the core features needed.",
			Lessons: "# Lessons",
			Extra:   map[string]string{"README.md": "# README"},
		},
		CIFiles:   []string{".github/workflows/ci.yml"},
		TodoCount: 3,
	}

	findings := Diagnose(s, DiagnoseOptions{TodoThreshold: 10})
	// Should have no errors or warnings for a well-configured project
	for _, f := range findings {
		if f.Severity == Error {
			t.Errorf("unexpected error finding: %s: %s", f.Rule, f.Message)
		}
	}
}

func TestSeverityString(t *testing.T) {
	tests := []struct {
		sev  Severity
		want string
	}{
		{Info, "info"},
		{Warning, "warning"},
		{Error, "error"},
		{Severity(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.sev.String(); got != tt.want {
			t.Errorf("Severity(%d).String() = %q, want %q", tt.sev, got, tt.want)
		}
	}
}
