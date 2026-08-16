package report

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
)

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("broken audit stream")
}

func TestVerboseSinkFormatsRelativePathAndSwallowsWriteErrors(t *testing.T) {
	root := filepath.Join(t.TempDir(), "repo")
	cwd := filepath.Join(root, "subdir")
	var got bytes.Buffer
	sink := NewVerboseSink(&got, cwd, root)
	sink.Add(engine.TraceEntry{File: "source.go", Rule: "audit", Cached: true, Outcome: "fail"})
	if want := "🔴 ../source.go [audit]\n"; got.String() != want {
		t.Fatalf("VerboseSink output = %q, want %q", got.String(), want)
	}

	NewVerboseSink(failingWriter{}, cwd, root).Add(engine.TraceEntry{File: "source.go", Rule: "audit", Outcome: "pass"})
}

func TestVerboseSinkSerializesConcurrentLines(t *testing.T) {
	root := t.TempDir()
	var got bytes.Buffer
	sink := NewVerboseSink(&got, root, root)
	const count = 64
	var adds sync.WaitGroup
	for i := range count {
		adds.Add(1)
		go func() {
			defer adds.Done()
			sink.Add(engine.TraceEntry{File: fmt.Sprintf("file-%02d.go", i), Rule: "audit", Outcome: "pass"})
		}()
	}
	adds.Wait()

	lines := strings.Split(strings.TrimSpace(got.String()), "\n")
	if len(lines) != count {
		t.Fatalf("VerboseSink wrote %d lines, want %d: %q", len(lines), count, got.String())
	}
	seen := make(map[string]bool, count)
	for _, line := range lines {
		seen[line] = true
	}
	for i := range count {
		want := fmt.Sprintf("🟢 file-%02d.go [audit]", i)
		if !seen[want] {
			t.Fatalf("VerboseSink output lacks complete line %q: %q", want, got.String())
		}
	}
}
