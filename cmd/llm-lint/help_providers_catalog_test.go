package main

import (
	"bytes"
	"strings"
	"testing"
)

// R-J6TT-EEVY
func TestRunHelpPrintsProviderCatalogWithoutInference(t *testing.T) {
	client := &fakeClient{}
	bindClient(t, client)

	var out, errOut bytes.Buffer
	code := run([]string{"--help"}, bytes.NewReader(nil), &out, &errOut, noEnv, t.TempDir())
	if code != 0 || errOut.Len() != 0 {
		t.Fatalf("run(--help) = code %d, stderr %q; want 0 and empty stderr", code, errOut.String())
	}
	for _, want := range []string{
		"Usage: llm-lint",
		"-c key=value",
		"providers:",
		"google\n",
		"auth=key (GEMINI_API_KEY)",
		"openai\n",
		"auth=sub (auth_file=~/.llm-lint/openai-auth.json)",
		"x-ai\n",
		"auth=sub (auth_file=~/.llm-lint/x-ai-auth.json)",
		"gemini-3.7-flash (reasoning ",
		"*",
		"bare model= derives its provider and must name a catalogued model",
		"explicit provider= accepts any model as pass-through",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("help output does not contain %q:\n%s", want, out.String())
		}
	}
	if client.calls != 0 {
		t.Fatalf("Judge calls = %d, want 0", client.calls)
	}
}
