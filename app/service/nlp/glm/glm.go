package glm

import (
	"catface/app/global/variable"
	"context"
	"errors"

	"github.com/yankeguo/zhipu"

)

// ChatWithGLM 封装了与GLM模型进行对话的逻辑
func Chat(message string) (string, error) {
	service := variable.GlmClient.ChatCompletion("glm-4-flash").
		AddMessage(zhipu.ChatCompletionMessage{
			Role:    "user",
			Content: message,
		})

	res, err := service.Do(context.Background())
	if err != nil {
		apiErrorCode := zhipu.GetAPIErrorCode(err)
		return "", errors.New(apiErrorCode) // 将字符串包装成 error 类型
	}

	return res.Choices[0].Message.Content, nil
}

// ChatStream 接收一个消息和一个通道，将流式响应发送到通道中
func ChatStream(message string, ch chan<- string) error {
	service := variable.GlmClient.ChatCompletion("glm-4-flash").
		AddMessage(zhipu.ChatCompletionMessage{Role: "user", Content: message}).
		SetStreamHandler(func(chunk zhipu.ChatCompletionResponse) error {
			content := chunk.Choices[0].Delta.Content
			if content != "" {
				ch <- content // 将内容发送到通道
			}
			return nil
		})

	// 执行服务调用
	_, err := service.Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}
