package redis

import (
	"time"

	"github.com/go-redis/redis"
)

// ClientOptions struct contains the options for connecting to redis
type ClientOptions struct {
	Host         string
	Port         string
	Password     string
	MaxRetries   int
	WriteTimeout time.Duration
	DB           int
}

// Client struct holds connection to redis
type Client struct {
	conn *redis.Client
}

// NewClient method will return a pointer to new client object
func NewClient(opts *ClientOptions) *Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         opts.Host + ":" + opts.Port,
		Password:     opts.Password,
		DB:           opts.DB,
		MaxRetries:   opts.MaxRetries,
		WriteTimeout: opts.WriteTimeout,
	})
	var client = &Client{conn: redisClient}
	return client
}

// GetConn returns a pointer to the underlying redis library
func (c *Client) GetConn() *redis.Client {
	return c.conn
}

// Close method closes the redis connection
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
