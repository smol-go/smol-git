package repository

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/smol-go/smol-git/internal/object"
)

type Repository struct {
	Path string
	refs map[string]string
}

func Init(path string) (*Repository, error) {
	gitDir := filepath.Join(path, ".git")
	dirs := []string{
		gitDir,
		filepath.Join(gitDir, "objects"),
		filepath.Join(gitDir, "refs"),
		filepath.Join(gitDir, "refs/heads"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	repo := &Repository{
		Path: path,
		refs: make(map[string]string),
	}

	return repo, nil
}

func (r *Repository) WriteObject(obj object.Object) (string, error) {
	hash := obj.Hash()
	objDir := filepath.Join(r.Path, ".git", "objects", hash[:2])
	objPath := filepath.Join(objDir, hash[2:])

	if err := os.MkdirAll(objDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create object directory: %w", err)
	}

	data, err := obj.Serialize()
	if err != nil {
		return "", fmt.Errorf("failed to serialize object: %w", err)
	}

	if err := os.WriteFile(objPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write object file: %w", err)
	}

	return hash, nil
}
