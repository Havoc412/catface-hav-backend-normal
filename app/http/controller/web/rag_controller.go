package web

import (
	"catface/app/global/errcode"
	"catface/app/global/variable"
	"catface/app/model_es"
	"catface/app/service/nlp"
	"catface/app/utils/response"
	"io"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Rag struct {
}

// v1 Http-POST 版本
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

func (r *Rag) HelpDetectCat(context *gin.Context) {
	// TODO
}
