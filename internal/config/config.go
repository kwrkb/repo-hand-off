package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds .handoff.yaml configuration.
type Config struct {
	Format        string   `yaml:"format"`
	Output        string   `yaml:"output"`
	Files         []string `yaml:"files"`
	Exclude       []string `yaml:"exclude"`
	Depth         int      `yaml:"depth"`
	TodoThreshold int      `yaml:"todo_threshold"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Format:        "markdown",
		Output:        "HANDOFF.md",
		Depth:         3,
		TodoThreshold: 10,
	}
}

// Load reads .handoff.yaml from dir. Missing file returns DefaultConfig, not an error.
func Load(dir string) (*Config, error) {
	path := filepath.Join(dir, ".handoff.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse .handoff.yaml: %w", err)
	}

	for _, pattern := range cfg.Exclude {
		if _, err := filepath.Match(pattern, ""); err != nil {
			return nil, fmt.Errorf("invalid exclude pattern %q in .handoff.yaml: %w", pattern, err)
		}
	}

	return cfg, nil
}
