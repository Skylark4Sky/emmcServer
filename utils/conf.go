package utils

import (
	"errors"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

type Mqtt struct {
	Host  string `yaml:"host"`
	Token string `yaml:"token"`
	Name  string `yaml:"name"`
	Pwsd  string `yaml:"pwsd"`
}

type Redis struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	Auth        string `yaml:"auth"`
	Accesstoken string `yaml:accesstoken`
	Expiredtime string `expiredtime`
}

type RabbitMQ struct {
	Host string `yaml:"host"`
	Name string `yaml:"name"`
	Pwsd string `yaml:"pwsd"`
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
}

type Log struct {
	Enabel   bool   `yaml:"enabel"`
	Filepath string `yaml:"filepath"`
	Filename string `yaml:"filename"`
}

type Service struct {
	Mqtt bool `yaml:"mqtt"`
	Web  bool `yaml:"web"`
}

type System struct {
	Service    Service `yaml:"service"`
	LogConfig  Log     `yaml:"log"`
	Timeformat string  `yaml:"timeformat"`
}

type Config struct {
	MqttConfig     Mqtt     `yaml:"mqtt"`
	RedisConfg     Redis    `yaml:"redis"`
	RabbitMQConfig RabbitMQ `yaml:"rabbitMQ"`
	MysqlConfig    Mysql    `yaml:"mysql"`
	WebConfig      Web      `yaml:"web"`
	SystemConfig   System   `yaml:"system"`
}

var config = &Config{}
var ErrConfString error

func init() {
	fp, err := ioutil.ReadFile("./config/conf.yml")
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

func MqttConf() *Mqtt {
	if ErrConfString != nil {
		return nil
	}
	return &(config.MqttConfig)
}

func MysqlConf() *Mysql {
	if ErrConfString != nil {
		return nil
	}
	return &(config.MysqlConfig)
}

func WebConf() *Web {
	if ErrConfString != nil {
		return nil
	}
	return &(config.WebConfig)
}

func SystemConf() *System {
	if ErrConfString != nil {
		return nil
	}
	return &(config.SystemConfig)
}

func LogConf() *Log {
	if ErrConfString != nil {
		return nil
	}
	return &(config.SystemConfig.LogConfig)
}
