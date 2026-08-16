package report

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"sort"

	"github.com/ikigenba/llm-lint/internal/engine"
	"github.com/ikigenba/llm-lint/internal/rules"
)

type Stats struct {
	Rules, Files, Pairs, Calls, CacheHits int
	InputTokens, OutputTokens             int64
	CostUSD                               float64
}

func Text(w io.Writer, cwd, root string, findings []engine.Finding) error {
	for _, finding := range sorted(findings) {
		name := finding.File
		path := finding.File
		if !filepath.IsAbs(path) {
			path = filepath.Join(root, filepath.FromSlash(path))
		}
		if rel, err := filepath.Rel(cwd, path); err == nil {
			name = rel
		}
		if _, err := fmt.Fprintf(w, "%s:%d: %s (%s)\n", name, finding.Line, finding.Explanation, finding.Rule); err != nil {
			return fmt.Errorf("report: %w", err)
		}
	}
	return nil
}

func JSON(w io.Writer, findings []engine.Finding) error {
	type findingJSON struct {
		Rule        string         `json:"rule"`
		Severity    rules.Severity `json:"severity"`
		File        string         `json:"file"`
		Line        int            `json:"line"`
		Evidence    string         `json:"evidence"`
		Explanation string         `json:"explanation"`
	}

	encoded := make([]findingJSON, 0, len(findings))
	for _, finding := range sorted(findings) {
		encoded = append(encoded, findingJSON{
			Rule:        finding.Rule,
			Severity:    finding.Severity,
			File:        finding.File,
			Line:        finding.Line,
			Evidence:    finding.Evidence,
			Explanation: finding.Explanation,
		})
	}
	if err := json.NewEncoder(w).Encode(encoded); err != nil {
		return fmt.Errorf("report: %w", err)
	}
	return nil
}

func StatsLine(w io.Writer, s Stats) {
	fmt.Fprintf(w, "llm-lint: %d rules, %d files, %d pairs, %d calls, %d cache hits, %s in / %s out tokens, $%.4f\n", s.Rules, s.Files, s.Pairs, s.Calls, s.CacheHits, tokenCount(s.InputTokens), tokenCount(s.OutputTokens), s.CostUSD)
}

func Verbose(w io.Writer, cwd, root string, entries []engine.TraceEntry) {
	ordered := append([]engine.TraceEntry(nil), entries...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].File != ordered[j].File {
			return ordered[i].File < ordered[j].File
		}
		return ordered[i].Rule < ordered[j].Rule
	})
	for _, entry := range ordered {
		name := entry.File
		path := entry.File
		if !filepath.IsAbs(path) {
			path = filepath.Join(root, filepath.FromSlash(path))
		}
		if rel, err := filepath.Rel(cwd, path); err == nil {
			name = rel
		}
		cacheStatus := "miss"
		if entry.Cached {
			cacheStatus = "hit"
		}
		fmt.Fprintf(w, "%s: %s %s %s\n", name, entry.Rule, cacheStatus, entry.Outcome)
	}
}

func sorted(findings []engine.Finding) []engine.Finding {
	ordered := append([]engine.Finding(nil), findings...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].File != ordered[j].File {
			return ordered[i].File < ordered[j].File
		}
		if ordered[i].Line != ordered[j].Line {
			return ordered[i].Line < ordered[j].Line
		}
		return ordered[i].Rule < ordered[j].Rule
	})
	return ordered
}

func tokenCount(n int64) string {
	if n >= 1000 {
		return fmt.Sprintf("%dk", n/1000)
	}
	return fmt.Sprintf("%d", n)
}

func ExitCode(findings []engine.Finding) int {
	for _, finding := range findings {
		if finding.Severity == rules.SeverityError {
			return 1
		}
	}
	return 0
}
