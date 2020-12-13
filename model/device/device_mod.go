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

const (
	COM_CHARGE_START_BIT     = 0x01 //下发
	COM_CHARGE_START_ACK_BIT = 0x02 //设备已执行
	COM_CHARGE_STOP_BIT      = 0x04 //下发
	COM_CHARGE_STOP_ACK_BIT  = 0x08 //设备已执行
	COM_CHARGE_ERROR_BIT     = 0x10 //0 正常 1 充电出错
)

type ModuleConnectLog struct {
	ID         uint64 `gorm:"primary_key;column:id;type:bigint(20) unsigned;not null" json:"id"`
	ModuleID   uint64 `gorm:"column:module_id;type:bigint(20) unsigned" json:"module_id"`     // 模组id
	AccessWay  uint8  `gorm:"column:access_way;type:tinyint(2) unsigned" json:"access_way"`   // 接入方式 1 GSM， 2，WIFI 3蓝牙
	ModuleSn   string `gorm:"column:module_sn;type:varchar(64)" json:"module_sn"`             // 模组序列号
	IP         string `gorm:"column:ip;type:varchar(16)" json:"ip"`                           // 连接时IP
	CreateTime int64  `gorm:"column:create_time;type:bigint(13) unsigned" json:"create_time"` // 连接时间
}

type ModuleInfo struct {
	ID            uint64 `gorm:"primary_key;column:id;type:bigint(20) unsigned;not null" json:"id"`   // 设备ID
	DeviceID      uint64 `gorm:"column:device_id;type:bigint(20) unsigned;not null" json:"device_id"` // 对应设备关系
	AccessWay     uint8  `gorm:"column:access_way;type:tinyint(2) unsigned" json:"access_way"`        // 接入方式
	ModuleSn      string `gorm:"column:module_sn;type:varchar(64)" json:"module_sn"`                  // 模组序列号
	ModuleVersion string `gorm:"column:module_version;type:varchar(32)" json:"module_version"`        // 模组固件版本
	CreateTime    int64  `gorm:"column:create_time;type:bigint(13) unsigned" json:"create_time"`      // 创建时间
	UpdateTime    int64  `gorm:"column:update_time;type:bigint(13) unsigned" json:"update_time"`      // 最后一次更新时间
}

type DeviceInfo struct {
	ID            uint64 `gorm:"primary_key;column:id;type:bigint(20) unsigned;not null" json:"id"` // 设备ID
	AccessWay     uint8  `gorm:"column:access_way;type:tinyint(2) unsigned" json:"access_way"`      // 当前接入方式
	DeviceSn      string `gorm:"column:device_sn;type:varchar(64)" json:"device_sn"`                // 设备序列号
	DeviceVersion string `gorm:"column:device_version;type:varchar(32)" json:"device_version"`      // 设备固件版本
	Remark        string `gorm:"column:remark;type:varchar(255)" json:"remark"`                     // 设备备注
	Localtion     string `gorm:"column:localtion;type:varchar(255)" json:"localtion"`               // 设备所在地址
	Type          int8   `gorm:"column:type;type:tinyint(2) unsigned" json:"type"`                  // 设备类型
	Status        int8   `gorm:"column:status;type:tinyint(2) unsigned" json:"status"`              // 设备状态
	Worker        int8   `gorm:"column:worker;type:tinyint(2) unsigned" json:"worker"`              // 设备状态
	CreateTime    int64  `gorm:"column:create_time;type:bigint(13) unsigned" json:"create_time"`    // 创建时间
	UpdateTime    int64  `gorm:"column:update_time;type:bigint(13) unsigned" json:"update_time"`    // 更新时间
}

