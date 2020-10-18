package device

import (
	. "GoServer/mqttPacket"
	. "GoServer/utils/log"
	. "GoServer/utils/string"
	. "GoServer/utils/threadWorker"
	. "GoServer/utils/time"
	M "github.com/eclipse/paho.mqtt.golang"
	"time"
)

var serverMap = make(map[string]interface{})

const (
	RECV_MQTT_MSG = 0
	SEND_MQTT_MSG = 1
)

type ExpiredMsg struct {
	Pattern string
	Chann string
	Message string
}

type MqMsg struct {
	Broker  string
	Topic   string
	Payload string
	Direction uint8
}

func GetMqttClient(brokerHost string) M.Client {
	broker := serverMap[brokerHost]
	if broker != nil {
		return broker.(M.Client)
	}
	return nil
}

func SetMqttClient(brokerHost string, handle interface{}) {
	if brokerHost != "" && handle != nil {
		serverMap[brokerHost] = handle
	}
}

func (expired *ExpiredMsg) ExecTask() error {
	DeviceExpiredMsgOps(expired.Pattern,expired.Chann,expired.Message)
	return nil
}

func (msg *MqMsg) ExecTask() error {
	switch msg.Direction {
	case RECV_MQTT_MSG: {
		ok, packet := MessageUnpack([]byte(msg.Payload))
		if ok && packet.Data != nil {
			deviceSN := GetDeviceSN(msg.Topic)
			//保存包数据入库
			SaveDeviceTransferDataOps(msg.Broker, deviceSN, packet)
			//处理包数据
			DeviceActBehaviorDataOps(packet, deviceSN, string(msg.Payload))
			MqttLog("[", msg.Broker, "] ===== ", packet.Json.ID, " =====>> ", msg.Topic, " time:", TimeFormat(time.Now()), "=========", GetGoroutineID(), GetWorkerQueueSize())
			MqttLog(packet.Data.(Protocol).Print())
		} else {
			MqttLog("analysis failed ->Topic:%s Payload:%s\n", msg.Topic, msg.Payload)
		}
	}
	case SEND_MQTT_MSG:
		Client := GetMqttClient(msg.Broker)
		if token := Client.Publish(msg.Topic,0,false,msg.Payload); token.Wait() && token.Error() != nil {
			MqttLog("Send MQ Message Error",msg)
			return token.Error()
		}
	}
	return nil
}

func (msg *MqMsg)Send(Broker, Topic string ,Payload string) {
	msg.Direction = SEND_MQTT_MSG
	msg.Broker = Broker
	msg.Topic = StringJoin([]interface{}{"/point_switch/",Topic})
	msg.Payload = Payload
	var work Job = msg
	InsertAsyncTask(work)
}
