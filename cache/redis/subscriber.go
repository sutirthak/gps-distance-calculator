package cache

import (
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

func (r *RedisInstance) Subscribe(channel string, callback func(*redis.Message)){
	pubsub := r.RedisClient.Subscribe(r.Ctx, channel)
	defer pubsub.Close()
	for {
		msg, err := pubsub.ReceiveMessage(r.Ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"message": "error in getting message from subscribe method",
			}).Error(err)
			return
		}
		callback(msg)
	}
}
