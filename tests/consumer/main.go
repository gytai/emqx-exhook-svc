package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

// DeviceAuth 设备认证信息
type DeviceAuth struct {
	Owner          string `json:"owner"`
	OrgId          int64  `json:"orgId"`
	DeviceId       string `json:"deviceId"`
	DeviceType     string `json:"deviceType"`
	DeviceProtocol string `json:"deviceProtocol"`
	ProductId      string `json:"productId"`
	RuleChainId    string `json:"ruleChainId"`
	Name           string `json:"name"`
	CreatedAt      int64  `json:"created_at"`
	ExpiredAt      int64  `json:"expired_at"`
}

// DeviceEventInfo 设备事件信息
type DeviceEventInfo struct {
	DeviceId   string      `json:"deviceId"`
	DeviceAuth *DeviceAuth `json:"deviceAuth"`
	Datas      string      `json:"datas"`
	RequestId  string      `json:"requestId"`
	Topic      string      `json:"topic"`
}

func main() {
	// 解析命令行参数
	var (
		rabbitURL    = flag.String("url", "amqp://guest:guest@localhost:5672/", "RabbitMQ连接URL")
		exchangeName = flag.String("exchange", "exchange_default_v1_devices_me_telemetry", "要订阅的交换机名称")
		queueName    = flag.String("queue", "", "队列名称（默认为随机生成）")
		consumerName = flag.String("consumer", "consumer1", "消费者名称")
	)
	flag.Parse()

	// 如果未指定队列名称，则使用消费者名称作为队列名前缀
	if *queueName == "" {
		*queueName = fmt.Sprintf("%s_queue", *consumerName)
	}

	log.Printf("启动消费者: %s, 连接到: %s, 订阅交换机: %s, 队列: %s\n",
		*consumerName, *rabbitURL, *exchangeName, *queueName)

	// 连接RabbitMQ
	conn, err := amqp.Dial(*rabbitURL)
	if err != nil {
		log.Fatalf("无法连接到RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 创建通道
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("无法创建通道: %v", err)
	}
	defer ch.Close()

	// 声明交换机
	err = ch.ExchangeDeclare(
		*exchangeName, // 交换机名称
		"fanout",      // 类型
		true,          // 持久化
		false,         // 自动删除
		false,         // 内部使用
		false,         // 不等待
		nil,           // 参数
	)
	if err != nil {
		log.Fatalf("无法声明交换机: %v", err)
	}

	// 声明队列
	q, err := ch.QueueDeclare(
		*queueName, // 队列名
		true,       // 持久化
		false,      // 自动删除
		false,      // 排他性
		false,      // 不等待
		nil,        // 参数
	)
	if err != nil {
		log.Fatalf("无法声明队列: %v", err)
	}

	// 将队列绑定到交换机
	err = ch.QueueBind(
		q.Name,        // 队列名
		"",            // 路由键（fanout模式忽略）
		*exchangeName, // 交换机名
		false,         // 不等待
		nil,           // 参数
	)
	if err != nil {
		log.Fatalf("无法绑定队列到交换机: %v", err)
	}

	// 消费消息
	msgs, err := ch.Consume(
		q.Name,        // 队列名
		*consumerName, // 消费者标签
		true,          // 自动确认
		false,         // 排他性
		false,         // 不等待
		false,         // 参数
		nil,           // 参数
	)
	if err != nil {
		log.Fatalf("无法注册消费者: %v", err)
	}

	// 创建一个通道来接收中断信号
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	log.Printf("消费者 %s 已启动，等待消息...\n", *consumerName)

	// 在goroutine中消费消息
	go func() {
		for msg := range msgs {
			log.Printf("收到消息: [%d] 字节\n", len(msg.Body))

			// 尝试解析JSON
			var eventInfo DeviceEventInfo
			if err := json.Unmarshal(msg.Body, &eventInfo); err != nil {
				log.Printf("JSON解析错误: %v\n", err)
				log.Printf("原始数据: %s\n", string(msg.Body))
				continue
			}

			// 打印消息详情
			log.Printf("Topic: %s\n", eventInfo.Topic)
			log.Printf("DeviceId: %s\n", eventInfo.DeviceId)
			log.Printf("Datas: %s\n", eventInfo.Datas)

			if eventInfo.DeviceAuth != nil {
				log.Printf("设备信息: 类型=%s, 名称=%s\n",
					eventInfo.DeviceAuth.DeviceType,
					eventInfo.DeviceAuth.Name)
			}

			log.Println("---------------------------------------------")
		}
	}()

	// 等待中断信号
	<-stop
	log.Println("关闭消费者...")
}
