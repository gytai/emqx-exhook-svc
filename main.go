package main

import (
	"log"
	"net"

	"emqx.io/grpc/exhook/kit"
	"emqx.io/grpc/exhook/netbase"
	pb "emqx.io/grpc/exhook/protobuf"
	"google.golang.org/grpc"
)

const (
	port = ":9000"
)

func main() {
	// redis连接
	client, err := kit.NewRedisClient("127.0.0.1", "123456", 6379, 0)
	if err != nil {
		log.Panic("Redis连接错误")
	} else {
		log.Println("Redis连接成功")
	}
	netbase.RedisDb = client

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}
	s := grpc.NewServer()

	// 使用hooks.go中定义的Server
	hookServer := netbase.NewHookServer()
	pb.RegisterHookProviderServer(s, hookServer)

	log.Println("gRPC服务运行在::9000")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("服务创建失败: %v", err)
	}
}
