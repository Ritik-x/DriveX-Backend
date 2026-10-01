package redis

import (
	"context"
	"encoding/json"
	"time"

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

func (c *Client) SetJSON(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) error {

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.RDB.Set(
		ctx,
		key,
		data,
		expiration,
	).Err()
}
func (c *Client) GetJSON(
	ctx context.Context,
	key string,
	result interface{},
) error {

	data, err := c.RDB.Get(ctx, key).Bytes()

	if err != nil {
		return err
	}

	return json.Unmarshal(data, result)
}