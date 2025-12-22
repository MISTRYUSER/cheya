package chat

import (
	"context"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// Generate 使用 ChatModel 生成完整回复
func Generate(ctx context.Context, chatModel model.ChatModel, messages []*schema.Message) (*schema.Message, error) {
	result, err := chatModel.Generate(ctx, messages)
	if err != nil {
		return nil, err
	}
	return result, nil
}

