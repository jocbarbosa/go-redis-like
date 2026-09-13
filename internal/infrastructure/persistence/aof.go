package persistence

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/jocbarbosa/go-redis-like/internal/domain/repository"
)

const defaultFilePermission = 0644

type AOF struct {
	filepath string
	file     *os.File
	mu       *sync.Mutex
}

func NewAOF(filepath string) (repository.PersistenceRepository, error) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR|os.O_WRONLY, defaultFilePermission)
	if err != nil {
		return nil, err
	}

	return &AOF{
		filepath: filepath,
		file:     file,
	}, nil
}

func (a *AOF) Append(ctx context.Context, cmd string, args []string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	line := cmd
	if len(args) > 0 {
		line = fmt.Sprintf("%s %s", line, strings.Join(args, " "))
	}

	line += "\n"

	_, err := a.file.WriteString(line)
	if err != nil {
		return err
	}

	return a.file.Sync()
}

func (a *AOF) Replay(ctx context.Context, store repository.KeyValueRepository) error {
	return nil
}

func (a *AOF) Close() {}
