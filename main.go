package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"recommend/config"
	"recommend/idl/gen/recommend"
	"recommend/model"
	"recommend/service"
)

func main() {
	if err := config.InitConfig(config.DefaultCfg); err != nil {
		panic(err)
	}
	if err := model.InitModel(); err != nil {
		panic(err)
	}
	service.InitService()

	var port int
	var ip string
	flag.IntVar(&port, "port", 50000, "port for the service")
	flag.StringVar(&ip, "ip", "127.0.0.1", "ip for the service")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	recommend.RegisterRecommenderServer(grpcServer, &RecommendServer{})
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
