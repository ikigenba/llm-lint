package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
)

// R-HG3Q-J731
func TestRunStatsFlagWritesSummaryOnlyToStderr(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		client := &fakeClient{usage: engine.Usage{Input: 1234, Output: 56, CostUSD: 0.0123}}
		bindClient(t, client)
		root := lintTree(t)
		if err := os.WriteFile(filepath.Join(root, "second_test.go"), []byte("sample\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		var out, errOut bytes.Buffer
		code := run([]string{"--stats", "--concurrency", "1"}, bytes.NewReader(nil), &out, &errOut, noEnv, root)
		want := "llm-lint: 1 rules, 2 files, 2 pairs, 2 calls, 1 cache hits, 1k in / 56 out tokens, $0.0123\n"
		if code != 0 || out.Len() != 0 || errOut.String() != want || client.calls != 1 {
			t.Fatalf("run() = code %d, stdout %q, stderr %q, provider calls %d; want 0, empty, %q, 1", code, out.String(), errOut.String(), client.calls, want)
		}
	})

	t.Run("absent", func(t *testing.T) {
		bindClient(t, &fakeClient{})
		var out, errOut bytes.Buffer
		if code := run(nil, bytes.NewReader(nil), &out, &errOut, noEnv, lintTree(t)); code != 0 {
			t.Fatalf("run() = %d, stderr %q", code, errOut.String())
		}
		if strings.Contains(errOut.String(), "llm-lint: 1 rules") {
			t.Fatalf("stderr contains stats without --stats: %q", errOut.String())
		}
	})
}
