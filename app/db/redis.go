package db

import (
	"github.com/go-redis/redis"
	"github.com/astaxie/beego"
)

var client *redis.Client

func GetRedisClient() *redis.Client {
	if client == nil {
		beego.Debug("new redis client")
		client = redis.NewClient(&redis.Options{
			Addr:        "localhost:6379",
			Password:    "",
			DB:          0,
			IdleTimeout: -1,
		})
	}
	return client
}
