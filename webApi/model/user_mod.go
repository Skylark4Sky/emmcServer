package model

import (
	. "GoServer/webApi/middleWare"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// UserAuth 用户授权表
type UserAuth struct {
	ID           int64  `gorm:"primary_key;column:id;type:bigint(20);not null" json:"-"`
	UId          int64  `gorm:"unique_index:only;index;column:uid;type:bigint(20) unsigned;not null" json:"uid"`               // 用户id
	IDentityType int8   `gorm:"unique_index:only;column:identity_type;type:tinyint(4) unsigned;not null" json:"identity_type"` // 1手机号 2邮箱 3用户名 4qq 5微信 6腾讯微博 7新浪微博
	IDentifier   string `gorm:"column:identifier;type:varchar(50);not null" json:"identifier"`                                 // 手机号 邮箱 用户名或第三方应用的唯一标识
	Certificate  string `gorm:"column:certificate;type:varchar(20);not null" json:"certificate"`                               // 密码凭证(站内的保存密码，站外的不保存或保存token)
	CreateTime   int    `gorm:"column:create_time;type:int(11) unsigned;not null" json:"create_time"`                          // 绑定时间
	UpdateTime   int    `gorm:"column:update_time;type:int(11) unsigned;not null" json:"update_time"`                          // 更新绑定时间
}

// UserBase 用户基础信息表
type UserBase struct {
	UId            int64  `gorm:"primary_key;column:uid;type:bigint(20);not null" json:"uid"`                      // 用户ID
	UserRole       int8   `gorm:"column:user_role;type:tinyint(2) unsigned;not null" json:"user_role"`             // 2正常用户 3禁言用户 4虚拟用户 5运营
	RegisterSource int8   `gorm:"column:register_source;type:tinyint(4) unsigned;not null" json:"register_source"` // 注册来源：1手机号 2邮箱 3用户名 4qq 5微信 6腾讯微博 7新浪微博
	UserName       string `gorm:"column:user_name;type:varchar(32);not null" json:"user_name"`                     // 用户账号，必须唯一
	UserPwsd       string `gorm:"column:user_pwsd;type:varchar(128);not null" json:"user_pwsd"`                    // 用户密码
	NickName       string `gorm:"column:nick_name;type:varchar(32);not null" json:"nick_name"`                     // 用户昵称
	Gender         int8   `gorm:"column:gender;type:tinyint(1) unsigned;not null" json:"gender"`                   // 用户性别 0-female 1-male
	Birthday       int64  `gorm:"column:birthday;type:bigint(20) unsigned;not null" json:"birthday"`               // 用户生日
	Signature      string `gorm:"column:signature;type:varchar(255);not null" json:"signature"`                    // 用户个人签名
	Mobile         string `gorm:"column:mobile;type:varchar(16);not null" json:"mobile"`                           // 手机号码(唯一)
	MobileBindTime int    `gorm:"column:mobile_bind_time;type:int(11) unsigned;not null" json:"mobile_bind_time"`  // 手机号码绑定时间
	Email          string `gorm:"column:email;type:varchar(100);not null" json:"email"`                            // 邮箱(唯一)
	EmailBindTime  int    `gorm:"column:email_bind_time;type:int(11) unsigned;not null" json:"email_bind_time"`    // 邮箱绑定时间
	Face           string `gorm:"column:face;type:varchar(255);not null" json:"face"`                              // 头像
	Face200        string `gorm:"column:face200;type:varchar(255);not null" json:"face200"`                        // 头像 200x200x80
	Srcface        string `gorm:"column:srcface;type:varchar(255);not null" json:"srcface"`                        // 原图头像
	CreateTime     int    `gorm:"column:create_time;type:int(11) unsigned;not null" json:"create_time"`            // 创建时间
	UpdateTime     int    `gorm:"column:update_time;type:int(11) unsigned;not null" json:"update_time"`            // 修改时间
	PushToken      string `gorm:"column:push_token;type:varchar(50);not null" json:"push_token"`                   // 用户设备push_token
}

// UserExtra 用户额外信息表
type UserExtra struct {
	UId           int64  `gorm:"primary_key;column:uid;type:bigint(20);not null" json:"uid"`            // 用户 ID
	Vendor        string `gorm:"column:vendor;type:varchar(64);not null" json:"vendor"`                 // 手机厂商：apple|htc|samsung，很少用
	ClientName    string `gorm:"column:client_name;type:varchar(50);not null" json:"client_name"`       // 客户端名称，如hjskang
	ClientVersion string `gorm:"column:client_version;type:varchar(50);not null" json:"client_version"` // 客户端版本号，如7.0.1
	OsName        string `gorm:"column:os_name;type:varchar(16);not null" json:"os_name"`               // 设备号:android|ios
	OsVersion     string `gorm:"column:os_version;type:varchar(16);not null" json:"os_version"`         // 系统版本号:2.2|2.3|4.0|5.1
	DeviceName    string `gorm:"column:device_name;type:varchar(32);not null" json:"device_name"`       // 设备型号，如:iphone6s、u880、u8800
	DeviceID      string `gorm:"column:device_id;type:varchar(128);not null" json:"device_id"`          // 设备ID
	IDfa          string `gorm:"column:idfa;type:varchar(50);not null" json:"idfa"`                     // 苹果设备的IDFA
	IDfv          string `gorm:"column:idfv;type:varchar(50);not null" json:"idfv"`                     // 苹果设备的IDFV
	Market        string `gorm:"column:market;type:varchar(20);not null" json:"market"`                 // 来源
	CreateTime    int    `gorm:"column:create_time;type:int(11) unsigned;not null" json:"create_time"`  // 添加时间
	UpdateTime    int    `gorm:"column:update_time;type:int(11) unsigned;not null" json:"update_time"`  // 更新时间
	Extend1       string `gorm:"column:extend1;type:varchar(100);not null" json:"extend1"`              // 扩展字段1
	Extend2       string `gorm:"column:extend2;type:varchar(100);not null" json:"extend2"`              // 扩展字段2
	Extend3       string `gorm:"column:extend3;type:varchar(100);not null" json:"extend3"`              // 扩展字段3
}

// UserInfoUpdate 用户注册日志表
type UserInfoUpdate struct {
	ID              int64  `gorm:"primary_key;column:id;type:bigint(20);not null" json:"-"`                     // 自增ID
	UId             int64  `gorm:"column:uid;type:bigint(20) unsigned;not null" json:"uid"`                     // 用户ID
	AttributeName   string `gorm:"column:attribute_name;type:varchar(30);not null" json:"attribute_name"`       // 属性名
	AttributeOldVal string `gorm:"column:attribute_old_val;type:varchar(30);not null" json:"attribute_old_val"` // 属性对应旧的值
	AttributeNewVal string `gorm:"column:attribute_new_val;type:varchar(30);not null" json:"attribute_new_val"` // 属性对应新的值
	UpdateTime      int    `gorm:"column:update_time;type:int(11);not null" json:"update_time"`                 // 修改时间
}

// UserLocation 用户定位表
type UserLocation struct {
	UId          int64   `gorm:"primary_key;column:uid;type:bigint(20) unsigned;not null" json:"uid"` // 用户ID
	CurrNation   string  `gorm:"column:curr_nation;type:varchar(10);not null" json:"curr_nation"`     // 所在地国
	CurrProvince string  `gorm:"column:curr_province;type:varchar(10);not null" json:"curr_province"` // 所在地省
	CurrCity     string  `gorm:"column:curr_city;type:varchar(10);not null" json:"curr_city"`         // 所在地市
	CurrDistrict string  `gorm:"column:curr_district;type:varchar(20);not null" json:"curr_district"` // 所在地地区
	Location     string  `gorm:"column:location;type:varchar(255);not null" json:"location"`          // 具体地址
	Longitude    float64 `gorm:"column:longitude;type:decimal(10,6)" json:"longitude"`                // 经度
	Latitude     float64 `gorm:"column:latitude;type:decimal(10,6)" json:"latitude"`                  // 纬度
	UpdateTime   int     `gorm:"column:update_time;type:int(11) unsigned" json:"update_time"`         // 修改时间
}

// UserLoginLog 登陆日志表
type UserLoginLog struct {
	ID         int64  `gorm:"primary_key;column:id;type:bigint(20);not null" json:"-"`
	UId        int64  `gorm:"index:idx_uid_type_time;column:uid;type:bigint(20) unsigned;not null" json:"uid"`   // 用户uid
	Type       int8   `gorm:"index:idx_uid_type_time;column:type;type:tinyint(3) unsigned;not null" json:"type"` // 登录方式 第三方/邮箱/手机等
	Command    int8   `gorm:"column:command;type:tinyint(3) unsigned;not null" json:"command"`                   // 操作类型 1登陆成功  2登出成功 3登录失败 4登出失败
	Version    string `gorm:"column:version;type:varchar(32);not null" json:"version"`                           // 客户端版本号
	Client     string `gorm:"column:client;type:varchar(20);not null" json:"client"`                             // 客户端
	DeviceID   string `gorm:"column:device_id;type:varchar(64);not null" json:"device_id"`                       // 登录时设备号
	Lastip     string `gorm:"column:lastip;type:varchar(32);not null" json:"lastip"`                             // 登录ip
	Os         string `gorm:"column:os;type:varchar(16);not null" json:"os"`                                     // 手机系统
	Osver      string `gorm:"column:osver;type:varchar(32);not null" json:"osver"`                               // 系统版本
	Text       string `gorm:"column:text;type:varchar(200);not null" json:"text"`
	CreateTime int    `gorm:"index:idx_uid_type_time;index;column:create_time;type:int(11) unsigned;not null" json:"create_time"` // 操作时间
}

// UserRegisterLog 用户注册日志表
type UserRegisterLog struct {
	ID             int64  `gorm:"primary_key;column:id;type:bigint(20);not null" json:"-"`                         // 自增ID
	UId            int64  `gorm:"column:uid;type:bigint(20) unsigned;not null" json:"uid"`                         // 用户ID
	RegisterMethod uint8  `gorm:"column:register_method;type:tinyint(2) unsigned;not null" json:"register_method"` // 注册方式1手机号 2邮箱 3用户名 4qq 5微信 6腾讯微博 7新浪微博
	RegisterTime   int    `gorm:"column:register_time;type:int(11);not null" json:"register_time"`                 // 注册时间
	RegisterIP     string `gorm:"column:register_ip;type:varchar(16);not null" json:"register_ip"`                 // 注册IP
	RegisterClient string `gorm:"column:register_client;type:varchar(16);not null" json:"register_client"`         // 注册客户端
}

type User struct {
	UserAuth        UserAuth
	UserBase        UserBase
	UserExtra       UserExtra
	UserInfoUpdate  UserInfoUpdate
	UserLocation    UserLocation
	UserLoginLog    UserLoginLog
	UserRegisterLog UserRegisterLog
}

type UserLogin struct {
	User    User
	Account string `json:"account"`
	Pwsd    string `json:"pwsd"`
}

func (m *UserLogin) Login(ip string) (*JwtObj, error) {
	if m.Pwsd == "" {
		return nil, errors.New("password is required")
	}
	entity := &m.User
	cond := fmt.Sprintf("email = '%s' or user_name = '%s' or mobile = '%s'", m.Account, m.Account, m.Account)
	err := DBInstance.Where(cond).First(&entity.UserBase).Error
	if err != nil {
		if IsRecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(entity.UserBase.UserPwsd), []byte(m.Pwsd)); err != nil {
		return nil, err
	}
	return JwtGenerateToken(m, entity.UserBase.UId), nil
}
