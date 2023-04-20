package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Legacynnn/Chatbot-GPT/goMicroService/configs"
	"github.com/Legacynnn/Chatbot-GPT/goMicroService/internal/infra/grpc/server"
	"github.com/Legacynnn/Chatbot-GPT/goMicroService/internal/infra/repositories"
	"github.com/Legacynnn/Chatbot-GPT/goMicroService/internal/infra/web"
	"github.com/Legacynnn/Chatbot-GPT/goMicroService/internal/infra/web/webservice"
	"github.com/Legacynnn/Chatbot-GPT/goMicroService/internal/useCases/chat/completion"
	completionstream "github.com/Legacynnn/Chatbot-GPT/goMicroService/internal/useCases/chat/completionStream"
	"github.com/sashabaranov/go-openai"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	conn, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	repo := repositories.NewChatRepositoryMySQL(conn)
	client := openai.NewClient(configs.OpenAIApiKey)

	chatConfig := completion.ChatCompletionConfigInputDTO{
		Model:                configs.Model,
		ModelMaxTokens:       configs.ModelMaxTokens,
		Temperature:          float32(configs.Temperature),
		TopP:                 float32(configs.TopP),
		N:                    configs.N,
		Stop:                 configs.Stop,
		MaxTokens:            configs.MaxTokens,
		InitialSystemMessage: configs.InitialChatMessage,
	}

	chatConfigStream := completion.ChatCompletionConfigInputDTO{
		Model:                configs.Model,
		ModelMaxTokens:       configs.ModelMaxTokens,
		Temperature:          float32(configs.Temperature),
		TopP:                 float32(configs.TopP),
		N:                    configs.N,
		Stop:                 configs.Stop,
		MaxTokens:            configs.MaxTokens,
		InitialSystemMessage: configs.InitialChatMessage,
	}

	usecase := completion.NewChatCompletionUseCase(repo, client)

	streamChannel := make(chan completionstream.ChatCompletionOutputDTO)
	usecaseStream := completionstream.NewChatCompletionUseCase(repo, client, streamChannel)

	fmt.Println("Starting gRPC server on port " + configs.GRPCServerPort)
	grpcServer := server.NewGRPCServer(
		*usecaseStream,
		chatConfigStream,
		configs.GRPCServerPort,
		configs.AuthToken,
		streamChannel,
	)
	go grpcServer.Start()

	webserver := webservice.NewWebServer(":" + configs.WebServerPort)
	webserverChatHandler := web.NewWebChatGPTHandler(*usecase, chatConfig, configs.AuthToken)
	webserver.AddHandler("/chat", webserverChatHandler.Handle)

	fmt.Println("Server running on port " + configs.WebServerPort)
	webserver.Start()
}
