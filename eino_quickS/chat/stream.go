package chat

import (
	"context"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// Stream 使用 ChatModel 生成流式回复
func Stream(ctx context.Context, chatModel model.ChatModel, messages []*schema.Message) (*schema.StreamReader[*schema.Message], error) {
	streamResult, err := chatModel.Stream(ctx, messages)
	if err != nil {
		return nil, err
	}
	return streamResult, nil
}

