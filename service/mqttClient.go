package Service

import (
	. "GoServer/handle/device"
	. "GoServer/middleWare/dataBases/redis"
	. "GoServer/utils/config"
	. "GoServer/utils/log"
	. "GoServer/utils/threadWorker"
	M "github.com/eclipse/paho.mqtt.golang"
	"net/url"
)

func init() {
	Redis().Subscribe(func(pattern, channel, message string) error {
		var work Job = &ExpiredMsg{Pattern: pattern, Chann:channel, Message:message}
		InsertAsyncTask(work)
		SystemLog("notifyCallback  pattern:",pattern," channel : ", channel, " message : ", message)
		return nil
	}, "__keyevent@0__:expired")
}

var MessageCb M.MessageHandler = func(client M.Client, msg M.Message) {
	rOps := client.OptionsReader()
	servers := rOps.Servers()
	broker := servers[0]
	var work Job = &MqMsg{Broker: broker.Host, Topic: msg.Topic(), Payload: string(msg.Payload()), Direction: RECV_MQTT_MSG}
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
