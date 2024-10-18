package web

import (
	"github.com/gin-gonic/gin"
)

type Encounters struct {
}

func (e *Encounters) Create(context *gin.Context) {
	// TODO 处理 Photos 文件，然后处理出 Avatar，并获取压缩后的 宽高，以及文件的存储路径。

	// Real Insert
	// if model.CreateEncounterFactory("").InsertDate(context) {
	// 	response.Success(context, consts.CurdStatusOkMsg, "")
	// } else {
	// 	response.Fail(context, consts.CurdCreatFailCode, consts.CurdCreatFailMsg+",新增错误", "")
	// }
}
