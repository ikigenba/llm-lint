package cache

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNormalRunPrunesOldEntriesAndKeepsYoung(t *testing.T) {
	// R-H6CJ-H15H
	dir := t.TempDir()
	old := filepath.Join(dir, "aa", "old.json")
	young := filepath.Join(dir, "bb", "young.json")
	for _, path := range []string{old, young} {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("[]"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	now := time.Date(2026, time.August, 15, 12, 0, 0, 0, time.UTC)
	if err := os.Chtimes(old, now.Add(-31*24*time.Hour), now.Add(-31*24*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(young, now.Add(-29*24*time.Hour), now.Add(-29*24*time.Hour)); err != nil {
		t.Fatal(err)
	}

	client := &CachingClient{Store: &Store{Dir: dir, Now: func() time.Time { return now }}, Next: &fakeClient{}}
	if _, _, err := client.Judge(context.Background(), testRule("one", "check"), "a.go", []byte("content")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Fatalf("old entry stat error = %v, want not exist", err)
	}
	if _, err := os.Stat(young); err != nil {
		t.Fatalf("young entry was removed: %v", err)
	}
}
