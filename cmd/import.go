package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kwrkb/repo-hand-off/internal/parser"
	"github.com/spf13/cobra"
)

var (
	importForce bool
	importFile  string
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Restore project files from HANDOFF.md",
	Long:  "Parse HANDOFF.md and write Vision, Plan, and Lessons sections back to their respective files.",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(importFile)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", importFile, err)
		}

		parsed, err := parser.ParseHandoffMarkdown(string(data))
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", importFile, err)
		}

		fileMap := map[string]string{
			"VISION.md":  parsed.Vision,
			"PLAN.md":    parsed.Plan,
			"LESSONS.md": parsed.Lessons,
		}

		// Add extra files
		for name, content := range parsed.Extra {
			fileMap[name] = content
		}

		// Sort keys for deterministic output
		keys := make([]string, 0, len(fileMap))
		for k := range fileMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, filename := range keys {
			content := fileMap[filename]
			if content == "" {
				logVerbose("Skipping %s (empty)", filename)
				continue
			}

			// Prevent path traversal: resolve and verify target stays under workDir
			target := filepath.Join(workDir, filename)
			resolved, err := filepath.Abs(target)
			if err != nil {
				return fmt.Errorf("failed to resolve path for %s: %w", filename, err)
			}
			if !strings.HasPrefix(resolved, workDir+string(filepath.Separator)) && resolved != workDir {
				logInfo("Skipping %s (path traversal detected)", filename)
				continue
			}

			if !importForce {
				f, err := os.OpenFile(resolved, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
				if err != nil {
					if os.IsExist(err) {
						logInfo("Skipping %s (already exists, use --force to overwrite)", filename)
						continue
					}
					return fmt.Errorf("failed to write %s: %w", filename, err)
				}
				if _, err := f.WriteString(content + "\n"); err != nil {
					f.Close()
					return fmt.Errorf("failed to write %s: %w", filename, err)
				}
				if err := f.Close(); err != nil {
					return fmt.Errorf("failed to close %s: %w", filename, err)
				}
			} else {
				if err := os.WriteFile(resolved, []byte(content+"\n"), 0644); err != nil {
					return fmt.Errorf("failed to write %s: %w", filename, err)
				}
			}
			logInfo("Wrote %s", filename)
		}

		return nil
	},
}

func init() {
	importCmd.Flags().BoolVar(&importForce, "force", false, "Overwrite existing files")
	importCmd.Flags().StringVarP(&importFile, "file", "f", "HANDOFF.md", "Input file path")
	rootCmd.AddCommand(importCmd)
}
