package report

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
)

// R-HDNX-RNLN
func TestTextRendersPathRelativeToInvokingDirectory(t *testing.T) {
	root := t.TempDir()
	cwd := filepath.Join(root, "tools")
	var out bytes.Buffer
	findings := []engine.Finding{{Rule: "clarity", File: "pkg/sample.go", Line: 7, Explanation: "make this clearer"}}

	if err := Text(&out, cwd, root, findings); err != nil {
		t.Fatal(err)
	}
	want := filepath.Join("..", "pkg", "sample.go") + ":7: make this clearer (clarity)\n"
	if out.String() != want {
		t.Fatalf("Text() = %q, want %q", out.String(), want)
	}
}
