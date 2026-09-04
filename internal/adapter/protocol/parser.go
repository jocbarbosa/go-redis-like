package protocol

import (
	"fmt"
	"strings"

	"github.com/jocbarbosa/go-redis-like/internal/domain/command"
)

type Command struct {
	Type command.Type
	Args []string
}

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseCommand(line string) (*Command, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("empty command")
	}

	parts := strings.Fields(line)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty line")
	}

	cmdType := command.Type(strings.ToUpper(parts[0]))
	if !cmdType.IsValid() {
		return nil, fmt.Errorf("invalid command : %s", parts[0])
	}

	cmd := &Command{
		Type: cmdType,
		Args: []string{},
	}

	if len(parts) > 1 {
		cmd.Args = parts[1:]
	}

	return cmd, nil
}

func (p *Parser) FormatResponse(result any) string {
	switch v := result.(type) {
	case string:
		return v
	case int, int64:
		return fmt.Sprintf("%d", v)
	case bool:
		if v {
			return p.FormatOK()
		}
		return "ERR operation failed"
	case error:
		return p.FormarError(v)
	default:
		return fmt.Sprintf("%v", result)
	}
}

func (p *Parser) FormatOK() string {
	return "OK"
}

func (p *Parser) FormarError(msg error) string {
	return fmt.Sprintf("ERR: %s", msg.Error())
}

func (p *Parser) FormatNil() string {
	return "nil"
}
