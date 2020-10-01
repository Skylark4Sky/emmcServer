package device

import (
//. "GoServer/utils"
)

type AccesswayType int8

const (
	UNKNOWN AccesswayType = iota
	GSM
	WIFI
	BLUETOOTH
)

type DeviceInfo struct {
	ID         int64  `json:"id" gorm:"pk autoincr comment('设备ID') BIGINT(20) 'id'"`
	UID        int64  `json:"uid" gorm:"default NULL comment('对应关系用户') BIGINT(20) 'uid'"`
	Accessway  int8   `json:"accessway" gorm:"default NULL comment('接入方式') TINYINT(1) 'accessWay'"`
	Type       int8   `json:"type" gorm:"default NULL comment('设备类型') TINYINT(1) 'type'"`
	DeviceSN   string `json:"devicesn" gorm:"default 'NULL' comment('设备序列号') VARCHAR(64) 'deviceSN'"`
	Remark     string `json:"remark" gorm:"default 'NULL' comment('设备备注') VARCHAR(255) 'remark'"`
	Localtion  string `json:"localtion" gorm:"default 'NULL' comment('所在位置') VARCHAR(255) 'localtion'"`
	Version    string `json:"version" gorm:"default 'NULL' comment('固件版本') VARCHAR(64) 'version'"`
	Createtime int64  `json:"createtime" gorm:"default NULL comment('创建时间') BIGINT(20) 'createTime'"`
	Updatetime int64  `json:"updatetime" gorm:"default NULL comment('最后一次更新时间') BIGINT(20) 'updateTime'"`
}

type DeviceTransferLog struct {
	ID           int64  `json:"id" gorm:"pk autoincr BIGINT(20) 'id'"`
	DeviceID     int64  `json:"device_id" gorm:"default NULL comment('设备ID') BIGINT(20) 'device_id'"`
	TransferID   int64  `json:"transfer_id" gorm:"default NULL comment('传输ID') BIGINT(20) 'transfer_id'"`
	TransferData string `json:"transfer_data" gorm:"default 'NULL' comment('传输数据base64') VARCHAR(512) 'transfer_data'"`
	Behavior     int64  `json:"behavior" gorm:"default NULL comment('传输行为') TINYINT(1) 'behavior'"`
	IP           string `json:"ip" gorm:"default 'NULL' comment('上报ID') VARCHAR(16) 'ip'"`
	CreateTime   int64  `json:"create_time" gorm:"default NULL comment('建立时间') INT(13) 'create_time'"`
}
