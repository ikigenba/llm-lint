package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ikigenba/llm-lint/internal/config"
)

// R-TIL4-XZAL
func TestRunHelpSeparatesSectionsWithOneBlankLine(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"--help"}, bytes.NewReader(nil), &out, &errOut, noEnv, t.TempDir()); code != 0 || errOut.Len() != 0 {
		t.Fatalf("run(--help) = code %d, stderr %q; want 0 and empty stderr", code, errOut.String())
	}

	help := out.String()
	if strings.Contains(help, "\n\n\n") {
		t.Fatalf("help contains more than one blank line between sections:\n%s", help)
	}
	blocks := strings.Split(strings.TrimSuffix(help, "\n"), "\n\n")
	wantHeaders := []string{
		"usage: llm-lint [flags] [path ...]",
		"flags:",
		"defaults:",
		"providers:",
	}
	for _, provider := range config.Providers() {
		wantHeaders = append(wantHeaders, string(provider))
	}
	wantHeaders = append(wantHeaders, "bare model= derives its provider and must name a catalogued model")
	if len(blocks) != len(wantHeaders) {
		t.Fatalf("help has %d blank-line-separated sections; want %d:\n%s", len(blocks), len(wantHeaders), help)
	}
	for i, block := range blocks {
		header, _, _ := strings.Cut(block, "\n")
		if header != wantHeaders[i] {
			t.Errorf("section %d starts with %q; want %q", i, header, wantHeaders[i])
		}
	}

	out.Reset()
	errOut.Reset()
	if code := run([]string{"--not-a-flag"}, bytes.NewReader(nil), &out, &errOut, noEnv, t.TempDir()); code != 2 {
		t.Fatalf("run(bad flag) = code %d; want 2", code)
	}
	stderr := errOut.String()
	if strings.Count(stderr, "usage: llm-lint [flags] [path ...]") != 1 || strings.Contains(stderr, "\n\n") || strings.Contains(stderr, "flags:") {
		t.Errorf("bad-flag stderr includes help sections or invalid usage separation:\n%s", stderr)
	}
}
