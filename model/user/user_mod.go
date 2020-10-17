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

const (
	UNKNOWN  uint8 = iota
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
	ID           uint64  `gorm:"primary_key;column:id;type:bigint(20);not null" json:"-"`
	UID          uint64 `gorm:"unique_index:only;index:idx_uid;column:uid;type:bigint(20) unsigned;not null" json:"uid"` // 用户id
	IdentityType  uint8   `gorm:"unique_index:only;column:identity_type;type:tinyint(2) unsigned" json:"identity_type"`    // 1手机号 2邮箱 3用户名 4qq 5微信 6腾讯微博 7新浪微博
	Identifier   string `gorm:"column:identifier;type:varchar(64)" json:"identifier"`                                    // 手机号 邮箱 用户名或第三方应用的唯一标识
	Certificate  string `gorm:"column:certificate;type:varchar(64)" json:"certificate"`                                  // 密码凭证(站内的保存密码，站外的不保存或保存token)
	CreateTime   int64  `gorm:"column:create_time;type:bigint(13) unsigned;not null" json:"create_time"`                    // 绑定时间
	UpdateTime   int64 `gorm:"column:update_time;type:bigint(13) unsigned;not null" json:"update_time"`                    // 更新绑定时间
}

// UserBase 用户基础信息表
type UserBase struct {
	UID            uint64  `gorm:"primary_key;column:uid;type:bigint(20) unsigned;not null" json:"-"`      // 用户ID
	UserRole       uint8   `gorm:"column:user_role;type:tinyint(2) unsigned" json:"user_role"`             // 2正常用户 3禁言用户 4虚拟用户 5运营
	RegisterSource uint8   `gorm:"column:register_source;type:tinyint(2) unsigned" json:"register_source"` // 注册来源：1手机号 2邮箱 3用户名 4qq 5微信 6腾讯微博 7新浪微博
	UserName       string `gorm:"column:user_name;type:varchar(32)" json:"user_name"`                     // 用户账号，必须唯一
	UserPwsd       string `gorm:"column:user_pwsd;type:varchar(64)" json:"user_pwsd"`
	NickName       string `gorm:"column:nick_name;type:varchar(32)" json:"nick_name"`                    // 用户昵称
	Gender         uint8   `gorm:"column:gender;type:tinyint(1) unsigned" json:"gender"`                  // 用户性别 0-female 1-male
	Birthday       int64 `gorm:"column:birthday;type:bigint(13) unsigned" json:"birthday"`              // 用户生日
	Signature      string `gorm:"column:signature;type:varchar(255)" json:"signature"`                   // 用户个人签名
	Mobile         string `gorm:"column:mobile;type:varchar(16)" json:"mobile"`                          // 手机号码(唯一)
	MobileBindTime int64    `gorm:"column:mobile_bind_time;type:bigint(13) unsigned" json:"mobile_bind_time"` // 手机号码绑定时间
	Email          string `gorm:"column:email;type:varchar(128)" json:"email"`                           // 邮箱(唯一)
	EmailBindTime  int64    `gorm:"column:email_bind_time;type:bigint(13) unsigned" json:"email_bind_time"`   // 邮箱绑定时间
	Face           string `gorm:"column:face;type:varchar(255)" json:"face"`                             // 头像
	Face200        string `gorm:"column:face200;type:varchar(255)" json:"face200"`                       // 头像 200x200x80
	Srcface        string `gorm:"column:srcface;type:varchar(255)" json:"srcface"`                       // 原图头像
	CreateTime     int64 `gorm:"column:create_time;type:bigint(13) unsigned;not null" json:"create_time"`  // 创建时间
	UpdateTime     int64    `gorm:"column:update_time;type:bigint(13) unsigned;not null" json:"update_time"`  // 修改时间
	PushToken      string `gorm:"column:push_token;type:varchar(64);not null" json:"push_token"`         // 用户设备push_token
}

// UserExtra 用户额外信息表
type UserExtra struct {
	UID           uint64  `gorm:"primary_key;column:uid;type:bigint(20);not null" json:"-"`             // 用户 ID
	Vendor        string `gorm:"column:vendor;type:varchar(64)" json:"vendor"`                         // 手机厂商：apple|htc|samsung，很少用
	Language      string `gorm:"column:language;type:varchar(32)" json:"language"`                     // 客户端语言设置
	ClientName    string `gorm:"column:client_name;type:varchar(50)" json:"client_name"`               // 客户端名称，如hjskang
	ClientVersion string `gorm:"column:client_version;type:varchar(50)" json:"client_version"`         // 客户端版本号，如7.0.1
	OsName        string `gorm:"column:os_name;type:varchar(16)" json:"os_name"`                       // 设备号:android|ios
	OsVersion     string `gorm:"column:os_version;type:varchar(16)" json:"os_version"`                 // 系统版本号:2.2|2.3|4.0|5.1
	DeviceName    string `gorm:"column:device_name;type:varchar(32)" json:"device_name"`               // 设备型号，如:iphone6s、u880、u8800
	DeviceID      string `gorm:"column:device_id;type:varchar(128)" json:"device_id"`                  // 设备ID
	IDfa          string `gorm:"column:idfa;type:varchar(50)" json:"idfa"`                             // 苹果设备的IDFA
	IDfv          string `gorm:"column:idfv;type:varchar(50)" json:"idfv"`                             // 苹果设备的IDFV
	Market        string `gorm:"column:market;type:varchar(20)" json:"market"`                         // 来源
	CreateTime    int64    `gorm:"column:create_time;type:bigint(13) unsigned;not null" json:"create_time"` // 添加时间
	UpdateTime    int64 `gorm:"column:update_time;type:bigint(13) unsigned" json:"update_time"`          // 更新时间
	Extend1       string `gorm:"column:extend1;type:varchar(100)" json:"extend1"`                      // 扩展字段1
	Extend2       string `gorm:"column:extend2;type:varchar(100)" json:"extend2"`                      // 扩展字段2
	Extend3       string `gorm:"column:extend3;type:varchar(100)" json:"extend3"`                      // 扩展字段3
}

