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
	ContextTimeoutEnabled bool

	// Maximum number of socket connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	PoolSize int

	// Minimum number of idle connections which is useful if it is slow to establish
	// a new connection.
	MinIdleConns int

	// Maximum number of idle connections.
	MaxIdleConns int

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActiveConns int

	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration

	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 30 minutes. -1 disables idle timeout check.
	ConnMaxIdleTime time.Duration

	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	ConnMaxLifetime time.Duration

	// TLS Config to use. When set, TLS will be negotiated.
	TLSConfig *tls.Config

	// Limiter interface used to implement circuit breaker or rate limiter.
	Limiter Limiter

	// Enables read only queries on slave/follower nodes.
	readOnly bool
}

// Limiter is the interface of a rate limiter or a circuit breaker.
type Limiter interface {
	// Allow returns nil if operation is allowed or an error otherwise.
	// If operation is allowed client must report the result of the operation
	// whether it is a success or a failure.
	Allow() error
	// ReportResult reports the result of the previously allowed operation.
	// nil indicates a success, non-nil error usually indicates a failure.
	ReportResult(result error)
}

func (opt *Options) init() {
	if opt.Addr == "" {
		opt.Addr = "localhost:6379"
	}
	if opt.Network == "" {
		if len(opt.Addr) > 0 && opt.Addr[0] == '/' {
			opt.Network = "unix"
		} else {
			opt.Network = "tcp"
		}
	}
	if opt.Protocol == 0 {
		opt.Protocol = 3
	}
	if opt.DB < 0 {
		opt.DB = 0
	}
	if opt.DialTimeout == 0 {
		opt.DialTimeout = 5 * time.Second
	}
	if opt.ReadTimeout == 0 {
		opt.ReadTimeout = 3 * time.Second
	}
	if opt.WriteTimeout == 0 {
		opt.WriteTimeout = opt.ReadTimeout
	}
	if opt.MaxRetries == 0 {
		opt.MaxRetries = 3
	}
	switch opt.MinRetryBackoff {
	case -1:
		opt.MinRetryBackoff = 0
	case 0:
		opt.MinRetryBackoff = 8 * time.Millisecond
	}
	switch opt.MaxRetryBackoff {
	case -1:
		opt.MaxRetryBackoff = 0
	case 0:
		opt.MaxRetryBackoff = 512 * time.Millisecond
	}
}
