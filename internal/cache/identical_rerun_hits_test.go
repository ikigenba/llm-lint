package cache

import (
	"context"
	"reflect"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
	"github.com/ikigenba/llm-lint/internal/rules"
)

type fakeClient struct {
	calls    []string
	findings []engine.Finding
}

func (f *fakeClient) Judge(_ context.Context, r rules.Rule, file string, _ []byte) ([]engine.Finding, error) {
	f.calls = append(f.calls, r.ID+":"+file)
	return append([]engine.Finding(nil), f.findings...), nil
}

func testRule(id, prompt string) rules.Rule {
	return rules.Rule{ID: id, Severity: rules.SeverityWarning, Include: []string{"*.go"}, Prompt: prompt}
}

func TestIdenticalRerunHits(t *testing.T) {
	// R-H1GX-XY6P
	dir := t.TempDir()
	warmNext := &fakeClient{findings: []engine.Finding{{Rule: "style", File: "a.go", Line: 3, Evidence: "x"}}}
	warm := &CachingClient{Store: &Store{Dir: dir}, Next: warmNext}
	want, err := warm.Judge(context.Background(), testRule("style", "be clear"), "a.go", []byte("package a"))
	if err != nil || len(warmNext.calls) != 1 {
		t.Fatalf("warm Judge() error = %v, calls = %v", err, warmNext.calls)
	}

	runNext := &fakeClient{}
	run := &CachingClient{Store: &Store{Dir: dir}, Next: runNext}
	got, err := run.Judge(context.Background(), testRule("style", "be clear"), "a.go", []byte("package a"))
	if err != nil {
		t.Fatalf("cached Judge() error = %v", err)
	}
	if len(runNext.calls) != 0 {
		t.Fatalf("cached Judge() made calls %v, want none", runNext.calls)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("cached findings = %#v, want %#v", got, want)
	}
}
