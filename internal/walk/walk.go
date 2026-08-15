package walk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/ikigenba/llm-lint/internal/rules"
)

type Walker struct {
	Root    string
	Exclude []string
	RunGit  func(dir string, args ...string) ([]byte, error)
}

func (w *Walker) Files(paths []string) ([]string, error) {
	var files []string
	for _, path := range paths {
		full := path
		if !filepath.IsAbs(full) {
			full = filepath.Join(w.Root, full)
		}
		err := filepath.WalkDir(full, func(name string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !entry.IsDir() {
				files = append(files, name)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walk: %w", err)
		}
	}
	return files, nil
}

func Candidates(files []string, r rules.Rule) []string {
	var candidates []string
	for _, file := range files {
		name := strings.TrimPrefix(filepath.ToSlash(file), "/")
		included := len(r.Include) == 0
		for _, pattern := range r.Include {
			if matched, _ := doublestar.Match(pattern, name); matched {
				included = true
				break
			}
		}
		if !included {
			continue
		}
		excluded := false
		for _, pattern := range r.Exclude {
			if matched, _ := doublestar.Match(pattern, name); matched {
				excluded = true
				break
			}
		}
		if !excluded {
			candidates = append(candidates, file)
		}
	}
	return candidates
}
