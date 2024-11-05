package main

import (
	"context"

	"github.com/yankeguo/zhipu"
)

var prompt = "请根据以下长文本生成一个合适的标题，不需要书名号，长度10字内："

// var content = "那天散步时，我遇见了一只瘦弱的流浪猫。我给了它食物，并带它去看病。之后，我决定收养它，给它一个家。现在，它是我生活中不可或缺的伙伴。"
var content = "因为猪皮脚崴了，带去医院拍了一下片子，无大碍，静养就好，最近这段时间不要太动他，让他慢慢恢复恢复就好。"

func main() {
	// 或者手动指定密钥
	client, err := zhipu.NewClient(zhipu.WithAPIKey("0cf510ebc01599dba2a593069c1bdfbc.nQBQ4skP8xBh7ijU"))

	service := client.ChatCompletion("glm-4-flash").
		AddMessage(zhipu.ChatCompletionMessage{
			Role:    "user",
			Content: prompt + content,
		})
	// 	.SetStreamHandler(func(chunk zhipu.ChatCompletionResponse) error {
	// 	println(chunk.Choices[0].Delta.Content)
	// 	return nil
	// })

	res, err := service.Do(context.Background())

	if err != nil {
		zhipu.GetAPIErrorCode(err) // get the API error code
	} else {
		println(res.Choices[0].Message.Content)
	}
}
