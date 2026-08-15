package engine

import (
	"context"
	"fmt"
	"io"

	"github.com/ikigenba/llm-lint/internal/rules"
)

type Finding struct {
	Rule        string
	Severity    rules.Severity
	File        string
	Line        int
	Evidence    string
	Explanation string
}

type Client interface {
	Judge(ctx context.Context, r rules.Rule, file string, content []byte) ([]Finding, error)
}

type Engine struct {
	Client      Client
	Concurrency int
}

type Stats struct {
	Rules, Files, Pairs, Calls int
	InputTokens, OutputTokens  int64
	CostUSD                    float64
}

type OpFailure struct{ Err error }

func (e *OpFailure) Error() string { return fmt.Sprintf("engine: operational failure: %v", e.Err) }
func (e *OpFailure) Unwrap() error { return e.Err }

func (e *Engine) Run(ctx context.Context, rs []rules.Rule, files map[string][]string, read func(path string) ([]byte, error), warn io.Writer) ([]Finding, Stats, error) {
	stats := Stats{Rules: len(rs)}
	seen := make(map[string]bool)
	var findings []Finding
	for _, rule := range rs {
		for _, file := range files[rule.ID] {
			stats.Pairs++
			if !seen[file] {
				stats.Files++
				seen[file] = true
			}
			content, err := read(file)
			if err != nil {
				return nil, stats, &OpFailure{Err: err}
			}
			stats.Calls++
			got, err := e.Client.Judge(ctx, rule, file, content)
			if err != nil {
				return nil, stats, &OpFailure{Err: err}
			}
			findings = append(findings, got...)
		}
	}
	_ = warn
	return findings, stats, nil
}
