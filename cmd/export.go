package cmd

import (
	"fmt"
	"os"

	"github.com/kwrkb/repo-hand-off/internal/collector"
	"github.com/kwrkb/repo-hand-off/internal/renderer"
	"github.com/spf13/cobra"
)

var (
	exportOutput string
	exportFormat string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export project state to HANDOFF.md",
	Long:  "Collect project state (Git info, project files, directory structure) and generate HANDOFF.md.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Apply config defaults for flags not explicitly set
		if !cmd.Flags().Changed("output") {
			exportOutput = cfg.Output
		}
		if !cmd.Flags().Changed("format") {
			exportFormat = cfg.Format
		}

		logVerbose("Collecting project state from %s", workDir)
		snapshot, err := collector.Collect(workDir, collectOpts())
		if err != nil {
			return fmt.Errorf("failed to collect project state: %w", err)
		}
		logVerbose("Collection complete")

		content := renderer.RenderHandoff(snapshot, exportFormat)

		if err := os.WriteFile(exportOutput, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", exportOutput, err)
		}

		logInfo("Exported project state to %s", exportOutput)
		logVerbose("Written %d bytes to %s", len(content), exportOutput)
		return nil
	},
}

func init() {
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "HANDOFF.md", "Output file path")
	exportCmd.Flags().StringVar(&exportFormat, "format", renderer.FormatMarkdown, "Output format: markdown or xml")
	rootCmd.AddCommand(exportCmd)
}
