package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"recommend/config"
	"recommend/idl/gen/recommend"
	"recommend/model"
	"recommend/service"

	"google.golang.org/grpc"
)

func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
			return
		}
	}()
	return handler(ctx, req)
}

func main() {
	if err := config.InitConfig(config.CfgFileMain); err != nil {
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
	var opts []grpc.ServerOption
	// 注册interceptor
	var interceptor grpc.UnaryServerInterceptor

	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	// 实例化grpc Server
	grpcServer := grpc.NewServer(opts...)
	recommend.RegisterRecommenderServer(grpcServer, &RecommendServer{})
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
