package nlp

import (
	"catface/app/global/variable"
	"catface/app/service/nlp/glm"
	"fmt"
	"strings"

	"github.com/yankeguo/zhipu"
)

func GenerateTitle(content string, client *zhipu.ChatCompletionService) string {
	message := variable.PromptsYml.GetString("Prompt.Title") + content
	title, _ := glm.Chat(message, client)
	return title
}

// ChatKnoledgeRAG 使用 RAG 模型进行知识问答
func ChatRAG(doc, query, mode string, ch chan<- string, client *zhipu.ChatCompletionService) error {
	// 读取配置文件中的 KnoledgeRAG 模板
	promptTemplate := variable.PromptsYml.GetString("Prompt.RAG." + mode)

	// 替换模板中的占位符
	message := strings.Replace(promptTemplate, "{question}", query, -1)
	message = strings.Replace(message, "{context}", doc, -1)

	// 调用聊天接口
	// err := glm.ChatStream(message, ch)
	err := glm.BufferedChatStream(message, ch, client)
	if err != nil {
		return fmt.Errorf("调用聊天接口失败: %w", err)
	}

	return nil
}
