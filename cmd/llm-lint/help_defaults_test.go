package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// R-4B8L-L2MN
func TestRunHelpPrintsCompiledDefaultsIndependentOfConfiguration(t *testing.T) {
	const defaults = "defaults: provider=google   model=gemini-3.7-flash   auth=key"
	cwd := t.TempDir()
	if err := os.WriteFile(filepath.Join(cwd, ".llm-lint.json"), []byte(`{"model":"gpt-5.4"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"--help"}, {"--help", "-c", "model=gpt-5.5"}} {
		var out, errOut bytes.Buffer
		if code := run(args, bytes.NewReader(nil), &out, &errOut, noEnv, cwd); code != 0 || errOut.Len() != 0 {
			t.Fatalf("run(%q) = code %d, stderr %q; want 0 and empty stderr", args, code, errOut.String())
		}
		if count := strings.Count(out.String(), defaults); count != 1 {
			t.Errorf("run(%q) contains compiled defaults %d times; want once:\n%s", args, count, out.String())
		}
	}
}
