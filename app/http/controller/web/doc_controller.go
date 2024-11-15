package web

import (
	"catface/app/global/consts"
	"catface/app/global/variable"
	"catface/app/model"
	"catface/app/utils/response"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type Docs struct {
}

func (d *Docs) Upload(context *gin.Context) {
	// TODO 1. 读取源文件，调用 py API 分块上传。
	path := context.GetString(consts.ValidatorPrefix + "path")
	filePath := filepath.Join(variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath"), variable.ConfigYml.GetString("FileUploadSetting.DocsRootPath"), path)
	_ = filePath

	// STAGE 2.
	if ok := model.CreateDocFactory("").InsertDocumentData(context); ok {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdCreatFailCode, consts.CurdCreatFailMsg, "上传文档错误")
	}
}
