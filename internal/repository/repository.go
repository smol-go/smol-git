package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/smol-go/smol-git/internal/blob"
	"github.com/smol-go/smol-git/internal/commit"
	"github.com/smol-go/smol-git/internal/index"
	"github.com/smol-go/smol-git/internal/tree"
	"github.com/smol-go/smol-git/pkg/types"
)

type Repository struct {
	Path  string
	index *index.Index
	refs  map[string]string
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

	idx := index.NewIndex()
	if err := idx.Write(filepath.Join(gitDir, "index")); err != nil {
		return nil, fmt.Errorf("failed to create index: %w", err)
	}

	repo := &Repository{
		Path:  path,
		index: idx,
		refs:  make(map[string]string),
	}

	return repo, nil
}

func (r *Repository) WriteObject(obj types.Object) (string, error) {
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

func Open(path string) (*Repository, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	gitDir := findGitDirectory(absPath)
	if gitDir == "" {
		return nil, fmt.Errorf("not a git repository")
	}

	idx, err := index.Read(filepath.Join(gitDir, "index"))
	if err != nil {
		return nil, fmt.Errorf("failed to read index: %w", err)
	}

	return &Repository{
		Path:  filepath.Dir(gitDir),
		index: idx,
		refs:  make(map[string]string),
	}, nil
}

func findGitDirectory(path string) string {
	for {
		gitPath := filepath.Join(path, ".git")
		if fi, err := os.Stat(gitPath); err == nil && fi.IsDir() {
			return gitPath
		}
		parent := filepath.Dir(path)
		if parent == path {
			return ""
		}
		path = parent
	}
}

func (r *Repository) Add(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	blob := blob.NewBlob(content)
	hash, err := r.WriteObject(blob)
	if err != nil {
		return fmt.Errorf("failed to write blob: %w", err)
	}

	relPath, err := filepath.Rel(r.Path, absPath)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}

	r.index.Add(relPath, hash)
	if err := r.index.Write(filepath.Join(r.Path, ".git", "index")); err != nil {
		return fmt.Errorf("failed to write index: %w", err)
	}

	return nil
}

func (r *Repository) Status() (string, error) {
	var sb strings.Builder

	_, err := r.readRef("HEAD")
	if err != nil {
		sb.WriteString("No commits yet\n\n")
	} else {
		sb.WriteString("On branch master\n")
	}

	staged := r.index.StagedFiles()
	if len(staged) > 0 {
		sb.WriteString("\nChanges to be committed:\n")
		for _, file := range staged {
			sb.WriteString(fmt.Sprintf("\tmodified: %s\n", file))
		}
	}

	return sb.String(), nil
}

func (r *Repository) Commit(message string) (string, error) {
	tree := tree.NewTree()
	for path, hash := range r.index.Entries {
		tree.AddEntry(path, hash)
	}

	treeHash, err := r.WriteObject(tree)
	if err != nil {
		return "", fmt.Errorf("failed to write tree: %w", err)
	}

	commit := commit.NewCommit(treeHash, message)

	if head, err := r.readRef("HEAD"); err == nil {
		commit.Parent = head
	}

	commitHash, err := r.WriteObject(commit)
	if err != nil {
		return "", fmt.Errorf("failed to write commit: %w", err)
	}

	if err := r.updateRef("HEAD", commitHash); err != nil {
		return "", fmt.Errorf("failed to update HEAD: %w", err)
	}

	return commitHash, nil
}

func (r *Repository) readRef(ref string) (string, error) {
	path := filepath.Join(r.Path, ".git", ref)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func (r *Repository) updateRef(ref, hash string) error {
	path := filepath.Join(r.Path, ".git", ref)
	return os.WriteFile(path, []byte(hash), 0644)
}
