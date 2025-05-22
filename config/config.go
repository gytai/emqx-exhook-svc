package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config 应用配置结构
type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`

	Redis struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`

	RabbitMQ struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		VHost    string `yaml:"vhost"`
	} `yaml:"rabbitmq"`
}

// LoadConfig 从配置文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件错误: %v", err)
	}

	// 解析YAML
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件错误: %v", err)
	}

	log.Println("配置加载成功")
	return config, nil
}
