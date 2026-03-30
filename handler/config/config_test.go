package config

import (
	"os"
	"testing"
)

func TestConfigLifecycle(t *testing.T) {
	// Setup a clean directory for testing to prevent modifying the actual config.yaml
	tempDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory to temp dir: %v", err)
	}

	// 1. CheckConfigFile should be false initially because there's no config.yaml
	if CheckConfigFile() {
		t.Error("CheckConfigFile should return false when no config.yaml exists")
	}

	// 2. LoadConfig should create the default config since it's missing, and then return it
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if len(cfg.Cameras) == 0 {
		t.Fatalf("Expected at least 1 camera in default config")
	}

	if cfg.Cameras[0].Name != "cam1" {
		t.Errorf("Expected first camera name to be 'cam1', got '%s'", cfg.Cameras[0].Name)
	}

	// 3. CheckConfigFile should now be true
	if !CheckConfigFile() {
		t.Error("CheckConfigFile should return true after LoadConfig creates default config")
	}

	// 4. ReadConfig should now work independently
	cfg2, err := ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig failed: %v", err)
	}

	if len(cfg2.Cameras) == 0 {
		t.Fatalf("Expected at least 1 camera in read config")
	}
}
