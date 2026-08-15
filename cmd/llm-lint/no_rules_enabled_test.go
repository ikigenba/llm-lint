package main

import (
	"bytes"
	"strings"
	"testing"
)

// R-GHYJ-TMBL
func TestRunWithNoEnabledRulesWarnsAndSucceeds(t *testing.T) {
	client := &fakeClient{}
	bindClient(t, client)
	var out, errOut bytes.Buffer
	code := run(nil, bytes.NewReader(nil), &out, &errOut, noEnv, t.TempDir())
	if code != 0 || out.Len() != 0 || client.calls != 0 || strings.TrimSpace(errOut.String()) != "llm-lint: no rules enabled" {
		t.Fatalf("run() = code %d, stdout %q, stderr %q, calls %d; want 0, empty stdout, warning, zero calls", code, out.String(), errOut.String(), client.calls)
	}
}
