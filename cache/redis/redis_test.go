package cache

import (
	"context"
	"testing"
)

var ctx = context.Background()

func TestLocalhostConnection(t *testing.T) {
	redisInstance := RedisInstance{Ctx: ctx}
	err := redisInstance.ConnectToRedis("localhost", "6379", "", 0)
	if err != nil {
		t.Errorf("localhost connection test case faild")
	}
}

func TestFMDPConnection(t *testing.T) {
	redisInstance := RedisInstance{Ctx: ctx}
	err := redisInstance.ConnectToRedis("fmdp-staging.ddnsfree.com", "6379", "", 0)
	if err != nil {
		t.Errorf("FMDP connection test case faild")
	}
}


