package ruleService

type ExecuteParamModel struct {
	OpenId    string `json:"open_id"`
	RuleId    string `json:"rule_id"`
	RuleType  int64  `json:"rule_type"`
	ImageName string `json:"image_name"`
	ImageUrl  string `json:"image_url"`
	Desc      string `json:"desc"`
	ExecuteType int64 `json:"execute_type"`
	OutputName string `json:"output_name"`
	OutputUrl string `json:"output_url"`
}
