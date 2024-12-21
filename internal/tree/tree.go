package tree

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/smol-go/smol-git/pkg/types"
)

type TreeEntry struct {
	Path string
	Hash string
}

type Tree struct {
	entries []TreeEntry
}

func NewTree() *Tree {
	return &Tree{entries: make([]TreeEntry, 0)}
}

func (t *Tree) Type() types.ObjectType {
	return types.TypeTree
}

func (t *Tree) AddEntry(path, hash string) {
	t.entries = append(t.entries, TreeEntry{Path: path, Hash: hash})
}

func (t *Tree) Serialize() ([]byte, error) {
	sort.Slice(t.entries, func(i, j int) bool {
		return t.entries[i].Path < t.entries[j].Path
	})

	var content strings.Builder
	for _, entry := range t.entries {
		content.WriteString(fmt.Sprintf("%s %s\n", entry.Hash, entry.Path))
	}

	data := content.String()
	header := fmt.Sprintf("%s %d\x00", t.Type(), len(data))
	return []byte(header + data), nil
}

func (t *Tree) Hash() string {
	data, _ := t.Serialize()
	hash := sha1.Sum(data)
	return hex.EncodeToString(hash[:])
}
