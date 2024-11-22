// //后端代码

// //注意 **我注释的代码，是不使用gin框架封装的Stream方法，也就是C.Stream(func())和C.ssevent(),只是C.Stream要改成for循环持续的从通道里面进行读，直到通道关闭，结束for循环**

package main

// import (
// 	"catface/app/service/nlp/glm"
// 	_ "catface/bootstrap"
// 	"fmt"
// 	"io"
// 	"testing"
// 	// "time"
//

// 	"github.com/gin-gonic/gin"

// )

// func SSE(c *gin.Context) {
// 	// 设置响应头，告诉前端适用event-stream事件流交互
// 	//c.Writer.Header().Set("Content-Type", "text/event-stream")
// 	//c.Writer.Header().Set("Cache-Control", "no-cache")
// 	//c.Writer.Header().Set("Connection", "keep-alive")

// 	// 判断是否支持sse
// 	//w := c.Writer
// 	//flusher, _ := w.(http.Flusher)
// 	query := c.Query("query")

// 	// 接收前端页面关闭连接通知
// 	closeNotify := c.Request.Context().Done()

// 	// 开启协程监听前端页面是否关闭了连接，关闭连接会触发此方法
// 	go func() {
// 		<-closeNotify
// 		fmt.Println("SSE关闭了")
// 		return
// 	}()

// 	//新建一个通道，用于数据接收和响应
// 	Chan := make(chan string)

// 	// 异步接收GPT响应，然后把响应的数据发送到通道Chan
// 	go func() {
// 		err := glm.ChatStream(query, Chan)
// 		if err != nil {
// 			fmt.Println("Error", err)
// 		}

// 		close(Chan)
// 	}()

// 	// gin框架封装的stream,会持续的调用这个func方法，记得返回true;返回false代表结束调用func方法
// 	c.Stream(func(w io.Writer) bool {
// 		select {
// 		case i, ok := <-Chan:
// 			if !ok {
// 				return false
// 			}
// 			c.SSEvent("chat", i) // c.SSEvent会自动修改响应头为事件流，并发送”test“事件流给前端监听”test“的回调方法
// 			//flusher.Flush() // 确保立即发送
// 			return true
// 		case <-closeNotify:
// 			fmt.Println("SSE关闭了")
// 			return false
// 		}
// 	})
// }

// func TestSSE(t *testing.T) {
// 	engine := gin.Default()
// 	// 设置跨域中间件
// 	engine.Use(func(context *gin.Context) {
// 		origin := context.GetHeader("Origin")
// 		// 允许 Origin 字段中的域发送请求
// 		context.Writer.Header().Add("Access-Control-Allow-Origin", origin) // 这边我的前端页面在63342，会涉及跨域，这个根据自己情况设置，或者直接设置为”*“，放行所有的
// 		// 设置预验请求有效期为 86400 秒
// 		context.Writer.Header().Set("Access-Control-Max-Age", "86400")
// 		// 设置允许请求的方法
// 		context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, U`PDATE, PATCH")
// 		// 设置允许请求的 Header
// 		context.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Apitoken")
// 		// 设置拿到除基本字段外的其他字段，如上面的Apitoken, 这里通过引用Access-Control-Expose-Headers，进行配置，效果是一样的。
// 		context.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Headers")
// 		// 配置是否可以带认证信息
// 		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		// OPTIONS请求返回200
// 		if context.Request.Method == "OPTIONS" {
// 			fmt.Println(context.Request.Header)
// 			context.AbortWithStatus(200)
// 		} else {
// 			context.Next()
// 		}
// 	})
// 	engine.GET("/admin/rag/default_talk", SSE) // TIP 记得适用get请求，我用post前端报404，资料说是SSE只支持get请求
// 	engine.Run(":20201")
// }
