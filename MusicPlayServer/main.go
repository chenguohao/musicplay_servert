package main

import (
	"MusicPlayServer/common"
	"MusicPlayServer/common/config"
	"MusicPlayServer/common/log"
	"MusicPlayServer/dao"
	"MusicPlayServer/http"
	"flag"
	"fmt"
	"os"
)

var envMode string
var ServerVersion = "1.0.01"

func init() {
	flag.StringVar(&envMode, "mode", common.ModeDebug, "Application run mode")

}

func main() {
	flag.Parse()
	fmt.Printf("============== %s ===============\n", ServerVersion)
	fmt.Printf("run mode: %s \n", envMode)

	err := config.LoadConfig(envMode)
	if err != nil {
		os.Exit(1)
	}

	err = log.InitLog(&config.ServerConfig)
	if err != nil {
		os.Exit(1)
	}

	dbClient, err := dao.InitDB(&config.ServerConfig)
	if err != nil {
		os.Exit(1)
	}
	dao.DBClient = dbClient

	http.StartHttpServer()
}
