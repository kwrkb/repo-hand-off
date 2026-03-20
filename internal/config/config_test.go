package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.Format != "markdown" {
		t.Errorf("Format = %q, want %q", cfg.Format, "markdown")
	}
	if cfg.Output != "HANDOFF.md" {
		t.Errorf("Output = %q, want %q", cfg.Output, "HANDOFF.md")
	}
	if cfg.Depth != 3 {
		t.Errorf("Depth = %d, want 3", cfg.Depth)
	}
}

func TestLoadPartialYAML(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".handoff.yaml"), []byte("format: xml\n"), 0644)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.Format != "xml" {
		t.Errorf("Format = %q, want %q", cfg.Format, "xml")
	}
	// Defaults preserved for unspecified fields
	if cfg.Output != "HANDOFF.md" {
		t.Errorf("Output = %q, want %q", cfg.Output, "HANDOFF.md")
	}
	if cfg.Depth != 3 {
		t.Errorf("Depth = %d, want 3", cfg.Depth)
	}
}

func TestLoadFullYAML(t *testing.T) {
	dir := t.TempDir()
	content := `format: xml
output: out.md
files:
  - NOTES.md
  - TODO.md
exclude:
  - "*.tmp"
depth: 5
`
	os.WriteFile(filepath.Join(dir, ".handoff.yaml"), []byte(content), 0644)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.Format != "xml" {
		t.Errorf("Format = %q, want %q", cfg.Format, "xml")
	}
	if cfg.Output != "out.md" {
		t.Errorf("Output = %q, want %q", cfg.Output, "out.md")
	}
	if len(cfg.Files) != 2 {
		t.Errorf("Files len = %d, want 2", len(cfg.Files))
	}
	if len(cfg.Exclude) != 1 {
		t.Errorf("Exclude len = %d, want 1", len(cfg.Exclude))
	}
	if cfg.Depth != 5 {
		t.Errorf("Depth = %d, want 5", cfg.Depth)
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".handoff.yaml"), []byte("{{invalid yaml"), 0644)

	_, err := Load(dir)
	if err == nil {
		t.Fatal("Load should return error for invalid YAML")
	}
	if !strings.Contains(err.Error(), "failed to parse .handoff.yaml") {
		t.Errorf("error should contain context, got: %v", err)
	}
}

func TestLoadBadExcludePattern(t *testing.T) {
	dir := t.TempDir()
	content := `exclude:
  - "[invalid"
`
	os.WriteFile(filepath.Join(dir, ".handoff.yaml"), []byte(content), 0644)

	_, err := Load(dir)
	if err == nil {
		t.Fatal("Load should return error for bad exclude pattern")
	}
	if !strings.Contains(err.Error(), "invalid exclude pattern") {
		t.Errorf("error should mention invalid pattern, got: %v", err)
	}
}

func TestLoadUnreadableFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".handoff.yaml")
	os.WriteFile(path, []byte("format: xml\n"), 0644)
	os.Chmod(path, 0000)
	t.Cleanup(func() { os.Chmod(path, 0644) })

	_, err := Load(dir)
	if err == nil {
		t.Fatal("Load should return error for unreadable file")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Format != "markdown" {
		t.Errorf("Format = %q, want %q", cfg.Format, "markdown")
	}
	if cfg.Output != "HANDOFF.md" {
		t.Errorf("Output = %q, want %q", cfg.Output, "HANDOFF.md")
	}
	if cfg.Depth != 3 {
		t.Errorf("Depth = %d, want 3", cfg.Depth)
	}
	if cfg.Files != nil {
		t.Error("Files should be nil by default")
	}
	if cfg.Exclude != nil {
		t.Error("Exclude should be nil by default")
	}
}
