package main

import (
	"bytes"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
	"github.com/ikigenba/llm-lint/internal/rules"
)

// R-FSCN-SFR0
func TestRunErrorFindingReturnsOne(t *testing.T) {
	bindClient(t, &fakeClient{findings: []engine.Finding{{Rule: "clarity", Severity: rules.SeverityError, File: "sample_test.go", Line: 1, Explanation: "unclear"}}})
	var out, errOut bytes.Buffer
	if code := run(nil, bytes.NewReader(nil), &out, &errOut, noEnv, lintTree(t)); code != 1 {
		t.Fatalf("run() = %d, want 1; stdout %q; stderr %q", code, out.String(), errOut.String())
	}
}
