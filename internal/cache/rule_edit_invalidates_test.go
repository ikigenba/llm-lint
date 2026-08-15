package cache

import (
	"context"
	"reflect"
	"testing"
)

func TestRuleEditInvalidatesOnlyEditedRule(t *testing.T) {
	// R-H2OU-BPXE
	dir := t.TempDir()
	rulesBefore := []struct{ id, prompt string }{{"one", "first"}, {"two", "second"}}
	warmNext := &fakeClient{}
	warm := &CachingClient{Store: &Store{Dir: dir}, Next: warmNext}
	for _, rule := range rulesBefore {
		if _, err := warm.Judge(context.Background(), testRule(rule.id, rule.prompt), "a.go", []byte("same")); err != nil {
			t.Fatal(err)
		}
	}

	next := &fakeClient{}
	cached := &CachingClient{Store: &Store{Dir: dir}, Next: next}
	for _, rule := range []struct{ id, prompt string }{{"one", "first"}, {"two", "edited"}} {
		if _, err := cached.Judge(context.Background(), testRule(rule.id, rule.prompt), "a.go", []byte("same")); err != nil {
			t.Fatal(err)
		}
	}
	if want := []string{"two:a.go"}; !reflect.DeepEqual(next.calls, want) {
		t.Fatalf("calls after rule edit = %v, want %v", next.calls, want)
	}
}
