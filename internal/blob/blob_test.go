package blob

import (
	"bytes"
	"testing"

	"github.com/smol-go/smol-git/pkg/types"
)

func TestBlobSerialization(t *testing.T) {
	content := []byte("test content")
	blob := NewBlob(content)

	if blob.Type() != types.TypeBlob {
		t.Errorf("Expected type %s, got %s", types.TypeBlob, blob.Type())
	}

	serialized, err := blob.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize blob: %v", err)
	}

	expectedHeader := []byte("blob 12\x00")
	if !bytes.HasPrefix(serialized, expectedHeader) {
		t.Error("Serialized blob has incorrect header")
	}

	if !bytes.HasSuffix(serialized, content) {
		t.Error("Serialized blob has incorrect content")
	}

	hash1 := blob.Hash()
	hash2 := blob.Hash()
	if hash1 != hash2 {
		t.Error("Hash is not consistent")
	}
}
