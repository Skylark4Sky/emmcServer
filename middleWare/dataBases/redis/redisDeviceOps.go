package redis

import (
//	"fmt"
)

const (
	DEVICEID = "deviceID:"
	RAWDATA  = "rawData:"
	COMDATA  = "comData:"
)

type DeviceStatus struct {
	Behavior     uint8  `json:behavior`
	Signal       int8   `json:signal`
	Worker       uint8  `json:worker`
	ProtoVersion uint16 `json:protoVersion`
}

func getDeviceIDKey(deviceSN string) string {
	return DEVICEID + deviceSN
}

func getRawDataKey(deviceSN string) string {
	return RAWDATA + deviceSN
}

func getComdDataKey(deviceSN string) string {
	return COMDATA + deviceSN
}

//更新Device令牌时间
func (c *Cacher) UpdateDeviceTokenExpiredTime(deviceSN string, status *DeviceStatus, time int64) {
	c.Set(deviceSN, status, time) //15分钟过期,正常1-2分钟后续数据就上来了
}

//插入Device对应令牌
func (c *Cacher) InitWithInsertDeviceIDToken(deviceSN string, deviceID int64) {
	c.Set(deviceSN, deviceID, 900)               //15分钟过期,正常1-2分钟后续数据就上来了
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
	c.HSet(getComdDataKey(deviceSN), comID, comData)
}
