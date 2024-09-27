package http

import (
	"MusicPlayServer/common/log"
	. "MusicPlayServer/controller"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"time"
)

func StartHttpServer() {
	router := gin.Default()
	router.Use(ginzap.Ginzap(log.Logger(), "2006-01-02 15:04:05", false))
	router.Use(ginzap.RecoveryWithZap(log.Logger(), true))
	log.Info("start http server")
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"OPTIONS", "GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Token"},
		AllowCredentials: false,
		MaxAge:           1 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			//TODO: only allow my domain
			return true
		},
	}))

	v1group := router.Group("/v1")
	{
		v1group.POST("/appleSign", AuthWithApple)
		v1group.POST("/createPlaylist", CreatePlayList)
		v1group.GET("/getPlaylist", GetPlaylist)
		v1group.POST("/updateProfile", ReqestUpdateProfile)
		v1group.POST("/updatePlaylist", RequestUpdatePlaylist)
		v1group.POST("/deletePlaylist", RequestDeletePlaylist)
		v1group.POST("/like", RequestLike)
		v1group.POST("/addPlayCount", RequestAddPlayCount)
		v1group.POST("/deleteAccount", DeleteAccount)
	}

	router.Run(":8011")
}
