package config

import (
	"errors"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"path/filepath"
)

type MqttOptions struct {
	Host  string `yaml:"host"`
	Token string `yaml:"token"`
	Name  string `yaml:"name"`
	Pwsd  string `yaml:"pwsd"`
}

type RedisOptions struct {
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	Auth           string `yaml:"auth"`
	MaxIdle        int    `yaml:"maxIdle"`
	MaxOpen        int    `yaml:"maxOpen"`
	ConnectTimeout int    `yaml:"connect_timeout"`
	ReadTimeout    int    `yaml:"read_timeout"`
	WriteTimeout   int    `yaml:"write_timeout"`
	IdleTimeout    int    `yaml:"idle_timeout"`
}

type MysqlOptions struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	Pwsd     string `yaml:"pwsd"`
	Basedata string `yaml:"basedata"`
	Debug    bool   `yaml:"debug"`
}

type WebOptions struct {
	Port string `yaml:"port"`
	Mode uint32 `yaml:"runMode"`
}

type LogOptions struct {
	Enabel   bool   `yaml:"enabel"`
	Mqttpath string `yaml:"mqttpath"`
	Webpath  string `yaml:"systempath"`
	Filename string `yaml:"filename"`
}

type ServiceOptions struct {
	Mqtt bool `yaml:"mqtt"`
	Web  bool `yaml:"web"`
}

type JwtOptions struct {
	AppSecret  string `yaml:"appSecret"`
	AppIss     string `yaml:"appIss"`
	ExpireTime uint32 `yaml:"expireTime"`
}

type WeAppOptions struct {
	CodeToSessURL string `yaml:"CodeToSessURL"`
	AppID         string `yaml:"AppID"`
	AppSecret     string `yaml:"AppSecret"`
}

type SystemOptions struct {
	Service    ServiceOptions `yaml:"service"`
	Timeformat string         `yaml:"timeformat"`
	Log        LogOptions     `yaml:"log"`
	Jwt        JwtOptions     `yaml:"jwt"`
	WeApp      WeAppOptions   `yaml:"weApp"`
}

type ConfigOptions struct {
	Mqtt   []MqttOptions `yaml:"mqtt"`
	Redis  RedisOptions  `yaml:redis`
	Mysql  MysqlOptions  `yaml:"mysql"`
	Web    WebOptions    `yaml:"web"`
	System SystemOptions `yaml:"system"`
}

var config = &ConfigOptions{}

var errOptionsString error

func init() {

	exePath, _ := filepath.Abs("./")
	exeFilePath := path.Join(exePath, "./utils/config/conf.yml")

	fp, err := ioutil.ReadFile(exeFilePath)
	if err != nil {
		errOptionsString = errors.New(fmt.Sprintf("yamlFile.Get err #%v ", err))
		return
	}

	err = yaml.Unmarshal(fp, config)
	if err != nil {
		errOptionsString = errors.New(fmt.Sprintf("yamlFile.Unmarshal err #%v ", err))
		return
	}
}

func chkOption(open error, option interface{}) (err error) {
	if open != nil && option == nil {
		err = errors.New("options read failed")
	}
	return
}

func GetConfig() (option *ConfigOptions, err error) {
	option = config
	err = chkOption(errOptionsString, option)
	return
}

func GetMqtt() (option []MqttOptions, err error) {
	option = config.Mqtt
	err = chkOption(errOptionsString, option)
	return
}

func GetRedis() (option *RedisOptions, err error) {
	option = &config.Redis
	err = chkOption(errOptionsString, option)
	return
}

func GetMysql() (option *MysqlOptions, err error) {
	option = &config.Mysql
	err = chkOption(errOptionsString, option)
	return
}

func GetWeb() (option *WebOptions, err error) {
	option = &config.Web
	err = chkOption(errOptionsString, option)
	return
}

func GetSystem() (option *SystemOptions, err error) {
	option = &config.System
	err = chkOption(errOptionsString, option)
	return
}

func GetLog() (option *LogOptions, err error) {
	option = &config.System.Log
	err = chkOption(errOptionsString, option)
	return
}

func GetJwt() (option *JwtOptions, err error) {
	option = &config.System.Jwt
	err = chkOption(errOptionsString, option)
	return
}

func GetWeApp() (option *WeAppOptions, err error) {
	option = &config.System.WeApp
	err = chkOption(errOptionsString, option)
	return
}
