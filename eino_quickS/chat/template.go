package chat

import (
	"context"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

// 把逻辑放在函数里
func CreateTemplate() ([]*schema.Message, error) {
	template := prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个{role},你要用{style}的语气回答问题,你的目标是帮助人解决心里问题"),

		// 插入需要的对话历史
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("问题:{question}:"),
	)

	// 用模板生成消息
	messages, err := template.Format(context.Background(), map[string]any{
		"role":     "程序员鼓励师",
		"style":    "积极、温暖且专业",
		"question": "我的代码一直报错，感觉好沮丧，该怎么办？",
		"chat_history": []*schema.Message{
			schema.UserMessage("你好"),
			schema.AssistantMessage("嘿！我是你的程序员鼓励师！...", nil),
			schema.UserMessage("我觉得自己写的代码太烂了"),
			schema.AssistantMessage("每个程序员都经历过这个阶段！...", nil),
		},
	})

	return messages, err
}
