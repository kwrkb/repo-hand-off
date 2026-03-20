package collector

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ProjectFiles holds the contents of key project files.
type ProjectFiles struct {
	Vision  string
	Plan    string
	Lessons string
	Readme  string
	Claude  string
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
	{"Readme", "README.md"},
	{"Claude", "CLAUDE.md"},
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
		case "Readme":
			pf.Readme = content
		case "Claude":
			pf.Claude = content
		}
	}

	if len(extraFiles) > 0 {
		pf.Extra = make(map[string]string)
		for _, f := range extraFiles {
			content, err := readFileIfExists(filepath.Join(dir, f))
			if err != nil {
				return nil, err
			}
			if content != "" {
				pf.Extra[f] = content
			}
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
