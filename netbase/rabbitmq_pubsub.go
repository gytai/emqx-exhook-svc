package netbase

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// DeclareExchange 声明交换机（持久化）
func (r *RabbitMQ) DeclareExchange(name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.channel.ExchangeDeclare(
		name,     // 交换机名称
		"fanout", // 类型：fanout表示广播给所有绑定的队列
		true,     // 持久化
		false,    // 自动删除
		false,    // 内部使用
		false,    // 不等待服务器响应
		nil,      // 额外参数
	)
}

// BindQueueToExchange 将队列绑定到交换机
func (r *RabbitMQ) BindQueueToExchange(queueName, exchangeName string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.channel.QueueBind(
		queueName,    // 队列名称
		"",           // 路由键（fanout模式下忽略）
		exchangeName, // 交换机名称
		false,        // 不等待服务器响应
		nil,          // 额外参数
	)
}

// PublishToExchange 发布消息到交换机（广播给所有绑定的队列）
func (r *RabbitMQ) PublishToExchange(exchangeName string, body []byte) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 确保交换机存在
	err := r.channel.ExchangeDeclare(
		exchangeName, // 交换机名称
		"fanout",     // 类型：fanout表示广播给所有绑定的队列
		true,         // 持久化
		false,        // 自动删除
		false,        // 内部使用
		false,        // 不等待服务器响应
		nil,          // 额外参数
	)
	if err != nil {
		return fmt.Errorf("声明交换机失败: %v", err)
	}

	return r.channel.Publish(
		exchangeName, // 交换机
		"",           // 路由键（fanout模式下忽略）
		false,        // 强制
		false,        // 立即
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // 持久化消息
			ContentType:  "application/json",
			Body:         body,
		},
	)
}
