package engine

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ikigenba/agentkit"
	"github.com/ikigenba/llm-lint/internal/rules"
)

const judgeSystemPrompt = `You are a lint rule judge. Judge only the supplied rule against the supplied file. Report every violation by calling report_violations exactly once. Call it with an empty violations array when the file is clean.`

type AgentkitClient struct {
	NewConversation func(system string, log io.Writer) (*agentkit.Conversation, error)
	Log             io.Writer
	Warn            io.Writer
	Context         int64
}

func (c *AgentkitClient) ContextWindow() int64 { return c.Context }

func (c *AgentkitClient) Judge(ctx context.Context, rule rules.Rule, file string, content []byte) ([]Finding, error) {
	if c.NewConversation == nil {
		panic("engine: AgentkitClient has nil NewConversation")
	}
	conversation, err := c.NewConversation(judgeSystemPrompt, c.Log)
	if err != nil {
		return nil, fmt.Errorf("engine: create conversation: %w", err)
	}
	defer conversation.Close()

	type violation struct {
		Line        int    `json:"line"`
		Evidence    string `json:"evidence"`
		Explanation string `json:"explanation"`
	}
	type report struct {
		Violations []violation `json:"violations"`
	}
	var reported []violation
	toolCalls := 0
	conversation.Tools = append(conversation.Tools, agentkit.NewTool(
		"report_violations",
		"Report the lines that violate the lint rule, or an empty list when clean.",
		func(_ context.Context, in report) (string, error) {
			toolCalls++
			reported = append(reported[:0], in.Violations...)
			return "violations recorded", nil
		},
	))

	userText := formatJudgeRequest(rule, file, content)
	for attempt := 0; attempt < 2; attempt++ {
		if attempt == 1 {
			userText = "You did not call report_violations. Call it exactly once now, using an empty violations array if the file is clean."
		}
		stream := conversation.Send(ctx, userText)
		for range stream.Events() {
		}
		if err := stream.Err(); err != nil {
			return nil, fmt.Errorf("engine: judge rule %s file %s: %w", rule.ID, file, err)
		}
		if toolCalls > 0 {
			break
		}
	}
	if toolCalls == 0 {
		return nil, fmt.Errorf("engine: judge rule %s file %s: report_violations was not called after retry", rule.ID, file)
	}

	lines := sourceLines(content)
	findings := make([]Finding, 0, len(reported))
	for _, item := range reported {
		if item.Line < 1 || item.Line > len(lines) {
			if c.Warn != nil {
				fmt.Fprintf(c.Warn, "engine: dropping out-of-range line %d for rule %s file %s\n", item.Line, rule.ID, file)
			}
			continue
		}
		findings = append(findings, Finding{
			Rule:        rule.ID,
			Severity:    rule.Severity,
			File:        file,
			Line:        item.Line,
			Evidence:    lines[item.Line-1],
			Explanation: item.Explanation,
		})
	}
	return findings, nil
}

func formatJudgeRequest(rule rules.Rule, file string, content []byte) string {
	var numbered strings.Builder
	for i, line := range sourceLines(content) {
		numbered.WriteString(strconv.Itoa(i + 1))
		numbered.WriteString(": ")
		numbered.WriteString(line)
		numbered.WriteByte('\n')
	}
	return fmt.Sprintf("Rule %s:\n%s\n\nFile %s:\n%s", rule.ID, rule.Prompt, file, numbered.String())
}

func sourceLines(content []byte) []string {
	if len(content) == 0 {
		return nil
	}
	lines := strings.Split(string(content), "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}
