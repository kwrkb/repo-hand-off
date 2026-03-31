package collector

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ProjectFiles holds the contents of key project files.
type ProjectFiles struct {
	Vision  string
	Plan    string
	Lessons string
	Extra   map[string]string
}

// knownFiles maps field names to file paths to look for.
var knownFiles = []struct {
	Name string
	Path string
}{
	{"Vision", "VISION.md"},
	{"Plan", "PLAN.md"},
	{"Lessons", "LESSONS.md"},
}

// autoExtraFiles are files automatically included in Extra if they exist.
var autoExtraFiles = []string{
	"README.md",
	"CLAUDE.md",
	"AGENTS.md",
	"GEMINI.md",
}

// CollectFiles reads key project files from the given directory.
func CollectFiles(dir string, extraFiles []string) (*ProjectFiles, error) {
	pf := &ProjectFiles{}
	for _, kf := range knownFiles {
		content, err := readFileIfExists(filepath.Join(dir, kf.Path))
		if err != nil {
			return nil, err
		}
		switch kf.Name {
		case "Vision":
			pf.Vision = content
		case "Plan":
			pf.Plan = content
		case "Lessons":
			pf.Lessons = content
		}
	}

	// Merge autoExtraFiles and extraFiles (deduplicated, order preserved)
	allExtra := make([]string, 0, len(autoExtraFiles)+len(extraFiles))
	seen := make(map[string]bool)
	for _, f := range append(autoExtraFiles, extraFiles...) {
		if !seen[f] {
			allExtra = append(allExtra, f)
			seen[f] = true
		}
	}

	pf.Extra = make(map[string]string)
	for _, f := range allExtra {
		content, err := readFileIfExists(filepath.Join(dir, f))
		if err != nil {
			return nil, err
		}
		if content != "" {
			pf.Extra[f] = content
		}
	}

	return pf, nil
}

func readFileIfExists(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(data), nil
}

// BuildDirTree builds a tree representation of the directory up to maxDepth.
func BuildDirTree(dir string, maxDepth int, exclude []string) (string, error) {
	if maxDepth <= 0 {
		maxDepth = 3
	}
	var lines []string
	lines = append(lines, filepath.Base(dir)+"/")

	err := buildTree(dir, dir, "", maxDepth, 0, exclude, &lines)
	if err != nil {
		return "", err
	}
	return strings.Join(lines, "\n"), nil
}

func buildTree(root, dir, prefix string, maxDepth, depth int, exclude []string, lines *[]string) error {
	if depth >= maxDepth {
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// Filter and sort: directories first, then files
	var dirs, files []fs.DirEntry
	for _, e := range entries {
		name := e.Name()
		if shouldSkip(name, exclude) {
			continue
		}
		if e.IsDir() {
			dirs = append(dirs, e)
		} else {
			files = append(files, e)
		}
	}

	sorted := append(dirs, files...)
	for i, e := range sorted {
		isLast := i == len(sorted)-1
		connector := "├── "
		childPrefix := "│   "
		if isLast {
			connector = "└── "
			childPrefix = "    "
		}

		name := e.Name()
		if e.IsDir() {
			name += "/"
		}
		*lines = append(*lines, fmt.Sprintf("%s%s%s", prefix, connector, name))

		if e.IsDir() {
			err := buildTree(root, filepath.Join(dir, e.Name()), prefix+childPrefix, maxDepth, depth+1, exclude, lines)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Pre-sorted list for binary search.
var skipNames = []string{
	".DS_Store", ".git", ".idea", ".next", ".venv", ".vscode",
	"__pycache__", "build", "dist", "node_modules", "target", "vendor",
}

var todoPattern = regexp.MustCompile(`(?i)\b(TODO|FIXME)\b`)

// CountTodos counts TODO/FIXME comments in source files under dir.
func CountTodos(dir string, exclude []string) (int, error) {
	count := 0
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if shouldSkip(d.Name(), exclude) {
				return filepath.SkipDir
			}
			return nil
		}
		if shouldSkip(d.Name(), exclude) {
			return nil
		}
		if !isTextFile(d.Name()) {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return nil // skip unreadable files
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			count += len(todoPattern.FindAllString(scanner.Text(), -1))
		}
		// Ignore scanner errors (e.g. token too long): count what we could read
		return nil
	})
	return count, err
}

// textExtensions lists file extensions considered as text/source files.
var textExtensions = map[string]bool{
	".go": true, ".py": true, ".js": true, ".ts": true, ".tsx": true, ".jsx": true,
	".rb": true, ".rs": true, ".java": true, ".kt": true, ".swift": true, ".c": true,
	".cpp": true, ".h": true, ".hpp": true, ".cs": true, ".php": true, ".sh": true,
	".bash": true, ".zsh": true, ".fish": true, ".yaml": true, ".yml": true,
	".toml": true, ".json": true, ".xml": true, ".html": true, ".css": true,
	".scss": true, ".less": true, ".sql": true, ".md": true, ".txt": true,
	".vue": true, ".svelte": true, ".dart": true, ".lua": true, ".r": true,
	".ex": true, ".exs": true, ".erl": true, ".hs": true, ".ml": true,
	".tf": true, ".proto": true, ".graphql": true, ".dockerfile": true,
}

func isTextFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	if textExtensions[ext] {
		return true
	}
	// Include common extensionless files
	lower := strings.ToLower(name)
	return lower == "makefile" || lower == "dockerfile" || lower == "rakefile" || lower == "gemfile"
}

// DetectCIFiles returns relative paths of CI configuration files found in dir.
func DetectCIFiles(dir string) []string {
	var files []string

	// GitHub Actions
	ghDir := filepath.Join(dir, ".github", "workflows")
	if entries, err := os.ReadDir(ghDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			ext := strings.ToLower(filepath.Ext(e.Name()))
			if ext == ".yml" || ext == ".yaml" {
				files = append(files, filepath.Join(".github", "workflows", e.Name()))
			}
		}
	}

	// GitLab CI
	for _, name := range []string{".gitlab-ci.yml", ".gitlab-ci.yaml"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			files = append(files, name)
			break
		}
	}

	return files
}

func shouldSkip(name string, exclude []string) bool {
	i := sort.SearchStrings(skipNames, name)
	if i < len(skipNames) && skipNames[i] == name {
		return true
	}
	for _, pattern := range exclude {
		if matched, _ := filepath.Match(pattern, name); matched {
			return true
		}
	}
	return false
}
