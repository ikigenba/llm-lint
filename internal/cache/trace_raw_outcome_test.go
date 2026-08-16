package cache

import (
	"context"
	"reflect"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
	"github.com/ikigenba/llm-lint/internal/rules"
	"github.com/ikigenba/llm-lint/internal/suppress"
)

func TestTraceRecordsRawVerdictOutcome(t *testing.T) {
	// R-H6G6-UFEQ
	rule := testRule("target-rule", "find target")
	finding := engine.Finding{Rule: rule.ID, Severity: rules.SeverityError, File: "ignored.go", Line: 1}
	trace := &engine.Trace{}
	failing := &CachingClient{Store: &Store{Dir: t.TempDir()}, Next: &fakeClient{findings: []engine.Finding{finding}}, Trace: trace}
	got, _, err := failing.Judge(context.Background(), rule, finding.File, []byte("work() // llm-lint:ignore target-rule\n"))
	if err != nil {
		t.Fatal(err)
	}
	filtered, err := suppress.Filter(got, func(string) ([]byte, error) {
		return []byte("work() // llm-lint:ignore target-rule\n"), nil
	})
	if err != nil || len(filtered) != 0 {
		t.Fatalf("Filter() = %#v, %v; want suppressed finding", filtered, err)
	}

	clean := &CachingClient{Store: &Store{Dir: t.TempDir()}, Next: &fakeClient{}, Trace: trace}
	if _, _, err := clean.Judge(context.Background(), rule, "clean.go", []byte("work()\n")); err != nil {
		t.Fatal(err)
	}
	want := []engine.TraceEntry{
		{File: "ignored.go", Rule: rule.ID, Outcome: "fail"},
		{File: "clean.go", Rule: rule.ID, Outcome: "pass"},
	}
	if !reflect.DeepEqual(trace.Entries(), want) {
		t.Fatalf("trace = %#v, want raw outcomes %#v", trace.Entries(), want)
	}
}
