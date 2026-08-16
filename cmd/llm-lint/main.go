package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/ikigenba/llm-lint/internal/cache"
	"github.com/ikigenba/llm-lint/internal/config"
	"github.com/ikigenba/llm-lint/internal/engine"
	"github.com/ikigenba/llm-lint/internal/report"
	"github.com/ikigenba/llm-lint/internal/rules"
	"github.com/ikigenba/llm-lint/internal/suppress"
	"github.com/ikigenba/llm-lint/internal/walk"
)

var version = "dev"

func main() {
	cwd, _ := os.Getwd()
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, os.Getenv, cwd))
}

type stringList []string

func (s *stringList) String() string { return strings.Join(*s, ",") }
func (s *stringList) Set(value string) error {
	*s = append(*s, value)
	return nil
}

var newClient = func(cfg *config.Config, errOut io.Writer) (engine.Client, error) {
	return &engine.AgentkitClient{
		NewConversation: cfg.NewConversation,
		Warn:            errOut,
	}, nil
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "Usage: llm-lint [options] [path...]")
	fmt.Fprintln(w, "Options: -c key=value, --rules path, --model id, --format text|json, --concurrency N, --no-cache, --stats, --list-rules, -V, --version, --help")
}

func run(args []string, in io.Reader, out, errOut io.Writer, getenv func(string) string, cwd string) int {
	fs := flag.NewFlagSet("llm-lint", flag.ContinueOnError)
	fs.SetOutput(errOut)
	fs.Usage = func() { usage(errOut) }
	var overrides, rulePaths stringList
	var model, format string
	var concurrency int
	var noCache, statsFlag, listRules, versionFlag, help bool
	fs.Var(&overrides, "c", "config override key=value (repeatable)")
	fs.Var(&rulePaths, "rules", "additional rule file or directory (repeatable)")
	fs.StringVar(&model, "model", "", "model id")
	fs.StringVar(&format, "format", "text", "output format: text or json")
	fs.IntVar(&concurrency, "concurrency", 8, "maximum in-flight inference calls")
	fs.BoolVar(&noCache, "no-cache", false, "skip cache reads")
	fs.BoolVar(&statsFlag, "stats", false, "print run statistics")
	fs.BoolVar(&listRules, "list-rules", false, "list known rules")
	fs.BoolVar(&versionFlag, "V", false, "print version")
	fs.BoolVar(&versionFlag, "version", false, "print version")
	fs.BoolVar(&help, "help", false, "print help")
	if err := fs.Parse(args); err != nil {
		usage(errOut)
		return 2
	}
	_ = in
	if help {
		usage(out)
		return 0
	}
	if versionFlag {
		fmt.Fprintln(out, version)
		return 0
	}
	if format != "text" && format != "json" {
		fmt.Fprintf(errOut, "llm-lint: unsupported format %q\n", format)
		return 2
	}
	if concurrency < 1 {
		fmt.Fprintln(errOut, "llm-lint: concurrency must be positive")
		return 2
	}
	if model != "" {
		overrides = append(overrides, "model="+model)
	}
	cfg, err := config.Load(cwd, overrides, getenv)
	if err != nil {
		fmt.Fprintf(errOut, "llm-lint: %v\n", err)
		if errors.Is(err, config.ErrAuth) {
			return 3
		}
		return 2
	}
	cfg.Rules = append(cfg.Rules, rulePaths...)
	allRules := rules.BuiltIns()
	localRules, err := rules.Load(cfg.Root, cfg.Rules)
	if err != nil {
		fmt.Fprintf(errOut, "llm-lint: %v\n", err)
		return 2
	}
	allRules = append(allRules, localRules...)
	selected, err := rules.Select(allRules, cfg.Enable, cfg.Disable)
	if err != nil {
		fmt.Fprintf(errOut, "llm-lint: %v\n", err)
		return 2
	}
	if listRules {
		enabled := make(map[string]bool, len(selected))
		for _, rule := range selected {
			enabled[rule.ID] = true
		}
		for _, rule := range allRules {
			state := "disabled"
			if enabled[rule.ID] {
				state = "enabled"
			}
			fmt.Fprintf(out, "%s\t%s\t%s\t%s\n", rule.ID, rule.Severity, state, rule.Description)
		}
		return 0
	}
	if len(selected) == 0 {
		fmt.Fprintln(errOut, "llm-lint: no rules enabled")
		return 0
	}

	paths := fs.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}
	walker := &walk.Walker{Root: cfg.Root, Exclude: cfg.Exclude, RunGit: runGit}
	files, err := walker.Files(paths)
	if err != nil {
		fmt.Fprintf(errOut, "llm-lint: %v\n", err)
		return 3
	}
	byRule := make(map[string][]string, len(selected))
	for _, rule := range selected {
		byRule[rule.ID] = walk.Candidates(files, rule)
	}
	client, err := newClient(cfg, errOut)
	if err != nil {
		fmt.Fprintf(errOut, "llm-lint: %v\n", err)
		return 3
	}
	var hits atomic.Int64
	client = &cache.CachingClient{Store: &cache.Store{}, Next: client, NoRead: noCache, Hits: &hits}
	runner := &engine.Engine{Client: client, Concurrency: concurrency}
	readFile := func(name string) ([]byte, error) {
		if !filepath.IsAbs(name) {
			name = filepath.Join(cfg.Root, filepath.FromSlash(name))
		}
		return os.ReadFile(name)
	}
	findings, engineStats, err := runner.Run(context.Background(), selected, byRule, readFile, errOut)
	if err != nil {
		fmt.Fprintf(errOut, "llm-lint: %v\n", err)
		return 3
	}
	findings, err = suppress.Filter(findings, readFile)
	if err != nil {
		fmt.Fprintf(errOut, "llm-lint: %v\n", err)
		return 3
	}
	if format == "json" {
		err = report.JSON(out, findings)
	} else if len(findings) > 0 {
		err = report.Text(out, cwd, cfg.Root, findings)
	}
	if err != nil {
		fmt.Fprintf(errOut, "llm-lint: %v\n", err)
		return 3
	}
	if statsFlag {
		report.StatsLine(errOut, report.Stats{Rules: engineStats.Rules, Files: engineStats.Files, Pairs: engineStats.Pairs, Calls: engineStats.Calls, CacheHits: int(hits.Load()), InputTokens: engineStats.InputTokens, OutputTokens: engineStats.OutputTokens, CostUSD: engineStats.CostUSD})
	}
	return report.ExitCode(findings)
}

func runGit(dir string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	return cmd.Output()
}
