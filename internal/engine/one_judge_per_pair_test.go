package engine

import (
	"context"
	"io"
	"sync"
	"testing"

	"github.com/ikigenba/llm-lint/internal/rules"
)

type judgeFunc func(context.Context, rules.Rule, string, []byte) ([]Finding, Usage, error)

func (f judgeFunc) Judge(ctx context.Context, rule rules.Rule, file string, content []byte) ([]Finding, Usage, error) {
	return f(ctx, rule, file, content)
}

func TestRunJudgesEveryRuleFilePairOnce(t *testing.T) {
	// R-GRPQ-VS95
	var mu sync.Mutex
	calls := make(map[string]int)
	client := judgeFunc(func(_ context.Context, rule rules.Rule, file string, _ []byte) ([]Finding, Usage, error) {
		mu.Lock()
		calls[rule.ID+"/"+file]++
		mu.Unlock()
		return nil, Usage{}, nil
	})
	rs := []rules.Rule{{ID: "one"}, {ID: "two"}}
	files := map[string][]string{"one": {"a.go", "b.go"}, "two": {"a.go", "b.go"}}
	engine := Engine{Client: client, Concurrency: 2}
	_, stats, err := engine.Run(context.Background(), rs, files, func(string) ([]byte, error) { return []byte("x"), nil }, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if stats.Calls != 4 || len(calls) != 4 {
		t.Fatalf("calls = %d across %d pairs, want 4 across 4", stats.Calls, len(calls))
	}
	for pair, count := range calls {
		if count != 1 {
			t.Errorf("calls for %s = %d, want 1", pair, count)
		}
	}
}
