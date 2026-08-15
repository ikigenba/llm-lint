package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
	"github.com/ikigenba/llm-lint/internal/rules"
)

func TestRunWiresTreeThroughPipelineToReport(t *testing.T) {
	client := &fakeClient{findings: []engine.Finding{{
		Rule:        "no-sleep-in-tests",
		Severity:    rules.SeverityError,
		File:        "sample_test.go",
		Line:        1,
		Evidence:    "sample",
		Explanation: "remove the sleep",
	}}}
	bindClient(t, client)
	root := lintTree(t)
	var out, errOut bytes.Buffer
	code := run(nil, bytes.NewReader(nil), &out, &errOut, noEnv, root)
	if code != 1 || client.calls != 1 || !strings.Contains(out.String(), "sample_test.go:1: remove the sleep (no-sleep-in-tests)") {
		t.Fatalf("run() = code %d, calls %d, stdout %q, stderr %q; pipeline did not report finding", code, client.calls, out.String(), errOut.String())
	}
}
