package user

import (
	. "GoServer/utils/string"
	. "GoServer/utils/time"
)

const (
	UNKnown_Gender = iota
	Female         //女
	Male           //男
)

type UserType int8

const (
	UNKNOWN  UserType = iota
	MOBILE            //手机
	EMAIL             //邮箱
	USERNAME          //用户名
	QQ                //QQ
	WECHAT            //微信
	WEIBO             //微博
)

const (
	LOGIN_SUCCEED  = 1
	LOGIN_FAILURED = 3
)

// UserAuth 用户授权表
type UserAuth struct {
	ID           int64  `json:"id" gorm:"pk autoincr BIGINT(20) 'id'"`
	UID          int64  `json:"uid" gorm:"not null default 0 comment('用户id') BIGINT(20) 'uid'"`
	IdentityType int8   `json:"identity_type" gorm:"default 1 comment('1手机号 2邮箱 3用户名 4qq 5微信 6腾讯微博 7新浪微博') TINYINT(1) 'identity_type'"`
	Identifier   string `json:"identifier" gorm:"default 'NULL' comment('手机号 邮箱 用户名或第三方应用的唯一标识') VARCHAR(50) 'identifier'"`
	Certificate  string `json:"certificate" gorm:"default 'NULL' comment('密码凭证(站内的保存密码，站外的不保存或保存token)') VARCHAR(20) 'certificate'"`
	CreateTime   int    `json:"create_time" gorm:"not null default 0 comment('绑定时间') INT(14) 'create_time'"`
	UpdateTime   int    `json:"update_time" gorm:"not null default 0 comment('更新绑定时间') INT(14) 'update_time'"`
}

// UserBase 用户基础信息表
type UserBase struct {
	UID            int64  `json:"uid" gorm:"not null pk autoincr comment('用户ID') BIGINT(20) 'uid'"`
	UserRole       int8   `json:"user_role" gorm:"default 2 comment('2正常用户 3禁言用户 4虚拟用户 5运营') TINYINT(1) 'user_role'"`
	RegisterSource int8   `json:"register_source" gorm:"default 0 comment('注册来源：1手机号 2邮箱 3用户名 4qq 5微信 6腾讯微博 7新浪微博') TINYINT(1) 'register_source'"`
	UserName       string `json:"user_name" gorm:"default 'NULL' comment('用户账号，必须唯一') VARCHAR(32) 'user_name'"`
	UserPwsd       string `json:"user_pwsd" gorm:"default 'NULL' VARCHAR(64) 'user_pwsd'"`
	NickName       string `json:"nick_name" gorm:"default 'NULL' comment('用户昵称') VARCHAR(32) 'nick_name'"`
	Gender         int8   `json:"gender" gorm:"default 0 comment('用户性别 0-female 1-male') TINYINT(1) 'gender'"`
	Birthday       int64  `json:"birthday" gorm:"default 0 comment('用户生日') BIGINT(20) 'birthday'"`
	Signature      string `json:"signature" gorm:"default 'NULL' comment('用户个人签名') VARCHAR(255) 'signature'"`
	Mobile         string `json:"mobile" gorm:"default 'NULL' comment('手机号码(唯一)') VARCHAR(16) 'mobile'"`
	MobileBindTime int64  `json:"mobile_bind_time" gorm:"default 0 comment('手机号码绑定时间') INT(11) 'mobile_bind_time'"`
	Email          string `json:"email" gorm:"default 'NULL' comment('邮箱(唯一)') VARCHAR(128) 'email'"`
	EmailBindTime  int64  `json:"email_bind_time" gorm:"default 0 comment('邮箱绑定时间') INT(11) 'email_bind_time'"`
	Face           string `json:"face" gorm:"default 'NULL' comment('头像') VARCHAR(255) 'face'"`
	Face200        string `json:"face200" gorm:"default 'NULL' comment('头像 200x200x80') VARCHAR(255) 'face200'"`
	Srcface        string `json:"srcface" gorm:"default 'NULL' comment('原图头像') VARCHAR(255) 'srcface'"`
	CreateTime     int    `json:"create_time" gorm:"not null comment('创建时间') INT(11) 'create_time'"`
	UpdateTime     int    `json:"update_time" gorm:"not null comment('修改时间') INT(11) 'update_time'"`
	PushToken      string `json:"push_token" gorm:"not null comment('用户设备push_token') VARCHAR(64) 'push_token'"`
}

