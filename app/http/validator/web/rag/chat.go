package rag

import (
	"catface/app/global/consts"
	"catface/app/http/controller/web"
	"catface/app/http/validator/core/data_transfer"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

// INFO 虽然起名为 Chat，但是默认就会去查询 知识库，也就是不作为一般的 LLM-chat 来使用。
type Chat struct {
	Query string `form:"query" json:"query" binding:"required"`
	Token string `form:"token" json:"token"` // UPDATE 暂时不想启用 user 的 token，就先单独处理。

	Mode   string `form:"mode" json:"mode"`
	CatsId string `form:"cats_id" json:"cats_id"` //
}

func (c Chat) CheckParams(context *gin.Context) {
	if err := context.ShouldBind(&c); err != nil {
		response.ValidatorError(context, err)
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(c, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "RAG CHAT 表单验证器json化失败", "")
	} else {
		(&web.Rag{}).ChatWebSocket(extraAddBindDataContext)
	}

}
