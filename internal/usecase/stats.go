package usecase

import (
	"context"
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/jocbarbosa/go-redis-like/internal/domain/repository"
)

type Stats struct {
	startedAt        time.Time
	totalCommands    atomic.Int64
	totalConnections atomic.Int64
	keyspace         repository.KeyValueRepository
}

func NewStats(keyspace repository.KeyValueRepository) *Stats {
	return &Stats{
		startedAt: time.Now(),
		keyspace:  keyspace,
	}
}

func (s *Stats) IncrementCommands() {
	s.totalCommands.Add(1)
}

func (s *Stats) IncrementConnections() {
	s.totalConnections.Add(1)
}

func (s *Stats) GetInfo(ctx context.Context) string {
	uptime := time.Since(s.startedAt).Seconds()
	commands := s.totalCommands.Load()
	connections := s.totalConnections.Load()
	dbSize := s.keyspace.Size(ctx)

	info := fmt.Sprintf(`# Server
			redis_version:redis-like-go/1.0.0
			os:%s
			uptime_in_seconds:%.0f
			uptime_in_days:%.0f

			# Clients
			connected_clients:%d
			total_connections_received:%d

			# Stats
			total_commands_processed:%d
			keyspace_hits:0
			keyspace_misses:0

			# Keyspace
			db0:keys=%d
			`,
		runtime.GOOS,
		uptime,
		uptime/86400,
		connections,
		connections,
		commands,
		dbSize,
	)

	return info
}
