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
	// STAGE 1. 插入 MySQL 记录。
	var doc_id int64
	ok := false
	if doc_id, ok = model.CreateDocFactory("").InsertDocumentData(context); !ok {
		response.Fail(context, consts.CurdCreatFailCode, consts.CurdCreatFailMsg, "上传文档错误")
	}

	// STAGE 2. 调用 python API
	path := context.GetString(consts.ValidatorPrefix + "path")
	filePath := filepath.Join(variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath"), variable.ConfigYml.GetString("FileUploadSetting.DocsRootPath"), path)

	// TODO
	_ = filePath
	_ = doc_id
	response.Success(context, consts.CurdStatusOkMsg, "")
}
