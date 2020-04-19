package conf

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
	Name string `yaml:"name"`
	Pwsd string `yaml:"pwsd"`
}

type Web struct {
	Port string `yaml:"port"`
}

type System struct {
	Timeformat string `yaml:"timeformat"`
}

type Config struct {
	MqttConfig   Mqtt   `yaml:"mqtt"`
	MysqlConfig  Mysql  `yaml:"mysql"`
	WebConfig    Web    `yaml:"web"`
	SystemConfig System `yaml:"system"`
}

var config = &Config{}
var ErrString error

func init() {
	fp, err := ioutil.ReadFile("./conf.yml")
	if err != nil {
		ErrString = errors.New(fmt.Sprintf("yamlFile.Get err #%v ", err))
		return
	}

	err = yaml.Unmarshal(fp, config)
	if err != nil {
		ErrString = errors.New(fmt.Sprintf("yamlFile.Unmarshal err #%v ", err))
		return
	}
}

func GetConfig() *Config {
	if ErrString != nil {
		return nil
	}
	return config
}

func (conf *Config) GetMqtt() *Mqtt {
	if ErrString != nil {
		return nil
	}
	return &(conf.MqttConfig)
}

func (conf *Config) GetMysql() *Mysql {
	if ErrString != nil {
		return nil
	}
	return &(conf.MysqlConfig)
}

func (conf *Config) GetWeb() *Web {
	if ErrString != nil {
		return nil
	}
	return &(conf.WebConfig)
}

func (conf *Config) GetSystem() *System {
	if ErrString != nil {
		return nil
	}
	return &(conf.SystemConfig)
}
