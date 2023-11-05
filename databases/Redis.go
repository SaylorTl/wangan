package databases

import (
	"github.com/go-redis/redis"
	"gopkg.in/ini.v1"
)

var Redis *redis.Client

func InitRedisDb(conf *ini.File) (err error) {
	// redis配置信息
	redisHostName := conf.Section("redis").Key("REDIS_HOST").String()
	redisPort := conf.Section("redis").Key("REDIS_PORT").String()
	redisPwd := conf.Section("redis").Key("REDISs_PASSWORD").String()
	Redis = redis.NewClient(&redis.Options{
		Addr:     redisHostName + ":" + redisPort,
		Password: redisPwd, // no password set
		DB:       0,        // use default DB
	})
	_, err = Redis.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
