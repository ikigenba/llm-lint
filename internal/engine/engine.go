package engine

import (
	"context"
	"fmt"
	"io"
	"sort"
	"sync"

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

type Usage struct {
	Input, Output int64
	CostUSD       float64
}

type Client interface {
	Judge(ctx context.Context, r rules.Rule, file string, content []byte) ([]Finding, Usage, error)
}

type TraceEntry struct {
	File, Rule string
	Cached     bool
	Outcome    string
}

type TraceSink interface {
	Add(TraceEntry)
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

func (e *Engine) Run(ctx context.Context, rs []rules.Rule, files map[string][]string, read func(path string) ([]byte, error), warn io.Writer) ([]Finding, Stats, error) {
	if e.Client == nil {
		panic("engine: nil Client")
	}
	if e.Concurrency < 1 {
		panic("engine: Concurrency must be at least 1")
	}
	if warn == nil {
		warn = io.Discard
	}

	type pair struct {
		rule rules.Rule
		file string
	}
	var pairs []pair
	stats := Stats{Rules: len(rs)}
	seen := make(map[string]bool)
	for _, rule := range rs {
		for _, file := range files[rule.ID] {
			stats.Pairs++
			if !seen[file] {
				stats.Files++
				seen[file] = true
			}
			pairs = append(pairs, pair{rule: rule, file: file})
		}
	}

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	jobs := make(chan pair)
	type result struct {
		findings []Finding
		usage    Usage
		called   bool
		err      error
	}
	results := make(chan result)
	var workers sync.WaitGroup
	workerCount := min(e.Concurrency, len(pairs))
	for range workerCount {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for p := range jobs {
				content, err := read(p.file)
				if err != nil {
					sendResult(runCtx, results, result{err: fmt.Errorf("read %s: %w", p.file, err)})
					continue
				}
				got, usage, err := e.Client.Judge(runCtx, p.rule, p.file, content)
				if err != nil {
					err = fmt.Errorf("rule %s file %s: %w", p.rule.ID, p.file, err)
				}
				sendResult(runCtx, results, result{findings: got, usage: usage, called: true, err: err})
			}
		}()
	}
	go func() {
		defer close(jobs)
		for _, p := range pairs {
			select {
			case jobs <- p:
			case <-runCtx.Done():
				return
			}
		}
	}()
	go func() {
		workers.Wait()
		close(results)
	}()

	var findings []Finding
	var firstErr error
	for got := range results {
		if got.called {
			stats.Calls++
			stats.InputTokens += got.usage.Input
			stats.OutputTokens += got.usage.Output
			stats.CostUSD += got.usage.CostUSD
		}
		if got.err != nil && firstErr == nil {
			firstErr = got.err
			cancel()
		}
		if firstErr == nil {
			findings = append(findings, got.findings...)
		}
	}
	if firstErr != nil {
		return nil, stats, fmt.Errorf("engine: operational failure: %w", firstErr)
	}
	if err := ctx.Err(); err != nil {
		return nil, stats, err
	}
	sort.SliceStable(findings, func(i, j int) bool {
		if findings[i].File != findings[j].File {
			return findings[i].File < findings[j].File
		}
		if findings[i].Line != findings[j].Line {
			return findings[i].Line < findings[j].Line
		}
		return findings[i].Rule < findings[j].Rule
	})
	return findings, stats, nil
}

func sendResult[T any](ctx context.Context, dst chan<- T, result T) {
	select {
	case dst <- result:
	case <-ctx.Done():
	}
}
