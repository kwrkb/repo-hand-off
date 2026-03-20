package cmd

import (
	"fmt"
	"os"

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
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		snapshot, err := collector.Collect(dir)
		if err != nil {
			return fmt.Errorf("failed to collect project state: %w", err)
		}

		output := renderer.RenderPrompt(snapshot, promptFormat)
		fmt.Print(output)
		return nil
	},
}

func init() {
	promptCmd.Flags().StringVar(&promptFormat, "format", renderer.FormatMarkdown, "Output format: markdown or xml")
	rootCmd.AddCommand(promptCmd)
}
