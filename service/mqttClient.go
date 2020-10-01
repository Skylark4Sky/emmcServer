package Service

import (
	. "GoServer/middleWare/dataBases/redis"
	. "GoServer/mqtt"
	. "GoServer/utils/config"
	. "GoServer/utils/log"
	. "GoServer/utils/threadWorker"
	. "GoServer/utils/time"
	"fmt"
	M "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"time"
)

type MqMsg struct {
	Broker  string
	Topic   string
	Payload []byte
}

func behaviorHandle(packet *Packet, cacheKey string, playload string) {
	switch packet.Json.Behavior {
	case GISUNLINK_CHARGEING, GISUNLINK_CHARGE_LEISURE:
		{
			rd := Redis().Get()
			defer rd.Close()

			var timeout int = 45
			if packet.Json.Behavior == GISUNLINK_CHARGE_LEISURE {
				timeout = 120
			}

			key := cacheKey[11:]

			_, err := rd.Do("SET", key, playload, "ex", timeout)
			if err != nil {
				SystemLog("set redis value", zap.String("cacheKey", cacheKey), zap.Error(err))
			}
		}
	}

}

func (msg *MqMsg) ExecTask() error {
	ok, packet := MessageHandler(msg.Payload)

	if ok && packet.JsonData != nil {
		behaviorHandle(packet, msg.Topic, string(msg.Payload))
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
	mqttOptions, _ := GetMqtt()
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
