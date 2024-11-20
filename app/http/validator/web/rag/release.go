package rag

import (
	"catface/app/global/consts"
	"catface/app/http/controller/web"
	"catface/app/http/validator/core/data_transfer"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

type Release struct {
	Token string `form:"token" json:"token" binding:"required"`
}

func (r Release) CheckParams(context *gin.Context) {
	if err := context.ShouldBind(&r); err != nil {
		response.ValidatorError(context, err)
		return
	}
	extraAddBindDataContext := data_transfer.DataAddContext(r, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "RAG RELEASE 表单验证器json化失败", "")
	} else {
		(&web.Rag{}).Release(extraAddBindDataContext)
	}
}
