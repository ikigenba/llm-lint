package cache

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ikigenba/llm-lint/internal/engine"
	"github.com/ikigenba/llm-lint/internal/rules"
)

type Store struct {
	Dir string
	Now func() time.Time
}

type CachingClient struct {
	Store  *Store
	Next   engine.Client
	NoRead bool
	Hits   *atomic.Int64
	Trace  engine.TraceSink

	pruneOnce sync.Once
	warnOnce  sync.Once
}

func (c *CachingClient) Judge(ctx context.Context, r rules.Rule, file string, content []byte) ([]engine.Finding, engine.Usage, error) {
	if c.Store == nil {
		return c.judgeNext(ctx, r, file, content)
	}

	dir := c.Store.dir()
	if dir == "" {
		return c.judgeNext(ctx, r, file, content)
	}
	now := c.Store.now()
	c.pruneOnce.Do(func() {
		if err := c.Store.prune(dir, now); err != nil {
			c.warn(err)
		}
	})

	ruleContent, err := json.Marshal(r)
	if err != nil {
		return c.judgeNext(ctx, r, file, content)
	}
	key := Key(ruleContent, content)
	path := filepath.Join(dir, key[:2], key+".json")
	if !c.NoRead {
		findings, ok, readErr := read(path, now)
		if readErr != nil {
			c.warn(readErr)
		}
		if ok {
			for i := range findings {
				findings[i].File = file
			}
			if c.Hits != nil {
				c.Hits.Add(1)
			}
			c.record(file, r.ID, findings, true)
			return findings, engine.Usage{}, nil
		}
	}

	findings, usage, err := c.judgeNext(ctx, r, file, content)
	if err != nil {
		return nil, usage, err
	}
	if err := write(path, findings); err != nil {
		c.warn(err)
	}
	return findings, usage, nil
}

func (c *CachingClient) judgeNext(ctx context.Context, r rules.Rule, file string, content []byte) ([]engine.Finding, engine.Usage, error) {
	findings, usage, err := c.Next.Judge(ctx, r, file, content)
	if err == nil {
		c.record(file, r.ID, findings, false)
	}
	return findings, usage, err
}

func (c *CachingClient) record(file, rule string, findings []engine.Finding, cached bool) {
	if c.Trace == nil {
		return
	}
	outcome := "pass"
	if len(findings) > 0 {
		outcome = "fail"
	}
	c.Trace.Add(engine.TraceEntry{File: file, Rule: rule, Cached: cached, Outcome: outcome})
}

func (c *CachingClient) warn(err error) {
	c.warnOnce.Do(func() {
		_, _ = fmt.Fprintf(os.Stderr, "cache: %v\n", err)
	})
}

func Key(rulePromptAndMeta, fileContent []byte) string {
	h := sha256.New()
	for _, part := range [][]byte{rulePromptAndMeta, fileContent} {
		var size [8]byte
		binary.BigEndian.PutUint64(size[:], uint64(len(part)))
		h.Write(size[:])
		h.Write(part)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func (s *Store) dir() string {
	if s.Dir != "" {
		return s.Dir
	}
	dir, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "llm-lint")
}

func (s *Store) now() time.Time {
	if s.Now != nil {
		return s.Now()
	}
	return time.Now()
}

func (s *Store) prune(dir string, now time.Time) error {
	if dir == "" {
		return nil
	}
	cutoff := now.Add(-30 * 24 * time.Hour)
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if os.IsNotExist(walkErr) {
				return nil
			}
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.ModTime().Before(cutoff) {
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
		return nil
	})
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func read(path string, now time.Time) ([]engine.Finding, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	var findings []engine.Finding
	if err := json.Unmarshal(data, &findings); err != nil {
		return nil, false, err
	}
	if findings == nil {
		findings = []engine.Finding{}
	}
	if err := os.Chtimes(path, now, now); err != nil {
		return findings, true, err
	}
	return findings, true, nil
}

func write(path string, findings []engine.Finding) error {
	if path == "" {
		return nil
	}
	if findings == nil {
		findings = []engine.Finding{}
	}
	data, err := json.Marshal(findings)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".entry-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if _, err = tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err = tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}
