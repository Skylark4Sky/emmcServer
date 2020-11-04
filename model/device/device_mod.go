package device

import (
	. "GoServer/utils/time"
)

const (
	UNKNOWN int8 = iota
	GSM
	WIFI
	BLUETOOTH
)

type ModuleConnectLog struct {
	ID         uint64 `gorm:"primary_key;column:id;type:bigint(20) unsigned;not null" json:"-"`
	ModuleID   uint64 `gorm:"column:module_id;type:bigint(20) unsigned" json:"module_id"`     // 模组id
	AccessWay  uint8  `gorm:"column:access_way;type:tinyint(2) unsigned" json:"access_way"`   // 接入方式 1 GSM， 2，WIFI 3蓝牙
	ModuleSn   string `gorm:"column:module_sn;type:varchar(64)" json:"module_sn"`             // 模组序列号
	IP         string `gorm:"column:ip;type:varchar(16)" json:"ip"`                           // 连接时IP
	CreateTime int64  `gorm:"column:create_time;type:bigint(13) unsigned" json:"create_time"` // 连接时间
}

type ModuleInfo struct {
	ID            uint64 `gorm:"primary_key;column:id;type:bigint(20) unsigned;not null" json:"-"`    // 设备ID
	DeviceID      uint64 `gorm:"column:device_id;type:bigint(20) unsigned;not null" json:"device_id"` // 对应设备关系
	AccessWay     uint8  `gorm:"column:access_way;type:tinyint(2) unsigned" json:"access_way"`        // 接入方式
	ModuleSn      string `gorm:"column:module_sn;type:varchar(64)" json:"module_sn"`                  // 模组序列号
	ModuleVersion string `gorm:"column:module_version;type:varchar(32)" json:"module_version"`        // 模组固件版本
	CreateTime    int64  `gorm:"column:create_time;type:bigint(13) unsigned" json:"create_time"`      // 创建时间
	UpdateTime    int64  `gorm:"column:update_time;type:bigint(13) unsigned" json:"update_time"`      // 最后一次更新时间
}

type DeviceInfo struct {
	ID            uint64 `gorm:"primary_key;column:id;type:bigint(20) unsigned;not null" json:"-"` // 设备ID
	AccessWay     uint8  `gorm:"column:access_way;type:tinyint(2) unsigned" json:"access_way"`     // 当前接入方式
	DeviceSn      string `gorm:"column:device_sn;type:varchar(64)" json:"device_sn"`               // 设备序列号
	DeviceVersion string `gorm:"column:device_version;type:varchar(32)" json:"device_version"`     // 设备固件版本
	Remark        string `gorm:"column:remark;type:varchar(255)" json:"remark"`                    // 设备备注
	Localtion     string `gorm:"column:localtion;type:varchar(255)" json:"localtion"`              // 设备所在地址
	Type          int8   `gorm:"column:type;type:tinyint(2) unsigned" json:"type"`                 // 设备类型
	CreateTime    int64  `gorm:"column:create_time;type:bigint(13) unsigned" json:"create_time"`   // 创建时间
	UpdateTime    int64  `gorm:"column:update_time;type:bigint(13) unsigned" json:"update_time"`   // 更新时间
}

type DeviceTransferLog struct {
	ID           uint64 `gorm:"primary_key;column:id;type:bigint(20) unsigned;not null" json:"-"`
	DeviceID     uint64 `gorm:"column:device_id;type:bigint(20) unsigned" json:"device_id"`         // 设备ID
	Behavior     uint8  `gorm:"column:behavior;type:tinyint(2)" json:"behavior"`                    // 传输行为
	DeviceSn     string `gorm:"column:device_sn;type:varchar(64)" json:"device_sn"`                 // 设备串号
	ServerNode   string `gorm:"column:server_node;type:varchar(32)" json:"server_node"`             // 服务节点
	TransferID   int64  `gorm:"column:transfer_id;type:bigint(13) unsigned" json:"transfer_id"`     // 传输ID
	TransferAct  string `gorm:"column:transfer_act;type:varchar(32)" json:"transfer_act"`           // 传输行为
	ComNum       uint8  `gorm:"column:com_num;type:tinyint(2) unsigned" json:"com_num"`             // 上报数据条数
	TransferData string `gorm:"column:transfer_data;type:varchar(512)" json:"transfer_data"`        // 传输数据base64
	TransferTime int64  `gorm:"column:transfer_time;type:bigint(13) unsigned" json:"transfer_time"` // 传输时间
	CreateTime   int64  `gorm:"column:create_time;type:bigint(13)" json:"create_time"`              // 建立时间
}

type CreateDeviceInfo struct {
	Module ModuleInfo
	Device DeviceInfo
	Log    ModuleConnectLog
}

func (log *ModuleConnectLog) Create(deviceID uint64, accessway uint8, moduleSN string, ip string) {
	log.ModuleID = deviceID
	log.AccessWay = accessway
	log.ModuleSn = moduleSN
	log.IP = ip
	log.CreateTime = GetTimestampMs()
}

func (module *ModuleInfo) Create(accessway uint8, moduleSN string, module_version string) {
	module.AccessWay = accessway
	module.ModuleSn = moduleSN
	module.ModuleVersion = module_version
	module.CreateTime = GetTimestampMs()
}

func (device *DeviceInfo) Create(accessway uint8, sn string, version string) {
	device.AccessWay = accessway
	device.DeviceSn = sn
	device.DeviceVersion = version
	device.CreateTime = GetTimestampMs()
}

func (module *ModuleInfo) Update(module_version string) {
	module.ModuleVersion = module_version
	module.UpdateTime = GetTimestampMs()
}

func (device *DeviceInfo) Update(id uint64, accessway uint8, version string, time int64) {
	device.ID = id
	device.AccessWay = accessway
	device.DeviceVersion = version
	device.UpdateTime = time
}

func (transfer *DeviceTransferLog) Create(transfer_id int64, act string, device_sn string, data string, serverIP string, behavior uint8, transferTime int64) {
	transfer.TransferID = transfer_id
	transfer.TransferAct = act
	transfer.DeviceSn = device_sn
	transfer.TransferData = data
	transfer.Behavior = behavior
	transfer.ServerNode = serverIP
	transfer.TransferTime = transferTime
	transfer.CreateTime = GetTimestampMs()
}
