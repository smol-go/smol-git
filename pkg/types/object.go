package types

type ObjectType string

const (
	TypeBlob   ObjectType = "blob"
	TypeTree   ObjectType = "tree"
	TypeCommit ObjectType = "commit"
)

type Object interface {
	Type() ObjectType
	Serialize() ([]byte, error)
	Hash() string
}
