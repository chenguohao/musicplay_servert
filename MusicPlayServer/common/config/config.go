package config

import (
	"MusicPlayServer/common"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Mode     string           `yaml:"env"`
	FileHost FileServerConfig `yaml:"fileServer"`
	Log      LogConfig        `yaml:"log"`
	Mysql    MysqlConfig      `yaml:mysql`
}
type FileServerConfig struct {
	Host string `yaml:host`
}

type LogConfig struct {
	Output string `yaml:output`
}

type MysqlConfig struct {
	DataSource string `yaml:datasource`
}

var ServerConfig Config

func LoadConfig(env string) error {
	cfgPath := "common/conf"
	cfgName := "config.yml"
	fname := cfgPath + "/" + cfgName

	switch env {
	case common.ModeTest:
		fname = cfgPath + "/test/" + cfgName
	case common.ModeProd:
		fname = cfgPath + "/prod/" + cfgName
	default:
		fname = cfgPath + "/" + cfgName
	}
	fmt.Printf("Load config from: %s \n", fname)

	return loadConfigFile(&ServerConfig, fname)
}

func loadConfigFile(config interface{}, fname string) error {
	if fname == "" {
		fmt.Printf("Config file is required. \n")
		return errors.New("miss config file")
	}

	data, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Printf("Read config file[%s] error: %v \n", fname, err)
		return err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		fmt.Printf("Unmarshal config file[%s] error: %v \n", fname, err)
		return err
	}

	fmt.Printf("Load config info: %#v \n", config)
	return nil
}
