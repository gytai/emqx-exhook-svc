package netbase

import (
	"context"
	"log"

	pb "emqx.io/grpc/exhook/protobuf"
)

// Server 是用于实现 emqx_exhook_v3 的服务器
type Server struct {
	pb.UnimplementedHookProviderServer
}

// HookProviderServer callbacks

func (s *Server) OnProviderLoaded(ctx context.Context, in *pb.ProviderLoadedRequest) (*pb.LoadedResponse, error) {
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

func (s *Server) OnProviderUnloaded(ctx context.Context, in *pb.ProviderUnloadedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *Server) OnClientConnect(ctx context.Context, in *pb.ClientConnectRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *Server) OnClientConnack(ctx context.Context, in *pb.ClientConnackRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *Server) OnClientConnected(ctx context.Context, in *pb.ClientConnectedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *Server) OnClientDisconnected(ctx context.Context, in *pb.ClientDisconnectedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *Server) OnClientAuthenticate(ctx context.Context, in *pb.ClientAuthenticateRequest) (*pb.ValuedResponse, error) {
	reply := &pb.ValuedResponse{}
	reply.Type = pb.ValuedResponse_STOP_AND_RETURN
	reply.Value = &pb.ValuedResponse_BoolResult{BoolResult: true}
	return reply, nil
}

func (s *Server) OnClientAuthorize(ctx context.Context, in *pb.ClientAuthorizeRequest) (*pb.ValuedResponse, error) {
	reply := &pb.ValuedResponse{}
	reply.Type = pb.ValuedResponse_STOP_AND_RETURN
	reply.Value = &pb.ValuedResponse_BoolResult{BoolResult: true}
	return reply, nil
}

func (s *Server) OnClientSubscribe(ctx context.Context, in *pb.ClientSubscribeRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *Server) OnClientUnsubscribe(ctx context.Context, in *pb.ClientUnsubscribeRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *Server) OnSessionCreated(ctx context.Context, in *pb.SessionCreatedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}
func (s *Server) OnSessionSubscribed(ctx context.Context, in *pb.SessionSubscribedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *Server) OnSessionUnsubscribed(ctx context.Context, in *pb.SessionUnsubscribedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *Server) OnSessionResumed(ctx context.Context, in *pb.SessionResumedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *Server) OnSessionDiscarded(ctx context.Context, in *pb.SessionDiscardedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *Server) OnSessionTakenover(ctx context.Context, in *pb.SessionTakenoverRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *Server) OnSessionTerminated(ctx context.Context, in *pb.SessionTerminatedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *Server) OnMessagePublish(ctx context.Context, in *pb.MessagePublishRequest) (*pb.ValuedResponse, error) {
	reply := &pb.ValuedResponse{}
	reply.Type = pb.ValuedResponse_STOP_AND_RETURN
	reply.Value = &pb.ValuedResponse_Message{Message: in.Message}

	datas := string(in.GetMessage().GetPayload())
	username := in.Message.Headers["username"]
	etoken := &DeviceAuth{}
	etoken.GetDeviceToken(username)

	data := &DeviceEventInfo{
		Datas:      datas,
		DeviceId:   etoken.DeviceId,
		DeviceAuth: etoken,
		Topic:      in.Message.Topic,
	}

	// 将消息发布到RabbitMQ交换机，支持多个消费者订阅
	err := publishToQueue("", data) // queueName参数已不使用，传空字符串
	if err != nil {
		log.Printf("发布消息到RabbitMQ失败: %v, topic: %s", err, in.Message.Topic)
	} else {
		log.Printf("成功发布消息到RabbitMQ交换机, topic: %s", in.Message.Topic)
	}

	return reply, nil
}

func (s *Server) OnMessageDelivered(ctx context.Context, in *pb.MessageDeliveredRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *Server) OnMessageDropped(ctx context.Context, in *pb.MessageDroppedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *Server) OnMessageAcked(ctx context.Context, in *pb.MessageAckedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

// NewHookServer 创建并返回一个新的Server实例
func NewHookServer() *Server {
	return &Server{}
}
