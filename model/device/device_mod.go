package device

import (
	. "GoServer/utils/time"
)

type AccesswayType int8

const (
	UNKNOWN AccesswayType = iota
	GSM
	WIFI
	BLUETOOTH
)

type DeviceConnectLog struct {
	ID         int64         `json:"id" xorm:"pk autoincr BIGINT(20) 'id'"`
	DeviceID   int64         `json:"device_id" xorm:"default NULL comment('设备id') BIGINT(20) 'device_id'"`
	AccessWay  AccesswayType `json:"access_way" xorm:"default NULL comment('接入方式 1 GSM， 2，WIFI 3蓝牙') TINYINT(1) 'access_way'"`
	ModuleSN   string        `json:"module_sn" xorm:"default 'NULL' comment('模组序列号') VARCHAR(64) 'module_sn'"`
	IP         string        `json:"ip" xorm:"default 'NULL' comment('连接时IP') VARCHAR(16) 'ip'"`
	CreateTime int64         `json:"create_time" xorm:"default NULL comment('连接时间') BIGINT(13) 'create_time'"`
}

type DeviceInfo struct {
	ID         int64         `json:"id" gorm:"pk autoincr comment('设备ID') BIGINT(20) 'id'"`
	UID        int64         `json:"uid" gorm:"default NULL comment('对应关系用户') BIGINT(20) 'uid'"`
	AccessWay  AccesswayType `json:"access_way" gorm:"default NULL comment('接入方式') TINYINT(1) 'access_way'"`
	Type       int8          `json:"type" gorm:"default NULL comment('设备类型') TINYINT(1) 'type'"`
	ModuleSN   string        `json:"module_sn" gorm:"default 'NULL' comment('模组序列号') VARCHAR(64) 'module_sn'"`
	DeviceSN   string        `json:"device_sn" gorm:"default 'NULL' comment('设备序列号') VARCHAR(64) 'device_sn'"`
	Remark     string        `json:"remark" gorm:"default 'NULL' comment('设备备注') VARCHAR(255) 'remark'"`
	Localtion  string        `json:"localtion" gorm:"default 'NULL' comment('所在位置') VARCHAR(255) 'localtion'"`
	Version    string        `json:"version" gorm:"default 'NULL' comment('固件版本') VARCHAR(64) 'version'"`
	CreateTime int64         `json:"create_time" gorm:"default NULL comment('创建时间') BIGINT(13) 'create_time'"`
	UpdateTime int64         `json:"update_time" gorm:"default NULL comment('最后一次更新时间') BIGINT(13) 'update_time'"`
}

type DeviceTransferLog struct {
	ID           int64  `json:"id" gorm:"pk autoincr BIGINT(20) 'id'"`
	TransferID   int64  `json:"transfer_id" gorm:"default NULL comment('传输ID') BIGINT(20) 'transfer_id'"`
	TransferAct  string `json:"transfer_act" gorm:"default NULL comment('设备行为') VARCHAR(32) 'transfer_act'"`
	DeviceSN     string `json:"device_sn" gorm:"default NULL comment('设备ID') VARCHAR(64) 'device_sn'"`
	ComNum       int64  `json:"com_num" gorm:"default NULL comment('上报条数') TINYINT(2) 'com_num'"`
	TransferData string `json:"transfer_data" gorm:"default 'NULL' comment('传输数据base64') VARCHAR(512) 'transfer_data'"`
	Behavior     int64  `json:"behavior" gorm:"default NULL comment('传输行为') TINYINT(1) 'behavior'"`
	ServerNode   string `json:"server_node" gorm:"default 'NULL' comment('上报服务节点') VARCHAR(32) 'server_node'"`
	TransferTime int64  `json:"transfer_time" gorm:"default NULL comment('建立时间') BIGINT(13) 'transfer_time'"`
	CreateTime   int64  `json:"create_time" gorm:"default NULL comment('建立时间') BIGINT(13) 'create_time'"`
}

func (log *DeviceConnectLog) Create(deviceID int64, accessway AccesswayType, moduleSN string, ip string) {
	log.DeviceID = deviceID
	log.AccessWay = accessway
	log.ModuleSN = moduleSN
	log.IP = ip
	log.CreateTime = GetTimestampMs()
}

func (device *DeviceInfo) Create(accessway AccesswayType, deviceType int8, moduleSN string, deviceSN string, version string) {
	device.AccessWay = accessway
	device.Type = deviceType
	device.ModuleSN = moduleSN
	device.DeviceSN = deviceSN
	device.Version = version
	device.CreateTime = GetTimestampMs()
}

func (transfer *DeviceTransferLog) Create(transfer_id int64, act string, device_sn string, data string, serverIP string, behavior int64, transferTime int64) {
	transfer.TransferID = transfer_id
	transfer.TransferAct = act
	transfer.DeviceSN = device_sn
	transfer.TransferData = data
	transfer.Behavior = behavior
	transfer.ServerNode = serverIP
	transfer.TransferTime = transferTime
	transfer.CreateTime = GetTimestampMs()
}
