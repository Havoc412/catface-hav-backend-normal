package glm

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/yankeguo/zhipu"
)

// ChatWithGLM 封装了与GLM模型进行对话的逻辑
func Chat(message string, client *zhipu.ChatCompletionService) (string, error) {
	service := client.AddMessage(zhipu.ChatCompletionMessage{
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
func ChatStream(message string, ch chan<- string, client *zhipu.ChatCompletionService) error {
	service := client.AddMessage(zhipu.ChatCompletionMessage{Role: "user", Content: message}).
		SetStreamHandler(func(chunk zhipu.ChatCompletionResponse) error {
			content := chunk.Choices[0].Delta.Content
			if content != "" {
				ch <- content // 将内容发送到通道
			}
			return nil
		})

	// Test
	//messages := client.GetMessages()
	//for id, message := range messages {
	//	variable.ZapLog.Info(fmt.Sprintf("message-%d", id+1), zap.String("message", message.(zhipu.ChatCompletionMessage).Role), zap.String("content", message.(zhipu.ChatCompletionMessage).Content))
	//}

	// 执行服务调用
	res, err := service.Do(context.Background())
	if err != nil {
		return err
	}
	// 增加 AI 回答的消息记录。
	client.AddMessage(zhipu.ChatCompletionMessage{Role: "assistant", Content: res.Choices[0].Message.Content})

	return nil
}

// 带缓冲机制的 ChatStream；计数 & 计时 双判定。
func BufferedChatStream(message string, ch chan<- string, client *zhipu.ChatCompletionService) error {
	bufferedCh := make(chan string)                // 带缓冲的通道，缓冲大小为10
	timer := time.NewTimer(500 * time.Millisecond) // 定时器，500毫秒

	go func() {
		err := ChatStream(message, bufferedCh, client)
		if err != nil {
			return
		}
		close(bufferedCh)
	}()

	var buffer strings.Builder
	for {
		select {
		case c, ok := <-bufferedCh:
			if !ok {
				if buffer.Len() > 0 {
					ch <- buffer.String()
				}
				return nil // 依靠这里停止函数。
			}
			buffer.WriteString(c)
			if buffer.Len() >= 10 {
				ch <- buffer.String()
				buffer.Reset()
				timer.Reset(500 * time.Millisecond)
			}
		case <-timer.C:
			if buffer.Len() > 0 {
				ch <- buffer.String()
				buffer.Reset()
			}
			timer.Reset(500 * time.Millisecond)
		}
	}
}