type DeviceCom struct {
	ID                   uint64  `gorm:"primary_key;column:id;type:bigint(20) unsigned ;not null" json:"-"`
	DeviceID             uint64  `gorm:"column:device_id;type:bigint(20) unsigned ;not null" json:"device_id"`               // 设备ID
	ChargeID             uint64  `gorm:"column:charge_id;type:bigint(20) unsigned ;not null" json:"charge_id"`               // 充电ID
	ComID                uint8   `gorm:"column:com_id;type:tinyint(2) unsigned ;not null" json:"com_id"`                     // 端口
	MaxEnergy            uint32  `gorm:"column:max_energy;type:int(10) unsigned " json:"max_energy"`                         // 最大使用电量
	MaxTime              uint32  `gorm:"column:max_time;type:int(10)" json:"max_time"`                                       // 最大使用时间
	MaxElectricity       uint32  `gorm:"column:max_electricity;type:int(10) unsigned " json:"max_electricity"`               // 最大使用电流
	UseEnergy            uint32  `gorm:"column:use_energy;type:int(10) unsigned " json:"use_energy"`                         // 已冲电量
	UseTime              uint32  `gorm:"column:use_time;type:int(10) unsigned " json:"use_time"`                             // 已冲时间
	MaxChargeElectricity uint32  `gorm:"column:max_charge_electricity;type:int(10) unsigned " json:"max_charge_electricity"` // 最大充电电流
	AveragePower         float64 `gorm:"column:average_power;type:decimal(10,0) unsigned " json:"average_power"`             // 平均功率
	MaxPower             float64 `gorm:"column:max_power;type:decimal(10,0) unsigned " json:"max_power"`                     // 最大功率
	State                uint32  `gorm:"column:state;type:int(10) unsigned ;not null" json:"state"`                          // 状态
	CreateTime           int64   `gorm:"column:create_time;type:bigint(13) unsigned " json:"create_time"`                    // 创建时间
	EndTime              int64   `gorm:"column:end_time;type:bigint(13) unsigned " json:"end_time"`                          // 结束时间
}

type DeviceTransferLog struct {
	ID           uint64 `gorm:"primary_key;column:id;type:bigint(20) unsigned;not null" json:"module_id`
	DeviceID     uint64 `gorm:"column:device_id;type:bigint(20) unsigned" json:"device_id"`         // 设备ID
	Behavior     uint8  `gorm:"column:behavior;type:tinyint(2)" json:"behavior"`                    // 传输行为
	DeviceSn     string `gorm:"column:device_sn;type:varchar(64)" json:"device_sn"`                 // 设备串号
	ServerNode   string `gorm:"column:server_node;type:varchar(32)" json:"server_node"`             // 服务节点
	TransferID   int64  `gorm:"column:transfer_id;type:bigint(13) unsigned" json:"transfer_id"`     // 传输ID
	TransferAct  string `gorm:"column:transfer_act;type:varchar(32)" json:"transfer_act"`           // 传输行为
	ComNum       uint8  `gorm:"column:com_num;type:tinyint(2) unsigned" json:"com_num"`             // 上报数据条数
	TransferData string `gorm:"column:transfer_data;type:varchar(512)" json:"transfer_data"`        // 传输数据base64
	PayloadData  string `gorm:"column:payload_data;type:varchar(2048)" json:"payload_data"`         // 对应数据
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

func (device *DeviceInfo) Create(accessway uint8, sn string, version string, status int8) {
	device.AccessWay = accessway
	device.DeviceSn = sn
	device.DeviceVersion = version
	device.Status = status
	device.CreateTime = GetTimestampMs()
}

func (module *ModuleInfo) Update(module_version string) {
	module.ModuleVersion = module_version
	module.UpdateTime = GetTimestampMs()
}

func (device *DeviceInfo) Update(id uint64, accessway uint8, version string, time int64, status int8) {
	device.ID = id
	device.AccessWay = accessway
	device.DeviceVersion = version
	device.UpdateTime = time
	device.Status = status
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

func (com *DeviceCom) Create(deviceID, chargeID uint64, comID uint8, maxEnergy, maxTime, maxElectricity uint32) {
	com.DeviceID = deviceID
	com.ChargeID = chargeID
	com.ComID = comID
	com.MaxEnergy = maxEnergy
	com.MaxTime = maxTime
	com.MaxElectricity = maxElectricity
	com.State = COM_CHARGE_START_BIT
	com.CreateTime = GetTimestampMs()
}

func (com *DeviceCom) ChangeValue(useEnergy, useTime, maxChargeElectricity uint32, averagePower, maxPower float64) {
	com.UseEnergy = useEnergy
	com.UseTime = useTime
	com.MaxChargeElectricity = maxChargeElectricity
	com.AveragePower = averagePower
	com.MaxPower = maxPower
}
