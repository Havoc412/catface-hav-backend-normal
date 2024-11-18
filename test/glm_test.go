package test

import (
	"catface/app/global/variable"
	_ "catface/bootstrap"
	"context"
	"testing"

	"github.com/yankeguo/zhipu"
)

func TestGlmMessageStore(t *testing.T) {
	glmClient, err := zhipu.NewClient(zhipu.WithAPIKey(variable.ConfigYml.GetString("Glm.ApiKey")))
	if err != nil {
		t.Fatal(err)
	}

	service := glmClient.ChatCompletion("glm-4-flash").AddMessage(zhipu.ChatCompletionMessage{
		Role:    "user",
		Content: "请你记一下我说的数字：2",
	})

	res, err := service.Do(context.Background())
	if err != nil {
		apiErrorCode := zhipu.GetAPIErrorCode(err)
		t.Fatal(apiErrorCode)
	}
	t.Log(res.Choices[0].Message.Content)

	messages := service.GetMessages()
	for _, message := range messages {
		t.Log(message.(zhipu.ChatCompletionMessage).Role, message.(zhipu.ChatCompletionMessage).Content)
	}

	service.AddMessage(zhipu.ChatCompletionMessage{
		Role:    "user",
		Content: "现在请你复述我刚才说的数字。",
	})
	res, err = service.Do(context.Background())
	if err != nil {
		apiErrorCode := zhipu.GetAPIErrorCode(err)
		t.Fatal(apiErrorCode)
	}

	messages = service.GetMessages()
	for _, message := range messages {
		t.Log(message.(zhipu.ChatCompletionMessage).Role, message.(zhipu.ChatCompletionMessage).Content)
	}

	t.Log(res.Choices[0].Message.Content)
}
