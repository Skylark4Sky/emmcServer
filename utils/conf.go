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

type Mysql struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Name string `yaml:"name"`
	Pwsd string `yaml:"pwsd"`
	Basedata string `yaml:"basedata"`
}

type Web struct {
	Port string `yaml:"port"`
}

type Log struct {
	Enabel bool `yaml:"enabel"`
	Filepath string `yaml:"filepath"`
	Filename string `yaml:"filename"`
}

type System struct {
	Timeformat string `yaml:"timeformat"`
	LogConfig Log `yaml:"log"`
}

type Config struct {
	MqttConfig   Mqtt   `yaml:"mqtt"`
	MysqlConfig  Mysql  `yaml:"mysql"`
	WebConfig    Web    `yaml:"web"`
	SystemConfig System `yaml:"system"`
}

var config = &Config{}
var ErrConfString error

func init() {
	fp, err := ioutil.ReadFile("./conf/conf.yml")
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

func (conf *Config) GetMqtt() *Mqtt {
	if ErrConfString != nil {
		return nil
	}
	return &(conf.MqttConfig)
}

func (conf *Config) GetMysql() *Mysql {
	if ErrConfString != nil {
		return nil
	}
	return &(conf.MysqlConfig)
}

func (conf *Config) GetWeb() *Web {
	if ErrConfString != nil {
		return nil
	}
	return &(conf.WebConfig)
}

func (conf *Config) GetSystem() *System {
	if ErrConfString != nil {
		return nil
	}
	return &(conf.SystemConfig)
}

func (system *System) GetLog() *Log {
	if ErrConfString != nil {
		return nil
	}
	return &(system.LogConfig)
}
