package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	. "GoServer/utils/config"
	"time"
	"unsafe"
)

type RedisSubscriberCallback func(pattern, channel, message string)

var redisNotify RedisSubscriber

type RedisSubscriber struct {
	client redis.PubSubConn
	cbMap  map[string]RedisSubscriberCallback
}

func init ()  {
	options, err := GetRedis()
	if err != nil {
		return
	}
	const healthCheckPeriod = time.Minute
	conn, err := redis.Dial("tcp",  options.Host+":"+options.Port,
		redis.DialPassword(options.Auth),
		redis.DialReadTimeout(healthCheckPeriod+10*time.Second),
		redis.DialWriteTimeout(10*time.Second))
	if err != nil {
		return
	}

	redisNotify.client = redis.PubSubConn{Conn: conn}
	redisNotify.cbMap = make(map[string]RedisSubscriberCallback)

	// Start a goroutine to receive notifications from the server.
	go func() {
		for {
			switch res := redisNotify.client.Receive().(type) {
			case error:
				continue;
			case redis.Message:
				pattern := (*string)(unsafe.Pointer(&res.Pattern))
				channel := (*string)(unsafe.Pointer(&res.Channel))
				message := (*string)(unsafe.Pointer(&res.Data))
				redisNotify.cbMap[*channel](*pattern, *channel, *message)
			case redis.Subscription:
				fmt.Printf("------------>%s: %s %d\n", res.Channel, res.Kind, res.Count)
			}
		}
	}()
	//return nil
}

func RedisNotifySubscribe(channel interface{}, cb RedisSubscriberCallback) {
	err := redisNotify.client.Subscribe(channel)
	if err != nil {
		fmt.Printf("redis Subscribe error.")
		return
	}
	redisNotify.cbMap[channel.(string)] = cb
}