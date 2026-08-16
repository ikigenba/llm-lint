package cache

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
)

type traceSink struct {
	mu      sync.Mutex
	entries []engine.TraceEntry
}

func (s *traceSink) Add(entry engine.TraceEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, entry)
}

func (s *traceSink) Entries() []engine.TraceEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]engine.TraceEntry(nil), s.entries...)
}

func TestTraceRecordsCachedHitAndMiss(t *testing.T) {
	// R-H58A-GNO1
	dir := t.TempDir()
	rule := testRule("style", "be clear")
	content := []byte("package a")
	warm := &CachingClient{Store: &Store{Dir: dir}, Next: &fakeClient{}}
	if _, _, err := warm.Judge(context.Background(), rule, "warm.go", content); err != nil {
		t.Fatal(err)
	}

	trace := &traceSink{}
	next := &fakeClient{}
	client := &CachingClient{Store: &Store{Dir: dir}, Next: next, Trace: trace}
	if _, _, err := client.Judge(context.Background(), rule, "hit.go", content); err != nil {
		t.Fatal(err)
	}
	if _, _, err := client.Judge(context.Background(), rule, "miss.go", []byte("package different")); err != nil {
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
