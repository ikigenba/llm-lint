package main

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ikigenba/llm-lint/internal/config"
	"github.com/ikigenba/llm-lint/internal/engine"
	"github.com/ikigenba/llm-lint/internal/rules"
)

type fakeClient struct {
	findings []engine.Finding
	err      error
	usage    engine.Usage
	calls    int
}

func (f *fakeClient) Judge(context.Context, rules.Rule, string, []byte) ([]engine.Finding, engine.Usage, error) {
	f.calls++
	return f.findings, f.usage, f.err
}

func bindClient(t *testing.T, client engine.Client) {
	t.Helper()
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	old := newClient
	newClient = func(*config.Config, io.Writer) (engine.Client, error) { return client, nil }
	t.Cleanup(func() { newClient = old })
}

func lintTree(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "sample_test.go"), []byte("sample\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	config := []byte(`{"disable":["boolean-state-machine"]}`)
	if err := os.WriteFile(filepath.Join(dir, ".llm-lint.json"), config, 0o600); err != nil {
		t.Fatal(err)
	}
	return dir
}

func noEnv(string) string { return "" }
