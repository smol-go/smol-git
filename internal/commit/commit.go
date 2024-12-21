package commit

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/smol-go/smol-git/pkg/types"
)

type Commit struct {
	TreeHash   string
	Parent     string
	Message    string
	Author     string
	CommitTime time.Time
}

func NewCommit(treeHash, message string) *Commit {
	return &Commit{
		TreeHash:   treeHash,
		Message:    message,
		Author:     "John Doe <john@example.com>",
		CommitTime: time.Now(),
	}
}

func (c *Commit) Type() types.ObjectType {
	return types.TypeCommit
}

func (c *Commit) Serialize() ([]byte, error) {
	var content strings.Builder
	content.WriteString(fmt.Sprintf("tree %s\n", c.TreeHash))
	if c.Parent != "" {
		content.WriteString(fmt.Sprintf("parent %s\n", c.Parent))
	}
	content.WriteString(fmt.Sprintf("author %s %d +0000\n", c.Author, c.CommitTime.Unix()))
	content.WriteString(fmt.Sprintf("committer %s %d +0000\n", c.Author, c.CommitTime.Unix()))
	content.WriteString(fmt.Sprintf("\n%s\n", c.Message))

	data := content.String()
	header := fmt.Sprintf("%s %d\x00", c.Type(), len(data))
	return []byte(header + data), nil
}

func (c *Commit) Hash() string {
	data, _ := c.Serialize()
	hash := sha1.Sum(data)
	return hex.EncodeToString(hash[:])
}
