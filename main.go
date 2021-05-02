package main

import (
	"gin_mani_crm_api/conf"
	"gin_mani_crm_api/userService"
	"gin_mani_crm_api/util"
	"github.com/gin-gonic/gin"
	logx "github.com/amoghe/distillog"
)

func main(){
	router := gin.Default()
	router.Use(util.Cors())
	router.Static("/static", "./static")

	// 用户相关
	user := router.Group("/v2")
	{
		user.POST("/login", userService.UserLogin)
	}

	configs := conf.GetConfig()
	logx.Infof("start mani_crm_api server...")
	router.Run(configs.Server.Port)
}
