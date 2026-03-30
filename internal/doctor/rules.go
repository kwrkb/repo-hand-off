package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kwrkb/repo-hand-off/internal/collector"
)

// defaultRules is the built-in set of diagnostic rules.
var defaultRules = []Rule{
	&VisionExists{},
	&PlanExists{},
	&ReadmeExists{},
	&LicenseExists{},
	&LessonsExists{},
	&CIExists{},
	&GitignoreExists{},
	&UncommittedChanges{},
	&TodoFixmeCount{Threshold: 10},
	&HandoffFreshness{MaxAge: 7 * 24 * time.Hour},
	&VisionNotEmpty{},
	&PlanNotEmpty{},
}

// minContentLength is the minimum length for a file to be considered non-empty.
const minContentLength = 50

// --- Rule 1: VisionExists ---

type VisionExists struct{}

func (r *VisionExists) Name() string { return "vision-exists" }
func (r *VisionExists) Run(s *collector.Snapshot) []Finding {
	if s.Files.Vision == "" {
		return []Finding{{
			Rule:     r.Name(),
			Severity: Error,
			Message:  "VISION.md が見つかりません",
			Action:   "VISION.md を作成し、プロジェクトの目的を記述してください",
		}}
	}
	return nil
}

// --- Rule 2: PlanExists ---

type PlanExists struct{}

func (r *PlanExists) Name() string { return "plan-exists" }
func (r *PlanExists) Run(s *collector.Snapshot) []Finding {
	if s.Files.Plan == "" {
		return []Finding{{
			Rule:     r.Name(),
			Severity: Error,
			Message:  "PLAN.md が見つかりません",
			Action:   "PLAN.md を作成し、実装計画を記述してください",
		}}
	}
	return nil
}

// --- Rule 3: ReadmeExists ---

type ReadmeExists struct{}

func (r *ReadmeExists) Name() string { return "readme-exists" }
func (r *ReadmeExists) Run(s *collector.Snapshot) []Finding {
	if s.Files.Extra["README.md"] == "" {
		return []Finding{{
			Rule:     r.Name(),
			Severity: Warning,
			Message:  "README.md が見つかりません",
			Action:   "README.md を作成し、プロジェクトの概要を記述してください",
		}}
	}
	return nil
}

// --- Rule 4: LicenseExists ---

type LicenseExists struct{}

func (r *LicenseExists) Name() string { return "license-exists" }
func (r *LicenseExists) Run(s *collector.Snapshot) []Finding {
	if s.WorkDir == "" {
		return nil
	}
	for _, name := range []string{"LICENSE", "LICENSE.md", "LICENSE.txt", "LICENCE"} {
		if _, err := os.Stat(filepath.Join(s.WorkDir, name)); err == nil {
			return nil
		}
	}
	return []Finding{{
		Rule:     r.Name(),
		Severity: Warning,
		Message:  "LICENSE ファイルが見つかりません",
		Action:   "LICENSE ファイルを追加してください",
	}}
}

// --- Rule 5: LessonsExists ---

type LessonsExists struct{}

func (r *LessonsExists) Name() string { return "lessons-exists" }
func (r *LessonsExists) Run(s *collector.Snapshot) []Finding {
	if s.Files.Lessons == "" {
		return []Finding{{
			Rule:     r.Name(),
			Severity: Info,
			Message:  "LESSONS.md が見つかりません",
			Action:   "LESSONS.md を作成し、開発で得た教訓を記録してください",
		}}
	}
	return nil
}

// --- Rule 6: CIExists ---

type CIExists struct{}

func (r *CIExists) Name() string { return "ci-exists" }
func (r *CIExists) Run(s *collector.Snapshot) []Finding {
	if len(s.CIFiles) == 0 {
		return []Finding{{
			Rule:     r.Name(),
			Severity: Warning,
			Message:  "CI 設定ファイルが見つかりません",
			Action:   ".github/workflows/ または .gitlab-ci.yml を追加してください",
		}}
	}
	return nil
}

// --- Rule 7: GitignoreExists ---

type GitignoreExists struct{}

