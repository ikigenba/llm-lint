package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
	"github.com/ikigenba/llm-lint/internal/rules"
)

// R-FTKK-67HP
func TestRunWarningFindingPrintsAndReturnsZero(t *testing.T) {
	bindClient(t, &fakeClient{findings: []engine.Finding{{Rule: "clarity", Severity: rules.SeverityWarning, File: "sample.txt", Line: 1, Explanation: "consider revising"}}})
	var out, errOut bytes.Buffer
	code := run(nil, bytes.NewReader(nil), &out, &errOut, noEnv, lintTree(t))
	if code != 0 || !strings.Contains(out.String(), "consider revising") {
		t.Fatalf("run() = code %d, stdout %q; want 0 and printed warning", code, out.String())
	}
}
