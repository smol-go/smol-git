package index

import (
	"encoding/json"
	"os"
	"sort"
)

type Index struct {
	Entries map[string]string
}

func NewIndex() *Index {
	return &Index{
		Entries: make(map[string]string),
	}
}

func Read(path string) (*Index, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return NewIndex(), nil
	}
	if err != nil {
		return nil, err
	}

	var idx Index
	if err := json.Unmarshal(data, &idx.Entries); err != nil {
		return nil, err
	}
	return &idx, nil
}

func (idx *Index) Write(path string) error {
	data, err := json.Marshal(idx.Entries)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (idx *Index) Add(path, hash string) {
	idx.Entries[path] = hash
}

func (idx *Index) StagedFiles() []string {
	var files []string
	for path := range idx.Entries {
		files = append(files, path)
	}
	sort.Strings(files)
	return files
}