func (r *GitignoreExists) Name() string { return "gitignore-exists" }
func (r *GitignoreExists) Run(s *collector.Snapshot) []Finding {
	if s.WorkDir == "" {
		return nil
	}
	if _, err := os.Stat(filepath.Join(s.WorkDir, ".gitignore")); err != nil {
		return []Finding{{
			Rule:     r.Name(),
			Severity: Info,
			Message:  ".gitignore が見つかりません",
			Action:   ".gitignore を作成し、追跡不要なファイルを除外してください",
		}}
	}
	return nil
}

// --- Rule 8: UncommittedChanges ---

type UncommittedChanges struct{}

func (r *UncommittedChanges) Name() string { return "uncommitted-changes" }
func (r *UncommittedChanges) Run(s *collector.Snapshot) []Finding {
	if s.Git.HasChanges {
		return []Finding{{
			Rule:     r.Name(),
			Severity: Warning,
			Message:  "未コミットの変更があります",
			Action:   "変更をコミットしてから handoff してください",
		}}
	}
	return nil
}

// --- Rule 9: TodoFixmeCount ---

type TodoFixmeCount struct {
	Threshold int
}

func (r *TodoFixmeCount) Name() string { return "todo-fixme-count" }
func (r *TodoFixmeCount) Run(s *collector.Snapshot) []Finding {
	threshold := r.Threshold
	if threshold <= 0 {
		threshold = 10
	}
	if s.TodoCount >= threshold {
		return []Finding{{
			Rule:     r.Name(),
			Severity: Info,
			Message:  fmt.Sprintf("TODO/FIXME が %d 件あります", s.TodoCount),
		}}
	}
	return nil
}

// --- Rule 10: HandoffFreshness ---

type HandoffFreshness struct {
	MaxAge time.Duration
}

func (r *HandoffFreshness) Name() string { return "handoff-freshness" }
func (r *HandoffFreshness) Run(s *collector.Snapshot) []Finding {
	if s.WorkDir == "" {
		return nil
	}
	path := filepath.Join(s.WorkDir, "HANDOFF.md")
	info, err := os.Stat(path)
	if err != nil {
		return []Finding{{
			Rule:     r.Name(),
			Severity: Warning,
			Message:  "HANDOFF.md が見つかりません",
			Action:   "handoff export を実行して HANDOFF.md を生成してください",
		}}
	}
	maxAge := r.MaxAge
	if maxAge <= 0 {
		maxAge = 7 * 24 * time.Hour
	}
	age := time.Since(info.ModTime())
	if age > maxAge {
		days := int(age.Hours() / 24)
		return []Finding{{
			Rule:     r.Name(),
			Severity: Warning,
			Message:  fmt.Sprintf("HANDOFF.md の最終更新が %d 日前です", days),
			Action:   "handoff export を実行して HANDOFF.md を更新してください",
		}}
	}
	return nil
}

// --- Rule 11: VisionNotEmpty ---

type VisionNotEmpty struct{}

func (r *VisionNotEmpty) Name() string { return "vision-not-empty" }
func (r *VisionNotEmpty) Run(s *collector.Snapshot) []Finding {
	if s.Files.Vision == "" {
		return nil // VisionExists handles missing case
	}
	if len(strings.TrimSpace(s.Files.Vision)) < minContentLength {
		return []Finding{{
			Rule:     r.Name(),
			Severity: Error,
			Message:  "VISION.md の内容が極端に短いです",
			Action:   "VISION.md にプロジェクトの目的・背景を十分に記述してください",
		}}
	}
	return nil
}

// --- Rule 12: PlanNotEmpty ---

type PlanNotEmpty struct{}

func (r *PlanNotEmpty) Name() string { return "plan-not-empty" }
func (r *PlanNotEmpty) Run(s *collector.Snapshot) []Finding {
	if s.Files.Plan == "" {
		return nil // PlanExists handles missing case
	}
	if len(strings.TrimSpace(s.Files.Plan)) < minContentLength {
		return []Finding{{
			Rule:     r.Name(),
			Severity: Error,
			Message:  "PLAN.md の内容が極端に短いです",
			Action:   "PLAN.md に実装計画を十分に記述してください",
		}}
	}
	return nil
}
