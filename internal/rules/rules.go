package rules

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
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

var validID = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

//go:embed builtin/*.md
var builtInFiles embed.FS

func Parse(name string, src []byte) (Rule, error) {
	id := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
	if !validID.MatchString(id) {
		return Rule{}, ruleError(name, "invalid rule id %q", id)
	}

	line, offset, ok := readLine(src, 0)
	if !ok || line != "---" {
		return Rule{}, ruleError(name, "malformed frontmatter: missing opening fence")
	}

	rule := Rule{ID: id, Exclude: []string{}}
	seen := make(map[string]bool)
	closed := false
	for offset <= len(src) {
		line, next, ok := readLine(src, offset)
		if !ok {
			break
		}
		offset = next
		if line == "---" {
			closed = true
			break
		}
		key, value, found := strings.Cut(line, ":")
		key = strings.TrimSpace(key)
		if !found || key == "" || seen[key] {
			return Rule{}, ruleError(name, "malformed frontmatter line %q", line)
		}
		seen[key] = true
		value = strings.TrimSpace(value)
		switch key {
		case "description":
			rule.Description = value
		case "severity":
			rule.Severity = Severity(value)
		case "include":
			if err := parseStringArray(value, &rule.Include); err != nil {
				return Rule{}, ruleError(name, "invalid include: %v", err)
			}
		case "exclude":
			if err := parseStringArray(value, &rule.Exclude); err != nil {
				return Rule{}, ruleError(name, "invalid exclude: %v", err)
			}
		default:
			return Rule{}, ruleError(name, "unknown frontmatter key %q", key)
		}
	}
	if !closed {
		return Rule{}, ruleError(name, "unterminated frontmatter")
	}
	if rule.Severity != SeverityError && rule.Severity != SeverityWarning {
		return Rule{}, ruleError(name, "invalid severity %q", rule.Severity)
	}
	if len(rule.Include) == 0 {
		return Rule{}, ruleError(name, "include must contain at least one pattern")
	}
	rule.Prompt = string(src[offset:])
	return rule, nil
}

func BuiltIns() []Rule {
	entries, err := fs.ReadDir(builtInFiles, "builtin")
	if err != nil {
		panic(fmt.Sprintf("rules: read embedded catalog: %v", err))
	}
	rules := make([]Rule, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		path := "builtin/" + entry.Name()
		src, err := builtInFiles.ReadFile(path)
		if err != nil {
			panic(fmt.Sprintf("rules: read embedded rule %s: %v", path, err))
		}
		rule, err := Parse(entry.Name(), src)
		if err != nil {
			panic(fmt.Sprintf("rules: parse embedded rule %s: %v", path, err))
		}
		rule.BuiltIn = true
		rules = append(rules, rule)
	}
	return rules
}

func Load(root string, paths []string) ([]Rule, error) {
	seen := make(map[string]string)
	for _, rule := range BuiltIns() {
		seen[rule.ID] = "built-in catalog"
	}

	var loaded []Rule
	for _, path := range paths {
		full := path
		if !filepath.IsAbs(full) {
			full = filepath.Join(root, full)
		}
		info, err := os.Stat(full)
		if err != nil {
			return nil, ruleError(path, "cannot read path: %v", err)
		}
		files := []string{full}
		if info.IsDir() {
			entries, err := os.ReadDir(full)
			if err != nil {
				return nil, ruleError(path, "cannot read directory: %v", err)
			}
			files = files[:0]
			for _, entry := range entries {
				if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
					files = append(files, filepath.Join(full, entry.Name()))
				}
			}
		}
		for _, file := range files {
			src, err := os.ReadFile(file)
			if err != nil {
				return nil, ruleError(file, "cannot read file: %v", err)
			}
			rule, err := Parse(file, src)
			if err != nil {
				return nil, err
			}
			if previous, exists := seen[rule.ID]; exists {
				return nil, ruleError(file, "rule id %q collides with %s", rule.ID, previous)
			}
			seen[rule.ID] = file
			loaded = append(loaded, rule)
		}
	}
	return loaded, nil
}

func Select(all []Rule, enable []string) ([]Rule, error) {
	known := make(map[string]Rule, len(all))
	for _, rule := range all {
		known[rule.ID] = rule
	}
	selected := make([]Rule, 0, len(enable))
	for _, id := range enable {
		rule, ok := known[id]
		if !ok {
			return nil, ruleError(id, "unknown enabled rule")
		}
		selected = append(selected, rule)
	}
	return selected, nil
}

func parseStringArray(value string, target *[]string) error {
	var parsed []string
	if err := json.Unmarshal([]byte(value), &parsed); err != nil {
		return err
	}
	if parsed == nil {
		return errors.New("must be a JSON string array")
	}
	*target = parsed
	return nil
}

func readLine(src []byte, offset int) (string, int, bool) {
	if offset >= len(src) {
		return "", offset, false
	}
	rest := src[offset:]
	if newline := strings.IndexByte(string(rest), '\n'); newline >= 0 {
		line := strings.TrimSuffix(string(rest[:newline]), "\r")
		return line, offset + newline + 1, true
	}
	return strings.TrimSuffix(string(rest), "\r"), len(src), true
}

func ruleError(name, format string, args ...any) error {
	reason := fmt.Sprintf(format, args...)
	return fmt.Errorf("%w: %s: %s", ErrRule, name, reason)
}
