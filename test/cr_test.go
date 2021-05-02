package test

import (
	"fmt"
	"gin_mani_crm_api/conf"
	"testing"
)

func TestAddr(t *testing.T){
	cf := conf.GetConfig()
	address := fmt.Sprintf("%s%s", cf.Server.Ip, cf.Server.Port)
	fmt.Println(address)
}
