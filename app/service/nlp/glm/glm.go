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
