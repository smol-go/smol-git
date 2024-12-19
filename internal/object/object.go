package object

import (
	"fmt"
)

const TypeBlob string = "blob"

type Blob struct {
	content []byte
}

func (b *Blob) Type() string {
	return TypeBlob
}

func (b *Blob) Serialize() ([]byte, error) {
	header := fmt.Sprintf("%s %d\x00", b.Type(), len(b.content))
	data := make([]byte, len(header)+len(b.content))
	copy(data, []byte(header))
	copy(data[len(header):], b.content)
	return data, nil
}
