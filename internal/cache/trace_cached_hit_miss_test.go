package cache

import (
	"context"
	"reflect"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
)

func TestTraceRecordsCachedHitAndMiss(t *testing.T) {
	// R-H58A-GNO1
	dir := t.TempDir()
	rule := testRule("style", "be clear")
	content := []byte("package a")
	warm := &CachingClient{Store: &Store{Dir: dir}, Next: &fakeClient{}}
	if _, err := warm.Judge(context.Background(), rule, "warm.go", content); err != nil {
		t.Fatal(err)
	}

	trace := &engine.Trace{}
	next := &fakeClient{}
	client := &CachingClient{Store: &Store{Dir: dir}, Next: next, Trace: trace}
	if _, err := client.Judge(context.Background(), rule, "hit.go", content); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Judge(context.Background(), rule, "miss.go", []byte("package different")); err != nil {
		t.Fatal(err)
	}
	want := []engine.TraceEntry{
		{File: "hit.go", Rule: rule.ID, Cached: true, Outcome: "pass"},
		{File: "miss.go", Rule: rule.ID, Cached: false, Outcome: "pass"},
	}
	if len(next.calls) != 1 || !reflect.DeepEqual(trace.Entries(), want) {
		t.Fatalf("next calls = %v, trace = %#v; want one miss and %#v", next.calls, trace.Entries(), want)
	}
}
