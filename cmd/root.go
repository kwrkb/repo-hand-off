package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "handoff",
	Short: "Export and share development state for seamless handoffs",
	Long:  "repo-hand-off preserves development context (code, plans, intent, lessons) and enables seamless handoffs between humans and AI.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
