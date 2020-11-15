package redis

import (
	. "GoServer/utils/log"
	. "GoServer/utils/string"
	"github.com/gomodule/redigo/redis"
)

const (
	REDIS_LOCK_DEFAULTIMEOUT = 10
)

type RedisLock struct {
	resource string
	token    string
	conn     *Cacher
	timeout  int
}

func (lock *RedisLock) tryLock() (ok bool, err error) {
	_, err = redis.String(lock.conn.Do("SET", lock.key(), lock.token, "EX", int(lock.timeout), "NX"))
	if err == redis.ErrNil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (lock *RedisLock) Unlock() (err error) {
	_, err = lock.conn.Do("del", lock.key())
	return
}

func (lock *RedisLock) key() string {
	return StringJoin([]interface{}{"redislock:", lock.resource})
}

func (lock *RedisLock) AddTimeout(ex_time int64) (ok bool, err error) {
	ttl_time, err := redis.Int64(lock.conn.Do("TTL", lock.key()))
	if err != nil {
		SystemLog("redis get failed:", err)
	}
	if ttl_time > 0 {
		_, err := redis.String(lock.conn.Do("SET", lock.key(), lock.token, "EX", int(ttl_time+ex_time)))
		if err == redis.ErrNil {
			return false, nil
		}
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

func TryLock(resource string, token string, defaulTimeout int) (lock *RedisLock, ok bool, err error) {
	return TryLockWithTimeout(resource, token, defaulTimeout)
}

func TryLockWithTimeout( resource string, token string, timeout int) (lock *RedisLock, ok bool, err error) {
	lock = &RedisLock{resource, token, Redis(), timeout}

	ok, err = lock.tryLock()

	if !ok || err != nil {
		lock = nil
	}

	return
}