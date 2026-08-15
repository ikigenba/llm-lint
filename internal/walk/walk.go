package walk

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/ikigenba/llm-lint/internal/rules"
)

const sniffSize = 8 * 1024

type Walker struct {
	Root    string
	Exclude []string
	RunGit  func(dir string, args ...string) ([]byte, error)
}

func (w *Walker) Files(paths []string) ([]string, error) {
	if len(paths) == 0 {
		paths = []string{"."}
	}

	root, err := filepath.Abs(w.Root)
	if err != nil {
		return nil, fmt.Errorf("walk: resolve root: %w", err)
	}
	gitPaths := make([]string, 0, len(paths))
	explicitFiles := make([]string, 0, len(paths))
	for _, path := range paths {
		full := path
		if !filepath.IsAbs(full) {
			full = filepath.Join(root, full)
		}
		info, statErr := os.Stat(full)
		if statErr != nil {
			return nil, fmt.Errorf("walk: inspect named path %q: %w", path, statErr)
		}
		rel, relErr := filepath.Rel(root, full)
		if relErr != nil {
			return nil, fmt.Errorf("walk: resolve named path %q: %w", path, relErr)
		}
		gitPaths = append(gitPaths, filepath.ToSlash(rel))
		if info.Mode().IsRegular() {
			explicitFiles = append(explicitFiles, rel)
		}
	}

	runGit := w.RunGit
	if runGit == nil {
		runGit = func(dir string, args ...string) ([]byte, error) {
			cmd := exec.Command("git", args...)
			cmd.Dir = dir
			return cmd.Output()
		}
	}

	var discovered []string
	if output, gitErr := runGit(root, "rev-parse", "--is-inside-work-tree"); gitErr == nil && strings.TrimSpace(string(output)) == "true" {
		args := []string{"ls-files", "--cached", "--others", "--exclude-standard", "-z", "--"}
		args = append(args, gitPaths...)
		output, gitErr = runGit(root, args...)
		if gitErr != nil {
			return nil, fmt.Errorf("walk: git ls-files: %w", gitErr)
		}
		for _, name := range bytes.Split(output, []byte{0}) {
			if len(name) > 0 {
				discovered = append(discovered, filepath.FromSlash(string(name)))
			}
		}
	} else {
		for _, path := range gitPaths {
			full := filepath.Join(root, filepath.FromSlash(path))
			walkErr := filepath.WalkDir(full, func(name string, entry os.DirEntry, entryErr error) error {
				if entryErr != nil {
					if name == full {
						return entryErr
					}
					return nil
				}
				if entry.IsDir() {
					if entry.Name() == ".git" {
						return filepath.SkipDir
					}
					return nil
				}
				rel, relErr := filepath.Rel(root, name)
				if relErr != nil {
					return relErr
				}
				discovered = append(discovered, rel)
				return nil
			})
			if walkErr != nil {
				return nil, fmt.Errorf("walk: traverse %q: %w", path, walkErr)
			}
		}
	}
	discovered = append(discovered, explicitFiles...)

	unique := make(map[string]struct{}, len(discovered))
	files := make([]string, 0, len(discovered))
	for _, file := range discovered {
		name := filepath.ToSlash(filepath.Clean(file))
		name = strings.TrimPrefix(name, "./")
		if _, exists := unique[name]; exists || matchesAny(w.Exclude, name) {
			continue
		}
		binary, readErr := binaryFile(filepath.Join(root, filepath.FromSlash(name)))
		if readErr != nil {
			if contains(explicitFiles, file) {
				return nil, fmt.Errorf("walk: read named path %q: %w", file, readErr)
			}
			continue
		}
		if binary {
			continue
		}
		unique[name] = struct{}{}
		files = append(files, name)
	}
	sort.Strings(files)
	return files, nil
}

func binaryFile(name string) (bool, error) {
	file, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer file.Close()
	buffer := make([]byte, sniffSize)
	n, err := io.ReadFull(file, buffer)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return false, err
	}
	return bytes.IndexByte(buffer[:n], 0) >= 0, nil
}

func contains(files []string, target string) bool {
	cleanTarget := filepath.Clean(target)
	for _, file := range files {
		if filepath.Clean(file) == cleanTarget {
			return true
		}
	}
	return false
}

func matchesAny(patterns []string, name string) bool {
	for _, pattern := range patterns {
		if matched, _ := doublestar.Match(filepath.ToSlash(pattern), name); matched {
			return true
		}
	}
	return false
}

func Candidates(files []string, r rules.Rule) []string {
	var candidates []string
	for _, file := range files {
		name := strings.TrimPrefix(filepath.ToSlash(file), "/")
		if (len(r.Include) == 0 || matchesAny(r.Include, name)) && !matchesAny(r.Exclude, name) {
			candidates = append(candidates, file)
		}
	}
	return candidates
}
