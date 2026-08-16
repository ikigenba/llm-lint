package engine

import (
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/ikigenba/llm-lint/internal/rules"
)

func TestOversizeSkipRecordsTraceOnlyWhenEnabled(t *testing.T) {
	// R-H40E-2VXC
	rule := rules.Rule{ID: "large-rule", Prompt: "inspect this file"}
	run := func(trace *Trace) (*budgetClient, []TraceEntry) {
		client := &budgetClient{}
		runner := Engine{Client: client, Concurrency: 1, Trace: trace}
		_, _, err := runner.Run(context.Background(), []rules.Rule{rule}, map[string][]string{rule.ID: {"large.go"}}, func(string) ([]byte, error) {
			return make([]byte, 100), nil
		}, io.Discard)
		if err != nil {
			t.Fatalf("Run() error = %v", err)
		}
		if trace == nil {
			return client, nil
		}
		return client, trace.Entries()
	}

	trace := &Trace{}
	client, entries := run(trace)
	want := []TraceEntry{{File: "large.go", Rule: rule.ID, Outcome: "skipped"}}
	if client.calls != 0 || !reflect.DeepEqual(entries, want) {
		t.Fatalf("calls = %d, trace = %#v; want zero calls and %#v", client.calls, entries, want)
	}
	client, entries = run(nil)
	if client.calls != 0 || len(entries) != 0 {
		t.Fatalf("nil trace: calls = %d, entries = %#v; want neither", client.calls, entries)
	}
}
