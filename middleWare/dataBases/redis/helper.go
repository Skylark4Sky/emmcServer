package redis

import (
	"github.com/gomodule/redigo/redis"
)

// Int is a helper that converts a command reply to an integer
func RedisInt(reply interface{}, err error) (int, error) {
	return redis.Int(reply, err)
}

// Int64 is a helper that converts a command reply to 64 bit integer
func RedisInt64(reply interface{}, err error) (int64, error) {
	return redis.Int64(reply, err)
}

// String is a helper that converts a command reply to a string
func RedisString(reply interface{}, err error) (string, error) {
	return redis.String(reply, err)
}

// Bool is a helper that converts a command reply to a boolean
func RedisBool(reply interface{}, err error) (bool, error) {
	return redis.Bool(reply, err)
}
