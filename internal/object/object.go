package object

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

const TypeBlob string = "blob"

type Blob struct {
	content []byte
}

func (b *Blob) Type() string {
	return TypeBlob
}

type Object interface {
	Type() string
	Serialize() ([]byte, error)
	Hash() string
}

func NewBlob(content []byte) *Blob {
	return &Blob{content: content}
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
