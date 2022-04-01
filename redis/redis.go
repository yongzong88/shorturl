package redis

import (
	"time"

	"github.com/go-redis/redis"
)

var (
	rdb *redis.Client
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
}

func IsUsed(id string) bool {
	_, err := rdb.Get(id).Result()
	if err != nil {
		return false
	}
	return true
}

func Set(textid string, target string, expired int) error {
	return rdb.Set(textid, target, time.Duration(expired)*time.Second).Err()
}

func Get(textid string) (string, error) {
	return rdb.Get(textid).Result()
}
