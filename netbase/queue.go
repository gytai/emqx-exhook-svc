package netbase

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// 根据topic获取交换机名称
func getExchangeNameByTopic(topic string) string {
	// 移除topic中可能包含的通配符
	cleanTopic := strings.Replace(topic, "#", "", -1)
	cleanTopic = strings.Replace(cleanTopic, "+", "", -1)

	// 将斜杠替换为下划线，避免命名问题
	exchangeName := strings.Replace(cleanTopic, "/", "_", -1)

	// 为不同类型的topic设置不同交换机
	if strings.HasPrefix(topic, "device/") {
		return "exchange_device_data_" + exchangeName
	} else if strings.HasPrefix(topic, "sensor/") {
		return "exchange_sensor_data_" + exchangeName
	} else if strings.HasPrefix(topic, "control/") {
		return "exchange_control_cmd_" + exchangeName
	} else if strings.HasPrefix(topic, "status/") {
		return "exchange_status_update_" + exchangeName
	}

	// 默认交换机
	return "exchange_default_" + exchangeName
}

// 发布消息到RabbitMQ交换机
func publishToQueue(queueName string, data *DeviceEventInfo) error {
	if RabbitMQClient == nil {
		return fmt.Errorf("RabbitMQ客户端未初始化")
	}

	// 使用topic作为交换机名称
	exchangeName := getExchangeNameByTopic(data.Topic)

	// 将数据序列化为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	// 发布消息到RabbitMQ交换机（消息会被标记为持久化）
	err = RabbitMQClient.PublishToExchange(exchangeName, jsonData)
	if err != nil {
		return fmt.Errorf("推送消息到RabbitMQ交换机失败: %v", err)
	}

	log.Printf("消息已推送到RabbitMQ交换机 %s，数据大小: %d字节", exchangeName, len(jsonData))
	return nil
}
