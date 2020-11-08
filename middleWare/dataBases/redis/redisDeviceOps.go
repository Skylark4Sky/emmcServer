package redis

import (
	mqtt "GoServer/mqttPacket"
	. "GoServer/utils/float64"
	. "GoServer/utils/string"
	"encoding/json"
	"strconv"
)

const (
	REDIS_DEVICE_TOEKN_KEY = "token:"
	REDIS_COM_DATA_KEY     = "comData:"
	REDIS_DEVICE_INFO_KEY  = "info:"
)

const (
	REDIS_DEVICE_ID_FIELD    		= "id"
	REDIS_DEVICE_DATA_TOTAL_FIELD   = "total"
	REDIS_DEVICE_STATUS			    = "status"
	REDIS_RAW_DATA_FIELD     		= "rawData"
)


type DeviceTatolInfo struct {
	UseEnergy      uint64 `json:"energy"`
	UseTime        uint64 `json:"time"`
	CurElectricity uint64 `json:"electricity"`
	CurPower       string `json:"c_power"`
	AveragePower   string `json:"a_power"`
	MaxPower       string `json:"m_power"`
	EnableCount    uint8  `json:"enable"`
}

type DeviceStatus struct {
	Behavior     uint8  `json:"behavior"`
	Signal       int8   `json:"signal"`
	Worker       uint8  `json:"worker"`
	ProtoVersion uint16 `json:"protoVersion"`
}

func GetDeviceTokenKey(deviceSN string) string {
	return StringJoin([]interface{}{REDIS_DEVICE_TOEKN_KEY, deviceSN})
}

func GetComdDataKey(deviceSN string) string {
	return StringJoin([]interface{}{REDIS_COM_DATA_KEY, deviceSN})
}

func GetDeviceInfoKey(deviceSN string) string {
	return StringJoin([]interface{}{REDIS_DEVICE_INFO_KEY, deviceSN})
}

//更新令牌时间
func (c *Cacher) UpdateDeviceTokenExpiredTime(deviceSN string, status *DeviceStatus, timeout int64) {
	c.Set(GetDeviceTokenKey(deviceSN), status, timeout) //15分钟过期,正常1-2分钟后续数据就上来了
}

//插入对应令牌
func (c *Cacher) InitWithInsertDeviceIDToken(deviceSN string, deviceID uint64) {
	c.Set(GetDeviceTokenKey(deviceSN), deviceID, 900) //15分钟过期,正常1-2分钟后续数据就上来了
	c.HSet(GetDeviceInfoKey(deviceSN), REDIS_DEVICE_ID_FIELD, deviceID) //插入设备ID
	c.HSet(GetDeviceInfoKey(deviceSN), REDIS_DEVICE_STATUS, 1) 	//设备在线状态
}

//取设备状态
func (c *Cacher) GetDeviceStatusFromRedis(deviceSN string) int {
	status, err := c.HGetInt(GetDeviceInfoKey(deviceSN), REDIS_DEVICE_STATUS)
	if err != nil {
		return 0
	}
	return status
}

//设置设备状态
func (c *Cacher) SetDeviceStatusFromRedis(deviceSN string, status int) {
	c.HSet(GetDeviceInfoKey(deviceSN), REDIS_DEVICE_STATUS, status) //全局设备列表
}

//取设备ID
func (c *Cacher) GetDeviceIDFromRedis(deviceSN string) uint64 {
	deviceID, err := c.HGetUint64(GetDeviceInfoKey(deviceSN), REDIS_DEVICE_ID_FIELD)
	if err != nil {
		return 0
	}
	return deviceID
}

// 获取某一端口数据
func (c *Cacher) GetDeviceComDataFormRedis(deviceSN string, comID uint8, comData interface{}) error {
	str, err := RedisString(c.HGet(GetComdDataKey(deviceSN), strconv.Itoa(int(comID))))
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(str), comData)
}

//更新端口数据
func (c *Cacher) UpdateDeviceComDataToRedis(deviceSN string, comID uint8, comData interface{}) {
	c.HSet(GetComdDataKey(deviceSN), strconv.Itoa(int(comID)), comData)
}

//更新工作状态数据统计
func (c *Cacher) updateDeviceTatolInfoToRedis(deviceSN string, infoData interface{}) {
	c.HSet(GetDeviceInfoKey(deviceSN),REDIS_DEVICE_DATA_TOTAL_FIELD, infoData)
}

//统计当前工作端口数量
func (c *Cacher) TatolWorkerByDevice(deviceSN string, comDataMap map[uint8]mqtt.ComData) uint8 {

	deviceInfo := &DeviceTatolInfo{
		UseEnergy:      0,
		UseTime:        0,
		CurElectricity: 0,
		CurPower:       "0w",
		AveragePower:   "0w",
		MaxPower:       "0w",
		EnableCount:    0,
	}

	var maxPower float64 = 0
	for _, comData := range comDataMap {
		if comData.Enable == 1 {
			deviceInfo.UseEnergy += uint64(comData.UseEnergy)
			deviceInfo.UseTime += uint64(comData.UseTime)
			deviceInfo.CurElectricity += uint64(comData.CurElectricity)
			maxPower += comData.MaxPower
			deviceInfo.EnableCount += 1
		}
	}

	if deviceInfo.EnableCount >= 1 {
		curPower := CalculateCurComPowerToString(CUR_VOLTAGE, uint32(deviceInfo.CurElectricity), 2)
		AveragePower := CalculateCurAverageComPowerToString(uint32(deviceInfo.UseEnergy), uint32(deviceInfo.UseTime), 2)
		MaxPower := GetPowerValue(maxPower, 2)
		deviceInfo.CurPower = StringJoin([]interface{}{curPower, "w"})
		deviceInfo.AveragePower = StringJoin([]interface{}{AveragePower, "w"})
		deviceInfo.MaxPower = StringJoin([]interface{}{MaxPower, "w"})
	}

	//更新统计数据
	c.updateDeviceTatolInfoToRedis(deviceSN, deviceInfo)
	return deviceInfo.EnableCount
}

//批量读端口数据
func BatchReadDeviceComDataiFromRedis(deviceSN string) map[uint8]mqtt.ComData {
	var maxCom int = 10

	conn := Redis().BatchStart()
	defer Redis().BatchEnd(conn)
	comList := make(map[uint8]mqtt.ComData)
	for i := 0; i < maxCom; i++ {
		Redis().BatchHGet(conn, GetComdDataKey(deviceSN), strconv.Itoa(i))
	}

	pipe_list, _ := RedisValues(Redis().BatchExec(conn))
	for _, v := range pipe_list {
		comDataString, _ := RedisString(v, nil)
		comData := mqtt.ComData{}
		err := json.Unmarshal([]byte(comDataString), &comData)
		if err == nil {
			comList[comData.Id] = comData
		}
	}

	return comList
}

//批量写端口数据
func BatchWriteDeviceComDataToRedis(deviceSN string, comList *mqtt.ComList, comOps func(comData *mqtt.ComData)) {
	if comList == nil {
		return
	}

	conn := Redis().BatchStart()
	defer Redis().BatchEnd(conn)

	for _, comID := range comList.ComID {
		var index uint8 = comID
		if comList.ComNum <= 5 {
			index = (comID % 5)
		}
		comData := (comList.ComPort[int(index)]).(mqtt.ComData)
		comData.Id = comID

		comOps(&comData)
		Redis().BatchHSet(conn, GetComdDataKey(deviceSN), strconv.Itoa(int(comID)), comData)
	}

	Redis().BatchExec(conn)
}
