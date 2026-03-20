package cmd

import (
	"fmt"
	"os"

	"github.com/kwrkb/repo-hand-off/internal/collector"
	"github.com/kwrkb/repo-hand-off/internal/config"
	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags.
var Version = "dev"

var (
	cfg     *config.Config
	workDir string
	quiet   bool
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:     "handoff",
	Short:   "Export and share development state for seamless handoffs",
	Long:    "repo-hand-off preserves development context (code, plans, intent, lessons) and enables seamless handoffs between humans and AI.",
	Version: Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		workDir, err = os.Getwd()
		if err != nil {
			return err
		}
		cfg, err = config.Load(workDir)
		return err
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Suppress informational messages")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed progress")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func collectOpts() collector.CollectOptions {
	return collector.CollectOptions{
		ExtraFiles: cfg.Files,
		Exclude:    cfg.Exclude,
		Depth:      cfg.Depth,
	}
}

func logInfo(format string, args ...any) {
	if !quiet {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}

func logVerbose(format string, args ...any) {
	if verbose {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}