// UserExtra 用户额外信息表
type UserExtra struct {
	UID           int64  `json:"uid" gorm:"not null pk comment('用户 ID') BIGINT(20) 'uid'"`
	Vendor        string `json:"vendor" gorm:"default 'NULL' comment('手机厂商：apple|htc|samsung，很少用') VARCHAR(64) 'vendor'"`
	ClientName    string `json:"client_name" gorm:"default 'NULL' comment('客户端名称，如hjskang') VARCHAR(50) 'client_name'"`
	ClientVersion string `json:"client_version" gorm:"default 'NULL' comment('客户端版本号，如7.0.1') VARCHAR(50) 'client_version'"`
	OsName        string `json:"os_name" gorm:"default 'NULL' comment('设备号:android|ios') VARCHAR(16) 'os_name'"`
	OsVersion     string `json:"os_version" gorm:"default 'NULL' comment('系统版本号:2.2|2.3|4.0|5.1') VARCHAR(16) 'os_version'"`
	DeviceName    string `json:"device_name" gorm:"default 'NULL' comment('设备型号，如:iphone6s、u880、u8800') VARCHAR(32) 'device_name'"`
	DeviceId      string `json:"device_id" gorm:"default 'NULL' comment('设备ID') VARCHAR(128) 'device_id'"`
	Idfa          string `json:"idfa" gorm:"default 'NULL' comment('苹果设备的IDFA') VARCHAR(50) 'idfa'"`
	Idfv          string `json:"idfv" gorm:"default 'NULL' comment('苹果设备的IDFV') VARCHAR(50) 'idfv'"`
	Market        string `json:"market" gorm:"default 'NULL' comment('来源') VARCHAR(20) 'market'"`
	CreateTime    int    `json:"create_time" gorm:"not null default 0 comment('添加时间') INT(11) 'create_time'"`
	UpdateTime    int    `json:"update_time" gorm:"not null default 0 comment('更新时间') INT(11) 'update_time'"`
	Extend1       string `json:"extend1" gorm:"default 'NULL' comment('扩展字段1') VARCHAR(100) 'extend1'"`
	Extend2       string `json:"extend2" gorm:"default 'NULL' comment('扩展字段2') VARCHAR(100) 'extend2'"`
	Extend3       string `json:"extend3" gorm:"default 'NULL' comment('扩展字段3') VARCHAR(100) 'extend3'"`
}

// UserInfoUpdate 用户注册日志表
type UserInfoUpdate struct {
	ID              int64  `json:"id" gorm:"pk autoincr comment('自增ID') BIGINT(20) 'id'"`
	UID             int64  `json:"uid" gorm:"not null comment('用户ID') BIGINT(20) 'uid'"`
	AttributeName   string `json:"attribute_name" gorm:"default 'NULL' comment('属性名') VARCHAR(30) 'attribute_name'"`
	AttributeOldVal string `json:"attribute_old_val" gorm:"default 'NULL' comment('属性对应旧的值') VARCHAR(30) 'attribute_old_val'"`
	AttributeNewVal string `json:"attribute_new_val" gorm:"default 'NULL' comment('属性对应新的值') VARCHAR(30) 'attribute_new_val'"`
	UpdateTime      int    `json:"update_time" gorm:"not null comment('修改时间') INT(11) 'update_time'"`
}

