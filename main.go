package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	config "emqx.io/grpc/exhook/config"

	"emqx.io/grpc/exhook/kit"
	"emqx.io/grpc/exhook/netbase"
	pb "emqx.io/grpc/exhook/protobuf"
	"google.golang.org/grpc"
)

func main() {
	// 加载配置
	conf, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 从配置获取Redis连接参数 (仍然保留Redis连接用于验证设备)
	client, err := kit.NewRedisClient(
		conf.Redis.Host,
		conf.Redis.Password,
		conf.Redis.Port,
		conf.Redis.DB,
	)
	if err != nil {
		log.Panic("Redis连接错误")
	} else {
		log.Println("Redis连接成功")
	}
	netbase.RedisDb = client

	// 初始化RabbitMQ连接
	rabbitMQConfig := netbase.RabbitMQConfig{
		Host:     conf.RabbitMQ.Host,
		Port:     conf.RabbitMQ.Port,
		Username: conf.RabbitMQ.Username,
		Password: conf.RabbitMQ.Password,
		VHost:    conf.RabbitMQ.VHost,
	}

	rabbitMQ, err := netbase.NewRabbitMQ(rabbitMQConfig)
	if err != nil {
		log.Fatalf("RabbitMQ连接失败: %v", err)
	}
	defer rabbitMQ.Close()
	netbase.RabbitMQClient = rabbitMQ

	// 从配置获取服务端口
	port := fmt.Sprintf(":%d", conf.Server.Port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}
	s := grpc.NewServer()

	// 使用hooks服务
	hookServer := netbase.NewHookServer()
	pb.RegisterHookProviderServer(s, hookServer)

	// 优雅关闭
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("接收到终止信号，正在关闭服务...")
		s.GracefulStop()
	}()

	log.Printf("gRPC服务运行在::%d", conf.Server.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("服务创建失败: %v", err)
	}
}
