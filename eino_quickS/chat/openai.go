package chat

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

// NewOpenAI 创建 ChatModel
// 支持多种国内免费模型，通过环境变量配置：
//   - SILICONFLOW: export LLM_BASE_URL=https://api.siliconflow.cn/v1 LLM_API_KEY=你的key LLM_MODEL=Qwen/Qwen2.5-7B-Instruct
//   - DEEPSEEK:    export LLM_BASE_URL=https://api.deepseek.com LLM_API_KEY=你的key LLM_MODEL=deepseek-chat
//   - OPENAI:      export LLM_BASE_URL=https://api.openai.com/v1 LLM_API_KEY=你的key LLM_MODEL=gpt-4o-mini
func NewOpenAI(ctx context.Context) (model.ToolCallingChatModel, error) {
	// 从环境变量读取配置，默认使用硅基流动
	baseURL := os.Getenv("LLM_BASE_URL")
	apiKey := os.Getenv("LLM_API_KEY")
	modelName := os.Getenv("LLM_MODEL")

	// 默认值：硅基流动免费模型
	if baseURL == "" {
		baseURL = "https://api.siliconflow.cn/v1"
	}
	if modelName == "" {
		modelName = "Qwen/Qwen2.5-7B-Instruct" // 硅基流动免费模型
	}

	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: baseURL,
		Model:   modelName,
		APIKey:  apiKey,
	})
	if err != nil {
		return nil, err
	}
	return chatModel, nil
}
