package redis

import (
	mqtt "GoServer/mqttPacket"
	. "GoServer/utils/float64"
	. "GoServer/utils/log"
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
	COM_ENABLE    uint8 = 1
	COM_DISENABLE uint8 = 0
)

//info
const (
	REDIS_INFO_DEVICE_ID_FIELD         = "id"
	REDIS_INFO_DEVICE_DATA_TOTAL_FIELD = "total"
	REDIS_INFO_DEVICE_STATUS_FIELD     = "status"
	REDIS_INFO_DEVICE_WORKER_FIELD     = "worker"
	REDIS_INFO_SYNC_UPDATE_FIELD       = "syncTime" //同步数据到Mysql
	REDIS_INFO_RAW_DATA_FIELD          = "rawData"
	REDIS_INFO_USER_ID_FIELD           = "userID"
)

//整机数据统计
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
	ID           uint64 `json:"id"`
	ProtoVersion uint16 `json:"protoVersion"`
}

//端口数据缓存统计
type CacheComData struct {
	mqtt.ComData
	MaxChargeElectricity uint32  `json:"maxChargeElectricity"` //最大电流
	CurPower             float64 `json:"c_power"`              //当前端口使用功率
	AveragePower         float64 `json:"a_power"`              //当前端口平均功率
	MaxPower             float64 `json:"m_power"`              //当前端口最高使用功率
	SyncTime             int64   `json:"sync_time"`            //当前端口同步数据到数据库
	WriteFlags           uint8   `write_flags`                 //写数据标志位
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
	status := DeviceStatus{
		Behavior:     0,
		Signal:       0,
		Worker:       0,
		ID:           deviceID,
		ProtoVersion: 0,
	}
	c.Set(GetDeviceTokenKey(deviceSN), status, 900)                          //15分钟过期,正常1-2分钟后续数据就上来了
	c.HSet(GetDeviceInfoKey(deviceSN), REDIS_INFO_DEVICE_ID_FIELD, deviceID) //插入设备ID
	c.HSet(GetDeviceInfoKey(deviceSN), REDIS_INFO_DEVICE_STATUS_FIELD, 1)    //设备在线状态
}

func (c *Cacher) GetDeviceWorkerFormRedis(deviceSN string) int {
	worker, err := c.HGetInt(GetDeviceInfoKey(deviceSN), REDIS_INFO_DEVICE_WORKER_FIELD)
	if err != nil {
		return 0
	}
	return worker
}

func (c *Cacher) SetDeviceWorkerToRedis(deviceSN string, worker int) {
	c.HSet(GetDeviceInfoKey(deviceSN), REDIS_INFO_DEVICE_WORKER_FIELD, worker)
}

func (c *Cacher) UpdateDeviceIDToRedisByDeviceSN(deviceSN string, deviceID uint64) {
	if deviceSN != "" && deviceID != 0 {
		c.InitWithInsertDeviceIDToken(deviceSN, deviceID)
	}
}

//取设备状态
func (c *Cacher) GetDeviceStatusFromRedis(deviceSN string) int {
	status, err := c.HGetInt(GetDeviceInfoKey(deviceSN), REDIS_INFO_DEVICE_STATUS_FIELD)
	if err != nil {
		return 0
	}
	return status
}

//设置设备状态
func (c *Cacher) SetDeviceStatusToRedis(deviceSN string, status int) {
	c.HSet(GetDeviceInfoKey(deviceSN), REDIS_INFO_DEVICE_STATUS_FIELD, status)
}

//设置设备同步时间
func (c *Cacher) SetDeviceSyncTimeToRedis(deviceSN string, syncTime int64) {
	c.HSet(GetDeviceInfoKey(deviceSN), REDIS_INFO_SYNC_UPDATE_FIELD, syncTime)
}

//取设备同步时间
func (c *Cacher) GetDeviceSyncTimeFromRedis(deviceSN string) int64 {
	time, err := c.HGetInt64(GetDeviceInfoKey(deviceSN), REDIS_INFO_SYNC_UPDATE_FIELD)
	if err != nil {
		return 0
	}
	return time
}

