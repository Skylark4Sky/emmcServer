package Service

import (
	. "GoServer/handle/device"
	. "GoServer/middleWare/dataBases/redis"
	deviceModel "GoServer/model/device"
	. "GoServer/mqtt"
	. "GoServer/utils/config"
	. "GoServer/utils/log"
	. "GoServer/utils/string"
	. "GoServer/utils/threadWorker"
	. "GoServer/utils/time"
	"encoding/json"
	"fmt"
	M "github.com/eclipse/paho.mqtt.golang"
	"github.com/gomodule/redigo/redis"
	"github.com/prometheus/common/log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var serverMap = make(map[string]interface{})

const (
	CHARGEING_TIME = 45
	LEISURE_TIME   = 120
)

type MqMsg struct {
	Broker  string
	Topic   string
	Payload []byte
}

func init() {
	RedisNotify.Subscribe("__keyevent@0__:expired", TestPubCallback)
}

func TestPubCallback(patter , chann, msg string){
	log.Debug( "TestPubCallback patter : " + patter + " channel : ", chann, " message : ", msg)
}

func behaviorHandle( packet *Packet, cacheKey string, playload string) {
	switch packet.Json.Behavior {
	case GISUNLINK_CHARGEING, GISUNLINK_CHARGE_LEISURE:
		{
			rd := Redis().Get()
			defer rd.Close()

			var timeout int = CHARGEING_TIME
			if packet.Json.Behavior == GISUNLINK_CHARGE_LEISURE {
				timeout = LEISURE_TIME
			}

			comList := packet.JsonData.(*ComList)

			deviceInfo := fmt.Sprintf("{\"status\"%d,\"signal\":%d,\"version\":%d}", comList.ComBehavior, int8(comList.Signal), comList.ComProtoVer)

			if SetRedisItem(rd, "HSET", cacheKey, "rawData", playload, "deviceInfo", deviceInfo) != nil {
				return
			}

			if SetRedisItem(rd, "expire", cacheKey, timeout) != nil {
				return
			}

			for _, comID := range comList.ComID {

				var index uint8 = comID
				if comList.ComNum <= 5 {
					index = (comID % 5)
				}

				comData := (comList.ComPort[int(index)]).(ComData)
				comData.Id = comID

				var jsonByte []byte
				var comIDString strings.Builder
				comIDString.WriteString("comPort")
				comIDString.WriteString(strconv.Itoa(int(comID)))

				jsonByte, err := json.Marshal(comData)
				if err == nil {
					SetRedisItem(rd, "HSET", cacheKey, comIDString.String(), string(jsonByte))
				}
			}
		}
	}
}

func saveTransferData(serverNode string, device_sn string, packet *Packet) {
	var comNum int64 = 0
	switch packet.Json.Behavior {
	case GISUNLINK_CHARGEING, GISUNLINK_CHARGE_LEISURE: //运行中,空闲中
		comList := packet.JsonData.(*ComList)
		comNum = int64(comList.ComNum)
		break
	case GISUNLINK_START_CHARGE,GISUNLINK_CHARGE_FINISH,GISUNLINK_CHARGE_NO_LOAD,GISUNLINK_CHARGE_BREAKDOWN: //开始,完成,空载,故障
		comList := packet.JsonData.(*ComList)
		comNum = int64(comList.ComNum)
		for _, comID := range comList.ComID {
			comNum = int64(comID)
		}
		break
	}

	rd := Redis().Get()
	defer rd.Close()

	var deviceID int64 = 0
	ItemValue := GetRedisItem(rd,"HGET",device_sn,"deviceID")
	
	if ItemValue != nil {
		deviceID ,_= redis.Int64(ItemValue,nil)
	}

	log := &deviceModel.DeviceTransferLog{
		TransferID:   int64(packet.Json.ID),
		DeviceID: deviceID,
		TransferAct:  packet.Json.Act,
		DeviceSN:     device_sn,
		ComNum:       comNum,
		TransferData: packet.Json.Data,
		Behavior:     int64(packet.Json.Behavior),
		ServerNode:   serverNode,
		TransferTime: int64(packet.Json.Ctime),
	}

	CreateDeviceTransferLog(log)

}

func (msg *MqMsg) ExecTask() error {
	ok, packet := MessageHandler(msg.Payload)

	if ok && packet.JsonData != nil {
		deviceSN := GetDeviceSN(msg.Topic)
		saveTransferData(msg.Broker, deviceSN, packet)
		behaviorHandle(packet, deviceSN, string(msg.Payload))
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

func GetMqttClient(brokerHost string) *M.Client {
	broker := serverMap[brokerHost]
	if broker != nil {
		return broker.(*M.Client)
	}
	return nil
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
			serverMap[key] = Client
		}
	}
	return nil
}
