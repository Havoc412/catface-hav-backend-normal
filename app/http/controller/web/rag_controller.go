package web

import (
	"catface/app/global/errcode"
	"catface/app/global/variable"
	"catface/app/model_es"
	"catface/app/service/nlp"
	"catface/app/utils/response"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Rag struct {
}

// v1 Http-POST 版本; chat 需要不使用 ch 的版本。
// func (r *Rag) Chat(context *gin.Context) {
// 	// 1. query embedding
// 	query := context.GetString(consts.ValidatorPrefix + "query")
// 	embedding, ok := nlp.GetEmbedding(query)
// 	if !ok {
// 		code := errcode.ErrPythonService
// 		response.Fail(context, code, errcode.ErrMsg[code], "")
// 		return
// 	}

// 	// 2. ES TopK
// 	docs, err := model_es.CreateDocESFactory().TopK(embedding, 1)
// 	if err != nil || len(docs) == 0 {
// 		variable.ZapLog.Error("ES TopK error", zap.Error(err))

// 		code := errcode.ErrNoDocFound
// 		response.Fail(context, code, errcode.ErrMsg[code], errcode.ErrMsgForUser[code])
// 	}

// 	// 3. LLM answer
// 	if answer, err := nlp.ChatKnoledgeRAG(docs[0].Content, query); err == nil {
// 		response.Success(context, consts.CurdStatusOkMsg, gin.H{
// 			"answer": answer,
// 		})
// 	} else {
// 		response.Fail(context, consts.CurdStatusOkCode, consts.CurdStatusOkMsg, "")
// 	}
// }

func (r *Rag) ChatSSE(context *gin.Context) {
	query := context.Query("query")

	// 1. query embedding
	embedding, ok := nlp.GetEmbedding(query)
	if !ok {
		code := errcode.ErrPythonService
		response.Fail(context, code, errcode.ErrMsg[code], "")
		return
	}

	// 2. ES TopK
	docs, err := model_es.CreateDocESFactory().TopK(embedding, 1)
	if err != nil || len(docs) == 0 {
		variable.ZapLog.Error("ES TopK error", zap.Error(err))

		code := errcode.ErrNoDocFound
		response.Fail(context, code, errcode.ErrMsg[code], errcode.ErrMsgForUser[code])
	}

	// UPDATE
	closeEventFromVue := context.Request.Context().Done()
	ch := make(chan string) // TIP 建立通道。

	// 3. LLM answer
	go func() {
		err := nlp.ChatKnoledgeRAG(docs[0].Content, query, ch)
		if err != nil {
			variable.ZapLog.Error("ChatKnoledgeRAG error", zap.Error(err))
		}
		close(ch)
	}()

	context.Stream(func(w io.Writer) bool {
		select {
		case c, ok := <-ch:
			if !ok {
				return false
			}
			context.SSEvent("chat", c)
			return true
		case <-closeEventFromVue:
			return false
		}
	})
}

var upgrader = websocket.Upgrader{ // TEST 测试，先写一个裸的 wss
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // info 在生产环境中可能需要更安全的检查
	},
}

func (r *Rag) ChatWebSocket(context *gin.Context) {
	query := context.Query("query")

	// 0. 协议升级
	ws, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		variable.ZapLog.Error("OnOpen error", zap.Error(err))
		response.Fail(context, errcode.ErrWebsocketUpgradeFail, errcode.ErrMsg[errcode.ErrWebsocketUpgradeFail], "")
		return
	}
	defer ws.Close()

	// 1. query embedding
	embedding, ok := nlp.GetEmbedding(query)
	if !ok {
		code := errcode.ErrServerDown
		err := ws.WriteMessage(websocket.TextMessage, []byte(errcode.ErrMsgForUser[code]))
		if err != nil {
			variable.ZapLog.Error("Failed to send error message via WebSocket", zap.Error(err))
		}
		return
	}

	// 2. ES TopK
	docs, err := model_es.CreateDocESFactory().TopK(embedding, 1)
	if err != nil || len(docs) == 0 {
		variable.ZapLog.Error("ES TopK error", zap.Error(err))

		code := errcode.ErrNoDocFound
		err := ws.WriteMessage(websocket.TextMessage, []byte(errcode.ErrMsgForUser[code]))
		if err != nil {
			variable.ZapLog.Error("Failed to send error message via WebSocket", zap.Error(err))
		}
		return
	}

	// 3.
	closeEventFromVue := context.Request.Context().Done() // 接收前端传来的中断信号。
	ch := make(chan string)                               // TIP 建立通道。

	go func() {
		err := nlp.ChatKnoledgeRAG(docs[0].Content, query, ch)
		if err != nil {
			variable.ZapLog.Error("ChatKnoledgeRAG error", zap.Error(err))
		}
		close(ch) // 这里 close，使得下方 for 结束。
	}()

	for {
		select {
		case c, ok := <-ch:
			if !ok {
				return
			}
			// variable.ZapLog.Info("ChatKnoledgeRAG", zap.String("c", c))
			err := ws.WriteMessage(websocket.TextMessage, []byte(c))
			if err != nil {
				return
			}
		case <-closeEventFromVue:
			return
		}
	}
}

func (r *Rag) HelpDetectCat(context *gin.Context) {
	// TODO
}
