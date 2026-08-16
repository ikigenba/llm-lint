package report

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
)

func TestVerboseSortsAndFormatsTraceEntries(t *testing.T) {
	// R-H7O3-875F
	root := filepath.Join(t.TempDir(), "repo")
	cwd := filepath.Join(root, "subdir")
	entries := []engine.TraceEntry{
		{File: "z.go", Rule: "beta", Outcome: "pass"},
		{File: "a.go", Rule: "zeta", Outcome: "pass"},
		{File: "a.go", Rule: "alpha", Cached: true, Outcome: "fail"},
	}
	var got bytes.Buffer
	Verbose(&got, cwd, root, entries)
	want := "../a.go: alpha hit fail\n../a.go: zeta miss pass\n../z.go: beta miss pass\n"
	if got.String() != want {
		t.Fatalf("Verbose() = %q, want %q", got.String(), want)
	}
}
