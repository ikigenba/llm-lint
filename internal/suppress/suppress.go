package suppress

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/ikigenba/llm-lint/internal/engine"
)

const marker = "llm-lint:ignore"

type directive struct {
	all   bool
	rules map[string]struct{}
}

func Filter(findings []engine.Finding, read func(path string) ([]byte, error)) ([]engine.Finding, error) {
	directivesByFile := make(map[string]map[int][]directive)
	for _, finding := range findings {
		if _, ok := directivesByFile[finding.File]; ok {
			continue
		}

		contents, err := read(finding.File)
		if err != nil {
			return nil, fmt.Errorf("suppress: read %q: %w", finding.File, err)
		}
		directivesByFile[finding.File] = parseDirectives(string(contents))
	}

	filtered := make([]engine.Finding, 0, len(findings))
	for _, finding := range findings {
		if !isSuppressed(finding.Rule, directivesByFile[finding.File][finding.Line]) {
			filtered = append(filtered, finding)
		}
	}
	return filtered, nil
}

func parseDirectives(contents string) map[int][]directive {
	byLine := make(map[int][]directive)
	for index, line := range strings.Split(contents, "\n") {
		markerIndex := strings.Index(line, marker)
		if markerIndex < 0 {
			continue
		}

		directive, ok := parseDirective(line[markerIndex+len(marker):])
		if !ok {
			continue
		}
		targetLine := index + 1
		if commentOnly(line[:markerIndex]) {
			targetLine++
		}
		byLine[targetLine] = append(byLine[targetLine], directive)
	}
	return byLine
}

func parseDirective(afterMarker string) (directive, bool) {
	tail := strings.TrimSpace(afterMarker)
	first, last := firstAndLastAlphanumeric(tail)
	if first < 0 {
		return directive{all: true}, true
	}
	if strings.Contains(tail[:first], ",") || strings.Contains(tail[last+1:], ",") {
		return directive{}, false
	}

	rules := make(map[string]struct{})
	for _, item := range strings.Split(tail[first:last+1], ",") {
		rule := strings.TrimSpace(item)
		if rule == "" || !validRuleID(rule) {
			return directive{}, false
		}
		rules[rule] = struct{}{}
	}
	return directive{rules: rules}, true
}

func firstAndLastAlphanumeric(value string) (int, int) {
	first, last := -1, -1
	for index, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if first < 0 {
				first = index
			}
			last = index + len(string(r)) - 1
		}
	}
	return first, last
}

func validRuleID(value string) bool {
	for _, r := range value {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' && r != '.' {
			return false
		}
	}
	return true
}

func commentOnly(beforeMarker string) bool {
	for _, r := range beforeMarker {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func isSuppressed(rule string, directives []directive) bool {
	for _, directive := range directives {
		if directive.all {
			return true
		}
		if _, ok := directive.rules[rule]; ok {
			return true
		}
	}
	return false
}
