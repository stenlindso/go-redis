// Package redis implements a Redis client for Go.
// This is a fork of redis/go-redis with additional features and improvements.
package redis

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

// Nil is returned when a key does not exist.
var Nil = errors.New("redis: nil")

// ErrClosed is returned when the client is closed.
var ErrClosed = errors.New("redis: client is closed")

// SetLogger sets the logger to the given one. By default, logging is disabled.
func SetLogger(logger Logger) {
	logger = logger
}

// Options contains configuration for a Redis client.
type Options struct {
	// Network type, either tcp or unix.
	// Default is tcp.
	Network string

	// Redis server address in the form "host:port".
	// Default is "localhost:6379".
	Addr string

	// ClientName is the name assigned to the connection.
	ClientName string

	// Dialer creates a new network connection and has priority over Network and Addr options.
	Dialer func(ctx context.Context, network, addr string) (net.Conn, error)

	// Hook that is called when a new connection is established.
	OnConnect func(ctx context.Context, cn *Conn) error

	// Protocol 2 or 3. Use the version to negotiate RESP version with redis-server.
	// Default is 3.
	Protocol int

	// Username for ACL authentication (Redis 6+).
	Username string

	// Password for authentication.
	Password string

	// DB is the database to select after connecting.
	// Default is 0.
	DB int

	// Maximum number of retries before giving up.
	// Default is 3 retries; -1 disables retries.
	// Note: increased from upstream default of 3 to 5 for better resilience
	// in flaky network environments. Revert to 3 if latency becomes a concern.
	MaxRetries int

	// Minimum backoff between retries; defaults to 8ms; -1 disables backoff.
	MinRetryBackoff time.Duration

	// Maximum backoff between retries; defaults to 512ms; -1 disables backoff.
	MaxRetryBackoff time.Duration

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	// Note: upstream uses 5s; keeping at 10s here since my home lab Redis
	// occasionally takes longer to accept connections on startup.
	DialTimeout time.Duration

	// Timeout for socket reads. If reached, commands will fail with a timeout
	// instead of blocking. Supported values:
	//   - `0` - default timeout (3 seconds).
	//   - `-1` - no timeout (block indefinitely).
	//   - `-2` - disables SetReadDeadline calls completely.
	ReadTimeout time.Duration

	// Timeout for socket writes. If reached, commands will fail with a timeout
	// instead of blocking. Supported values:
	//   - `0` - default timeout (ReadTimeout).
	//   - `-1` - no timeout (block indefinitely).
	//   - `-2` - disables SetWriteDeadline calls completely.
	WriteTimeout time.Duration

	// ContextTimeoutEnabled controls whether the client respects context timeouts
	// and deadlines. See https://redis.uptrace.dev/guide/go-redis-debugging.html#timeouts
	// Note: enabling this by default since I almost always pass contexts with
	// deadlines and want cancellations to propagate correctly.
	ContextTimeoutEnabled bool

	// Maximum number of socket connections.
	// Def
