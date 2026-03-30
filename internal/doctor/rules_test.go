package doctor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kwrkb/repo-hand-off/internal/collector"
)

func TestVisionExists(t *testing.T) {
	tests := []struct {
		name    string
		vision  string
		wantLen int
	}{
		{"missing", "", 1},
		{"present", "# Vision\nSome content here that is long enough to pass checks.", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &collector.Snapshot{Files: collector.ProjectFiles{Vision: tt.vision}}
			got := (&VisionExists{}).Run(s)
			if len(got) != tt.wantLen {
				t.Errorf("VisionExists findings = %d, want %d", len(got), tt.wantLen)
			}
			if tt.wantLen > 0 && got[0].Severity != Error {
				t.Errorf("severity = %v, want Error", got[0].Severity)
			}
		})
	}
}

func TestPlanExists(t *testing.T) {
	tests := []struct {
		name    string
		plan    string
		wantLen int
	}{
		{"missing", "", 1},
		{"present", "# Plan\nSome content here that is long enough to pass checks.", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &collector.Snapshot{Files: collector.ProjectFiles{Plan: tt.plan}}
			got := (&PlanExists{}).Run(s)
			if len(got) != tt.wantLen {
				t.Errorf("PlanExists findings = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestReadmeExists(t *testing.T) {
	tests := []struct {
		name    string
		extra   map[string]string
		wantLen int
	}{
		{"missing", map[string]string{}, 1},
		{"present", map[string]string{"README.md": "# README"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &collector.Snapshot{Files: collector.ProjectFiles{Extra: tt.extra}}
			got := (&ReadmeExists{}).Run(s)
			if len(got) != tt.wantLen {
				t.Errorf("ReadmeExists findings = %d, want %d", len(got), tt.wantLen)
			}
			if tt.wantLen > 0 && got[0].Severity != Warning {
				t.Errorf("severity = %v, want Warning", got[0].Severity)
			}
		})
	}
}

func TestLicenseExists(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		dir := t.TempDir()
		s := &collector.Snapshot{WorkDir: dir}
		got := (&LicenseExists{}).Run(s)
		if len(got) != 1 {
			t.Fatalf("findings = %d, want 1", len(got))
		}
		if got[0].Severity != Warning {
			t.Errorf("severity = %v, want Warning", got[0].Severity)
		}
	})

	t.Run("present", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "LICENSE"), []byte("MIT"), 0644)
		s := &collector.Snapshot{WorkDir: dir}
		got := (&LicenseExists{}).Run(s)
		if len(got) != 0 {
			t.Errorf("findings = %d, want 0", len(got))
		}
	})

	t.Run("no workdir", func(t *testing.T) {
		s := &collector.Snapshot{}
		got := (&LicenseExists{}).Run(s)
		if len(got) != 0 {
			t.Errorf("findings = %d, want 0", len(got))
		}
	})
}

func TestLessonsExists(t *testing.T) {
	s := &collector.Snapshot{Files: collector.ProjectFiles{}}
	got := (&LessonsExists{}).Run(s)
	if len(got) != 1 {
		t.Fatalf("findings = %d, want 1", len(got))
	}
	if got[0].Severity != Info {
		t.Errorf("severity = %v, want Info", got[0].Severity)
	}
}

