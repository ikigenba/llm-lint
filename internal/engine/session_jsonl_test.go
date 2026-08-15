package engine

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/ikigenba/agentkit"
	"github.com/ikigenba/llm-lint/internal/rules"
)

func TestSessionJSONLContainsCallRecords(t *testing.T) {
	// R-GXT8-SMYM
	log, path, err := OpenSessionLog(t.TempDir(), "session-test")
	if err != nil {
		t.Fatal(err)
	}
	input, _ := json.Marshal(map[string]any{"violations": []any{}})
	provider := &scriptedProvider{responses: []agentkit.Message{
		{Role: agentkit.RoleAssistant, Blocks: []agentkit.Block{agentkit.ToolUseBlock{ID: "call-1", Name: "report_violations", Input: input}}},
		assistantText("done"),
	}}
	client := AgentkitClient{Log: log, NewConversation: func(system string, writer io.Writer) (*agentkit.Conversation, error) {
		return &agentkit.Conversation{Provider: provider, Model: "fake", System: system, Log: writer}, nil
	}}
	if _, err := client.Judge(t.Context(), rules.Rule{ID: "rule-log", Prompt: "judge"}, "logged.go", []byte("source")); err != nil {
		t.Fatal(err)
	}
	if err := log.Close(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if !strings.Contains(text, `"type":"turn_start"`) || !strings.Contains(text, `"type":"message"`) || !strings.Contains(text, `"type":"tool_use"`) {
		t.Fatalf("session log does not contain the call's turn, response, and tool records: %s", text)
	}
}
