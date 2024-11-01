package main

import (
	"context"

	"github.com/yankeguo/zhipu"

)

func main() {
	// 或者手动指定密钥
	client, err := zhipu.NewClient(zhipu.WithAPIKey("0cf510ebc01599dba2a593069c1bdfbc.nQBQ4skP8xBh7ijU"))

	service := client.ChatCompletion("glm-4-flash").
		AddMessage(zhipu.ChatCompletionMessage{
			Role:    "user",
			Content: "你好",
		}).SetStreamHandler(func(chunk zhipu.ChatCompletionResponse) error {
			println(chunk.Choices[0].Delta.Content)
			return nil
		})

	res, err := service.Do(context.Background())

	if err != nil {
		zhipu.GetAPIErrorCode(err) // get the API error code
	} else {
		println(res.Choices[0].Message.Content)
	}
}
