package device

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/middleWare/dataBases/redis"
	. "GoServer/model"
	. "GoServer/model/device"
	mqtt "GoServer/mqtt"
	. "GoServer/utils/log"
	. "GoServer/utils/respond"
	. "GoServer/utils/time"
	"encoding/json"
	"fmt"
	M "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
)

const (
	CHARGEING_TIME = 150 //忙时150秒
	LEISURE_TIME   = 300 //空闲时300秒
)

var serverMap = make(map[string]interface{})

type RequestData struct {
	AccessWay     AccesswayType `form:"type" json:"type" binding:"required"`
	ModuleSN      string        `form:"module_sn" json:"module_sn" binding:"required"`
	ModuleVersion string        `form:"module_version" json:"module_version" binding:"required"`
	DeviceSN      string        `form:"device_sn" json:"device_sn" binding:"required"`
	DeviceVersion string        `form:"device_version" json:"device_version" binding:"required"`
	Token         string        `form:"token" json:"token" binding:"required"`
}

type FirmwareInfo struct {
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

func GetMqttClient(brokerHost string) *M.Client {
	broker := serverMap[brokerHost]
	if broker != nil {
		return broker.(*M.Client)
	}
	return nil
}

func SetMqttClient(brokerHost string, handle interface{}) {
	if brokerHost != "" && handle != nil {
		serverMap[brokerHost] = handle
	}
}

func createConnectLog(ctx *gin.Context, device_id int64, accessway AccesswayType, moduleSN string) {
	var log ModuleConnectLog
	log.Create(device_id, accessway, moduleSN, ctx.ClientIP())
	CreateAsyncSQLTask(ASYNC_MODULE_CONNECT_LOG, log)
}

func (data *RequestData) Connect(ctx *gin.Context) interface{} {
	var module ModuleInfo
	err := ExecSQL().Where("module_sn = ?", data.ModuleSN).First(&module).Error
	var hasRecord = true

	if err != nil {
		if IsRecordNotFound(err) {
			hasRecord = false
		} else {
			return CreateMessage(SYSTEM_ERROR, err)
		}
	}

	if !hasRecord {
		//创建对应关系 第一次连接仅建立关系，版本检测下次连接时在进行处理
		var info CreateDeviceInfo
		curTimestampMs := GetTimestampMs()
		info.Module.Create(data.AccessWay, data.ModuleSN, data.ModuleVersion)
		info.Device.Create(data.AccessWay, data.DeviceSN, data.DeviceVersion)
		info.Log.Create(0, data.AccessWay, data.ModuleSN, ctx.ClientIP())
		info.Module.CreateTime = curTimestampMs
		info.Device.CreateTime = curTimestampMs
		CreateAsyncSQLTask(ASYNC_DEV_AND_MODULE_CREATE, info)
		return CreateMessage(SUCCESS, nil)
	} else {
		module.Update(data.ModuleVersion)
		CreateAsyncSQLTask(ASYNC_UP_MODULE_VERSION, module)
		var device DeviceInfo
		device.Update(module.DeviceID, data.AccessWay, data.DeviceVersion, module.UpdateTime)
		device.DeviceSn = data.DeviceSN
		CreateAsyncSQLTask(ASYNC_UP_DEVICE_VERSION, device)
		// 检测并返回固件版本
		// 返回版本升级格式
		//	data := &FirmwareInfo{
		//		URL:  "http://www.gisunlink.com/GiSunLink.ota.bin",
		//		Size: 476448,
		//	}
		createConnectLog(ctx, module.ID, data.AccessWay, data.ModuleSN)
		return CreateMessage(SUCCESS, nil)
	}
	return CreateMessage(SUCCESS, nil)
}

func createDeviceTransferLog(transfer *DeviceTransferLog) {
	if transfer == nil {
		return
	}

	transfer.CreateTime = GetTimestampMs()
	if err := ExecSQL().Create(&transfer).Error; err != nil {
		SystemLog("CreateDeviceTransferLog", err)
	}
}

//保存上报数据入库
func SaveDeviceTransferData(serverNode string, device_sn string, packet *mqtt.Packet) {
	var comNum int64 = 0
	switch packet.Json.Behavior {
	case mqtt.GISUNLINK_CHARGEING, mqtt.GISUNLINK_CHARGE_LEISURE: //运行中,空闲中
		comList := packet.JsonData.(*mqtt.ComList)
		comNum = int64(comList.ComNum)
		break
	case mqtt.GISUNLINK_START_CHARGE, mqtt.GISUNLINK_CHARGE_FINISH, mqtt.GISUNLINK_CHARGE_NO_LOAD, mqtt.GISUNLINK_CHARGE_BREAKDOWN: //开始,完成,空载,故障
		comList := packet.JsonData.(*mqtt.ComList)
		comNum = int64(comList.ComNum)
		for _, comID := range comList.ComID {
			comNum = int64(comID)
		}
		break
	}

	log := &DeviceTransferLog{
		TransferID:   int64(packet.Json.ID),
		DeviceID:     Redis().GetDeviceIDFromRedis(device_sn, "deviceID"),
		TransferAct:  packet.Json.Act,
		DeviceSN:     device_sn,
		ComNum:       comNum,
		TransferData: packet.Json.Data,
		Behavior:     int64(packet.Json.Behavior),
		ServerNode:   serverNode,
		TransferTime: int64(packet.Json.Ctime),
	}
	createDeviceTransferLog(log)
}

func DeviceActBehaviorDataAnalysis(packet *mqtt.Packet, cacheKey string, playload string) {
	switch packet.Json.Behavior {
	case mqtt.GISUNLINK_CHARGEING, mqtt.GISUNLINK_CHARGE_LEISURE:
		{
			//写入原始数据
			Redis().UpdateDeviceRawDataToRedis(cacheKey, playload)

			//循环写入端口数据
			comList := packet.JsonData.(*mqtt.ComList)
			for _, comID := range comList.ComID {
				var index uint8 = comID
				if comList.ComNum <= 5 {
					index = (comID % 5)
				}
				comData := (comList.ComPort[int(index)]).(mqtt.ComData)
				comData.Id = comID
				fmt.Println("--------CurPower",comData.CurPower,"AveragePower",comData.AveragePower,"MaxPower",comData.MaxPower)
				//更新端口值
				Redis().UpdateDeviceComDataToRedis(cacheKey, comID, comData)
			}

			deviceStatus := &DeviceStatus{
				Behavior: comList.ComBehavior,
				Signal:   int8(comList.Signal),
				Worker: Redis().TatolWorkerByDevice(cacheKey, func(comDataString string) (enable bool, useEnergy uint32, useTime uint32, curElectricity uint32) {
					comData := &mqtt.ComData{}
					err := json.Unmarshal([]byte(comDataString), comData)
					if err != nil {
						return false, 0, 0, 0
					}

					if comData.Enable == 1 {
						return true, comData.UseEnergy, comData.UseTime, comData.CurElectricity
					}

					return false, 0, 0, 0
				}),
				ProtoVersion: comList.ComProtoVer,
			}

			var timeout int64 = CHARGEING_TIME
			if packet.Json.Behavior == mqtt.GISUNLINK_CHARGE_LEISURE {
				timeout = LEISURE_TIME
			}

			Redis().UpdateDeviceTokenExpiredTime(cacheKey, deviceStatus, timeout)

		}
	}
}
