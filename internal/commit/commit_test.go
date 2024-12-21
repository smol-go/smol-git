package commit

import (
	"strings"
	"testing"

	"github.com/smol-go/smol-git/pkg/types"
)

func TestCommit(t *testing.T) {
	t.Run("Basic Commit", func(t *testing.T) {
		commit := NewCommit("tree-hash", "Initial commit")
		if commit.Type() != types.TypeCommit {
			t.Errorf("Expected type %s, got %s", types.TypeCommit, commit.Type())
		}

		data, err := commit.Serialize()
		if err != nil {
			t.Fatalf("Failed to serialize commit: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "tree-hash") {
			t.Error("Serialized commit missing tree hash")
		}
		if !strings.Contains(content, "Initial commit") {
			t.Error("Serialized commit missing message")
		}
	})

	t.Run("Commit with Parent", func(t *testing.T) {
		commit := NewCommit("tree-hash", "Second commit")
		commit.Parent = "parent-hash"

		data, err := commit.Serialize()
		if err != nil {
			t.Fatalf("Failed to serialize commit: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "parent parent-hash") {
			t.Error("Serialized commit missing parent hash")
		}
	})

	t.Run("Consistent Hash", func(t *testing.T) {
		commit := NewCommit("tree-hash", "Test commit")
		hash1 := commit.Hash()
		hash2 := commit.Hash()

		if hash1 != hash2 {
			t.Error("Commit hash should be consistent")
		}
	})
}
