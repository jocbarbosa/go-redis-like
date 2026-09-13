package persistence

import "github.com/jocbarbosa/go-redis-like/internal/domain/repository"

type AOFProviderOption struct {
	EnableAOF bool
	FilePath  string
}

func NewAOFProvider(opt AOFProviderOption) (repository.PersistenceRepository, error) {
	if !opt.EnableAOF {
		return nil, nil
	}

	return NewAOF(opt.FilePath)
}
