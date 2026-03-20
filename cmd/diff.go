package cmd

import (
	"fmt"
	"os"

	"github.com/kwrkb/repo-hand-off/internal/collector"
	"github.com/kwrkb/repo-hand-off/internal/differ"
	"github.com/kwrkb/repo-hand-off/internal/parser"
	"github.com/spf13/cobra"
)

var diffFile string

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare HANDOFF.md sections against current project files",
	Long:  "Parse HANDOFF.md and compare each section against the current file contents, showing which sections have changed.",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(diffFile)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", diffFile, err)
		}

		parsed, err := parser.ParseHandoffMarkdown(string(data))
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", diffFile, err)
		}

		// Merge extra file names from parsed HANDOFF into collect options
		// so differ can compare them against current files on disk
		opts := collectOpts()
		for name := range parsed.Extra {
			opts.ExtraFiles = append(opts.ExtraFiles, name)
		}

		logVerbose("Collecting current project state from %s", workDir)
		snapshot, err := collector.Collect(workDir, opts)
		if err != nil {
			return fmt.Errorf("failed to collect project state: %w", err)
		}

		diffs := differ.Compare(parsed, &snapshot.Files)

		fmt.Printf("%-15s %s\n", "Section", "Status")
		fmt.Printf("%-15s %s\n", "-------", "------")
		for _, d := range diffs {
			fmt.Printf("%-15s %s\n", d.Name, d.Status)
		}

		return nil
	},
}

func init() {
	diffCmd.Flags().StringVarP(&diffFile, "file", "f", "HANDOFF.md", "Input file path")
	rootCmd.AddCommand(diffCmd)
}
