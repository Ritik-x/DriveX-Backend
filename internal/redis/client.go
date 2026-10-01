package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	RDB *redis.Client

}





func NewRedisClient (addr string ) *Client{
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &Client{
		RDB: rdb,
	}
}

func (c *Client) Ping(ctx context.Context) error {
	return c.RDB.Ping(ctx).Err()
}