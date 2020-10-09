package redis

import (
	. "GoServer/utils/string"
	"strconv"
)

const (
	STATUSKEY   = "status:"
	DEVICEIDKEY = "ID:"
	RAWDATAKEY  = "rawData:"
	COMDATAKEY  = "comData:"
)

type DeviceStatus struct {
	Behavior     uint8  `json:behavior`
	Signal       int8   `json:signal`
	Worker       uint8  `json:worker`
	ProtoVersion uint16 `json:protoVersion`
}

func getStatusKey(deviceSN string) string {
	return StringJoin([]interface{}{STATUSKEY, deviceSN})
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

//更新Device令牌时间
func (c *Cacher) UpdateDeviceTokenExpiredTime(deviceSN string, status *DeviceStatus, timeout int64) {
	c.Set(getStatusKey(deviceSN), status, timeout) //15分钟过期,正常1-2分钟后续数据就上来了
}

//插入Device对应令牌
func (c *Cacher) InitWithInsertDeviceIDToken(deviceSN string, deviceID int64) {
	c.Set(getStatusKey(deviceSN), deviceID, 900) //15分钟过期,正常1-2分钟后续数据就上来了
	c.Set(getDeviceIDKey(deviceSN), deviceID, 0) //全局设备列表
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

//统计当前设备工作端口数量
func (c *Cacher) TatolWorkerByDevice(deviceSN string, isEnable func(comDataString string) bool) uint8 {
	var maxCom int = 10
	var worker uint8 = 0

	for i := 0; i < maxCom; i++ {
		str, err := RedisString(c.HGet(getComdDataKey(deviceSN), strconv.Itoa(i)))
		if err != nil {
			continue
		}
		if isEnable(str) {
			worker += 1
		}
	}

	return worker
}
