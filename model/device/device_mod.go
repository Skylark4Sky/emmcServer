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

type ModuleConnectLog struct {
	ID         int64         `json:"id" gorm:"pk autoincr BIGINT(20) 'id'"`
	ModuleID   int64         `json:"module_id" gorm:"default NULL comment('模组id') BIGINT(20) 'module_id'"`
	AccessWay  AccesswayType `json:"access_way" gorm:"default NULL comment('接入方式 1 GSM， 2，WIFI 3蓝牙') TINYINT(1) 'access_way'"`
	ModuleSN   string        `json:"module_sn" gorm:"default 'NULL' comment('模组序列号') VARCHAR(64) 'module_sn'"`
	IP         string        `json:"ip" gorm:"default 'NULL' comment('连接时IP') VARCHAR(16) 'ip'"`
	CreateTime int64         `json:"create_time" gorm:"default NULL comment('连接时间') BIGINT(13) 'create_time'"`
}

type ModuleInfo struct {
	ID            int64         `json:"id" gorm:"pk autoincr comment('设备ID') BIGINT(20) 'id'"`
	DeviceID      int64         `json:"device_id" gorm:"default NULL comment('对应设备关系') BIGINT(20) 'device_id'"`
	AccessWay     AccesswayType `json:"access_way" gorm:"default NULL comment('接入方式') TINYINT(1) 'access_way'"`
	ModuleSN      string        `json:"module_sn" gorm:"default 'NULL' comment('模组序列号') VARCHAR(64) 'module_sn'"`
	ModuleVersion string        `json:"module_version" gorm:"default 'NULL' comment('模组固件版本') VARCHAR(32) 'module_version'"`
	CreateTime    int64         `json:"create_time" gorm:"default NULL comment('创建时间') BIGINT(13) 'create_time'"`
	UpdateTime    int64         `json:"update_time" gorm:"default NULL comment('最后一次更新时间') BIGINT(13) 'update_time'"`
}

type DeviceInfo struct {
	ID            int64  `json:"id" gorm:"pk autoincr comment('设备ID') BIGINT(20) 'id'"`
	Type          int64  `json:"type" gorm:"default 0 comment('设备类型') TINYINT(2) 'type'"`
	DeviceSn      string `json:"device_sn" gorm:"default 'NULL' comment('设备序列号') VARCHAR(64) 'device_sn'"`
	DeviceVersion string `json:"device_version" gorm:"default 'NULL' comment('设备固件版本') VARCHAR(32) 'device_version'"`
	Remark        string `json:"remark" gorm:"default 'NULL' comment('设备备注') VARCHAR(255) 'remark'"`
	Localtion     string `json:"localtion" gorm:"default 'NULL' comment('设备所在地址') VARCHAR(255) 'localtion'"`
	CreateTime    int64  `json:"create_time" gorm:"default NULL comment('创建时间') BIGINT(13) 'create_time'"`
	UpdateTime    int64  `json:"update_time" gorm:"default NULL comment('最后一次更新时间') BIGINT(13) 'update_time'"`
}

type DeviceTransferLog struct {
	ID           int64  `json:"id" gorm:"pk autoincr BIGINT(20) 'id'"`
	DeviceID     int64  `json:"device_id" gorm:"default 0 comment('对应设备关系') BIGINT(20) 'device_id'"`
	TransferID   int64  `json:"transfer_id" gorm:"default 0 comment('传输ID') BIGINT(20) 'transfer_id'"`
	TransferAct  string `json:"transfer_act" gorm:"default 'NULL' comment('设备行为') VARCHAR(32) 'transfer_act'"`
	DeviceSN     string `json:"device_sn" gorm:"default 'NULL' comment('设备ID') VARCHAR(64) 'device_sn'"`
	ComNum       int64  `json:"com_num" gorm:"default 0 comment('上报条数') TINYINT(2) 'com_num'"`
	TransferData string `json:"transfer_data" gorm:"default 'NULL' comment('传输数据base64') VARCHAR(512) 'transfer_data'"`
	Behavior     int64  `json:"behavior" gorm:"default 0 comment('传输行为') TINYINT(1) 'behavior'"`
	ServerNode   string `json:"server_node" gorm:"default 'NULL' comment('上报服务节点') VARCHAR(32) 'server_node'"`
	TransferTime int64  `json:"transfer_time" gorm:"default 0 comment('建立时间') BIGINT(13) 'transfer_time'"`
	CreateTime   int64  `json:"create_time" gorm:"default 0 comment('建立时间') BIGINT(13) 'create_time'"`
}

type CreateDeviceInfo struct {
	Module ModuleInfo
	Device DeviceInfo
	Log    ModuleConnectLog
}

func (log *ModuleConnectLog) Create(deviceID int64, accessway AccesswayType, moduleSN string, ip string) {
	log.ModuleID = deviceID
	log.AccessWay = accessway
	log.ModuleSN = moduleSN
	log.IP = ip
	log.CreateTime = GetTimestampMs()
}

func (module *ModuleInfo) Create(accessway AccesswayType, moduleSN string, module_version string) {
	module.AccessWay = accessway
	module.ModuleSN = moduleSN
	module.ModuleVersion = module_version
	module.CreateTime = GetTimestampMs()
}

func (device *DeviceInfo) Create(sn string, version string) {
	device.DeviceSn = sn
	device.DeviceVersion = version
	device.CreateTime = GetTimestampMs()
}

func (module *ModuleInfo) Update(module_version string) {
	module.ModuleVersion = module_version
	module.UpdateTime = GetTimestampMs()
}

func (device *DeviceInfo) Update(id int64, version string, time int64) {
	device.ID = id
	device.DeviceVersion = version
	device.UpdateTime = time
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
