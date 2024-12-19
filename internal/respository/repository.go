package respository

import (
	"fmt"
	"os"
	"path/filepath"
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
