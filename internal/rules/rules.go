package rules

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

type Rule struct {
	ID          string
	Description string
	Severity    Severity
	Include     []string
	Exclude     []string
	Prompt      string
	BuiltIn     bool
}

var ErrRule = errors.New("rules: invalid rule")

//go:embed builtin/*.md
var builtInFiles embed.FS

func Parse(name string, src []byte) (Rule, error) {
	id := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
	if id == "" || len(src) == 0 {
		return Rule{}, fmt.Errorf("%w: %s is empty", ErrRule, name)
	}
	return Rule{ID: id, Description: "built-in lint rule", Severity: SeverityError, Include: []string{"**/*"}, Prompt: string(src)}, nil
}

func BuiltIns() []Rule {
	src, err := builtInFiles.ReadFile("builtin/clarity.md")
	if err != nil {
		panic("rules: embedded built-in missing")
	}
	rule, err := Parse("clarity.md", src)
	if err != nil {
		panic("rules: embedded built-in invalid")
	}
	rule.BuiltIn = true
	rule.Description = "flags unclear language"
	return []Rule{rule}
}

func Load(root string, paths []string) ([]Rule, error) {
	var loaded []Rule
	for _, path := range paths {
		full := path
		if !filepath.IsAbs(full) {
			full = filepath.Join(root, full)
		}
		info, err := os.Stat(full)
		if err != nil {
			return nil, fmt.Errorf("%w: %s: %v", ErrRule, path, err)
		}
		if info.IsDir() {
			continue
		}
		src, err := os.ReadFile(full)
		if err != nil {
			return nil, fmt.Errorf("%w: %s: %v", ErrRule, path, err)
		}
		rule, err := Parse(path, src)
		if err != nil {
			return nil, err
		}
		loaded = append(loaded, rule)
	}
	return loaded, nil
}

func Select(all []Rule, enable []string) ([]Rule, error) {
	if len(enable) == 0 {
		return append([]Rule(nil), all...), nil
	}
	wanted := make(map[string]bool, len(enable))
	for _, id := range enable {
		wanted[id] = true
	}
	var selected []Rule
	for _, rule := range all {
		if wanted[rule.ID] {
			selected = append(selected, rule)
			delete(wanted, rule.ID)
		}
	}
	for id := range wanted {
		return nil, fmt.Errorf("%w: unknown rule %q", ErrRule, id)
	}
	return selected, nil
}
