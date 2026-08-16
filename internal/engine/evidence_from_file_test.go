package engine

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/ikigenba/agentkit"
	"github.com/ikigenba/llm-lint/internal/rules"
)

func TestAgentkitClientUsesActualSourceEvidenceAndDropsBeyondEOF(t *testing.T) {
	// R-GSXN-9JZU
	input, err := json.Marshal(map[string]any{"violations": []map[string]any{
		{"line": 2, "evidence": "invented", "explanation": "bad call"},
		{"line": 9, "evidence": "invented", "explanation": "outside"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	provider := &scriptedProvider{responses: []agentkit.Message{
		{Role: agentkit.RoleAssistant, Blocks: []agentkit.Block{agentkit.ToolUseBlock{ID: "call-1", Name: "report_violations", Input: input}}},
		assistantText("done"),
	}}
	var warning bytes.Buffer
	client := AgentkitClient{
		NewConversation: func(system string, log io.Writer) (*agentkit.Conversation, error) {
			return &agentkit.Conversation{Provider: provider, Model: "fake", System: system, Log: log}, nil
		},
		Warn: &warning,
	}
	rule := rules.Rule{ID: "bad-api", Severity: rules.SeverityWarning, Prompt: "flag bad calls"}
	findings, _, err := client.Judge(t.Context(), rule, "main.go", []byte("first\nactual second\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("findings = %#v, want one", findings)
	}
	got := findings[0]
	if got.Line != 2 || got.Evidence != "actual second" || got.Explanation != "bad call" || got.Severity != rules.SeverityWarning {
		t.Fatalf("finding = %#v", got)
	}
	if !strings.Contains(warning.String(), "line 9") {
		t.Fatalf("warning = %q, want dropped line number", warning.String())
	}
}
