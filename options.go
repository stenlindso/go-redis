package redis

import (
	"context"
	"crypto/tls"
	"net"
	"time"
)

// Options keeps the settings to setup redis connection.
type Options struct {
	// The network type, either tcp or unix.
	// Default is tcp.
	Network string

	// host:port address.
	Addr string

	// ClientName will execute the `CLIENT SETNAME ClientName` command for each conn.
	ClientName string

	// Dialer creates new network connection and has priority over
	// Network and Addr options.
	Dialer func(ctx context.Context, network, addr string) (net.Conn, error)

	// Hook that is called when a new connection is established.
	OnConnect func(ctx context.Context, cn *Conn) error

	// Protocol 2 or 3. Use the version to negotiate RESP version with redis-server.
	// Default is 3.
	Protocol int

	// Use the specified username to authenticate the current connection
	// with one of the connections defined in the ACL list when connecting
	// to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
	Username string

	// Optional password. Must match the password specified in the
	// requirepass server configuration option (if connecting to a Redis 5.0 instance, or lower),
	// or the User Password when connecting to a Redis 6.0 instance, or greater,
	// that is using the Redis ACL system.
	Password string

	// CredentialsProvider allows the username and password to be updated
	// before reconnecting. It should return the current username and password.
	CredentialsProvider func() (username string, password string)

	// Database to be selected after connecting to the server.
	DB int

	// Maximum number of retries before giving up.
	// Default is 3 retries; -1 (not 0) disables retries.
	MaxRetries int

	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	MinRetryBackoff time.Duration

	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff time.Duration

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration

	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Supported values:
	//   - `0` - default timeout (3 seconds)
	//   - `-1` - no timeout (block indefinitely)
	//   - `-2` - disables SetReadDeadline calls completely
	ReadTimeout time.Duration

	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking. Supported values:
	//   - `0` - default timeout (ReadTimeout)
	//   - `-1` - no timeout (block indefinitely)
	//   - `-2` - disables SetWriteDeadline calls completely
	WriteTimeout time.Duration

	// ContextTimeoutEnabled controls whether the client respects context timeouts and deadlines.
	// Default is false (for backwards compatibility).
	// NOTE(personal): I prefer this enabled by default for new projects — remember to set
	// this to true when initializing clients, since it won't be changed here to avoid
	// breaking existing behavior.
	ContextTimeoutEnabled bool

	// Maximum number of socket connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	PoolSize int

	// Minimum number of idle connect
