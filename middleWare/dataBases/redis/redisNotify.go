package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
	"unsafe"
	. "GoServer/utils/config"
	. "GoServer/utils/log"
)

type PSubscribeCallback func(pattern, channel, message string)

var RedisNotify RedisSubscriber

type RedisSubscriber struct {
	client redis.PubSubConn
	cbMap  map[string]PSubscribeCallback
}

func init() {

	options, err := GetRedis()
	if err != nil {
		return
	}

	c, err := redis.Dial("tcp", options.Host+":"+options.Port,
		redis.DialPassword(options.Auth),
		redis.DialConnectTimeout(time.Duration(options.ConnectTimeout)*time.Second),
		redis.DialReadTimeout(time.Duration(options.ReadTimeout)*time.Second),
		redis.DialWriteTimeout(time.Duration(options.WriteTimeout)*time.Second))

	if err != nil {
		SystemLog("redis RedisNotify connect fail")
		return
	}

	RedisNotify.client = redis.PubSubConn{c}
	RedisNotify.cbMap = make(map[string]PSubscribeCallback)

	go func() {
		for {
			fmt.Printf("wait...")
			switch res := RedisNotify.client.Receive().(type) {
			case redis.Message:
				pattern := (*string)(unsafe.Pointer(&res.Pattern))
				channel := (*string)(unsafe.Pointer(&res.Channel))
				message := (*string)(unsafe.Pointer(&res.Data))
				RedisNotify.cbMap[*channel](*pattern, *channel, *message)
				//pattern := res.Pattern
				//channel := string(res.Channel)
				//message := string(res.Data)
				//if (pattern == "__keyspace@0__:blog*"){
				//	switch  message {
				//	case "set":
				//		// do something
				//		fmt.Println("set", channel)
				//	case "del":
				//		fmt.Println("del", channel)
				//	case "expire":
				//		fmt.Println("expire", channel)
				//	case "expired":
				//		fmt.Println("expired", channel)
				//	}
				//}
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", res.Channel, res.Kind, res.Count)
			case error:
				//				log.Error("error handle...")
				continue
			}
		}
	}()
}

func (redisNotify *RedisSubscriber) Subscribe(channel interface{}, cb PSubscribeCallback) {
	err := redisNotify.client.PSubscribe(channel)
	if err != nil {
		fmt.Printf("redis Subscribe error.")
		return
	}
	redisNotify.cbMap[channel.(string)] = cb
}