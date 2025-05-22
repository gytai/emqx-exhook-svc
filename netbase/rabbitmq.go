package netbase

import (
	"fmt"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQClient RabbitMQ客户端
var RabbitMQClient *RabbitMQ

// RabbitMQ 连接和通道封装
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	mutex   sync.Mutex
}

// RabbitMQConfig RabbitMQ配置
type RabbitMQConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	VHost    string
}

// NewRabbitMQ 创建新的RabbitMQ连接
func NewRabbitMQ(config RabbitMQConfig) (*RabbitMQ, error) {
	// 构建连接URL
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.VHost)

	// 连接RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("连接RabbitMQ失败: %v", err)
	}

	// 创建通道
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("创建通道失败: %v", err)
	}

	log.Println("RabbitMQ连接成功")
	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}, nil
}

// Close 关闭连接
func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// DeclareQueue 声明队列（持久化）
func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.channel.QueueDeclare(
		name,  // 队列名称
		true,  // 持久化
		false, // 自动删除
		false, // 排他性
		false, // 不等待服务器响应
		nil,   // 额外参数
	)
}

// PublishMessage 发布消息到队列（持久化消息）
func (r *RabbitMQ) PublishMessage(queueName string, body []byte) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 确保队列存在
	_, err := r.DeclareQueue(queueName)
	if err != nil {
		return err
	}

	return r.channel.Publish(
		"",        // 交换机
		queueName, // 路由键（队列名）
		false,     // 强制
		false,     // 立即
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // 持久化消息
			ContentType:  "application/json",
			Body:         body,
		},
	)
}
