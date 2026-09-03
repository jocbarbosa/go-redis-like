package storage

import (
	"context"
	"path/filepath"
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

// Set sets a new item in the mapping
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

// Get returns an item when is available based on the key
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

// Del deletes an item based on a given key
func (s *Store) Del(ctx context.Context, key string) int {
	if ctx.Err() != nil {
		return 0
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)

	return 1
}

// Expire sets an expiration TTL for an item key
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

// TTL sets a TTL to a given item key
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

// Persist removes TTL and transform an item to persistent
func (s *Store) Persist(ctx context.Context, key string) bool {
	if ctx.Err() != nil {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.data[key]
	if !exists {
		return false
	}

	item.ExpiresAt = nil

	return true

}

// Keys return all keys that matches a given pattern
func (s *Store) Keys(ctx context.Context, pattern string) []string {
	if ctx.Err() != nil {
		return []string{}
	}

	s.mu.RLock()
	defer s.mu.Unlock()

	now := time.Now().Unix()
	var matches []string

	for key, item := range s.data {
		if item.IsExpired(now) {
			continue
		}

		if matchPattern(key, pattern) {
			matches = append(matches, key)
		}
	}

	return matches

}

// Exists returns if a given key exists
func (s *Store) Exists(ctx context.Context, key string) bool {
	if ctx.Err() != nil {
		return false
	}

	s.mu.RLock()
	defer s.mu.Unlock()

	item, exists := s.data[key]
	if !exists {
		return false
	}

	return !item.IsExpired(time.Now().Unix())
}

// Size returns the number of all items
// this ignore if the item is expired or not
func (s *Store) Size(ctx context.Context) int {
	if ctx.Err() != nil {
		return 0
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

// StartCleanup cleans up all expired keys
func (s *Store) StartCleanup(ctx context.Context, intervalMs int64) {
	interval := time.Duration(intervalMs) * time.Millisecond
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				now := time.Now().Unix()
				for k, v := range s.data {
					if v.IsExpired(now) {
						s.Del(ctx, k)
					}
				}
			case <-ctx.Done():
				return
			case <-s.stopCleanup:
				return
			}
		}

	}()
}

func (s *Store) StopCleanup() {
	close(s.stopCleanup)
}

func matchPattern(key, pattern string) bool {
	if pattern == "*" {
		return true
	}

	matched, err := filepath.Match(pattern, key)
	if err != nil {
		return false
	}

	return matched

}
