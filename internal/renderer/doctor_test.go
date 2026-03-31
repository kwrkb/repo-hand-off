package renderer

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/kwrkb/repo-hand-off/internal/doctor"
)

func TestRenderDoctorTextNoFindings(t *testing.T) {
	out := RenderDoctorText(nil, "my-repo")
	if !strings.Contains(out, "All checks passed") {
		t.Errorf("expected 'All checks passed', got: %s", out)
	}
	if !strings.Contains(out, "my-repo") {
		t.Errorf("expected repo name in output, got: %s", out)
	}
}

func TestRenderDoctorText(t *testing.T) {
	findings := []doctor.Finding{
		{Rule: "vision-exists", Severity: doctor.Error, Message: "VISION.md が見つかりません", Action: "VISION.md を作成してください"},
		{Rule: "ci-exists", Severity: doctor.Warning, Message: "CI 設定ファイルが見つかりません", Action: ".github/workflows/ を追加してください"},
		{Rule: "todo-fixme-count", Severity: doctor.Info, Message: "TODO/FIXME が 12 件あります"},
	}

	out := RenderDoctorText(findings, "test-repo")

	if !strings.Contains(out, "[ERROR  ]") {
		t.Error("expected [ERROR  ] label")
	}
	if !strings.Contains(out, "[WARNING]") {
		t.Error("expected [WARNING] label")
	}
	if !strings.Contains(out, "[INFO   ]") {
		t.Error("expected [INFO   ] label")
	}
	if !strings.Contains(out, "→ VISION.md を作成してください") {
		t.Error("expected action line")
	}
	if !strings.Contains(out, "Summary: 1 error, 1 warning, 1 info") {
		t.Errorf("expected summary line, got: %s", out)
	}
}

func TestRenderDoctorJSON(t *testing.T) {
	findings := []doctor.Finding{
		{Rule: "vision-exists", Severity: doctor.Error, Message: "VISION.md が見つかりません", Action: "VISION.md を作成してください"},
		{Rule: "todo-fixme-count", Severity: doctor.Info, Message: "TODO/FIXME が 5 件あります"},
	}

	out, err := RenderDoctorJSON(findings)
	if err != nil {
		t.Fatalf("RenderDoctorJSON failed: %v", err)
	}

	var parsed struct {
		Findings []struct {
			Rule     string `json:"rule"`
			Severity string `json:"severity"`
			Message  string `json:"message"`
			Action   string `json:"action"`
		} `json:"findings"`
		Summary struct {
			Error   int `json:"error"`
			Warning int `json:"warning"`
			Info    int `json:"info"`
		} `json:"summary"`
	}

	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("JSON parse failed: %v", err)
	}

	if len(parsed.Findings) != 2 {
		t.Errorf("findings count = %d, want 2", len(parsed.Findings))
	}
	if parsed.Summary.Error != 1 {
		t.Errorf("error count = %d, want 1", parsed.Summary.Error)
	}
	if parsed.Summary.Info != 1 {
		t.Errorf("info count = %d, want 1", parsed.Summary.Info)
	}
	// Action should be omitted when empty
	if parsed.Findings[1].Action != "" {
		t.Errorf("expected empty action for todo finding, got: %q", parsed.Findings[1].Action)
	}
}

func TestCountErrors(t *testing.T) {
	findings := []doctor.Finding{
		{Severity: doctor.Error},
		{Severity: doctor.Warning},
		{Severity: doctor.Error},
		{Severity: doctor.Info},
	}
	if got := CountErrors(findings); got != 2 {
		t.Errorf("CountErrors = %d, want 2", got)
	}
}
