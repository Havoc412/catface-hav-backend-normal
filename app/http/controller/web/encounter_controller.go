package web

import (
	"catface/app/global/consts"
	"catface/app/global/errcode"
	"catface/app/global/variable"
	"catface/app/http/validator/core/data_transfer"
	"catface/app/model"
	"catface/app/service/encounter/curd"
	"catface/app/service/upload_file"
	"catface/app/utils/response"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Encounters struct {
}

func (e *Encounters) Create(context *gin.Context) {
	photos := data_transfer.GetStringSlice(context, "photos")
	if len(photos) > 0 {
		userId := strconv.Itoa(int(context.GetFloat64(consts.ValidatorPrefix + "user_id")))
		avatar := photos[0]
		avatarWidth := variable.ConfigYml.GetFloat64("FileUploadSetting.AvatarWidth")

		srcPath := filepath.Join(variable.BasePath, variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath"), "encounterPhotos", "hum_"+userId, avatar)
		dstPath := filepath.Join(variable.BasePath, variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath"), "encounterAvatar", avatar)
		avatarHeight, err := upload_file.ResizeImage(srcPath, dstPath, int(avatarWidth))
		if err != nil {
			response.Fail(context, consts.FilesUploadFailCode, consts.FilesUploadFailMsg, "")
			return
		}
		context.Set(consts.ValidatorPrefix+"avatar", avatar)
		context.Set(consts.ValidatorPrefix+"avatar_height", float64(avatarHeight))
		context.Set(consts.ValidatorPrefix+"avatar_width", float64(avatarWidth))
	}
	// 将 Array 转化为 string 类型
	animals_id := data_transfer.GetFloat64Slice(context, "animals_id")
	if res, err := data_transfer.ConvertSliceToString(animals_id); err == nil {
		context.Set(consts.ValidatorPrefix+"animals_id", res)
	} else {
		response.Fail(context, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, "")
		return
	}
	if res, err := data_transfer.ConvertSliceToString(photos); err == nil {
		context.Set(consts.ValidatorPrefix+"photos", res)
	} else {
		response.Fail(context, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, "")
		return
	}
	// Real Insert
	if model.CreateEncounterFactory("").InsertDate(context) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdCreatFailCode, consts.CurdCreatFailMsg+",新增错误", "")
	}
}

func (e *Encounters) List(context *gin.Context) {
	num := context.GetFloat64(consts.ValidatorPrefix + "num")
	skip := context.GetFloat64(consts.ValidatorPrefix + "skip")
	user_id := context.GetFloat64(consts.ValidatorPrefix + "user_id")

	encounters := curd.CreateEncounterCurdFactory().List(int(num), int(skip), int(user_id))
	if encounters != nil {
		response.Success(context, consts.CurdStatusOkMsg, encounters)
	} else {
		response.Fail(context, errcode.ErrDataNoFound, errcode.ErrMsg[errcode.ErrDataNoFound], "")
	}
}

func (e *Encounters) Detail(context *gin.Context) {
	encounterId := context.Param("encounter_id")

	encounters := curd.CreateEncounterCurdFactory().Detail(encounterId)
	if encounters != nil {
		response.Success(context, consts.CurdStatusOkMsg, encounters)
	} else {
		response.Fail(context, errcode.ErrDataNoFound, errcode.ErrMsg[errcode.ErrDataNoFound], "")
	}
}
