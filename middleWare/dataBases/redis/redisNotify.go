package redis

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	//. "GoServer/utils/config"
	"time"
)

type PSubscribeCallback func(pattern, channel, message string)

var RedisNotify RedisSubscriber

type RedisSubscriber struct {
	client redis.PubSubConn
	cbMap  map[string]PSubscribeCallback
}
//
//func init() {
//
//	options, err := GetRedis()
//	if err != nil {
//		return
//	}
//
//	c, err := redis.Dial("tcp", options.Host+":"+options.Port,
//		redis.DialPassword(options.Auth),
//		redis.DialConnectTimeout(time.Duration(options.ConnectTimeout)*time.Second),
//		redis.DialReadTimeout(time.Duration(options.ReadTimeout)*time.Second),
//		redis.DialWriteTimeout(time.Duration(options.WriteTimeout)*time.Second))
//
//	if err != nil {
//		SystemLog("redis RedisNotify connect fail")
//		return
//	}
//
//	RedisNotify.client = redis.PubSubConn{c}
//	RedisNotify.cbMap = make(map[string]PSubscribeCallback)
//
//	go func() {
//		for {
//			//fmt.Printf("wait...")
//			switch res := RedisNotify.client.Receive().(type) {
//			case redis.Message:
//				pattern := (*string)(unsafe.Pointer(&res.Pattern))
//				channel := (*string)(unsafe.Pointer(&res.Channel))
//				message := (*string)(unsafe.Pointer(&res.Data))
//				RedisNotify.cbMap[*channel](*pattern, *channel, *message)
//				//pattern := res.Pattern
//				//channel := string(res.Channel)
//				//message := string(res.Data)
//				//if (pattern == "__keyspace@0__:blog*"){
//				//	switch  message {
//				//	case "set":
//				//		// do something
//				//		fmt.Println("set", channel)
//				//	case "del":
//				//		fmt.Println("del", channel)
//				//	case "expire":
//				//		fmt.Println("expire", channel)
//				//	case "expired":
//				//		fmt.Println("expired", channel)
//				//	}
//				//}
//			case redis.Subscription:
//				fmt.Printf("------------>%s: %s %d\n", res.Channel, res.Kind, res.Count)
//			case error:
//				//				log.Error("error handle...")
//				continue
//			}
//		}
//	}()
//}

func listenPubSubChannels(ctx context.Context, redisServerAddr string,
	onMessage func(channel string, data []byte) error,
	channels ...string) error {
	// A ping is set to the server with this period to test for the health of
	// the connection and server.
	const healthCheckPeriod = time.Minute

	c, err := redis.Dial("tcp", redisServerAddr,
		// Read timeout on server should be greater than ping period.
		redis.DialPassword("xyCbwbCcRyjfAHAP"),
		redis.DialReadTimeout(healthCheckPeriod+10*time.Second),
		redis.DialWriteTimeout(10*time.Second))
	if err != nil {
		return err
	}
	defer c.Close()

	psc := redis.PubSubConn{Conn: c}

	if err := psc.Subscribe(redis.Args{}.AddFlat(channels)...); err != nil {
		return err
	}

	done := make(chan error, 1)

	// Start a goroutine to receive notifications from the server.
	go func() {
		for {
			switch n := psc.Receive().(type) {
			case error:
				done <- n
				return
			case redis.Message:
				if err := onMessage(n.Channel, n.Data); err != nil {
					done <- err
					return
				}
			case redis.Subscription:
				switch n.Count {
				case len(channels):

				case 0:
					// Return from the goroutine when all channels are unsubscribed.
					done <- nil
					return
				}
			}
		}
	}()

	ticker := time.NewTicker(healthCheckPeriod)
	defer ticker.Stop()
loop:
	for err == nil {
		select {
		case <-ticker.C:
			// Send ping to test health of connection and server. If
			// corresponding pong is not received, then receive on the
			// connection will timeout and the receive goroutine will exit.
			if err = psc.Ping(""); err != nil {
				break loop
			}
		case <-ctx.Done():
			break loop
		case err := <-done:
			// Return error from the receive goroutine.
			return err
		}
	}

	// Signal the receiving goroutine to exit by unsubscribing from all channels.
	psc.Unsubscribe()

	// Wait for goroutine to complete.
	return <-done
}

// This example shows how receive pubsub notifications with cancelation and
// health checks.
func ExamplePubSubConn() {
	ctx, cancel := context.WithCancel(context.Background())

	err := listenPubSubChannels(ctx,
		"172.16.0.8:6379",
		func(channel string, message []byte) error {
			fmt.Printf("channel: %s, message: %s\n", channel, message)

			// For the purpose of this example, cancel the listener's context
			// after receiving last message sent by publish().
			if string(message) == "goodbye" {
				cancel()
			}
			return nil
		},
		"c1", "__keyevent@0__:expired")

	if err != nil {
		fmt.Println(err)
		return
	}

	// Output:
	// channel: c1, message: hello
	// channel: c2, message: world
	// channel: c1, message: goodbye
}


func (redisNotify *RedisSubscriber) Subscribe(channel interface{}, cb PSubscribeCallback) {
	ExamplePubSubConn()
	//err := redisNotify.client.PSubscribe(channel)
	//if err != nil {
	//	fmt.Printf("redis Subscribe error.")
	//	return
	//}
	//redisNotify.cbMap[channel.(string)] = cb
}