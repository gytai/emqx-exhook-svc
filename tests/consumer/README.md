# RabbitMQ 测试消费者

这是一个用于测试 EMQX ExHook 服务将消息发送到 RabbitMQ 的测试消费者程序。

## 功能

- 连接到 RabbitMQ 服务器
- 创建一个队列并绑定到指定的交换机
- 消费从交换机接收的消息
- 解析并显示消息内容

## 使用方法

### 编译

```bash
cd tests/consumer
go build -o consumer
```

### 运行

基本用法：

```bash
./consumer
```

默认情况下，消费者将使用以下参数：
- RabbitMQ 连接 URL: `amqp://guest:guest@localhost:5672/`
- 交换机名称: `exchange_device_data_device_test`
- 消费者名称: `consumer1`
- 队列名称: `consumer1_queue`

### 自定义参数

您可以通过命令行参数自定义消费者行为：

```bash
./consumer -url=amqp://user:password@rabbitmq-host:5672/ -exchange=exchange_name -consumer=consumer_name -queue=queue_name
```

参数说明：
- `-url`: RabbitMQ 服务器连接 URL
- `-exchange`: 要订阅的交换机名称
- `-consumer`: 消费者标识名称
- `-queue`: 要创建的队列名称（如果不指定，将使用消费者名称加上后缀 "_queue"）

## 测试多个消费者

要测试多个消费者接收同一条消息，可以启动多个消费者实例，使它们订阅相同的交换机：

```bash
# 终端 1
./consumer -consumer=consumer1 -exchange=exchange_device_data_device_test

# 终端 2
./consumer -consumer=consumer2 -exchange=exchange_device_data_device_test
```

当 EMQX 发布消息到指定主题时，两个消费者都会收到相同的消息。

## 配合 EMQX 测试

1. 确保 EMQX 和 RabbitMQ 服务都已启动
2. 启动 exhook-svc 服务
3. 启动一个或多个消费者
4. 使用 MQTT 客户端发布消息到 EMQX 的主题，例如 `device/test`
5. 观察消费者输出的消息内容 