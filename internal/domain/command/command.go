package command

type Type string

const (
	SET     Type = "SET"
	GET     Type = "GET"
	DEL     Type = "DEL"
	EXPIRE  Type = "EXPIRE"
	TTL     Type = "TTL"
	PERSIST Type = "PERSIST"
	QUIT    Type = "QUIT"

	KEYS   Type = "KEYS"
	EXISTS Type = "EXISTS"
	PING   Type = "PING"
	INFO   Type = "INFO"
)

// IsValid validates if the command is acceptable
func (t Type) IsValid() bool {
	switch t {
	case SET, GET, DEL, EXPIRE, TTL, PERSIST, QUIT, KEYS, EXISTS, PING, INFO:
		return true
	default:
		return false
	}
}

// String returns a string of [Type]
func (t Type) String() string {
	return string(t)
}

// IsWriteCommand validates if the command is a write command
func (t Type) IsWriteCommand() bool {
	switch t {
	case SET, DEL, EXPIRE, PERSIST:
		return true
	default:
		return false
	}
}
