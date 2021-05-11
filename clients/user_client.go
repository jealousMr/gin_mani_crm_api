package clients

import (
	"fmt"
	"gin_mani_crm_api/conf"
	pb_mani "gin_mani_crm_api/pb"
	"google.golang.org/grpc"
	"log"
)

var userClient pb_mani.GinUserServiceClient

func GetUserClient() (pb_mani.GinUserServiceClient, error) {
	if userClient != nil {
		return userClient, nil
	}
	var err error
	userClient, err = connectUser()
	return userClient, err
}

func connectUser() (pb_mani.GinUserServiceClient, error) {
	cf := conf.GetConfig()
	address := fmt.Sprintf("%s%s",cf.Server.Ip, cf.Client.User)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("user service client connect err: ", err)
	}
	c := pb_mani.NewGinUserServiceClient(conn)
	return c, err
}
