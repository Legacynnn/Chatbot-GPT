package service

import (
	"github.com/Legacynnn/Chatbot-GPT/goMicroService/internal/infra/grpc/pb"
	"github.com/Legacynnn/Chatbot-GPT/goMicroService/internal/useCases/chat/completion"
	completionstream "github.com/Legacynnn/Chatbot-GPT/goMicroService/internal/useCases/chat/completionStream"
)

type ChatService struct {
	pb.UnimplementedChatServiceServer
	ChatCompletionStreamUseCase completion.ChatCompletionUseCase
	ChatConfigStream            completion.ChatCompletionConfigInputDTO
	StreamChannel               chan completion.ChatCompletionOutputDTO
}

func NewChatService(chatCompletionStreamUseCase completion.ChatCompletionUseCase, chatConfigStream completion.ChatCompletionConfigInputDTO, streamChannel chan completion.ChatCompletionOutputDTO) *ChatService {
	return &ChatService{
		ChatCompletionStreamUseCase: chatCompletionStreamUseCase,
		ChatConfigStream:            chatConfigStream,
		StreamChannel:               streamChannel,
	}
}

func (c *ChatService) ChatStream(req *pb.ChatRequest, stream pb.ChatService_ChatStreamServer) error {
	chatConfig := completionstream.ChatCompletionConfigInputDTO{
		Model:                c.ChatConfigStream.Model,
		ModelMaxTokens:       c.ChatConfigStream.ModelMaxTokens,
		Temperature:          c.ChatConfigStream.Temperature,
		TopP:                 c.ChatConfigStream.TopP,
		N:                    c.ChatConfigStream.N,
		Stop:                 c.ChatConfigStream.Stop,
		MaxTokens:            c.ChatConfigStream.MaxTokens,
		InitialSystemMessage: c.ChatConfigStream.InitialSystemMessage,
	}
	input := completionstream.ChatCompletionInputDTO{
		UserMessage: req.GetUserMessage(),
		UserID:      req.GetUserId(),
		ChatID:      req.GetChatId(),
		Config:      chatConfig,
	}

	ctx := stream.Context()
	go func() {
		for msg := range c.StreamChannel {
			stream.Send(&pb.ChatResponse{
				ChatId:  msg.ChatID,
				UserId:  msg.UserID,
				Content: msg.Content,
			})
		}
	}()

	_, err := c.ChatCompletionStreamUseCase.Execute(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
