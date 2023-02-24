package cache

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/redis/go-redis/v9"
)

func (r *RedisInstance) ConnectToRedis(host, port, password string, db int) error {
	adress := fmt.Sprintf("%s:%s", host, port)
	r.RedisClient = redis.NewClient(&redis.Options{
		Addr:     adress,
		Password: password,
		DB:       db,
	})

	pong, err := r.RedisClient.Ping(r.Ctx).Result()
	if err != nil {
		return err

	}
	log.Info(pong)
	return nil
}
