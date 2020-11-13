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
	MysqlConfig `yaml:"mysql"`
}

// MysqlConfig 数据库配置信息
type MysqlConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
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
