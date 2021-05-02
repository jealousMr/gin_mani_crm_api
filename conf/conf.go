package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Configs struct {
	Server struct{
		Port string `yaml:"port"`
		Ip string `yaml:"ip"`
	}
	Client struct{
		User string `yaml:"user"`
		Center string `yaml:"center"`
		Engine string `yaml:"engine"`
	}
}

func GetConfig() *Configs{
	config := &Configs{}
	content, err := ioutil.ReadFile("/home/xyl/src/gin_mani_crm_api/conf/meta.yaml")
	if err != nil {
		log.Fatalf("解析config.yaml读取错误: %v", err)
	}
	if yaml.Unmarshal(content, &config) != nil {
		log.Fatalf("解析config.yaml出错: %v", err)
	}
	return config
}