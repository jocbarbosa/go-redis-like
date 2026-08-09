# go-redis-like

An in-memory key-value database written in Go, inspired by Redis. It speaks a
simple plain-text protocol over TCP, supports key expiration, and can optionally
persist every write to an Append-Only File so data survives a restart.

## Features

### 11 Commands

**Basic:** `SET`, `GET`, `DEL`, `EXPIRE`, `TTL`, `PERSIST`

**Advanced:** `KEYS`, `EXISTS`, `PING`, `INFO`, `QUIT`

### Thread-Safe

All store operations are guarded by a `sync.RWMutex`, so concurrent clients can
read in parallel while writes are serialized.

### Background Expiration

A background goroutine periodically sweeps expired keys, so memory is reclaimed
even for keys that are never read again.

### AOF Persistence (optional)

Write commands are appended to an Append-Only File. On startup the file is
replayed automatically to rebuild the dataset.

### Simple TCP Protocol

Plain text, one command per line. Works with standard tools such as `netcat` and
`telnet` — no special client library required.

## Getting Started

### Requirements

- Go 1.21 or newer

### Build and run

```bash
go build -o go-redis-like .
./go-redis-like
```

### Connect

```bash
nc localhost 6379
```

```
PING
PONG
SET name joao
OK
GET name
joao
```

## Commands

| Command | Description |
| --- | --- |
| `SET key value` | Stores a value under `key`, overwriting any existing value. |
| `GET key` | Returns the value of `key`, or nil if it does not exist. |
| `DEL key` | Removes `key` and returns whether it existed. |
| `EXPIRE key seconds` | Sets a time-to-live on `key`. |
| `TTL key` | Returns the remaining time-to-live of `key` in seconds. |
| `PERSIST key` | Removes the expiration from `key`, making it permanent. |
| `KEYS` | Lists the keys currently stored. |
| `EXISTS key` | Reports whether `key` exists. |
| `PING` | Health check; replies `PONG`. |
| `INFO` | Returns server statistics. |
| `QUIT` | Closes the connection. |

## How It Works

- **Server** — accepts TCP connections and handles each client in its own
  goroutine.
- **Store** — an in-memory map protected by a `sync.RWMutex`, holding values plus
  optional expiration timestamps.
- **Expiration** — keys are checked lazily on access and swept periodically by a
  background cleanup goroutine.
- **AOF** — when enabled, write commands are appended to disk and replayed at
  startup to restore state.
