package nlp

import (
	"catface/app/global/variable"
	"catface/app/service/nlp/glm"
	"fmt"
	"strings"
)

func GenerateTitle(content string) string {
	message := variable.PromptsYml.GetString("Prompt.Title") + content
	title, _ := glm.Chat(message)
	return title
}

// ChatKnoledgeRAG 使用 RAG 模型进行知识问答
func ChatKnoledgeRAG(doc, query string, ch chan<- string) error {
	// 读取配置文件中的 KnoledgeRAG 模板
	promptTemplate := variable.PromptsYml.GetString("Prompt.KnoledgeRAG")

	// 替换模板中的占位符
	message := strings.Replace(promptTemplate, "{question}", query, -1)
	message = strings.Replace(message, "{context}", doc, -1)

	// 调用聊天接口
	// err := glm.ChatStream(message, ch)
	err := glm.BufferedChatStream(message, ch)
	if err != nil {
		return fmt.Errorf("调用聊天接口失败: %w", err)
	}

	return nil
}
