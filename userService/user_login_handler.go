package userService

import (
	"gin_mani_crm_api/clients"
	pb_mani "gin_mani_crm_api/pb"
	"gin_mani_crm_api/util"
	logx "github.com/amoghe/distillog"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserLogin(c *gin.Context) {
	userInfo, err := bindUserRegisterParam(c)
	if err != nil {
		logx.Errorf("UserLogin bindUserRegisterParam error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.BIND_PARAM_ERROR})
		return
	}
	userClient, err := clients.GetUserClient()
	if err != nil {
		logx.Errorf("UserLogin get userClient error:%v", err)
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.CLIENT_ERROR})
		return
	}
	_, err = userClient.AddAndUpdateUserInfo(c, &pb_mani.AddAndUpdateUserInfoReq{
		UserInfo: userInfo,
	})
	if err != nil {
		logx.Errorf("UserLogin AddAndUpdateUserInfo error:%v", err)
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.RPC_ERROR})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": util.SUCCESS})
	return
}

func bindUserRegisterParam(c *gin.Context) (*pb_mani.UserInfo, error) {
	obj := make(map[string]pb_mani.UserInfo)
	if err := c.BindJSON(&obj); err != nil {
		return nil, err
	}
	param := obj["user_info"]
	return &param, nil
}
