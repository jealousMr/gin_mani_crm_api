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

func GetHomeList(c *gin.Context) {
	pn := c.DefaultQuery("page_no", "1")
	ps := c.DefaultQuery("page_size", "6")
	pageNo, _ := strconv.ParseInt(pn, 10, 64)
	pageSize, _ := strconv.ParseInt(ps, 10, 64)
	center_rpc, err := clients.GetCenterClient()
	if err != nil {
		logx.Errorln("GetSwiperList GetCenterClient error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.CLIENT_ERROR})
		return
	}
	resp, err := center_rpc.GetRuleByRuleType(c, &pb_mani.GetRuleByRuleTypeReq{
		RuleType: pb_mani.RuleType_default_all,
		PageSize: pageSize,
		PageNo:   pageNo,
	})
	if err != nil {
		logx.Errorf("GetHomeList GetRuleByRuleType error:%v", err)
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.RPC_ERROR})
		return
	}
	urls := make([]*pb_mani.CrmUrl,0)
	for _,r := range resp.RuleList{
		urls = append(urls,&pb_mani.CrmUrl{
			Tag: r.RuleId,
			Url: r.RuleConfig.SourceUrl,
			Name: r.RuleConfig.SourceName,
			CrmType: pb_mani.CrmType_sys_crm,
		})
	}
	engine_rpc,err := clients.GetEngineClient()
	if err != nil {
		logx.Errorln("GetHomeList GetEngineClient error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.CLIENT_ERROR})
		return
	}
	urlResp,err := engine_rpc.FileUriToCrm(c,&pb_mani.FileUriToCrmReq{
		CrmList: urls,
		FileAction: pb_mani.FileAction_default_all_action,
	})
	if err != nil{
		logx.Errorln("GetHomeList FileUriToCrm error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.RPC_ERROR})
		return
	}
	for _,r := range resp.RuleList{
		r.RuleConfig.SourceUrl = urlResp.TagUrlMap[r.RuleId]
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  util.SUCCESS,
		"data": resp.RuleList,
		"page": resp.Page,
	})
	return
}
