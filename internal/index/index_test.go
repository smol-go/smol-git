package index

import (
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestIndex(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gogit-index-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("New Index", func(t *testing.T) {
		idx := NewIndex()
		if len(idx.Entries) != 0 {
			t.Error("New index should be empty")
		}
	})

	t.Run("Add and List Files", func(t *testing.T) {
		idx := NewIndex()

		files := map[string]string{
			"file1.txt":     "hash1",
			"file2.txt":     "hash2",
			"dir/file3.txt": "hash3",
		}

		for path, hash := range files {
			idx.Add(path, hash)
		}

		staged := idx.StagedFiles()
		sort.Strings(staged)

		expected := []string{"file1.txt", "file2.txt", "dir/file3.txt"}
		sort.Strings(expected)

		if !reflect.DeepEqual(staged, expected) {
			t.Errorf("Expected staged files %v, got %v", expected, staged)
		}
	})

	t.Run("Remove Files", func(t *testing.T) {
		idx := NewIndex()

		idx.Add("file1.txt", "hash1")
		if !idx.IsStaged("file1.txt") {
			t.Error("File should be staged")
		}

		idx.Remove("file1.txt")
		if idx.IsStaged("file1.txt") {
			t.Error("File should not be staged after removal")
		}
	})
}
