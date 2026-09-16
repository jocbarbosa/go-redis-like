package usecase

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jocbarbosa/go-redis-like/internal/adapter/protocol"
	"github.com/jocbarbosa/go-redis-like/internal/domain/command"
	"github.com/jocbarbosa/go-redis-like/internal/domain/repository"
)

type CommandHandler struct {
	store   repository.KeyValueRepository
	persist repository.PersistenceRepository
	parser  *protocol.Parser
	stats   *Stats
}

func NewCommandHandler(
	parser *protocol.Parser,
	store repository.KeyValueRepository,
	persist repository.PersistenceRepository,
	stats *Stats) *CommandHandler {
	return &CommandHandler{
		store:   store,
		persist: persist,
		parser:  parser,
		stats:   stats,
	}
}

func (h *CommandHandler) ExecuteCommand(ctx context.Context, cmd *protocol.Command) string {
	h.stats.IncrementCommands()

	switch cmd.Type {
	case command.SET:
		return h.handleSet(ctx, cmd.Args)
	case command.GET:
		return h.handleGet(ctx, cmd.Args)
	case command.DEL:
		return h.handleDelete(ctx, cmd.Args)
	case command.EXISTS:
		return h.handleExists(ctx, cmd.Args)
	case command.EXPIRE:
		return h.handleExpire(ctx, cmd.Args)
	case command.TTL:
		return h.handleTTL(ctx, cmd.Args)
	case command.KEYS:
		return h.handleKeys(ctx, cmd.Args)
	case command.PING:
		return h.handlePing(ctx, cmd.Args)
	case command.INFO:
		return h.handleInfo(ctx, cmd.Args)
	default:
		h.parser.FormatError(fmt.Errorf("unknow command type"))
	}

	return ""
}

func (h *CommandHandler) handleSet(ctx context.Context, args []string) string {
	if len(args) < 2 {
		return h.parser.FormatError(fmt.Errorf("invalid number of arguments for SET command: %d", len(args)))
	}

	key := args[0]
	value := strings.Join(args[1:], " ")

	h.store.Set(ctx, key, value)

	if h.persist != nil {
		h.persist.Append(ctx, command.SET.String(), args)
	}

	return h.parser.FormatOK()
}
func (h *CommandHandler) handleGet(ctx context.Context, args []string) string {
	if len(args) < 1 {
		return h.parser.FormatError(fmt.Errorf("invalid number of arguments for GET command: %d", len(args)))
	}

	value, found := h.store.Get(ctx, args[0])
	if found {
		return value
	}

	return h.parser.FormatNil()
}

func (h *CommandHandler) handleDelete(ctx context.Context, args []string) string {
	if len(args) < 1 {
		return h.parser.FormatError(fmt.Errorf("invalid number of arguments for DEL command: %d", len(args)))
	}

	count := 0
	for _, v := range args {
		count += h.store.Del(ctx, v)
	}

	if h.persist != nil && count != 0 {
		for _, v := range args {
			h.persist.Append(ctx, command.DEL.String(), []string{v})
		}
	}

	return h.parser.FormatResponse(count)
}

func (h *CommandHandler) handleExists(ctx context.Context, args []string) string {
	if len(args) < 1 {
		return h.parser.FormatError(fmt.Errorf("invalid number of arguments for EXISTS command: %d", len(args)))
	}

	count := 0
	for _, key := range args {
		if h.store.Exists(ctx, key) {
			count++
		}
	}

	return h.parser.FormatResponse(count)
}

func (h *CommandHandler) handleExpire(ctx context.Context, args []string) string {
	if len(args) < 1 {
		return h.parser.FormatError(fmt.Errorf("invalid number of arguments for EXPIRE command: %d", len(args)))
	}

	key := args[0]
	seconds, err := strconv.Atoi(args[1])
	if err != nil {
		return h.parser.FormatError(fmt.Errorf("the second value must be an integer"))
	}

	success := h.store.Expire(ctx, key, seconds)
	if success {
		if h.persist != nil {
			h.persist.Append(ctx, command.EXPIRE.String(), []string{key})
		}
		return h.parser.FormatOK()
	}

	return h.parser.FormatResponse(success)
}

func (h *CommandHandler) handleExpireAt(ctx context.Context, args []string) string {
	if len(args) < 1 {
		return h.parser.FormatError(fmt.Errorf("invalid number of arguments for EXISTS command: %d", len(args)))
	}
	return ""
}

func (h *CommandHandler) handleTTL(ctx context.Context, args []string) string {
	if len(args) < 1 {
		return h.parser.FormatError(fmt.Errorf("invalid number of arguments for TTL command: %d", len(args)))
	}
	return ""
}

func (h *CommandHandler) handleKeys(ctx context.Context, args []string) string {
	if len(args) < 1 {
		return h.parser.FormatError(fmt.Errorf("invalid number of arguments for KEYS command: %d", len(args)))
	}
	return ""
}

func (h *CommandHandler) handlePing(ctx context.Context, args []string) string {
	if len(args) < 1 {
		return h.parser.FormatError(fmt.Errorf("invalid number of arguments for PING command: %d", len(args)))
	}
	return ""
}

func (h *CommandHandler) handleInfo(ctx context.Context, args []string) string {
	if len(args) < 1 {
		return h.parser.FormatError(fmt.Errorf("invalid number of arguments for INFO command: %d", len(args)))
	}
	return ""
}
