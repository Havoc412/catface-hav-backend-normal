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
	// TODO 这里还需要处理一下历史记录？
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
		(&web.Rag{}).Chat(extraAddBindDataContext)
	}

}
