package cmd

import (
	"fmt"

	"github.com/kwrkb/repo-hand-off/internal/collector"
	"github.com/kwrkb/repo-hand-off/internal/renderer"
	"github.com/spf13/cobra"
)

var promptFormat string

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Generate an AI-ready prompt from project state",
	Long:  "Collect project state and output an AI-ready context prompt to stdout.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Apply config default for format if not explicitly set
		if !cmd.Flags().Changed("format") {
			promptFormat = cfg.Format
		}

		logVerbose("Collecting project state from %s", workDir)
		snapshot, err := collector.Collect(workDir, collectOpts())
		if err != nil {
			return fmt.Errorf("failed to collect project state: %w", err)
		}
		logVerbose("Collection complete")

		output, err := renderer.RenderPrompt(snapshot, promptFormat)
		if err != nil {
			return err
		}
		fmt.Print(output)
		return nil
	},
}

func init() {
	promptCmd.Flags().StringVar(&promptFormat, "format", renderer.FormatMarkdown, "Output format: markdown or xml")
	rootCmd.AddCommand(promptCmd)
}
