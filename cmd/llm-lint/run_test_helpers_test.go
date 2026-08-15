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
	calls    int
}

func (f *fakeClient) Judge(context.Context, rules.Rule, string, []byte) ([]engine.Finding, error) {
	f.calls++
	return f.findings, f.err
}

func bindClient(t *testing.T, client engine.Client) {
	t.Helper()
	old := newClient
	newClient = func(*config.Config, io.Writer) (engine.Client, error) { return client, nil }
	t.Cleanup(func() { newClient = old })
}

func lintTree(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "sample.txt"), []byte("sample\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return dir
}

func noEnv(string) string { return "" }