//取设备ID
func (c *Cacher) GetDeviceIDFromRedis(deviceSN string) uint64 {
	deviceID, err := c.HGetUint64(GetDeviceInfoKey(deviceSN), REDIS_INFO_DEVICE_ID_FIELD)
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
	c.HSet(GetDeviceInfoKey(deviceSN), REDIS_INFO_DEVICE_DATA_TOTAL_FIELD, infoData)
}

//统计整机数据
func (c *Cacher) TatolWorkerByDevice(deviceSN string, comDataMap map[uint8]CacheComData) uint8 {

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
		if comData.Enable == COM_ENABLE {
			deviceInfo.UseEnergy += uint64(comData.UseEnergy)
			deviceInfo.UseTime += uint64(comData.UseTime)
			deviceInfo.CurElectricity += uint64(comData.CurElectricity)
			maxPower += comData.MaxPower
			deviceInfo.EnableCount += 1
		}
	}

	if deviceInfo.EnableCount >= 1 {
		curPower := CalculateCurComPowerToString(CUR_VOLTAGE, uint32(deviceInfo.CurElectricity), 5)
		AveragePower := CalculateCurAverageComPowerToString(uint32(deviceInfo.UseEnergy), uint32(deviceInfo.UseTime), 5)
		MaxPower := GetPowerValue(maxPower, 5)
		deviceInfo.CurPower = StringJoin([]interface{}{curPower, "w"})
		deviceInfo.AveragePower = StringJoin([]interface{}{AveragePower, "w"})
		deviceInfo.MaxPower = StringJoin([]interface{}{MaxPower, "w"})
	}

	//更新统计数据
	c.updateDeviceTatolInfoToRedis(deviceSN, deviceInfo)
	return deviceInfo.EnableCount
}

//批量读端口数据
func BatchReadDeviceComDataiFromRedis(deviceSN string) map[uint8]CacheComData {
	var maxCom int = 10

	conn := Redis().BatchStart()
	defer Redis().BatchEnd(conn)
	comList := make(map[uint8]CacheComData)
	for i := 0; i < maxCom; i++ {
		Redis().BatchHGet(conn, GetComdDataKey(deviceSN), strconv.Itoa(i))
	}

	pipe_list, _ := RedisValues(Redis().BatchExec(conn))
	for _, v := range pipe_list {
		comDataString, _ := RedisString(v, nil)
		comData := CacheComData{}
		err := json.Unmarshal([]byte(comDataString), &comData)
		if err == nil {
			comList[comData.Id] = comData
		}
	}

	return comList
}

//批量写端口数据
func BatchWriteDeviceComDataToRedis(deviceSN string, comList *mqtt.ComList, comOps func(comData *mqtt.ComData) *CacheComData) {
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

		cacheComData := comOps(&comData)
		if err := Redis().BatchHSet(conn, GetComdDataKey(deviceSN), strconv.Itoa(int(comID)), cacheComData); err != nil {
			SystemLog("BatchWriteDeviceComDataToRedis err: ", err)
		}
	}

	Redis().BatchExec(conn)
}

//批量读Token
func BatchReadDeviceTokenFromRedis(deviceMap map[uint64]string) ([]interface{}, []interface{}) {
	conn := Redis().BatchStart()
	defer Redis().BatchEnd(conn)

	onLine := make([]interface{}, 0)
	workInLine := make([]interface{}, 0)

	for _, deviceSN := range deviceMap {
		Redis().BatchGet(conn, GetDeviceTokenKey(deviceSN))
	}

	token_list, _ := RedisValues(Redis().BatchExec(conn))

	for _, v := range token_list {
		statusValue, _ := RedisString(v, nil)
		status := DeviceStatus{}
		err := json.Unmarshal([]byte(statusValue), &status)
		if err == nil {
			if _, ok := deviceMap[status.ID]; ok {
				switch status.Behavior {
				case mqtt.GISUNLINK_CHARGE_LEISURE:
					onLine = append(onLine, status)
					break
				case mqtt.GISUNLINK_CHARGEING:
					workInLine = append(workInLine, status)
					break
				}
				delete(deviceMap, status.ID)
			}
		}
	}
	return onLine, workInLine
}
