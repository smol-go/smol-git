package tree

import (
	"strings"
	"testing"

	"github.com/smol-go/smol-git/pkg/types"
)

func TestTree(t *testing.T) {
	t.Run("Empty Tree", func(t *testing.T) {
		tree := NewTree()
		if tree.Type() != types.TypeTree {
			t.Errorf("Expected type %s, got %s", types.TypeTree, tree.Type())
		}
	})

	t.Run("Add Entries", func(t *testing.T) {
		tree := NewTree()
		tree.AddEntry("file1.txt", "hash1")
		tree.AddEntry("file2.txt", "hash2")

		data, err := tree.Serialize()
		if err != nil {
			t.Fatalf("Failed to serialize tree: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "file1.txt") || !strings.Contains(content, "hash1") {
			t.Error("Serialized tree missing entry")
		}
	})

	t.Run("Consistent Hash", func(t *testing.T) {
		tree1 := NewTree()
		tree1.AddEntry("file1.txt", "hash1")
		tree1.AddEntry("file2.txt", "hash2")

		tree2 := NewTree()
		tree2.AddEntry("file2.txt", "hash2")
		tree2.AddEntry("file1.txt", "hash1")

		if tree1.Hash() != tree2.Hash() {
			t.Error("Trees with same entries in different order should have same hash")
		}
	})
}