func TestCIExists(t *testing.T) {
	tests := []struct {
		name    string
		ciFiles []string
		wantLen int
	}{
		{"missing", nil, 1},
		{"present", []string{".github/workflows/ci.yml"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &collector.Snapshot{CIFiles: tt.ciFiles}
			got := (&CIExists{}).Run(s)
			if len(got) != tt.wantLen {
				t.Errorf("CIExists findings = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestGitignoreExists(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		dir := t.TempDir()
		s := &collector.Snapshot{WorkDir: dir}
		got := (&GitignoreExists{}).Run(s)
		if len(got) != 1 {
			t.Fatalf("findings = %d, want 1", len(got))
		}
		if got[0].Severity != Info {
			t.Errorf("severity = %v, want Info", got[0].Severity)
		}
	})

	t.Run("present", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("*.o\n"), 0644)
		s := &collector.Snapshot{WorkDir: dir}
		got := (&GitignoreExists{}).Run(s)
		if len(got) != 0 {
			t.Errorf("findings = %d, want 0", len(got))
		}
	})
}

func TestUncommittedChanges(t *testing.T) {
	tests := []struct {
		name       string
		hasChanges bool
		wantLen    int
	}{
		{"clean", false, 0},
		{"dirty", true, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &collector.Snapshot{Git: collector.GitInfo{HasChanges: tt.hasChanges}}
			got := (&UncommittedChanges{}).Run(s)
			if len(got) != tt.wantLen {
				t.Errorf("UncommittedChanges findings = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestTodoFixmeCount(t *testing.T) {
	tests := []struct {
		name      string
		count     int
		threshold int
		wantLen   int
	}{
		{"below threshold", 5, 10, 0},
		{"at threshold", 10, 10, 1},
		{"above threshold", 15, 10, 1},
		{"zero threshold uses default", 15, 0, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &collector.Snapshot{TodoCount: tt.count}
			r := &TodoFixmeCount{Threshold: tt.threshold}
			got := r.Run(s)
			if len(got) != tt.wantLen {
				t.Errorf("TodoFixmeCount findings = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestHandoffExists(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		dir := t.TempDir()
		s := &collector.Snapshot{WorkDir: dir}
		got := (&HandoffExists{}).Run(s)
		if len(got) != 1 {
			t.Fatalf("findings = %d, want 1", len(got))
		}
		if !strings.Contains(got[0].Message, "見つかりません") {
			t.Errorf("message = %q, want contains '見つかりません'", got[0].Message)
		}
	})

	t.Run("present", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "HANDOFF.md"), []byte("# HANDOFF"), 0644)
		s := &collector.Snapshot{WorkDir: dir}
		got := (&HandoffExists{}).Run(s)
		if len(got) != 0 {
			t.Errorf("findings = %d, want 0", len(got))
		}
	})

	t.Run("custom output path", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "out.md"), []byte("# HANDOFF"), 0644)
		s := &collector.Snapshot{WorkDir: dir}
		got := (&HandoffExists{OutputPath: "out.md"}).Run(s)
		if len(got) != 0 {
			t.Errorf("findings = %d, want 0", len(got))
		}
	})

	t.Run("custom output path missing", func(t *testing.T) {
		dir := t.TempDir()
		s := &collector.Snapshot{WorkDir: dir}
		got := (&HandoffExists{OutputPath: "out.md"}).Run(s)
		if len(got) != 1 {
			t.Fatalf("findings = %d, want 1", len(got))
		}
		if !strings.Contains(got[0].Message, "out.md") {
			t.Errorf("message = %q, want contains 'out.md'", got[0].Message)
		}
	})
}

func TestVisionNotEmpty(t *testing.T) {
	tests := []struct {
		name    string
		vision  string
		wantLen int
	}{
		{"missing (skip)", "", 0},
		{"too short", "# Vision", 1},
		{"sufficient", "# Vision\nThis is a project that aims to solve a very important problem.", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &collector.Snapshot{Files: collector.ProjectFiles{Vision: tt.vision}}
			got := (&VisionNotEmpty{}).Run(s)
			if len(got) != tt.wantLen {
				t.Errorf("VisionNotEmpty findings = %d, want %d", len(got), tt.wantLen)
			}
			if tt.wantLen > 0 && got[0].Severity != Error {
				t.Errorf("severity = %v, want Error", got[0].Severity)
			}
		})
	}
}

func TestPlanNotEmpty(t *testing.T) {
	tests := []struct {
		name    string
		plan    string
		wantLen int
	}{
		{"missing (skip)", "", 0},
		{"too short", "# Plan", 1},
		{"sufficient", "# Plan\nPhase 1: Set up the project and implement core features.", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &collector.Snapshot{Files: collector.ProjectFiles{Plan: tt.plan}}
			got := (&PlanNotEmpty{}).Run(s)
			if len(got) != tt.wantLen {
				t.Errorf("PlanNotEmpty findings = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}
