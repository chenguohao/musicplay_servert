package common

import (
	_ "embed" // 引入 embed 包
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
)

// 使用 go:embed 将文件打包进二进制
//
//go:embed conf/config.yml
var defaultConfig []byte

//go:embed conf/test/config.yml
var testConfig []byte

//go:embed conf/prod/config.yml
var prodConfig []byte

type Config struct {
	Mode     string           `yaml:"env"`
	FileHost FileServerConfig `yaml:"fileServer"`
	Log      LogConfig        `yaml:"log"`
	Mysql    MysqlConfig      `yaml:"mysql"`
}

type FileServerConfig struct {
	Host string `yaml:"host"`
}

type LogConfig struct {
	Output string `yaml:"output"`
}

type MysqlConfig struct {
	DataSource string `yaml:"datasource"`
}

var ServerConfig Config

func LoadConfig(env string) error {
	var configData []byte

	switch env {
	case ModeTest:
		configData = testConfig
	case ModeProd:
		configData = prodConfig
	default:
		configData = defaultConfig
	}

	if len(configData) == 0 {
		fmt.Println("Config data is empty")
		return errors.New("miss config data")
	}

	return loadConfigData(&ServerConfig, configData)
}

func loadConfigData(config interface{}, data []byte) error {
	if len(data) == 0 {
		fmt.Println("Config data is required")
		return errors.New("miss config data")
	}

	err := yaml.Unmarshal(data, config)
	if err != nil {
		fmt.Printf("Unmarshal config data error: %v \n", err)
		return err
	}

	return nil
}
