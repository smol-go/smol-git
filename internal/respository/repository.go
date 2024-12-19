package respository

type Repository struct {
	Path string
	refs map[string]string
}

func Init(path string) (*Repository, error) {
	repo := &Repository{
		Path: path,
		refs: make(map[string]string),
	}

	return repo, nil
}
