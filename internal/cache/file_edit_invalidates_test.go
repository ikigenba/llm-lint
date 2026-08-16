package cache

import (
	"context"
	"reflect"
	"testing"
)

func TestFileEditInvalidatesOnlyEditedFile(t *testing.T) {
	// R-H3WQ-PHO3
	dir := t.TempDir()
	rule := testRule("one", "check")
	warm := &CachingClient{Store: &Store{Dir: dir}, Next: &fakeClient{}}
	for _, file := range []struct{ name, content string }{{"a.go", "old a"}, {"b.go", "same b"}} {
		if _, _, err := warm.Judge(context.Background(), rule, file.name, []byte(file.content)); err != nil {
			t.Fatal(err)
		}
	}

	next := &fakeClient{}
	cached := &CachingClient{Store: &Store{Dir: dir}, Next: next}
	for _, file := range []struct{ name, content string }{{"a.go", "new a"}, {"b.go", "same b"}} {
		if _, _, err := cached.Judge(context.Background(), rule, file.name, []byte(file.content)); err != nil {
			t.Fatal(err)
		}
	}
	if want := []string{"one:a.go"}; !reflect.DeepEqual(next.calls, want) {
		t.Fatalf("calls after file edit = %v, want %v", next.calls, want)
	}
}
