package config

import (
	"io/ioutil"
	"log"
	"runtime"
	"sync"

	"gopkg.in/yaml.v2"
)

// YamlConfig 全部配置信息
type YamlConfig struct {
	AppConfig     `yaml:"app"`
	MysqlConfig   `yaml:"mysql"`
	JWTConfig     `yaml:"jwt"`
	RedisConfig   `yaml:"redis"`
	MessageConfig `yaml:"message"`
}

// AppConfig 系统配置信息
type AppConfig struct {
	Name  string `yaml:"name"`
	Debug bool   `yaml:"debug"`
	Port  string `yaml:"port"`
}

// MysqlConfig 数据库配置信息
type MysqlConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// JWTConfig Jwt配置信息
type JWTConfig struct {
	ExpireTime int    `yaml:"expire_time"`
	Salt       string `yaml:"salt"`
}

// RedisConfig Redis配置信息
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database int    `yaml:"database"`
	Password string `yaml:"password"`
}

// MessageConfig 短信配置信息
type MessageConfig struct {
	ServerURL string `yaml:"server_url"`
	Account   string `yaml:"account"`
	Password  string `yaml:"password"`
	Signature string `yaml:"signature"`
}

var (
	// Conf 全局配置
	Conf YamlConfig
	once sync.Once
)

func init() {
	once.Do(func() {
		_, filename, _, _ := runtime.Caller(0)

		// 配置文件路径
		confFile := filename[:len(filename)-7] + "../.env.yaml"
		file, err := ioutil.ReadFile(confFile)
		if err != nil {
			log.Fatalln("配置文件读取失败")
		}
		yaml.Unmarshal(file, &Conf)
	})
}
