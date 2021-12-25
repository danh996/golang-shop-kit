package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(address, username, passsword string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Username: username,
		Password: passsword, // no password set
		DB:       0,         // use default DB
	})
	if cmd := client.Ping(context.Background()); cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return client, nil
}
