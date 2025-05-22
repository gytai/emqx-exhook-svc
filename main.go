package main

import (
	"context"
	"emqx.io/grpc/exhook/kit"
	"emqx.io/grpc/exhook/netbase"
	pb "emqx.io/grpc/exhook/protobuf"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":9000"
)

// server is used to implement emqx_exhook_v3.s *server
type server struct {
	pb.UnimplementedHookProviderServer
}

// HookProviderServer callbacks

func (s *server) OnProviderLoaded(ctx context.Context, in *pb.ProviderLoadedRequest) (*pb.LoadedResponse, error) {
	hooks := []*pb.HookSpec{
		&pb.HookSpec{Name: "client.connect"},
		&pb.HookSpec{Name: "client.connack"},
		&pb.HookSpec{Name: "client.connected"},
		&pb.HookSpec{Name: "client.disconnected"},
		&pb.HookSpec{Name: "client.authenticate"},
		&pb.HookSpec{Name: "client.authorize"},
		&pb.HookSpec{Name: "client.subscribe"},
		&pb.HookSpec{Name: "client.unsubscribe"},
		&pb.HookSpec{Name: "session.created"},
		&pb.HookSpec{Name: "session.subscribed"},
		&pb.HookSpec{Name: "session.unsubscribed"},
		&pb.HookSpec{Name: "session.resumed"},
		&pb.HookSpec{Name: "session.discarded"},
		&pb.HookSpec{Name: "session.takenover"},
		&pb.HookSpec{Name: "session.terminated"},
		&pb.HookSpec{Name: "message.publish"},
		&pb.HookSpec{Name: "message.delivered"},
		&pb.HookSpec{Name: "message.acked"},
		&pb.HookSpec{Name: "message.dropped"},
	}
	return &pb.LoadedResponse{Hooks: hooks}, nil
}

func (s *server) OnProviderUnloaded(ctx context.Context, in *pb.ProviderUnloadedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *server) OnClientConnect(ctx context.Context, in *pb.ClientConnectRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *server) OnClientConnack(ctx context.Context, in *pb.ClientConnackRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *server) OnClientConnected(ctx context.Context, in *pb.ClientConnectedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *server) OnClientDisconnected(ctx context.Context, in *pb.ClientDisconnectedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *server) OnClientAuthenticate(ctx context.Context, in *pb.ClientAuthenticateRequest) (*pb.ValuedResponse, error) {
	reply := &pb.ValuedResponse{}
	reply.Type = pb.ValuedResponse_STOP_AND_RETURN
	reply.Value = &pb.ValuedResponse_BoolResult{BoolResult: true}
	return reply, nil
}

func (s *server) OnClientAuthorize(ctx context.Context, in *pb.ClientAuthorizeRequest) (*pb.ValuedResponse, error) {
	reply := &pb.ValuedResponse{}
	reply.Type = pb.ValuedResponse_STOP_AND_RETURN
	reply.Value = &pb.ValuedResponse_BoolResult{BoolResult: true}
	return reply, nil
}

func (s *server) OnClientSubscribe(ctx context.Context, in *pb.ClientSubscribeRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *server) OnClientUnsubscribe(ctx context.Context, in *pb.ClientUnsubscribeRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *server) OnSessionCreated(ctx context.Context, in *pb.SessionCreatedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}
func (s *server) OnSessionSubscribed(ctx context.Context, in *pb.SessionSubscribedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *server) OnSessionUnsubscribed(ctx context.Context, in *pb.SessionUnsubscribedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *server) OnSessionResumed(ctx context.Context, in *pb.SessionResumedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *server) OnSessionDiscarded(ctx context.Context, in *pb.SessionDiscardedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *server) OnSessionTakenover(ctx context.Context, in *pb.SessionTakenoverRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *server) OnSessionTerminated(ctx context.Context, in *pb.SessionTerminatedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *server) OnMessagePublish(ctx context.Context, in *pb.MessagePublishRequest) (*pb.ValuedResponse, error) {
	reply := &pb.ValuedResponse{}
	reply.Type = pb.ValuedResponse_STOP_AND_RETURN
	reply.Value = &pb.ValuedResponse_Message{Message: in.Message}

	datas := string(in.GetMessage().GetPayload())
	username := in.Message.Headers["username"]
	etoken := &netbase.DeviceAuth{}
	etoken.GetDeviceToken(username)

	data := &netbase.DeviceEventInfo{
		Datas:      datas,
		DeviceId:   etoken.DeviceId,
		DeviceAuth: etoken,
		Topic:      in.Message.Topic,
	}

	// todo 发送消息到消息队列
	log.Println(data)

	return reply, nil
}

func (s *server) OnMessageDelivered(ctx context.Context, in *pb.MessageDeliveredRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *server) OnMessageDropped(ctx context.Context, in *pb.MessageDroppedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *server) OnMessageAcked(ctx context.Context, in *pb.MessageAckedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func main() {
	client, err := kit.NewRedisClient("127.0.0.1", "123456", 6379, 0)
	if err != nil {
		log.Panic("Redis连接错误")
	} else {
		log.Println("Redis连接成功")
	}
	netbase.RedisDb = client

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHookProviderServer(s, &server{})
	log.Println("Started gRPC server on ::9000")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