// UserInfoUpdate 用户注册日志表
type UserInfoUpdate struct {
	ID              uint64  `gorm:"primary_key;column:id;type:bigint(20);not null" json:"-"`              // 自增ID
	UID             uint64 `gorm:"column:uid;type:bigint(20) unsigned;not null" json:"uid"`              // 用户ID
	AttributeName   string `gorm:"column:attribute_name;type:varchar(30)" json:"attribute_name"`         // 属性名
	AttributeOldVal string `gorm:"column:attribute_old_val;type:varchar(30)" json:"attribute_old_val"`   // 属性对应旧的值
	AttributeNewVal string `gorm:"column:attribute_new_val;type:varchar(30)" json:"attribute_new_val"`   // 属性对应新的值
	UpdateTime      int64    `gorm:"column:update_time;type:bigint(13) unsigned;not null" json:"update_time"` // 修改时间
}

// UserLocation 用户定位表
type UserLocation struct {
	UID          uint64   `gorm:"primary_key;column:uid;type:bigint(20) unsigned;not null" json:"-"` // 用户ID
	CurrNation   string  `gorm:"column:curr_nation;type:varchar(10)" json:"curr_nation"`            // 所在地国
	CurrProvince string  `gorm:"column:curr_province;type:varchar(10)" json:"curr_province"`        // 所在地省
	CurrCity     string  `gorm:"column:curr_city;type:varchar(10)" json:"curr_city"`                // 所在地市
	CurrDistrict string  `gorm:"column:curr_district;type:varchar(20)" json:"curr_district"`        // 所在地地区
	Location     string  `gorm:"column:location;type:varchar(255)" json:"location"`                 // 具体地址
	Longitude    float64 `gorm:"column:longitude;type:decimal(10,6)" json:"longitude"`              // 经度
	Latitude     float64 `gorm:"column:latitude;type:decimal(10,6)" json:"latitude"`                // 纬度
	UpdateTime   int64  `gorm:"column:update_time;type:bigint(13) unsigned" json:"update_time"`       // 修改时间
}

// UserLoginLog 登陆日志表
type UserLoginLog struct {
	ID         uint64  `gorm:"primary_key;column:id;type:bigint(20);not null" json:"-"`
	UID        uint64 `gorm:"column:uid;type:bigint(20) unsigned;not null" json:"uid"` // 用户uid
	Type       uint8   `gorm:"column:type;type:tinyint(2) unsigned" json:"type"`        // 登录方式 第三方/邮箱/手机等
	Command    uint8   `gorm:"column:command;type:tinyint(2) unsigned" json:"command"`  // 操作类型 1登陆成功 2登出成功 3登录失败 4登出失败
	Version    string `gorm:"column:version;type:varchar(32)" json:"version"`          // 客户端版本号
	Client     string `gorm:"column:client;type:varchar(20)" json:"client"`            // 客户端
	DeviceID   string `gorm:"column:device_id;type:varchar(64)" json:"device_id"`      // 登录时设备号
	Lastip     string `gorm:"column:lastip;type:varchar(32)" json:"lastip"`            // 登录ip
	Os         string `gorm:"column:os;type:varchar(16)" json:"os"`                    // 手机系统
	Osver      string `gorm:"column:osver;type:varchar(32)" json:"osver"`              // 系统版本
	Text       string `gorm:"column:text;type:varchar(200)" json:"text"`
	CreateTime int64    `gorm:"column:create_time;type:bigint(13) unsigned;not null" json:"create_time"` // 操作时间
}

// UserRegisterLog 用户注册日志表
type UserRegisterLog struct {
	ID             uint64  `gorm:"primary_key;column:id;type:bigint(20);not null" json:"-"`                // 自增ID
	UID            uint64 `gorm:"column:uid;type:bigint(20) unsigned;not null" json:"uid"`                // 用户ID
	RegisterMethod uint8   `gorm:"column:register_method;type:tinyint(1) unsigned" json:"register_method"` // 注册方式1手机号 2邮箱 3用户名 4qq 5微信 6腾讯微博 7新浪微博
	RegisterTime   int64    `gorm:"column:register_time;type:bigint(13);not null" json:"register_time"`        // 注册时间
	RegisterIP     string `gorm:"column:register_ip;type:varchar(16);not null" json:"register_ip"`        // 注册IP
	RegisterClient string `gorm:"column:register_client;type:varchar(16)" json:"register_client"`         // 注册客户端
}

type CreateUserInfo struct {
	Auth     UserAuth //第三方授权登录需要
	Base     UserBase
	Extra    UserExtra
	Location UserLocation
	Log      UserRegisterLog
}

func (auth *UserAuth) Create(userID uint64, IDentityType uint8, IDentifier string, Certificate string) {
	auth.UID = userID
	auth.IdentityType = IDentityType
	auth.Identifier = IDentifier
	auth.Certificate = Certificate
	auth.CreateTime = GetTimestamp()
}

func (login *UserLoginLog) Create(ip string, Command uint8, loginType uint8, userID uint64) {
	login.UID = userID
	login.Type = loginType
	login.CreateTime = GetTimestamp()
	login.Command = Command
	login.Lastip = ip
}

func (m *UserBase) CreateByDefaultInfo(userType uint8) {
	m.NickName = RandomDigitAndLetters(12)
	m.Gender = UNKnown_Gender
	m.UserRole = 2
	m.RegisterSource = userType
	m.CreateTime = GetTimestamp()
}
