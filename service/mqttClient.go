package Service

import (
	. "GoServer/dataBases/redis"
	. "GoServer/mqtt"
	. "GoServer/utils"
	"fmt"
	M "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	//	"reflect"
	"time"
)

type MqMsg struct {
	Broker  string
	Topic   string
	Payload []byte
}

func (msg *MqMsg) ExecTask() error {
	ok, packet := MessageHandler(msg.Payload)
	rd := Redis().Get()
	_, err := rd.Do("SET", msg.Topic, string(msg.Payload))
	if err != nil {
		WebLog("lPop websocket user msg from queue failed", zap.String("cacheKey", msg.Topic), zap.Error(err))
	}

	if ok && packet.JsonData != nil {
		MqttLog("[", msg.Broker, "] =========>>", msg.Topic, " time:", TimeFormat(time.Now()), "=========", GetGoroutineID(), GetWorkerQueueSize())
		MqttLog(packet.JsonData.(Protocol).Print())
	} else {
		fmt.Printf("analysis failed ->Topic:%s Payload:%s\n", msg.Topic, msg.Payload)
	}
	return nil
}

var MessageCb M.MessageHandler = func(client M.Client, msg M.Message) {
	rOps := client.OptionsReader()
	servers := rOps.Servers()
	broker := servers[0]
	var work Job = &MqMsg{Broker: broker.Host, Topic: msg.Topic(), Payload: msg.Payload()}
	InsertAsynTask(work)
}

func StartMqttService() error {
	mqttOptions ,_ := GetMqtt()
	for _, mqtt := range mqttOptions {
		opts := M.NewClientOptions().AddBroker(mqtt.Host)
		opts.SetClientID(mqtt.Token)
		opts.SetUsername(mqtt.Name)
		opts.SetPassword(mqtt.Pwsd)
		opts.SetAutoReconnect(true)
		opts.SetDefaultPublishHandler(MessageCb)

		Client := M.NewClient(opts)
		if token := Client.Connect(); token.Wait() && token.Error() != nil {
			return token.Error()
		}

		if token := Client.Subscribe("/#", 0, nil); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}
	return nil
}
