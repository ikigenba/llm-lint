package rules

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// R-GFIR-22U7
func TestLoadRejectsBuiltInAndLocalIDCollisions(t *testing.T) {
	t.Run("built-in", func(t *testing.T) {
		root := t.TempDir()
		writeTestRule(t, filepath.Join(root, "no-sleep-in-tests.md"), "collision")
		_, err := Load(root, []string{"no-sleep-in-tests.md"})
		assertCollisionError(t, err, "no-sleep-in-tests")
	})

	t.Run("local", func(t *testing.T) {
		root := t.TempDir()
		first := filepath.Join(root, "first")
		second := filepath.Join(root, "second")
		if err := os.Mkdir(first, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(second, 0o700); err != nil {
			t.Fatal(err)
		}
		writeTestRule(t, filepath.Join(first, "duplicate.md"), "first")
		writeTestRule(t, filepath.Join(second, "duplicate.md"), "second")
		_, err := Load(root, []string{"first", "second"})
		assertCollisionError(t, err, "duplicate")
	})
}

func assertCollisionError(t *testing.T, err error, id string) {
	t.Helper()
	if !errors.Is(err, ErrRule) || !strings.Contains(err.Error(), id) || !strings.Contains(err.Error(), "collides") {
		t.Fatalf("Load() error = %v, want ErrRule naming collision %q", err, id)
	}
}
