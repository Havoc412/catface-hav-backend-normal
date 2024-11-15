package upload_files

import (
	"catface/app/global/consts"
	"catface/app/global/variable"
	"catface/app/http/controller/web"
	"catface/app/http/validator/core/data_transfer"
	"catface/app/utils/files"
	"catface/app/utils/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UpFiles struct {
	DirName string `form:"dir_name" json:"dir_name"`
}

// 文件上传公共模块表单参数验证器
func (u UpFiles) CheckParams(context *gin.Context) {
	// 1.基本的验证规则没有通过
	if err := context.ShouldBind(&u); err != nil {
		response.ValidatorError(context, err)
		return
	}
	// 该函数主要是将本结构体的字段（成员）按照 consts.ValidatorPrefix+ json标签对应的 键 => 值 形式绑定在上下文，便于下一步（控制器）可以直接通过 context.Get(键) 获取相关值
	extraAddBindDataContext := data_transfer.DataAddContext(u, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "upload Files 表单验证器json化失败", "")
		return
	}

	// 2. File 内容的验证；
	tmpFile, err := context.FormFile(variable.ConfigYml.GetString("FileUploadSetting.UploadFileField")) //  file 是一个文件结构体（文件对象）
	//获取文件发生错误，可能上传了空文件等
	if err != nil {
		response.Fail(context, consts.FilesUploadFailCode, consts.FilesUploadFailMsg, err.Error())
		return
	}
	if tmpFile.Size == 0 {
		response.Fail(context, consts.FilesUploadMoreThanMaxSizeCode, consts.FilesUploadIsEmpty, "")
		return
	}

	//超过系统设定的最大值：32M，tmpFile.Size 的单位是 bytes 和我们定义的文件单位M 比较，就需要将我们的单位*1024*1024(即2的20次方)，一步到位就是 << 20
	sizeLimit := variable.ConfigYml.GetInt64("FileUploadSetting.Size")
	if tmpFile.Size > sizeLimit<<20 {
		response.Fail(context, consts.FilesUploadMoreThanMaxSizeCode, consts.FilesUploadMoreThanMaxSizeMsg+strconv.FormatInt(sizeLimit, 10)+"M", "")
		return
	}
	//不允许的文件mime类型
	var isPass bool
	var mimeType string
	if fp, err := tmpFile.Open(); err == nil {
		mimeType = files.GetFilesMimeByFp(fp)

		for _, value := range variable.ConfigYml.GetStringSlice("FileUploadSetting.AllowMimeType") {
			if strings.ReplaceAll(value, " ", "") == strings.ReplaceAll(mimeType, " ", "") {
				isPass = true
				break
			}
		}
		_ = fp.Close()
	} else {
		response.ErrorSystem(context, consts.ServerOccurredErrorMsg, "")
		return
	}
	//凡是存在相等的类型，通过验证，调用控制器
	if !isPass {
		response.Fail(context, consts.FilesUploadMimeTypeFailCode, consts.FilesUploadMimeTypeFailMsg, gin.H{
			"mime_type": mimeType,
		})
	} else {
		(&web.Upload{}).StartUpload(context)
	}
}
