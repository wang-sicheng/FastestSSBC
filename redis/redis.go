package redis

import (
	"context"
	"github.com/cloudflare/cfssl/log"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{})

//初始化
func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

//set
func SetIntoRedis(key string, value string) error {
	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
	return err
}

//get
func GetFromRedis(key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Errorf("the key:%s does not exist\n", key)
		return "", nil
	} else if err != nil {
		panic(err)
	} else {
		return val, nil
	}
}
