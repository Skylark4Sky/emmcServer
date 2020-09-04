package utils

import (
	"errors"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"path/filepath"
)

type Mqtt struct {
	Host  string `yaml:"host"`
	Token string `yaml:"token"`
	Name  string `yaml:"name"`
	Pwsd  string `yaml:"pwsd"`
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	Pwsd     string `yaml:"pwsd"`
	Basedata string `yaml:"basedata"`
	Debug    bool   `yaml:"debug"`
}

type Web struct {
	Port string `yaml:"port"`
	Mode uint32 `yaml:"runMode"`
}

type Log struct {
	Enabel   bool   `yaml:"enabel"`
	Filepath string `yaml:"filepath"`
	Filename string `yaml:"filename"`
}

type ServiceConf struct {
	Mqtt bool `yaml:"mqtt"`
	Web  bool `yaml:"web"`
}

type System struct {
	Service    ServiceConf `yaml:"service"`
	Timeformat string      `yaml:"timeformat"`
	LogConfig  Log         `yaml:"log"`
}

type Config struct {
	MqttConfig   []Mqtt `yaml:"mqtt"`
	MysqlConfig  Mysql  `yaml:"mysql"`
	WebConfig    Web    `yaml:"web"`
	SystemConfig System `yaml:"system"`
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

func GetMqtt() []Mqtt {
	if ErrConfString != nil {
		return nil
	}
	return config.MqttConfig
}

func GetMysql() *Mysql {
	if ErrConfString != nil {
		return nil
	}
	return &(config.MysqlConfig)
}

func GetWeb() *Web {
	if ErrConfString != nil {
		return nil
	}
	return &(config.WebConfig)
}

func GetSystem() *System {
	if ErrConfString != nil {
		return nil
	}
	return &(config.SystemConfig)
}

func GetLog() *Log {
	if ErrConfString != nil {
		return nil
	}
	return &(config.SystemConfig.LogConfig)
}
