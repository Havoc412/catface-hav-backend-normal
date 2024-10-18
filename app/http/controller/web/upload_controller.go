package web

import (
	"catface/app/global/consts"
	"catface/app/global/variable"
	"catface/app/service/upload_file"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

type Upload struct {
}

//	文件上传是一个独立模块，给任何业务返回文件上传后的存储路径即可。
//
// 开始上传
func (u *Upload) StartUpload(context *gin.Context) {
	// TODO 如果之后要存储到 Linux 服务器上特殊路径下，就需要修改这里。
	dir_name := context.GetString(consts.ValidatorPrefix + "dir_name")
	// 检查 dir_name 是否以 / 结尾，如果不是则补充一个 /
	if len(dir_name) > 0 && dir_name[len(dir_name)-1] != '/' {
		dir_name += "/"
	}

	savePath := variable.BasePath + variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath") + dir_name
	if r, finnalSavePath := upload_file.Upload(context, savePath); r == true {
		response.Success(context, consts.CurdStatusOkMsg, finnalSavePath)
	} else {
		response.Fail(context, consts.FilesUploadFailCode, consts.FilesUploadFailMsg, "")
	}
}
