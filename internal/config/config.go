package config

// Config holds application configuration.
type Config struct {
	Output  string
	Quiet   bool
	Verbose bool
	Format  string
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Output: "HANDOFF.md",
		Format: "markdown",
	}
}
