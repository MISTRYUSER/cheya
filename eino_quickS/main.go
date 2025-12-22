package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/xuewentao/cheya/eino_quickS/chat"
)

func main() {
	ctx := context.Background()

	// 创建 ChatModel
	fmt.Println("🚀 正在初始化 AI 助手...")
	chatModel, err := chat.NewOpenAI(ctx)
	if err != nil {
		fmt.Printf("❌ 初始化失败: %v\n", err)
		return
	}
	fmt.Println("✅ AI 助手已就绪！")
	fmt.Println()
	fmt.Println("💬 程序员鼓励师 - 交互式对话")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("输入你的问题，按 Enter 发送")
	fmt.Println("输入 'exit' 或 'quit' 退出")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 对话历史
	chatHistory := []*schema.Message{}

	// 系统提示
	systemPrompt := schema.SystemMessage("你是一个程序员鼓励师。你需要用积极、温暖且专业的语气回答问题。你的目标是帮助程序员保持积极乐观的心态，提供技术建议的同时也要关注他们的心理健康。")

	reader := bufio.NewReader(os.Stdin)

	for {
		// 显示输入提示
		fmt.Print("🧑 You: ")

		// 读取用户输入
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("❌ 读取输入失败: %v\n", err)
			continue
		}

		// 去除换行符和空格
		input = strings.TrimSpace(input)

		// 检查退出命令
		if input == "exit" || input == "quit" || input == "q" {
			fmt.Println("\n👋 再见！记住，你是最棒的程序员！")
			break
		}

		// 跳过空输入
		if input == "" {
			continue
		}

		// 构建消息列表
		messages := []*schema.Message{systemPrompt}
		messages = append(messages, chatHistory...)
		messages = append(messages, schema.UserMessage(input))

		// 显示 AI 回复提示
		fmt.Print("🤖 AI: ")

		// 使用流式输出
		streamResult, err := chatModel.Stream(ctx, messages)
		if err != nil {
			fmt.Printf("\n❌ 请求失败: %v\n", err)
			continue
		}

		// 收集完整回复
		fullResponse := streamResponse(streamResult)
		fmt.Println() // 换行

		// 更新对话历史
		chatHistory = append(chatHistory, schema.UserMessage(input))
		chatHistory = append(chatHistory, schema.AssistantMessage(fullResponse, nil))

		// 限制历史长度（保留最近 10 轮对话）
		if len(chatHistory) > 20 {
			chatHistory = chatHistory[len(chatHistory)-20:]
		}

		fmt.Println()
	}
}

// streamResponse 处理流式输出并返回完整内容
func streamResponse(sr *schema.StreamReader[*schema.Message]) string {
	defer sr.Close()

	var fullContent strings.Builder

	for {
		message, err := sr.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("\n❌ 接收失败: %v", err)
			break
		}
		// 打印并收集内容
		fmt.Print(message.Content)
		fullContent.WriteString(message.Content)
	}

	return fullContent.String()
}

// 保留原有函数供参考
func createOpenAIChatModel(ctx context.Context) model.ChatModel {
	chatModel, err := chat.NewOpenAI(ctx)
	if err != nil {
		fmt.Printf("Failed to create OpenAI ChatModel: %v\n", err)
		os.Exit(1)
	}
	return chatModel
}
