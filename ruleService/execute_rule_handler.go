package ruleService

import (
	"gin_mani_crm_api/clients"
	pb_mani "gin_mani_crm_api/pb"
	"gin_mani_crm_api/util"
	logx "github.com/amoghe/distillog"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ExecuteRuleV1(c *gin.Context) {
	param, err := bindExecuteV1Param(c)
	if err != nil {
		logx.Errorf("ExecuteRuleV1 bindExecuteV1Param error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.BIND_PARAM_ERROR})
		return
	}

	center_rpc, err := clients.GetCenterClient()
	if err != nil {
		logx.Errorln("ExecuteRuleV1 GetCenterClient error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.CLIENT_ERROR})
		return
	}
	ruleResp, err := center_rpc.GetRuleByCondition(c, &pb_mani.GetRuleByConditionReq{
		RuleId: param.RuleId,
	})
	if err != nil {
		logx.Errorln("ExecuteRuleV1 GetRuleByCondition error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.RPC_ERROR})
		return
	}
	rule := ruleResp.Rules[param.RuleId]
	rule.User = param.OpenId
	if param.Desc != ""{
		rule.RuleConfig.DescText = param.Desc
	}
	if param.RuleType != 0{
		rule.RuleType = pb_mani.RuleType(param.RuleType)
	}

	util.GoParallel(c, func() {
		execReq := &pb_mani.ExecuteRuleReq{
			Rule: rule,
		}
		_, err := center_rpc.ExecuteRule(c, execReq)
		if err != nil {
			logx.Errorln("ExecuteRuleV1 GetRuleByCondition error")
		}
	})
	c.JSON(http.StatusOK, gin.H{
		"msg":  util.SUCCESS,
	})
	return

}

func bindExecuteV1Param(c *gin.Context) (*ExecuteParamModel, error) {
	obj := make(map[string]ExecuteParamModel)
	if err := c.BindJSON(&obj); err != nil {
		return nil, err
	}
	param := obj["model_info"]
	return &param, nil
}
