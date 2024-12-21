package blob

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/smol-go/smol-git/pkg/types"
)

type Blob struct {
	content []byte
}

func NewBlob(content []byte) *Blob {
	return &Blob{content: content}
}

func (b *Blob) Type() types.ObjectType {
	return types.TypeBlob
}

func (b *Blob) Serialize() ([]byte, error) {
	header := fmt.Sprintf("%s %d\x00", b.Type(), len(b.content))
	data := make([]byte, len(header)+len(b.content))
	copy(data, []byte(header))
	copy(data[len(header):], b.content)
	return data, nil
}

func (b *Blob) Hash() string {
	data, _ := b.Serialize()
	hash := sha1.Sum(data)
	return hex.EncodeToString(hash[:])
}
