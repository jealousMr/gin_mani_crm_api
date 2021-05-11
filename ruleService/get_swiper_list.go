package ruleService

import (
	"gin_mani_crm_api/clients"
	pb_mani "gin_mani_crm_api/pb"
	"gin_mani_crm_api/util"
	logx "github.com/amoghe/distillog"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

const SysUser = "sys"

func GetSwiperList(c *gin.Context){
	center_rpc,err := clients.GetCenterClient()
	if err != nil {
		logx.Errorln("GetSwiperList GetCenterClient error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.CLIENT_ERROR})
		return
	}
	resp,err := center_rpc.GetRuleByCondition(c,&pb_mani.GetRuleByConditionReq{
		User: SysUser,
		RuleType: pb_mani.RuleType_default_all,
	})
	if err != nil{
		logx.Errorf("GetSwiperList GetRuleByCondition error:%v",err)
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.RPC_ERROR})
		return
	}
	rules := make([]*pb_mani.Rule,0)
	ruleList := make([]*pb_mani.Rule,4)
	urlList := make([]*pb_mani.CrmUrl,0)
	for _,rule := range resp.Rules{
		rules  =append(rules,rule)
	}
	for i := 0; i < 4; i++ {
		rd := rand.New(rand.NewSource(time.Now().UnixNano()))
		index := rd.Intn(len(rules))
		ruleList[i] = rules[index]
		urlList = append(urlList,&pb_mani.CrmUrl{
			Tag: rules[index].RuleId,
			Url: rules[index].RuleConfig.SourceUrl,
			Name: rules[index].RuleConfig.SourceName,
			CrmType: pb_mani.CrmType_sys_crm,
		})
	}
	engine_rpc,err := clients.GetEngineClient()
	if err != nil {
		logx.Errorln("GetSwiperList GetEngineClient error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.CLIENT_ERROR})
		return
	}
	urlResp,err := engine_rpc.FileUriToCrm(c,&pb_mani.FileUriToCrmReq{
		CrmList: urlList,
		FileAction: pb_mani.FileAction_default_all_action,
	})
	if err != nil{
		logx.Errorln("GetSwiperList FileUriToCrm error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.RPC_ERROR})
		return
	}
	for _,r := range ruleList{
		r.RuleConfig.SourceUrl = urlResp.TagUrlMap[r.RuleId]
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": util.SUCCESS,
		"data":ruleList,
	})
	return
}