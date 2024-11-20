package web

import (
	"catface/app/global/consts"
	"catface/app/global/errcode"
	"catface/app/global/variable"
	"catface/app/model"
	"catface/app/model_es"
	"catface/app/service/nlp"
	"catface/app/service/rag/curd"
	"catface/app/utils/llm_factory"
	"catface/app/utils/micro_service"
	"catface/app/utils/response"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Rag struct {
}

func (r *Rag) Release(context *gin.Context) {
	token := context.GetString(consts.ValidatorPrefix + "token")
	if ok := variable.GlmClientHub.ReleaseOneGlmClient(token); ok {
		variable.ZapLog.Info("释放一个 GLM Client",
			zap.String("token", token),
			zap.String("当前空闲连接数", strconv.Itoa(variable.GlmClientHub.Idle)))
	} else {
		variable.ZapLog.Warn("尝试释放一个 GLM Client，但是 token 无效",
			zap.String("当前空闲连接数", strconv.Itoa(variable.GlmClientHub.Idle)))
	}

	response.Success(context, consts.CurdStatusOkMsg, "")
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
	token := context.Query("token")

	// 0-1. 测试 python
	if !micro_service.TestLinkPythonService() {
		code := errcode.ErrPythonService
		response.Fail(context, code, errcode.ErrMsg[code], "")
		return
	}

	// 0-2. 获取一个 GLM Client
	if token == "" {
		token = variable.SnowFlake.GetIdAsString()
	}
	client, ercode := variable.GlmClientHub.GetOneGlmClient(token, llm_factory.GlmModeKnowledgeHub)
	if ercode != 0 {
		response.Fail(context, ercode, errcode.ErrMsg[ercode], errcode.ErrMsgForUser[ercode])
		return
	}
	defer variable.GlmClientHub.UnavtiveOneGlmClient(token) // INFO ws 结束时，取消 Avtive 的占用。

	// 1. query embedding
	embedding, ok := nlp.GetEmbedding([]string{query})
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
		err := nlp.ChatKnoledgeRAG(docs[0].Content, query, ch, client)
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
	token := context.Query("token")

	if token == "" {
		token = variable.SnowFlake.GetIdAsString()
	}

	// 0-1. 协议升级
	ws, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		variable.ZapLog.Error("OnOpen error", zap.Error(err))
		response.Fail(context, errcode.ErrWebsocketUpgradeFail, errcode.ErrMsg[errcode.ErrWebsocketUpgradeFail], "")
		return
	}
	defer ws.Close()

	// 0-2. 测试 Python 微服务是否启动
	if !micro_service.TestLinkPythonService() {
		code := errcode.ErrPythonServierDown
		err := ws.WriteMessage(websocket.TextMessage, model.CreateNlpWebSocketResult("", errcode.ErrMsgForUser[code]).JsonMarshal())
		if err != nil {
			variable.ZapLog.Error("Failed to send error message via WebSocket", zap.Error(err))
		}
		return
	}

	// 0-3. 从 GLM_HUB 中获取一个可用的 glm client;
	clientInfo, ercode := variable.GlmClientHub.GetOneGlmClientInfo(token, llm_factory.GlmModeKnowledgeHub)
	if ercode != 0 {
		variable.ZapLog.Error("GetOneGlmClient error", zap.Error(err))
		err := ws.WriteMessage(websocket.TextMessage, model.CreateNlpWebSocketResult("", errcode.ErrMsgForUser[ercode]).JsonMarshal())
		if err != nil {
			variable.ZapLog.Error("Failed to send error message via WebSocket", zap.Error(err))
		}
		return
	}
	defer variable.GlmClientHub.UnavtiveOneGlmClient(token) // INFO ws 结束时，取消 Avtive 的占用。

	// 1. query embedding
	clientInfo.AddQuery(query)
	embedding, ok := nlp.GetEmbedding(clientInfo.UserQuerys)
	if !ok {
		code := errcode.ErrPythonServierDown
		err := ws.WriteMessage(websocket.TextMessage, model.CreateNlpWebSocketResult("", errcode.ErrMsgForUser[code]).JsonMarshal())
		if err != nil {
			variable.ZapLog.Error("Failed to send error message via WebSocket", zap.Error(err))
		}
		return
	}

	// 2. ES TopK // TODO 这里需要特化选取不同知识库的文档；目前是依靠显式的路由。
	docs, err := curd.CreateDocCurdFactory().TopK(embedding, 1)
	if err != nil || len(docs) == 0 {
		variable.ZapLog.Error("ES TopK error", zap.Error(err))

		code := errcode.ErrNoDocFound
		err := ws.WriteMessage(websocket.TextMessage, model.CreateNlpWebSocketResult("", errcode.ErrMsgForUser[code]).JsonMarshal())
		if err != nil {
			variable.ZapLog.Error("Failed to send error message via WebSocket", zap.Error(err))
		}
		return
	}

	// STAGE websocket 的 defer 关闭函数，但是需要 ES 拿到的 doc—id
	defer func() { // UPDATE 临时"持久化"方案，之后考虑结合 jwt 维护的 token 处理。
		// 0. 传递参考资料的信息
		docMsg := model.CreateNlpWebSocketResult(docs[0].Type, docs)
		err := ws.WriteMessage(websocket.TextMessage, docMsg.JsonMarshal())
		if err != nil {
			variable.ZapLog.Error("Failed to send doc message via WebSocket", zap.Error(err))
		}

		// 1. 传递 token 信息； // UPDATE 临时方案
		tokenMsg := model.CreateNlpWebSocketResult("token", token)
		err = ws.WriteMessage(websocket.TextMessage, tokenMsg.JsonMarshal())
		if err != nil {
			variable.ZapLog.Error("Failed to send token message via WebSocket", zap.Error(err))
		}
		// ws.Close()  // 在上面调用了 defer；// TIP defer 的“栈”性质。
	}()

	// 3.
	closeEventFromVue := context.Request.Context().Done() // 接收前端传来的中断信号。
	ch := make(chan string)                               // TIP 建立通道。

	go func() {
		err := nlp.ChatKnoledgeRAG(docs[0].Content, query, ch, clientInfo.Client)
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
			err := ws.WriteMessage(websocket.TextMessage, model.CreateNlpWebSocketResult("", c).JsonMarshal())
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
