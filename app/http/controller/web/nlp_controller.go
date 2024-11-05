package web

import (
	"catface/app/global/consts"
	"catface/app/utils/nlp"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

type Nlp struct {
}

func (n *Nlp) Title(context *gin.Context) {
	content := context.GetString(consts.ValidatorPrefix + "content")

	newTitle := nlp.GenerateTitle(content)
	if newTitle != "" {
		response.Success(context, consts.CurdStatusOkMsg, gin.H{"title": newTitle})
	} else {
		response.Fail(context, consts.CurdStatusOkCode, consts.CurdStatusOkMsg, "")
	}
}
