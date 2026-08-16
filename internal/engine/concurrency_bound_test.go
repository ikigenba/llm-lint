package engine

import (
	"context"
	"io"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ikigenba/llm-lint/internal/rules"
)

func TestRunBoundsConcurrentJudges(t *testing.T) {
	// R-GVDG-13H8
	const limit = 3
	var active, maximum atomic.Int64
	client := judgeFunc(func(_ context.Context, _ rules.Rule, _ string, _ []byte) ([]Finding, Usage, error) {
		current := active.Add(1)
		for old := maximum.Load(); current > old && !maximum.CompareAndSwap(old, current); old = maximum.Load() {
		}
		time.Sleep(15 * time.Millisecond)
		active.Add(-1)
		return nil, Usage{}, nil
	})
	files := map[string][]string{"rule": {"a", "b", "c", "d", "e", "f"}}
	engine := Engine{Client: client, Concurrency: limit}
	if _, _, err := engine.Run(context.Background(), []rules.Rule{{ID: "rule"}}, files, func(string) ([]byte, error) { return []byte("x"), nil }, io.Discard); err != nil {
		t.Fatal(err)
	}
	if got := maximum.Load(); got <= 1 || got > limit {
		t.Fatalf("maximum in flight = %d, want > 1 and <= %d", got, limit)
	}
}
