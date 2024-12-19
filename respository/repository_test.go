package respository

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gogit-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	repo, err := Init(tmpDir)
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	if repo == nil {
		t.Error("Initialized repository is nil")
	}

	gitDir := filepath.Join(tmpDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Error("Git directory was not created")
	}

	requiredDirs := []string{
		"objects",
		"refs",
		"refs/heads",
	}

	for _, dir := range requiredDirs {
		path := filepath.Join(gitDir, dir)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Required directory %s was not created", dir)
		}
	}
}
