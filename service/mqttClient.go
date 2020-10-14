package Service

import (
	. "GoServer/handle/device"
	. "GoServer/middleWare/dataBases/redis"
	. "GoServer/mqtt"
	. "GoServer/utils/config"
	. "GoServer/utils/log"
	. "GoServer/utils/string"
	. "GoServer/utils/threadWorker"
	. "GoServer/utils/time"
	"fmt"
	M "github.com/eclipse/paho.mqtt.golang"
	"net/url"
	"time"
)

type MQOffLine struct {
	key string
}

type MqMsg struct {
	Broker  string
	Topic   string
	Payload []byte
}

func init() {
	Redis().Subscribe(func(chann string, msg []byte) error {

		var work Job = &MQOffLine{key: string(msg)}
		InsertAsyncTask(work)

		SystemLog("notifyCallback  channel : ", chann, " message : ", msg)
		return nil
	}, "__keyevent@0__:expired")
}

func (offline *MQOffLine) ExecTask() {

}

func (msg *MqMsg) ExecTask() error {
	ok, packet := MessageHandler(msg.Payload)
	if ok && packet.JsonData != nil {
		deviceSN := GetDeviceSN(msg.Topic)
		SaveDeviceTransferData(msg.Broker, deviceSN, packet)
		DeviceActBehaviorDataAnalysis(packet, deviceSN, string(msg.Payload))
		MqttLog("[", msg.Broker, "] ===== ", packet.Json.ID, " =====>> ", msg.Topic, " time:", TimeFormat(time.Now()), "=========", GetGoroutineID(), GetWorkerQueueSize())
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
	InsertAsyncTask(work)
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

		URL, err := url.Parse(mqtt.Host)

		if err == nil {
			key := URL.Host
			SetMqttClient(key, Client)
		}
	}
	return nil
}
