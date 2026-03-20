package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds .handoff.yaml configuration.
type Config struct {
	Format  string   `yaml:"format"`
	Output  string   `yaml:"output"`
	Files   []string `yaml:"files"`
	Exclude []string `yaml:"exclude"`
	Depth   int      `yaml:"depth"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Format: "markdown",
		Output: "HANDOFF.md",
		Depth:  3,
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
		return nil, err
	}
	return cfg, nil
}
