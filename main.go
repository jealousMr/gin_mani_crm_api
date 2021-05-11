package main

import (
	"gin_mani_crm_api/conf"
	"gin_mani_crm_api/ruleService"
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
		user.GET("/user_task",userService.GetUserTaskPageList)
	}

	// rule
	rule := router.Group("/v2")
	{
		rule.GET("/swiper_list",ruleService.GetSwiperList)
		rule.GET("/home_list",ruleService.GetHomeList)
		rule.GET("/promote_list",ruleService.GetPromotePageList)
		rule.POST("/execute_v1",ruleService.ExecuteRuleV1)
	}

	configs := conf.GetConfig()
	logx.Infof("start mani_crm_api server...")
	router.Run(configs.Server.Port)
}
