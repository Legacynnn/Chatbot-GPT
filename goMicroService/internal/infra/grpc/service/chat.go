package service

import (
	"github.com/Legacynnn/Chatbot-GPT/goMicroService/internal/infra/grpc/pb"
	"github.com/Legacynnn/Chatbot-GPT/goMicroService/internal/useCases/chat/completion"
	completionstream "github.com/Legacynnn/Chatbot-GPT/goMicroService/internal/useCases/chat/completionStream"
)

type ChatService struct {
	pb.UnimplementedChatServiceServer
	ChatCompletionStreamUseCase completionstream.ChatCompletionUseCase
	ChatConfigStream            completionstream.ChatCompletionConfigInputDTO
	StreamChannel               chan completionstream.ChatCompletionOutputDTO
}

func NewChatService(chatCompletionStreamUseCase completionstream.ChatCompletionUseCase, chatConfigStream completion.ChatCompletionConfigInputDTO, streamChannel chan completionstream.ChatCompletionOutputDTO) *ChatService {
	return &ChatService{
		ChatCompletionStreamUseCase: chatCompletionStreamUseCase,
		ChatConfigStream:            completionstream.ChatCompletionConfigInputDTO(chatConfigStream),
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
	// If error in stream, put completionStream
	input := completionstream.ChatCompletionInputDTO{
		UserMessage: req.GetUserMessage(),
		UserID:      req.GetUserId(),
		ChatID:      req.GetChatId(),
		// If have error in stream, back in this
		Config: chatConfig,
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
