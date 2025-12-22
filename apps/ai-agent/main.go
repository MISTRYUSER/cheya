package main

import (
    "log"
    "net"
    "google.golang.org/grpc"
    aiv1 "github.com/xuewentao/cheya/api/ai/v1"
 	"google.golang.org/grpc/reflection"
	"github.com/xuewentao/cheya/apps/ai-agent/server"
	"github.com/xuewentao/cheya/apps/ai-agent/llm"
)	

func main() {
	//1.listen port 50055
	lis, err := net.Listen("tcp",":50055")
	if err != nil {
		log.Fatalf("Failure listen:%v",err)
	}

	//create LLM client
	llmClient := llm.NewLLMClient (
		"https://api.siliconflow.cn/v1",
        "sk-lixunuekoigiqayuoimeggpzeqlrijcrvhgqmjgarwcqsehy",  // 稍后用环境变量
        "Qwen/Qwen2.5-7B-Instruct",
	)

	//3.create gRPC container
	grpcServer := grpc.NewServer()

	//3.register AIService
	aiv1.RegisterAIServiceServer(grpcServer,&server.AIServer{
		LLMClient: llmClient,
	})
	//you should register reflection 
	reflection.Register(grpcServer)
	//4.bringup
	log.Println("AI Agent server startup: 50055")
	if err := grpcServer.Serve(lis);err != nil {
		log.Fatalf("启动失败: %v", err)
    }

}