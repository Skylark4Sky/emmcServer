package redis

import (
	mqtt "GoServer/mqtt"
	. "GoServer/utils/float64"
	. "GoServer/utils/string"
	"encoding/json"
	"strconv"
)

const (
	DEVICETOEKNKEY = "token:"
	DEVICEIDKEY    = "ID:"
	RAWDATAKEY     = "rawData:"
	COMDATAKEY     = "comData:"
	DEVICEINFOKEY  = "info:"
)

type DeviceTatolInfo struct {
	UseEnergy      uint64 `json:"energy"`
	UseTime        uint64 `json:"time"`
	CurElectricity uint64 `json:"electricity"`
	CurPower       string `json:"c_power"`
	AveragePower   string `json:"a_power"`
	MaxPower   	   string `json:"m_power"`
	EnableCount    uint8  `json:"enable"`
}

type DeviceStatus struct {
	Behavior     uint8  `json:"behavior"`
	Signal       int8   `json:"signal"`
	Worker       uint8  `json:"worker"`
	ProtoVersion uint16 `json:"protoVersion"`
}

func getDeviceTokenKey(deviceSN string) string {
	return StringJoin([]interface{}{DEVICETOEKNKEY, deviceSN})
}

func getDeviceIDKey(deviceSN string) string {
	return StringJoin([]interface{}{DEVICEIDKEY, deviceSN})
}

func getRawDataKey(deviceSN string) string {
	return StringJoin([]interface{}{RAWDATAKEY, deviceSN})
}

func getComdDataKey(deviceSN string) string {
	return StringJoin([]interface{}{COMDATAKEY, deviceSN})
}

func getDeviceInfoKey(deviceSN string) string {
	return StringJoin([]interface{}{DEVICEINFOKEY, deviceSN})
}

//更新令牌时间
func (c *Cacher) UpdateDeviceTokenExpiredTime(deviceSN string, status *DeviceStatus, timeout int64) {
	c.Set(getDeviceTokenKey(deviceSN), status, timeout) //15分钟过期,正常1-2分钟后续数据就上来了
}

//插入对应令牌
func (c *Cacher) InitWithInsertDeviceIDToken(deviceSN string, deviceID int64) {
	c.Set(getDeviceTokenKey(deviceSN), deviceID, 900) //15分钟过期,正常1-2分钟后续数据就上来了
	c.Set(getDeviceIDKey(deviceSN), deviceID, 0)      //全局设备列表
}

//取设备ID
func (c *Cacher) GetDeviceIDFromRedis(deviceSN string, key_field string) int64 {
	deviceID, err := c.GetInt64(getDeviceIDKey(deviceSN))
	if err != nil {
		return 0
	}
	return deviceID
}

//更新原始上报数据
func (c *Cacher) UpdateDeviceRawDataToRedis(deviceSN string, rawData string) {
	c.Set(getRawDataKey(deviceSN), rawData, 0)
}

// 获取某一端口数据
func (c *Cacher) GetDeviceComDataFormRedis(deviceSN string, comID uint8, comData interface{}) error {
	str, err := RedisString(c.HGet(getComdDataKey(deviceSN), strconv.Itoa(int(comID))))
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(str), comData)
}

//更新端口数据
func (c *Cacher) UpdateDeviceComDataToRedis(deviceSN string, comID uint8, comData interface{}) {
	c.HSet(getComdDataKey(deviceSN), strconv.Itoa(int(comID)), comData)
}

//更新工作状态数据统计
func (c *Cacher) updateDeviceTatolInfoToRedis(deviceSN string, infoData interface{}) {
	c.Set(getDeviceInfoKey(deviceSN), infoData, 0)
}

//统计当前工作端口数量
func (c *Cacher) TatolWorkerByDevice(deviceSN string, comDataMap map[uint8]mqtt.ComData) uint8 {

	deviceInfo := &DeviceTatolInfo{
		UseEnergy:      0,
		UseTime:        0,
		CurElectricity: 0,
		CurPower:       "0w",
		AveragePower:   "0w",
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
		MaxPower := GetPowerValue(maxPower,2)
		deviceInfo.CurPower = StringJoin([]interface{}{curPower, "w"})
		deviceInfo.AveragePower = StringJoin([]interface{}{AveragePower, "w"})
		deviceInfo.MaxPower = StringJoin([]interface{}{MaxPower, "w"})
	}

	//更新统计数据
	c.updateDeviceTatolInfoToRedis(deviceSN, deviceInfo)
	return deviceInfo.EnableCount
}

//批量读端口数据
func BatchReadDeviceComData(deviceSN string) map[uint8]mqtt.ComData {
	var maxCom int = 10

	conn := Redis().BatchStart()
	defer Redis().BatchEnd(conn)
	comList := make(map[uint8]mqtt.ComData)
	for i := 0; i < maxCom; i++  {
		Redis().BatchHGet(conn,getComdDataKey(deviceSN), strconv.Itoa(i))
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
func BatchWriteDeviceComData(deviceSN string,comList *mqtt.ComList, comOps func(comData *mqtt.ComData)) {
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
		Redis().BatchHSet(conn,getComdDataKey(deviceSN), strconv.Itoa(int(comID)),comData)
	}

	Redis().BatchExec(conn)
}