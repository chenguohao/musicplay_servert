package main

import (
	"MusicPlayServer/common"
	"MusicPlayServer/common/log"
	"MusicPlayServer/dao"
	"MusicPlayServer/http"
	_ "embed"
	"flag"
	"fmt"
	"os"
)

//go:embed common/conf/config.yml
var configFile []byte
var envMode string
var ServerVersion = "1.0.01"

func init() {
	flag.StringVar(&envMode, "mode", common.ModeDebug, "Application run mode")

}

func main() {
	flag.Parse()
	fmt.Printf("============== %s ===============\n", ServerVersion)
	fmt.Printf("run mode: %s \n", envMode)

	err := common.LoadConfig(envMode)
	if err != nil {
		os.Exit(1)
	}

	err = log.InitLog(&common.ServerConfig)
	if err != nil {
		os.Exit(1)
	}

	dbClient, err := dao.InitDB(&common.ServerConfig)
	if err != nil {
		os.Exit(1)
	}
	dao.DBClient = dbClient

	http.StartHttpServer()
}
