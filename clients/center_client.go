package clients

import (
	"fmt"
	"gin_mani_crm_api/conf"
	pb_mani "gin_mani_crm_api/pb"
	"google.golang.org/grpc"
	"log"
)

var centerClient pb_mani.GinCenterServiceClient

func GetCenterClient() (pb_mani.GinCenterServiceClient, error) {
	if centerClient != nil {
		return centerClient, nil
	}
	var err error
	centerClient, err = connectCenter()
	return centerClient, err
}

func connectCenter() (pb_mani.GinCenterServiceClient, error) {
	cf := conf.GetConfig()
	address := fmt.Sprintf("%s%s", cf.Server.Ip, cf.Server.Port)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("center service client connect err: ", err)
	}
	c := pb_mani.NewGinCenterServiceClient(conn)
	return c, err
}
