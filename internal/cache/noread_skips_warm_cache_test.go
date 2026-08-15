package cache

import (
	"context"
	"reflect"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
)

func TestNoReadSkipsWarmCacheButWrites(t *testing.T) {
	// R-H54N-39ES
	dir := t.TempDir()
	rule := testRule("one", "check")
	warm := &CachingClient{Store: &Store{Dir: dir}, Next: &fakeClient{findings: []engine.Finding{{Evidence: "old"}}}}
	if _, err := warm.Judge(context.Background(), rule, "a.go", []byte("same")); err != nil {
		t.Fatal(err)
	}

	refreshNext := &fakeClient{findings: []engine.Finding{{Evidence: "new"}}}
	refresh := &CachingClient{Store: &Store{Dir: dir}, Next: refreshNext, NoRead: true}
	if _, err := refresh.Judge(context.Background(), rule, "a.go", []byte("same")); err != nil {
		t.Fatal(err)
	}
	if len(refreshNext.calls) != 1 {
		t.Fatalf("NoRead calls = %v, want one", refreshNext.calls)
	}

	next := &fakeClient{}
	normal := &CachingClient{Store: &Store{Dir: dir}, Next: next}
	got, err := normal.Judge(context.Background(), rule, "a.go", []byte("same"))
	if err != nil {
		t.Fatal(err)
	}
	if len(next.calls) != 0 {
		t.Fatalf("normal run calls = %v, want none", next.calls)
	}
	if want := []engine.Finding{{File: "a.go", Evidence: "new"}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("normal findings = %#v, want %#v", got, want)
	}
}
