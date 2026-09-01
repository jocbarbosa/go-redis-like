package storage

import (
	"context"
	"sync"
	"time"

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
func (s *Store) Set(ctx context.Context, key, value string) {
	if ctx.Err() != nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = &entity.Item{
		Value:     value,
		ExpiresAt: nil,
	}

}

func (s *Store) Get(ctx context.Context, key string) (string, bool) {
	if ctx.Err() != nil {
		return "", false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.data[key]
	if !exists {
		return "", false
	}

	if item.IsExpired(time.Now().Unix()) {
		return "", false
	}

	return item.Value, true

}

func (s *Store) Del(ctx context.Context, key string) int {
	if ctx.Err() != nil {
		return 0
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)

	return 1
}

func (s *Store) Expire(ctx context.Context, key string, seconds int) bool {
	if ctx.Err() != nil {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.data[key]
	if !exists {
		return false
	}

	expiresAt := time.Now().Add(time.Second * time.Duration(seconds)).Unix()

	item.ExpiresAt = &expiresAt

	return true

}

func (s *Store) TTL(ctx context.Context, key string) int64 {
	if ctx.Err() != nil {
		return -1
	}

	s.mu.RLock()
	defer s.mu.Unlock()

	item, exists := s.data[key]
	if !exists {
		return -1
	}

	if item.ExpiresAt == nil {
		return -1
	}

	remaining := *item.ExpiresAt - time.Now().Unix()

	if remaining <= 0 {
		return -1
	}

	return remaining

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
