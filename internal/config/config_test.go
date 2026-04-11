package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := LoadConfig()

	// Test default pagination
	if cfg.Pagination != "numeric" && cfg.Pagination != "dots" {
		t.Errorf("Expected pagination to be 'numeric' or 'dots', got %s", cfg.Pagination)
	}

	// Test default behavioral flags
	if cfg.AutoClear != false {
		t.Errorf("Expected default AutoClear to be false, got %v", cfg.AutoClear)
	}
	if cfg.AutoExit != false {
		t.Errorf("Expected default AutoExit to be false, got %v", cfg.AutoExit)
	}

	// Test default views exist
	if len(cfg.Views) == 0 {
		t.Error("Expected default views to be populated")
	}

	// Test default theme exists
	if cfg.Theme.Primary == "" {
		t.Error("Expected primary theme color to be set")
	}
}
