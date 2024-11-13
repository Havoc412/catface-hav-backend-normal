package web

import (
	"catface/app/global/consts"
	"catface/app/global/errcode"
	"catface/app/model_es"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

type Knowledge struct {
}

func (k *Knowledge) RandomList(context *gin.Context) {
	num := context.GetFloat64(consts.ValidatorPrefix + "num")

	knowledgeList, err := model_es.CreateKnowledgeESFactory().RandomDocuments(int(num))
	if err != nil {
		code := errcode.ErrKnowledgeRandomList
		response.Fail(context, code, errcode.ErrMsg[code], errcode.ErrMsgForUser[code])
	} else {
		response.Success(context, consts.CurdStatusOkMsg, knowledgeList)
	}

}
