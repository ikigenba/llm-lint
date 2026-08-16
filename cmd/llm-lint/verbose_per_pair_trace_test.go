package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
	"github.com/ikigenba/llm-lint/internal/rules"
)

type perPairClient struct {
	mu sync.Mutex
}

func (c *perPairClient) Judge(_ context.Context, rule rules.Rule, file string, _ []byte) ([]engine.Finding, engine.Usage, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if rule.ID != "no-sleep-in-tests" {
		return nil, engine.Usage{}, nil
	}
	return []engine.Finding{{Rule: rule.ID, Severity: rule.Severity, File: file, Line: 1, Evidence: "source", Explanation: "found issue"}}, engine.Usage{}, nil
}

func verboseLintTree(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	for name, content := range map[string]string{
		"alpha_test.go": "alpha\n",
		"beta_test.go":  "beta\n",
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

func TestRunVerboseWritesPerPairTraceOnlyToStderr(t *testing.T) {
	// R-H1KL-BCFY
	bindClient(t, &perPairClient{})
	root, rulePath := verboseLintTree(t)
	var out, errOut bytes.Buffer
	code := run([]string{"--verbose", "--rules", rulePath}, bytes.NewReader(nil), &out, &errOut, noEnv, root)
	if code != 1 {
		t.Fatalf("run() code = %d, stderr = %q", code, errOut.String())
	}
	wantTrace := []string{
		"alpha_test.go: extra-audit miss pass",
		"alpha_test.go: no-sleep-in-tests miss fail",
		"beta_test.go: extra-audit miss pass",
		"beta_test.go: no-sleep-in-tests miss fail",
	}
	if got := strings.Split(strings.TrimSpace(errOut.String()), "\n"); !equalStrings(got, wantTrace) {
		t.Fatalf("verbose stderr = %#v, want %#v", got, wantTrace)
	}
	if strings.Count(out.String(), "found issue") != 2 || strings.Contains(out.String(), " miss ") {
		t.Fatalf("findings stdout = %q; want two findings and no trace", out.String())
	}

	out.Reset()
	errOut.Reset()
	if code := run([]string{"--no-cache", "--rules", rulePath}, bytes.NewReader(nil), &out, &errOut, noEnv, root); code != 1 {
		t.Fatalf("non-verbose run() code = %d, stderr = %q", code, errOut.String())
	}
	if errOut.Len() != 0 {
		t.Fatalf("non-verbose stderr = %q, want no audit trace", errOut.String())
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
