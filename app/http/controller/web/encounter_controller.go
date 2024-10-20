package web

import (
	"catface/app/global/variable"
	"catface/app/http/validator/core/data_transfer"

	"github.com/gin-gonic/gin"
)

type Encounters struct {
}

func (e *Encounters) Create(context *gin.Context) {
	// TODO 处理 Photos 文件，然后处理出 Avatar，并获取压缩后的 宽高，以及文件的存储路径。
	photos := data_transfer.GetStringSlice(context, "photos")
	if len(photos) > 0 {
		avatar := photos[0]
		avatarWidth := variable.ConfigYml.GetFloat64("FileUploadSetting.AvatarWidth")

	}
	// Real Insert
	// if model.CreateEncounterFactory("").InsertDate(context) {
	// 	response.Success(context, consts.CurdStatusOkMsg, "")
	// } else {
	// 	response.Fail(context, consts.CurdCreatFailCode, consts.CurdCreatFailMsg+",新增错误", "")
	// }
}
