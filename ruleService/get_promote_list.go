package ruleService

import (
	"gin_mani_crm_api/clients"
	pb_mani "gin_mani_crm_api/pb"
	"gin_mani_crm_api/util"
	logx "github.com/amoghe/distillog"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetPromotePageList(c *gin.Context){
	pn := c.DefaultQuery("page_no", "1")
	ps := c.DefaultQuery("page_size", "6")
	pageNo, _ := strconv.ParseInt(pn, 10, 64)
	pageSize, _ := strconv.ParseInt(ps, 10, 64)
	center_rpc, err := clients.GetCenterClient()
	if err != nil {
		logx.Errorln("GetPromotePageList GetCenterClient error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.CLIENT_ERROR})
		return
	}
	resp, err := center_rpc.GetRuleByRuleType(c, &pb_mani.GetRuleByRuleTypeReq{
		RuleType: pb_mani.RuleType_default_image,
		PageSize: pageSize,
		PageNo:   pageNo,
	})
	if err != nil {
		logx.Errorf("GetPromotePageList GetRuleByRuleType error:%v", err)
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.RPC_ERROR})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  util.SUCCESS,
		"data": resp.RuleList,
		"page": resp.Page,
	})
	return
}