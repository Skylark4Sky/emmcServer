package redis

import (
	. "GoServer/utils/float64"
	. "GoServer/utils/string"
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
	UseEnergy      uint64 `json:energy`
	UseTime        uint64 `json:time`
	CurElectricity uint64 `json:electricity`
	CurPower       string `json:power`
	EnableCount    uint8  `json:enable`
}

type DeviceStatus struct {
	Behavior     uint8  `json:behavior`
	Signal       int8   `json:signal`
	Worker       uint8  `json:worker`
	ProtoVersion uint16 `json:protoVersion`
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

//更新端口数据
func (c *Cacher) UpdateDeviceComDataToRedis(deviceSN string, comID uint8, comData interface{}) {
	c.HSet(getComdDataKey(deviceSN), strconv.Itoa(int(comID)), comData)
}

//更新工作状态数据统计
func (c *Cacher) updateDeviceTatolInfoToRedis(deviceSN string, infoData interface{}) {
	c.Set(getDeviceInfoKey(deviceSN), infoData, 0)
}

//统计当前工作端口数量
func (c *Cacher) TatolWorkerByDevice(deviceSN string, analysisComStatus func(comDataString string) (enable bool, useEnergy uint32, useTime uint32,curElectricity uint32)) uint8 {
	var maxCom int = 10

	deviceInfo := &DeviceTatolInfo{
		UseEnergy:      0,
		UseTime:        0,
		CurElectricity: 0,
		CurPower:       "0 w",
		EnableCount:    0,
	}

	for i := 0; i < maxCom; i++ {
		str, err := RedisString(c.HGet(getComdDataKey(deviceSN), strconv.Itoa(i)))
		if err != nil {
			continue
		}

		enable, useEnergy, useTime, curElectricity := analysisComStatus(str)
		if enable {
			deviceInfo.UseEnergy += uint64(useEnergy)
			deviceInfo.UseTime += uint64(useTime)
			deviceInfo.CurElectricity += uint64(curElectricity)
			deviceInfo.EnableCount += 1
		}
	}

	if deviceInfo.EnableCount >= 1 {
		curPower := strconv.FormatFloat(CalculateComPower(CUR_VOLTAGE,uint32(deviceInfo.CurElectricity),2), 'f', 2, 64)
		deviceInfo.CurPower = StringJoin([]interface{}{curPower, " w"})
	}

	//更新统计数据
	c.updateDeviceTatolInfoToRedis(deviceSN, deviceInfo)
	return deviceInfo.EnableCount
}
