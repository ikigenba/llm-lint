package engine

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// OpenSessionLog creates the session JSONL in an injected agentkit directory.
func OpenSessionLog(dir, sessionID string) (io.WriteCloser, string, error) {
	if dir == "" || sessionID == "" {
		return nil, "", fmt.Errorf("engine: session log directory and id are required")
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, "", fmt.Errorf("engine: create session log directory: %w", err)
	}
	path := filepath.Join(dir, sessionID+".jsonl")
	log, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, "", fmt.Errorf("engine: create session log: %w", err)
	}
	return log, path, nil
}
