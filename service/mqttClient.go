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

type ExpiredMsg struct {
	Pattern string
	Chann string
	Message string
}

type MqMsg struct {
	Broker  string
	Topic   string
	Payload []byte
}

func init() {
	Redis().Subscribe(func(pattern, channel, message string) error {
		var work Job = &ExpiredMsg{Pattern:pattern, Chann:channel, Message:message}
		InsertAsyncTask(work)
		SystemLog("notifyCallback  pattern:",pattern," channel : ", channel, " message : ", message)
		return nil
	}, "__keyevent@0__:expired")
}

func (expired *ExpiredMsg) ExecTask() error{
	DeviceExpiredMsgOps(expired.Pattern,expired.Chann,expired.Message)
	return nil
}

func (msg *MqMsg) ExecTask() error {
	ok, packet := MessageUnpack(msg.Payload)
	if ok && packet.Data != nil {
		deviceSN := GetDeviceSN(msg.Topic)
		//保存包数据入库
		SaveDeviceTransferDataOps(msg.Broker, deviceSN, packet)
		//处理包数据
		DeviceActBehaviorDataOps(packet, deviceSN, string(msg.Payload))
		MqttLog("[", msg.Broker, "] ===== ", packet.Json.ID, " =====>> ", msg.Topic, " time:", TimeFormat(time.Now()), "=========", GetGoroutineID(), GetWorkerQueueSize())
		MqttLog(packet.Data.(Protocol).Print())
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
