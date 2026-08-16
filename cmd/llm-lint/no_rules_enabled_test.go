package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// R-2P4O-HGWM
func TestRunWithEveryActiveRuleDisabledWarnsAndSucceeds(t *testing.T) {
	client := &fakeClient{}
	bindClient(t, client)
	dir := t.TempDir()
	config := []byte(`{"disable":["boolean-state-machine","no-sleep-in-tests"]}`)
	if err := os.WriteFile(filepath.Join(dir, ".llm-lint.json"), config, 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	code := run(nil, bytes.NewReader(nil), &out, &errOut, noEnv, dir)
	if code != 0 || out.Len() != 0 || client.calls != 0 || strings.TrimSpace(errOut.String()) != "llm-lint: no rules enabled" {
		t.Fatalf("run() = code %d, stdout %q, stderr %q, calls %d; want 0, empty stdout, warning, zero calls", code, out.String(), errOut.String(), client.calls)
	}
}
