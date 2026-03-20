package cmd

import (
	"fmt"
	"os"

	"github.com/kwrkb/repo-hand-off/internal/collector"
	"github.com/kwrkb/repo-hand-off/internal/renderer"
	"github.com/spf13/cobra"
)

var exportOutput string

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export project state to HANDOFF.md",
	Long:  "Collect project state (Git info, project files, directory structure) and generate HANDOFF.md.",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		snapshot, err := collector.Collect(dir)
		if err != nil {
			return fmt.Errorf("failed to collect project state: %w", err)
		}

		content := renderer.RenderHandoff(snapshot)

		if err := os.WriteFile(exportOutput, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", exportOutput, err)
		}

		fmt.Fprintf(os.Stderr, "Exported project state to %s\n", exportOutput)
		return nil
	},
}

func init() {
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "HANDOFF.md", "Output file path")
	rootCmd.AddCommand(exportCmd)
}
