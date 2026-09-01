package storage

import (
	"context"
	"sync"

	"github.com/jocbarbosa/go-redis-like/internal/domain/entity"
	"github.com/jocbarbosa/go-redis-like/internal/domain/repository"
)

type Store struct {
	data        map[string]*entity.Item
	mu          sync.RWMutex
	stopCleanup chan struct{}
}

func NewStore() repository.KeyValueRepository {
	return &Store{
		data:        make(map[string]*entity.Item),
		stopCleanup: make(chan struct{}),
	}
}
func (s *Store) Set(ctx context.Context, key, value string) {}

func (s *Store) Get(ctx context.Context, key string) (string, bool) {
	return "", false
}

func (s *Store) Del(ctx context.Context, key string) int {
	return 0
}

func (s *Store) Expire(ctx context.Context, key string, seconds int) bool {
	return false
}

func (s *Store) TTL(ctx context.Context, key string) int64 {
	return 0
}

func (s *Store) Persist(ctx context.Context, key string) bool {
	return false
}

func (s *Store) Keys(ctx context.Context, pattern string) []string {
	return []string{}
}

func (s *Store) Exists(ctx context.Context, key string) bool {
	return false
}

func (s *Store) Size(ctx context.Context) int {
	return 0
}

func (s *Store) StartCleanup(intervalMs int64) {}

func (s *Store) StopCleanup() {}
