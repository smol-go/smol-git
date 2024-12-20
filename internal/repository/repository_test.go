package repository

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRepository(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "smolgit-repo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("Initialize Repository", func(t *testing.T) {
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

		indexPath := filepath.Join(gitDir, "index")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			t.Error("Index file was not created")
		}
	})

	t.Run("Open Repository", func(t *testing.T) {
		_, err := Init(tmpDir)
		if err != nil {
			t.Fatalf("Failed to initialize repository: %v", err)
		}

		repo, err := Open(tmpDir)
		if err != nil {
			t.Fatalf("Failed to open repository: %v", err)
		}

		if repo.Path != tmpDir {
			t.Errorf("Expected repository path %s, got %s", tmpDir, repo.Path)
		}
	})

	t.Run("Invalid Repository", func(t *testing.T) {
		invalidDir := filepath.Join(tmpDir, "invalid")
		if err := os.MkdirAll(invalidDir, 0755); err != nil {
			t.Fatalf("Failed to create invalid directory: %v", err)
		}

		_, err := Open(invalidDir)
		if err == nil {
			t.Error("Expected error when opening invalid repository")
		}
	})
}
