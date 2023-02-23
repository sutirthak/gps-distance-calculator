package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func Subscriber(redisClient *redis.Client, ctx context.Context, channel string, callback  func(*redis.Message)) {
	pubsub := redisClient.Subscribe(ctx, channel)
	defer pubsub.Close()
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		callback (msg)
	}
}
