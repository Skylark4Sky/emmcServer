package Service

import (
	. "GoServer/packet"
	. "GoServer/utils"
	"fmt"
	M "github.com/eclipse/paho.mqtt.golang"
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
	if ok && packet.JsonData != nil {
		PrintInfo("[", msg.Broker, "] =========>>", msg.Topic, " time:", time.Now().Format(GetSystem().Timeformat), "=========", GetGoroutineID(), GetWorkerQueueSize())
		PrintInfo(packet.JsonData.(Protocol).Print())
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
	for _, mqtt := range GetMqtt() {

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
