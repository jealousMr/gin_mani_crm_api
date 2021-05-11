package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
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
	cpath := fmt.Sprintf("%s/src/gin_mani_crm_api/conf/meta.yaml",os.Getenv("GOPATH"))
	content, err := ioutil.ReadFile(cpath)
	if err != nil {
		log.Fatalf("解析config.yaml读取错误: %v", err)
	}
	if yaml.Unmarshal(content, &config) != nil {
		log.Fatalf("解析config.yaml出错: %v", err)
	}
	return config
}