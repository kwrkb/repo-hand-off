package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kwrkb/repo-hand-off/internal/collector"
	"github.com/kwrkb/repo-hand-off/internal/doctor"
	"github.com/kwrkb/repo-hand-off/internal/renderer"
	"github.com/spf13/cobra"
)

var (
	doctorFormat string
	doctorStrict bool
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose handoff readiness of the project",
	Long:  "Check whether the project has the essential files and state needed for a successful handoff.",
	RunE: func(cmd *cobra.Command, args []string) error {
		logVerbose("Collecting project state from %s", workDir)
		snapshot, err := collector.Collect(workDir, collectOpts())
		if err != nil {
			return fmt.Errorf("failed to collect project state: %w", err)
		}
		logVerbose("Collection complete, running diagnostics")

		findings := doctor.Diagnose(snapshot, doctor.DiagnoseOptions{
			TodoThreshold: cfg.TodoThreshold,
			OutputPath:    cfg.Output,
		})
		logVerbose("Diagnostics complete: %d findings", len(findings))

		var output string
		switch doctorFormat {
		case renderer.FormatJSON:
			output, err = renderer.RenderDoctorJSON(findings)
			if err != nil {
				return err
			}
		default:
			repoName := filepath.Base(workDir)
			output = renderer.RenderDoctorText(findings, repoName)
		}

		fmt.Print(output)

		if doctorStrict && renderer.CountErrors(findings) > 0 {
			// Diagnostic output is already printed; exit directly to avoid
			// cobra printing a redundant error message.
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	doctorCmd.Flags().StringVar(&doctorFormat, "format", renderer.FormatText, "Output format: text or json")
	doctorCmd.Flags().BoolVar(&doctorStrict, "strict", false, "Exit with code 1 if any errors are found")
	rootCmd.AddCommand(doctorCmd)
}