// UserLocation 用户定位表
type UserLocation struct {
	UID          int64  `json:"uid" gorm:"not null pk comment('用户ID') BIGINT(20) 'uid'"`
	CurrNation   string `json:"curr_nation" gorm:"default 'NULL' comment('所在地国') VARCHAR(10) 'curr_nation'"`
	CurrProvince string `json:"curr_province" gorm:"default 'NULL' comment('所在地省') VARCHAR(10) 'curr_province'"`
	CurrCity     string `json:"curr_city" gorm:"default 'NULL' comment('所在地市') VARCHAR(10) 'curr_city'"`
	CurrDistrict string `json:"curr_district" gorm:"default 'NULL' comment('所在地地区') VARCHAR(20) 'curr_district'"`
	Location     string `json:"location" gorm:"default 'NULL' comment('具体地址') VARCHAR(255) 'location'"`
	Longitude    string `json:"longitude" gorm:"default NULL comment('经度') DECIMAL(10,6) 'longitude'"`
	Latitude     string `json:"latitude" gorm:"default NULL comment('纬度') DECIMAL(10,6) 'latitude'"`
	UpdateTime   int    `json:"update_time" gorm:"not null default 0 comment('修改时间') INT(11) 'update_time'"`
}

// UserLoginLog 登陆日志表
type UserLoginLog struct {
	ID         int64  `json:"id" gorm:"pk autoincr BIGINT(20) 'id'"`
	UID        int64  `json:"uid" gorm:"not null default 0 comment('用户uid') BIGINT(20) 'uid'"`
	Type       int8   `json:"type" gorm:"default 1 comment('登录方式 第三方/邮箱/手机等') TINYINT(1) 'type'"`
	Command    int8   `json:"command" gorm:"default 1 comment('操作类型 1登陆成功 2登出成功 3登录失败 4登出失败') TINYINT(1) 'command'"`
	Version    string `json:"version" gorm:"default '1.0' comment('客户端版本号') VARCHAR(32) 'version'"`
	Client     string `json:"client" gorm:"default 'dabaozha' comment('客户端') VARCHAR(20) 'client'"`
	DeviceId   string `json:"device_id" gorm:"default 'NULL' comment('登录时设备号') VARCHAR(64) 'device_id'"`
	Lastip     string `json:"lastip" gorm:"default 'NULL' comment('登录ip') VARCHAR(32) 'lastip'"`
	Os         string `json:"os" gorm:"default 'NULL' comment('手机系统') VARCHAR(16) 'os'"`
	Osver      string `json:"osver" gorm:"default 'NULL' comment('系统版本') VARCHAR(32) 'osver'"`
	Text       string `json:"text" gorm:"default 'NULL' VARCHAR(200) 'text'"`
	CreateTime int    `json:"create_time" gorm:"not null default 0 comment('操作时间') INT(11) 'create_time'"`
}

// UserRegisterLog 用户注册日志表
type UserRegisterLog struct {
	ID             int64  `json:"id" gorm:"pk autoincr comment('自增ID') BIGINT(20) 'id'"`
	UID            int64  `json:"uid" gorm:"not null comment('用户ID') BIGINT(20) 'uid'"`
	RegisterMethod int8   `json:"register_method" gorm:"default NULL comment('注册方式1手机号 2邮箱 3用户名 4qq 5微信 6腾讯微博 7新浪微博') TINYINT(1) 'register_method'"`
	RegisterTime   int    `json:"register_time" gorm:"not null comment('注册时间') INT(11) 'register_time'"`
	RegisterIP     string `json:"register_ip" gorm:"not null default '' comment('注册IP') VARCHAR(16) 'register_ip'"`
	RegisterClient string `json:"register_client" gorm:"default 'NULL' comment('注册客户端') VARCHAR(16) 'register_client'"`
}

func (auth *UserAuth) Create(userID int64, IDentityType int8, IDentifier string, Certificate string) {
	auth.UID = userID
	auth.IdentityType = IDentityType
	auth.Identifier = IDentifier
	auth.Certificate = Certificate
	auth.CreateTime = GetTimestamp()
}

func (login *UserLoginLog) Create(ip string, Command int8, loginType UserType, userID int64) {
	login.UID = userID
	login.Type = int8(loginType)
	login.CreateTime = GetTimestamp()
	login.Command = Command
	login.Lastip = ip
}

func (m *UserBase) CreateByDefaultInfo(userType UserType) {
	m.NickName = RandomDigitAndLetters(12)
	m.Gender = UNKnown_Gender
	m.UserRole = 1
	m.RegisterSource = int8(userType)
	m.CreateTime = GetTimestamp()
}
