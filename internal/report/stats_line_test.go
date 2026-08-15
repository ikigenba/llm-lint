package report

import (
	"bytes"
	"testing"
)

func TestStatsLineUsesHumanSummary(t *testing.T) {
	var out bytes.Buffer
	StatsLine(&out, Stats{Rules: 2, Files: 41, Pairs: 82, Calls: 7, CacheHits: 75, InputTokens: 118000, OutputTokens: 2000, CostUSD: 0.0113})
	want := "llm-lint: 2 rules, 41 files, 82 pairs, 7 calls, 75 cache hits, 118k in / 2k out tokens, $0.0113\n"
	if out.String() != want {
		t.Fatalf("StatsLine() = %q, want %q", out.String(), want)
	}
}
