package report

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"github.com/ikigenba/llm-lint/internal/engine"
	"github.com/ikigenba/llm-lint/internal/rules"
)

type Stats struct {
	Rules, Files, Pairs, Calls, CacheHits int
	InputTokens, OutputTokens             int64
	CostUSD                               float64
}

func Text(w io.Writer, cwd, root string, findings []engine.Finding) error {
	for _, finding := range findings {
		name := finding.File
		if rel, err := filepath.Rel(cwd, filepath.Join(root, finding.File)); err == nil {
			name = rel
		}
		if _, err := fmt.Fprintf(w, "%s:%d: %s %s: %s\n", name, finding.Line, finding.Severity, finding.Rule, finding.Explanation); err != nil {
			return fmt.Errorf("report: %w", err)
		}
	}
	return nil
}

func JSON(w io.Writer, findings []engine.Finding) error {
	if err := json.NewEncoder(w).Encode(findings); err != nil {
		return fmt.Errorf("report: %w", err)
	}
	return nil
}

func StatsLine(w io.Writer, s Stats) {
	fmt.Fprintf(w, "rules=%d files=%d pairs=%d calls=%d cache_hits=%d input_tokens=%d output_tokens=%d cost_usd=%.6f\n", s.Rules, s.Files, s.Pairs, s.Calls, s.CacheHits, s.InputTokens, s.OutputTokens, s.CostUSD)
}

func ExitCode(findings []engine.Finding) int {
	for _, finding := range findings {
		if finding.Severity == rules.SeverityError {
			return 1
		}
	}
	return 0
}
