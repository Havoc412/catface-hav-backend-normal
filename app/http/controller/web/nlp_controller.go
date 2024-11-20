package web

import (
	"catface/app/global/consts"
	"catface/app/global/errcode"
	"catface/app/global/variable"
	"catface/app/service/nlp"
	"catface/app/utils/llm_factory"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

type Nlp struct {
}

func (n *Nlp) Title(context *gin.Context) {
	content := context.GetString(consts.ValidatorPrefix + "content")

	tempGlmKey := variable.SnowFlake.GetIdAsString()
	client, ercode := variable.GlmClientHub.GetOneGlmClient(tempGlmKey, llm_factory.GlmModeSimple)
	if ercode > 0 {
		response.Fail(context, ercode, errcode.ErrMsg[ercode], errcode.ErrMsgForUser[ercode])
	}
	defer variable.GlmClientHub.UnavtiveOneGlmClient(tempGlmKey)
	defer variable.GlmClientHub.ReleaseOneGlmClient(tempGlmKey) // 临时使用，用完就释放。

	newTitle := nlp.GenerateTitle(content, client)
	if newTitle != "" {
		response.Success(context, consts.CurdStatusOkMsg, gin.H{"title": newTitle})
	} else {
		response.Fail(context, consts.CurdStatusOkCode, consts.CurdStatusOkMsg, "")
	}
}
