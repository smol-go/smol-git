package respository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smol-go/smol-git/internal/object"
)

func TestInit(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "smolgit-test-*")
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

func TestWriteObject(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "smolgit-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	repo, err := Init(tmpDir)
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	content := []byte("test content")
	blob := object.NewBlob(content)
	hash, err := repo.WriteObject(blob)
	if err != nil {
		t.Fatalf("Failed to write object: %v", err)
	}

	objectPath := filepath.Join(tmpDir, ".git", "objects", hash[:2], hash[2:])
	if _, err := os.Stat(objectPath); os.IsNotExist(err) {
		t.Error("Object file was not created")
	}

	data, err := os.ReadFile(objectPath)
	if err != nil {
		t.Fatalf("Failed to read object file: %v", err)
	}

	expectedData, _ := blob.Serialize()
	if string(data) != string(expectedData) {
		t.Error("Object content does not match expected content")
	}
}
