# emqx-exhook-svr-go

这是一个使用 Go 语言编写的 ExHook 服务。

## 重要说明
如果emqx 版本大于等于5.9.0 则需要使用protobuf-v3以上版本

## 前提条件

- [Go](https://golang.org)（支持任意一个最新的三个主要版本）
- [Protocol buffer](https://developers.google.com/protocol-buffers) **编译器**, `protoc`
  安装说明请参阅
  [Protocol Buffer 编译器安装指南](https://grpc.io/docs/protoc-installation/)。
- **Go 插件** 用于协议编译器：
    - 使用以下命令安装 Go 的协议编译器插件:
    ```
    export GO111MODULE=on  # 启用模块模式
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```

    - 更新你的 PATH，以便 protoc 编译器可以找到这些插件:
    ```
    export PATH="$PATH:$(go env GOPATH)/bin"
    ```


## 运行

尝试编译 [*.proto](file:///Users/taiguangyin/Desktop/emqx-extension-examples-master/exhook-svr-go/protobuf/exhook.proto) 文件：

```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    protobuf/exhook.proto
```


运行服务器：
```
go run main.go
```
