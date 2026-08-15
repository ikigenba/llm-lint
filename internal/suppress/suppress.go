package suppress

import "github.com/ikigenba/llm-lint/internal/engine"

func Filter(findings []engine.Finding, read func(path string) ([]byte, error)) ([]engine.Finding, error) {
	_ = read
	return append([]engine.Finding(nil), findings...), nil
}
