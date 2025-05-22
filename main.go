package main

import (
	"fmt"
	"log"
	"net"

	config "emqx.io/grpc/exhook/config"

	"emqx.io/grpc/exhook/kit"
	"emqx.io/grpc/exhook/netbase"
	pb "emqx.io/grpc/exhook/protobuf"
	"google.golang.org/grpc"
)

func main() {
	// 加载配置
	config, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 从配置获取Redis连接参数
	client, err := kit.NewRedisClient(
		config.Redis.Host,
		config.Redis.Password,
		config.Redis.Port,
		config.Redis.DB,
	)
	if err != nil {
		log.Panic("Redis连接错误")
	} else {
		log.Println("Redis连接成功")
	}
	netbase.RedisDb = client

	// 从配置获取服务端口
	port := fmt.Sprintf(":%d", config.Server.Port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}
	s := grpc.NewServer()

	// 使用hooks服务
	hookServer := netbase.NewHookServer()
	pb.RegisterHookProviderServer(s, hookServer)

	log.Printf("gRPC服务运行在::%d", config.Server.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("服务创建失败: %v", err)
	}
}
