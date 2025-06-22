package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Defaults(t *testing.T) {
	os.Clearenv()
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.WatchFolder != "./watched" {
		t.Errorf("expected ./watched, got %s", cfg.WatchFolder)
	}
}
