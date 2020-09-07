package utils

import (
	"errors"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"path/filepath"
)

type MqttConf struct {
	Host  string `yaml:"host"`
	Token string `yaml:"token"`
	Name  string `yaml:"name"`
	Pwsd  string `yaml:"pwsd"`
}

type MysqlConf struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	Pwsd     string `yaml:"pwsd"`
	Basedata string `yaml:"basedata"`
	Debug    bool   `yaml:"debug"`
}

type WebConf struct {
	Port string `yaml:"port"`
	Mode uint32 `yaml:"runMode"`
}

type LogConf struct {
	Enabel   bool   `yaml:"enabel"`
	Filepath string `yaml:"filepath"`
	Filename string `yaml:"filename"`
}

type ServiceConf struct {
	Mqtt bool `yaml:"mqtt"`
	Web  bool `yaml:"web"`
}

type JwtConf struct {
	AppSecret  string `yaml:"appSecret"`
	AppIss     string `yaml:"appIss"`
	ExpireTime uint32 `yaml:"expireTime"`
}

type WeAppConf struct {
	CodeToSessURL string `yaml:"CodeToSessURL"`
	AppID         string `yaml:"AppID"`
	AppSecret     string `yaml:"AppSecret"`
}

type SystemConf struct {
	Service    ServiceConf `yaml:"service"`
	Timeformat string      `yaml:"timeformat"`
	Log        LogConf     `yaml:"log"`
	Jwt        JwtConf     `yaml:"jwt"`
	WeApp      WeAppConf   `yaml:"weApp"`
}

type Config struct {
	Mqtt   []MqttConf `yaml:"mqtt"`
	Mysql  MysqlConf  `yaml:"mysql"`
	Web    WebConf    `yaml:"web"`
	System SystemConf `yaml:"system"`
}

var config = &Config{}
var ErrConfString error

func init() {

	exePath, _ := filepath.Abs("./")
	exeFilePath := path.Join(exePath, "./config/conf.yml")

	fp, err := ioutil.ReadFile(exeFilePath)
	if err != nil {
		ErrConfString = errors.New(fmt.Sprintf("yamlFile.Get err #%v ", err))
		return
	}

	err = yaml.Unmarshal(fp, config)
	if err != nil {
		ErrConfString = errors.New(fmt.Sprintf("yamlFile.Unmarshal err #%v ", err))
		return
	}
}

func GetConfig() *Config {
	if ErrConfString != nil {
		return nil
	}
	return config
}

func GetMqtt() []MqttConf {
	if ErrConfString != nil {
		return nil
	}
	return config.Mqtt
}

func GetMysql() *MysqlConf {
	if ErrConfString != nil {
		return nil
	}
	return &(config.Mysql)
}

func GetWeb() *WebConf {
	if ErrConfString != nil {
		return nil
	}
	return &(config.Web)
}

func GetSystem() *SystemConf {
	if ErrConfString != nil {
		return nil
	}
	return &(config.System)
}

func GetLog() *LogConf {
	if ErrConfString != nil {
		return nil
	}
	return &(config.System.Log)
}

func GetJwt() *JwtConf {
	if ErrConfString != nil {
		return nil
	}
	return &(config.System.Jwt)
}

func GetWeApp() *WeAppConf {
	if ErrConfString != nil {
		return nil
	}
	return &(config.System.WeApp)
}