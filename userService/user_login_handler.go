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
	user_rpc, err := clients.GetUserClient()
	if err != nil {
		logx.Errorf("UserLogin get userClient error:%v", err)
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.CLIENT_ERROR})
		return
	}
	_, err = user_rpc.AddAndUpdateUserInfo(c, &pb_mani.AddAndUpdateUserInfoReq{
		UserInfo: &pb_mani.UserInfo{
			UserId: userInfo.OpenId,
			NickName: userInfo.NickName,
			AvatarUrl: userInfo.AvatarUrl,
			Gender: userInfo.Gender,
			Country: userInfo.Country,
			UserState: pb_mani.UserState_user_state_valid,
		},
	})
	if err != nil {
		logx.Errorf("UserLogin AddAndUpdateUserInfo error:%v", err)
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.RPC_ERROR})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": util.SUCCESS})
	return
}

func bindUserRegisterParam(c *gin.Context) (*UserLoginParam, error) {
	obj := make(map[string]UserLoginParam)
	if err := c.BindJSON(&obj); err != nil {
		return nil, err
	}
	param := obj["user_info"]
	return &param, nil
}

type UserLoginParam struct {
	OpenId string `json:"open_id"`
	NickName string `json:"nick_name"`
	AvatarUrl string `json:"avatar_url"`
	Gender int64 `json:"gender"`
	Country string `json:"country"`
}