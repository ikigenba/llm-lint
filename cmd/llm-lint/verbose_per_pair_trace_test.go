package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
	"github.com/ikigenba/llm-lint/internal/rules"
)

type perPairClient struct {
	blockBeta <-chan struct{}
}

func (c *perPairClient) Judge(_ context.Context, rule rules.Rule, file string, _ []byte) ([]engine.Finding, engine.Usage, error) {
	if c.blockBeta != nil && strings.Contains(file, "beta_test.go") {
		<-c.blockBeta
	}
	if rule.ID != "no-sleep-in-tests" {
		return nil, engine.Usage{}, nil
	}
	return []engine.Finding{{Rule: rule.ID, Severity: rule.Severity, File: file, Line: 1, Evidence: "source", Explanation: "found issue"}}, engine.Usage{}, nil
}

func verboseLintTree(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	for name, content := range map[string]string{
		"alpha_test.go": "alpha " + root + "\n",
		"beta_test.go":  "beta " + root + "\n",
	} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	rulePath := filepath.Join(root, "extra-audit.md")
	rule := "---\ndescription: Extra audit\nseverity: warning\ninclude: [\"**/*_test.go\"]\n---\nJudge the file.\n"
	if err := os.WriteFile(rulePath, []byte(rule), 0o600); err != nil {
		t.Fatal(err)
	}
	return root, rulePath
}

type lineCapture struct {
	mu     sync.Mutex
	buf    bytes.Buffer
	writes chan struct{}
}

func (w *lineCapture) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	n, err := w.buf.Write(p)
	select {
	case w.writes <- struct{}{}:
	default:
	}
	return n, err
}

func (w *lineCapture) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.String()
}

func TestRunVerboseWritesPerPairTraceOnlyToStderr(t *testing.T) {
	// R-H1KL-BCFY
	root, rulePath := verboseLintTree(t)
	release := make(chan struct{})
	bindClient(t, &perPairClient{blockBeta: release})
	var out bytes.Buffer
	errOut := &lineCapture{writes: make(chan struct{}, 8)}
	done := make(chan int, 1)
	go func() {
		done <- run([]string{"--verbose", "--concurrency", "4", "--rules", rulePath}, bytes.NewReader(nil), &out, errOut, noEnv, root)
	}()
	select {
	case <-errOut.writes:
		select {
		case code := <-done:
			t.Fatalf("run() completed with code %d before blocked pairs were released", code)
		default:
		}
	case code := <-done:
		t.Fatalf("run() completed with code %d before streaming a trace line", code)
	}
	close(release)
	code := <-done
	if code != 1 {
		t.Fatalf("run() code = %d, stderr = %q", code, errOut.String())
	}
	wantTrace := []string{
		"🟢 alpha_test.go extra-audit",
		"🔴 alpha_test.go no-sleep-in-tests",
		"🟢 beta_test.go extra-audit",
		"🔴 beta_test.go no-sleep-in-tests",
	}
	gotTrace := strings.Split(strings.TrimSpace(errOut.String()), "\n")
	sort.Strings(gotTrace)
	sort.Strings(wantTrace)
	if !equalStrings(gotTrace, wantTrace) {
		t.Fatalf("verbose stderr = %#v, want %#v", gotTrace, wantTrace)
	}
	if strings.Count(out.String(), "found issue") != 2 || strings.Contains(out.String(), "🟢") || strings.Contains(out.String(), "🔴") {
		t.Fatalf("findings stdout = %q; want two findings and no trace", out.String())
	}

	verboseOut := out.String()
	var plainOut, plainErr bytes.Buffer
	bindClient(t, &perPairClient{})
	if code := run([]string{"--no-cache", "--concurrency", "4", "--rules", rulePath}, bytes.NewReader(nil), &plainOut, &plainErr, noEnv, root); code != 1 {
		t.Fatalf("non-verbose run() code = %d, stderr = %q", code, plainErr.String())
	}
	if plainErr.Len() != 0 || plainOut.String() != verboseOut {
		t.Fatalf("non-verbose stdout/stderr = %q/%q, want unchanged stdout %q and no audit trace", plainOut.String(), plainErr.String(), verboseOut)
	}
}

func equalStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
