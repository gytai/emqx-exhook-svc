module emqx.io/grpc/exhook

go 1.24

replace emqx.io/grpc/exhook => ./

require (
	github.com/go-redis/redis/v8 v8.11.5
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.71.0
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
)
