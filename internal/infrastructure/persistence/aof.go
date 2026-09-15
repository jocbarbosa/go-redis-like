package persistence

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/jocbarbosa/go-redis-like/internal/domain/command"
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
	if ctx.Err() != nil {
		return ctx.Err()
	}

	file, err := os.Open(a.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		cmd := strings.ToUpper(fields[0])
		args := fields[1:]

		switch command.Type(cmd) {
		case command.SET:
			if len(args) < 2 {
				continue
			}

			key := args[0]
			value := strings.Join(args[1:], " ")
			store.Set(ctx, key, value)
		case command.EXPIRE:
			if len(args) < 2 {
				continue
			}

			key := args[0]
			seconds, err := strconv.Atoi(args[1])
			if err != nil {
				continue
			}

			store.Expire(ctx, key, seconds)
		case command.DEL:
			if len(args) < 1 {
				continue
			}

			store.Del(ctx, args[0])
		default:
		}

	}

	if scanner.Err() != nil {
		return fmt.Errorf("error reading AOF file %w", err)
	}

	return nil
}

func (a *AOF) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.file != nil {
		a.file.Close()
	}

	return nil
}
