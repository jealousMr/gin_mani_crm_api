package userService

import (
	"context"
	"fmt"
	"gin_mani_crm_api/clients"
	pb_mani "gin_mani_crm_api/pb"
	"gin_mani_crm_api/util"
	logx "github.com/amoghe/distillog"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
)

// default get all

func GetUserTaskPageList(c *gin.Context) {
	openId := c.DefaultQuery("open_id", "")

	user_rpc, err := clients.GetUserClient()
	if err != nil {
		logx.Errorln("GetUserTaskPageList GetUserClient error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.CLIENT_ERROR})
		return
	}
	userResp, err := user_rpc.QueryUserInfoByIds(c, &pb_mani.QueryUserInfoByIdsReq{
		IdList: []string{openId},
	})
	if err != nil {
		logx.Errorln("GetUserTaskPageList QueryUserInfoByIds error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.RPC_ERROR})
		return
	}
	if userResp.Users == nil || len(userResp.Users) == 0 {
		logx.Errorln("GetUserTaskPageList invalid user error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.INVALID_USER})
		return
	}

	center_rpc, err := clients.GetCenterClient()
	if err != nil {
		logx.Errorln("GetUserTaskPageList GetCenterClient error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.CLIENT_ERROR})
		return
	}
	ruleResp, err := center_rpc.GetRuleByCondition(c, &pb_mani.GetRuleByConditionReq{
		User: openId,
	})
	if err != nil {
		logx.Errorln("GetUserTaskPageList GetRuleByCondition error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.RPC_ERROR})
		return
	}
	ruleMap := ruleResp.Rules
	ruleIdList := make([]string, 0)
	for _, r := range ruleResp.Rules {
		ruleIdList = append(ruleIdList, r.RuleId)
	}

	engine_rpc, err := clients.GetEngineClient()
	if err != nil {
		logx.Errorln("GetUserTaskPageList GetEngineClient error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.CLIENT_ERROR})
		return
	}
	taskMapResp, err := engine_rpc.GetTaskByRules(c, &pb_mani.GetTaskByRulesReq{
		RuleList: ruleIdList,
	})
	if err != nil {
		logx.Errorln("GetUserTaskPageList GetTaskByRules error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.RPC_ERROR})
		return
	}
	taskMap := taskMapResp.RuleTaskMap
	userTaskList := make([]*UserTaskModel, 0)

	for _, rule := range ruleMap {
		if _, ok := taskMap[rule.RuleId]; ok {
			userTaskList = append(userTaskList, &UserTaskModel{
				OpenId:       openId,
				TaskId:       taskMap[rule.RuleId].TaskId,
				RuleId:       rule.RuleId,
				RuleType:     int64(rule.RuleType),
				ImageName:    rule.RuleConfig.SourceName,
				ImageUrl:     rule.RuleConfig.SourceUrl,
				Desc:         rule.RuleConfig.DescText,
				ExecuteType:  int64(rule.ExecuteType),
				ExecuteState: int64(taskMap[rule.RuleId].ExecuteState),
				OutputName:   taskMap[rule.RuleId].OutputName,
				OutputUrl:    taskMap[rule.RuleId].OutputUrl,
			})
		}
	}
	uts, err := changeUrlToCrm(c, userTaskList, engine_rpc)
	if err != nil {
		logx.Errorln("GetUserTaskPageList changeUrlToCrm error")
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": util.RPC_ERROR})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  util.SUCCESS,
		"data": uts,
	})

	return
}
func changeUrlToCrm(ctx context.Context, uList []*UserTaskModel, engine_rpc pb_mani.GinEngineServiceClient) (uts []*UserTaskModel, err error) {
	uts = make([]*UserTaskModel, 0)
	defaultAllList := make([]*pb_mani.CrmUrl, 0)
	openAllList := make([]*pb_mani.CrmUrl, 0)
	dImageList := make([]*pb_mani.CrmUrl, 0)

	ruleToUserTask := make(map[string]*UserTaskModel)
	for _, u := range uList {
		ruleToUserTask[u.RuleId] = u
		switch pb_mani.RuleType(u.RuleType) {
		case pb_mani.RuleType_default_all:
			defaultAllList = append(defaultAllList, &pb_mani.CrmUrl{
				Tag:     fmt.Sprintf("%s-source-da", u.RuleId),
				Url:     u.ImageUrl,
				Name:    u.ImageName,
				CrmType: pb_mani.CrmType_user_crm,
			})
			defaultAllList = append(defaultAllList, &pb_mani.CrmUrl{
				Tag:     fmt.Sprintf("%s-output-da", u.RuleId),
				Url:     u.OutputUrl,
				Name:    u.OutputName,
				CrmType: pb_mani.CrmType_user_crm,
			})
			break
		case pb_mani.RuleType_open_all:
			openAllList = append(openAllList, &pb_mani.CrmUrl{
				Tag:     fmt.Sprintf("%s-source-oa", u.RuleId),
				Url:     u.ImageUrl,
				Name:    u.ImageName,
				CrmType: pb_mani.CrmType_user_crm,
			})
			openAllList = append(openAllList, &pb_mani.CrmUrl{
				Tag:     fmt.Sprintf("%s-output-oa", u.RuleId),
				Url:     u.OutputUrl,
				Name:    u.OutputName,
				CrmType: pb_mani.CrmType_user_crm,
			})
			break
		case pb_mani.RuleType_default_image:
			dImageList = append(dImageList, &pb_mani.CrmUrl{
				Tag:     fmt.Sprintf("%s-source-di", u.RuleId),
				Url:     u.ImageUrl,
				Name:    u.ImageName,
				CrmType: pb_mani.CrmType_user_crm,
			})
			dImageList = append(dImageList, &pb_mani.CrmUrl{
				Tag:     fmt.Sprintf("%s-output-di", u.RuleId),
				Url:     u.OutputUrl,
				Name:    u.OutputName,
				CrmType: pb_mani.CrmType_user_crm,
			})
			break
		}
	}

	var wg sync.WaitGroup

	wg.Add(1)
	util.GoParallel(ctx, func() {
		defer wg.Done()
		if len(defaultAllList) > 0 {
			urlResp, err := engine_rpc.FileUriToCrm(ctx, &pb_mani.FileUriToCrmReq{
				CrmList:    defaultAllList,
				FileAction: pb_mani.FileAction_default_all_action,
			})
			if err != nil {
				logx.Errorf("changeUrlToCrm FileUriToCrm defaultAllList error:%v", err)
				return
			}
			for _, da := range defaultAllList {
				ruleId := strings.Split(da.Tag, "-")[0]
				ruleToUserTask[ruleId].ImageUrl = urlResp.TagUrlMap[fmt.Sprintf("%s-source-da", ruleId)]
				ruleToUserTask[ruleId].OutputUrl = urlResp.TagUrlMap[fmt.Sprintf("%s-output-da", ruleId)]
			}
		}
	})

	wg.Add(1)
	util.GoParallel(ctx, func() {
		defer wg.Done()
		if len(openAllList) > 0 {
			urlResp, err := engine_rpc.FileUriToCrm(ctx, &pb_mani.FileUriToCrmReq{
				CrmList:    openAllList,
				FileAction: pb_mani.FileAction_open_all_action,
			})
			if err != nil {
				logx.Errorf("changeUrlToCrm FileUriToCrm openAllList error:%v", err)
				return
			}
			for _, da := range openAllList {
				ruleId := strings.Split(da.Tag, "-")[0]
				ruleToUserTask[ruleId].ImageUrl = urlResp.TagUrlMap[fmt.Sprintf("%s-source-oa", ruleId)]
				ruleToUserTask[ruleId].OutputUrl = urlResp.TagUrlMap[fmt.Sprintf("%s-output-oa", ruleId)]
			}
		}
	})

	wg.Add(1)
	util.GoParallel(ctx, func() {
		defer wg.Done()
		if len(dImageList) > 0 {
			urlResp, err := engine_rpc.FileUriToCrm(ctx, &pb_mani.FileUriToCrmReq{
				CrmList:    dImageList,
				FileAction: pb_mani.FileAction_default_image_action,
			})
			if err != nil {
				logx.Errorf("changeUrlToCrm FileUriToCrm dImageList error:%v", err)
				return
			}
			for _, da := range dImageList {
				ruleId := strings.Split(da.Tag, "-")[0]
				ruleToUserTask[ruleId].ImageUrl = urlResp.TagUrlMap[fmt.Sprintf("%s-source-di", ruleId)]
				ruleToUserTask[ruleId].OutputUrl = urlResp.TagUrlMap[fmt.Sprintf("%s-output-di", ruleId)]
			}
		}
	})
	wg.Wait()
	for _,u := range ruleToUserTask{
		uts = append(uts,u)
	}
	return uts, nil
}
