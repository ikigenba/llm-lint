package main

import (
	"bytes"
	"strings"
	"testing"
)

// R-FX89-BIPS
func TestRunListRulesDescribesKnownRulesWithoutInference(t *testing.T) {
	client := &fakeClient{}
	bindClient(t, client)
	var out, errOut bytes.Buffer
	code := run([]string{"--list-rules"}, bytes.NewReader(nil), &out, &errOut, noEnv, t.TempDir())
	listing := strings.TrimSpace(out.String())
	for _, field := range []string{"no-sleep-in-tests", "error", "enabled", "fixed-duration sleeps"} {
		if !strings.Contains(listing, field) {
			t.Fatalf("rule listing %q does not contain %q", listing, field)
		}
	}
	if code != 0 || client.calls != 0 {
		t.Fatalf("run() = code %d, calls %d, listing %q; want 0 and zero calls", code, client.calls, listing)
	}
}

// R-HIJJ-AQKF
func TestRunListRulesIncludesNoSleepInTests(t *testing.T) {
	var out, errOut bytes.Buffer
	code := run([]string{"--list-rules"}, bytes.NewReader(nil), &out, &errOut, noEnv, t.TempDir())
	if code != 0 || !strings.Contains(out.String(), "no-sleep-in-tests\terror\tenabled\t") {
		t.Fatalf("run(--list-rules) = %d, output %q; want embedded no-sleep-in-tests listing", code, out.String())
	}
}

// R-4IQ0-3CA4
func TestRunListRulesIncludesBooleanStateMachine(t *testing.T) {
	var out, errOut bytes.Buffer
	code := run([]string{"--list-rules"}, bytes.NewReader(nil), &out, &errOut, noEnv, t.TempDir())
	if code != 0 || !strings.Contains(out.String(), "boolean-state-machine\terror\tenabled\t") {
		t.Fatalf("run(--list-rules) = %d, output %q; want embedded boolean-state-machine listing", code, out.String())
	}
}
