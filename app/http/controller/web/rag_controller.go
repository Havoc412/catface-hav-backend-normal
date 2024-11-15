package web

import (
	"catface/app/global/consts"
	"catface/app/global/errcode"
	"catface/app/global/variable"
	"catface/app/model_es"
	"catface/app/service/nlp"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Rag struct {
}

func (r *Rag) Chat(context *gin.Context) {
	// 1. query embedding
	query := context.GetString(consts.ValidatorPrefix + "query")
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

	// 3. LLM answer
	if answer, err := nlp.ChatKnoledgeRAG(docs[0].Content, query); err == nil {
		response.Success(context, consts.CurdStatusOkMsg, gin.H{
			"answer": answer,
		})
	} else {
		response.Fail(context, consts.CurdStatusOkCode, consts.CurdStatusOkMsg, "")
	}
}

func (r *Rag) HelpDetectCat(context *gin.Context) {
	// TODO
}
