package redis

import (
	. "GoServer/utils"
	"github.com/gomodule/redigo/redis"
	"sync"
	"time"
)

const (
	redisIndex = "redisIndex"
)

var (
	redisPool sync.Map
	mutex sync.Mutex
)

func newRedisClient(index string, re *RedisOptions) error {
	pool := &redis.Pool {
		MaxIdle:   re.MaxIdle,
		MaxActive: re.MaxOpen,
		Dial: func() (redis.Conn, error) {

			c, err := redis.Dial("tcp", re.Host+":"+re.Port,
				redis.DialPassword(re.Auth),
				redis.DialConnectTimeout(time.Duration(re.ConnectTimeout)*time.Second),
				redis.DialReadTimeout(time.Duration(re.ReadTimeout)*time.Second),
				redis.DialWriteTimeout(time.Duration(re.WriteTimeout)*time.Second))

			if err != nil {
				zaplog.Info("redis connect fail")
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		IdleTimeout: 300 * time.Second,
		Wait:        true,
	}

	conn := pool.Get()
	defer conn.Close()

	if _, err := conn.Do("PING"); err != nil {
		zaplog.Info("redis ping fail")
		return err
	} else {
		redisPool.Store(index, pool)
		return nil
	}
}

func Redis() *redis.Pool {
	if redisLoad, ok := redisPool.Load(redisIndex); ok {
		redisPool := redisLoad.(*redis.Pool)
		return redisPool
	}

	//加锁防并发初始化
	mutex.Lock()
	defer mutex.Unlock()

	options, err := GetRedis()
	if err != nil {
		panic(err)
	}

	if err := newRedisClient(redisIndex, options); err != nil {
		panic(err)
	}

	poolLoad, _ := redisPool.Load(redisIndex)
	redisPool := poolLoad.(*redis.Pool)
	return redisPool
}