package config

import "testing"

func TestLoadAbsentConfigUsesDefaults(t *testing.T) {
	// R-G0VY-GTXV
	cwd := t.TempDir()
	cfg, err := Load(cwd, nil, func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Root != cwd || len(cfg.Enable) != 0 || len(cfg.Disable) != 0 {
		t.Fatalf("Root = %q, Enable = %#v, Disable = %#v", cfg.Root, cfg.Enable, cfg.Disable)
	}
	if cfg.Model.ModelID != defaultModel {
		t.Fatalf("ModelID = %q, want %q", cfg.Model.ModelID, defaultModel)
	}
}
