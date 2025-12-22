package server

import (
	"io"

	"github.com/xuewentao/cheya/apps/ai-agent/llm"
	aiv1 "github.com/xuewentao/cheya/api/ai/v1"
	_ "google.golang.org/grpc"
)
type AIServer struct {
	aiv1.UnimplementedAIServiceServer
	LLMClient *llm.LLMClient
}
func (s *AIServer) Chat(
    req *aiv1.ChatRequest,
    stream aiv1.AIService_ChatServer,
) error {
	//0. create context
	ctx := stream.Context()

	//1.create messages(System + history + current message)
	messages := []llm.Message {
		{Role: "system",Content: "你是 CheYa 车联网平台的智能助手"},
	}
	//add history message
	for _, h := range req.History {
		messages = append(messages, llm.Message{
			Role: 		h.Role,
			Content:	h.Content,
		})
	}

	//add user message
	messages = append(messages,llm.Message{
		Role: 		"user",
		Content: 	req.UserMessage,
	})

	//2. Call LLM Stream
	llmStream,err := s.LLMClient.Stream(ctx,messages)
	if err != nil {
		return err
	}
	defer llmStream.Close()

	//forward the response
	for {
		chunk, err := llmStream.Recv()
		if err == io.EOF{
			break//if recv done break
		}
		if err != nil {
			return err
		}
		//send to client
		stream.Send(&aiv1.ChatResponse{
			AiResponse: chunk.Content,
			Done: 		false,
		})
	}
	//4.send finish flagb
	stream.Send(&aiv1.ChatResponse{Done: true})
	return nil
}