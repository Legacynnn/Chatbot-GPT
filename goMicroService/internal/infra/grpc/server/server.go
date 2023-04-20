package server

import (
	"net"

	"github.com/Legacynnn/Chatbot-GPT/goMicroService/internal/infra/grpc/pb"
	"github.com/Legacynnn/Chatbot-GPT/goMicroService/internal/infra/grpc/service"
	"github.com/Legacynnn/Chatbot-GPT/goMicroService/internal/useCases/chat/completion"
	completionstream "github.com/Legacynnn/Chatbot-GPT/goMicroService/internal/useCases/chat/completionStream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	ChatCompletionStreamUseCase completionstream.ChatCompletionUseCase
	ChatConfigStream            completionstream.ChatCompletionConfigInputDTO
	ChatService                 service.ChatService
	Port                        string
	AuthToken                   string
	StreamChannel               chan completionstream.ChatCompletionOutputDTO
}

func NewGRPCServer(chatCompletionStreamUseCase completionstream.ChatCompletionUseCase, chatConfigStream completion.ChatCompletionConfigInputDTO, port, authToken string, streamChannel chan completionstream.ChatCompletionOutputDTO) *GRPCServer {
	chatService := service.NewChatService(chatCompletionStreamUseCase, chatConfigStream, streamChannel)
	return &GRPCServer{
		ChatCompletionStreamUseCase: chatCompletionStreamUseCase,
		ChatConfigStream:            completionstream.ChatCompletionConfigInputDTO(chatConfigStream),
		Port:                        port,
		AuthToken:                   authToken,
		StreamChannel:               streamChannel,
		ChatService:                 *chatService,
	}
}

func (g *GRPCServer) AuthInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	token := md.Get("authorization")
	if len(token) == 0 {
		return status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	if token[0] != g.AuthToken {
		return status.Error(codes.Unauthenticated, "authorization token is invalid")
	}

	return handler(srv, ss)
}

func (g *GRPCServer) Start() error {
	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(g.AuthInterceptor),
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterChatServiceServer(grpcServer, &g.ChatService)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+g.Port)
	if err != nil {
		panic(err.Error())
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err.Error())
	}

	return nil
}
