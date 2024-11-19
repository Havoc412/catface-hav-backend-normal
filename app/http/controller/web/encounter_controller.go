package web

import (
	"catface/app/global/consts"
	"catface/app/global/errcode"
	"catface/app/global/variable"
	"catface/app/http/validator/core/data_transfer"
	"catface/app/model"
	"catface/app/model_es"
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

		srcPath := filepath.Join(variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath"), "encounterPhotos", "hum_"+userId, avatar)
		dstPath := filepath.Join(variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath"), "encounterAvatar", avatar)
		avatarHeight, err := upload_file.ResizeImage(srcPath, dstPath, int(avatarWidth))
		if err != nil {
			response.Fail(context, consts.FilesUploadFailCode, consts.FilesUploadFailMsg, "")
			return
		}
		context.Set(consts.ValidatorPrefix+"avatar", avatar)
		context.Set(consts.ValidatorPrefix+"avatar_height", float64(avatarHeight))
		context.Set(consts.ValidatorPrefix+"avatar_width", float64(avatarWidth))
	}
	poi := context.GetStringMap(consts.ValidatorPrefix + "poi")
	if poi != nil {
		// 感觉这里就是获取信息之后，然后解析后再存储，方便后续 Model 直接绑定到数据。
		latitude := poi["latitude"].(float64)
		longitude := poi["longitude"].(float64)
		context.Set(consts.ValidatorPrefix+"latitude", latitude)
		context.Set(consts.ValidatorPrefix+"longitude", longitude)
	}
	extra := context.GetStringMap(consts.ValidatorPrefix + "extra")
	var tags []string
	if extra != nil {
		context.Set(consts.ValidatorPrefix+"topics", extra["topics"])
		tags = data_transfer.GetStringSlice(context, "topics")
		context.Set(consts.ValidatorPrefix+"tags_list", tags)
		context.Set(consts.ValidatorPrefix+"tags", tags) // INFO 这里字段没有直接匹配上。
	}
	// STAGE - 2: Make []string to String
	if res, err := data_transfer.ConvertSliceToString(photos); err == nil {
		context.Set(consts.ValidatorPrefix+"photos", res)
	} else {
		response.Fail(context, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, "")
		return
	}
	if res, err := data_transfer.ConvertSliceToString(tags); err == nil {
		context.Set(consts.ValidatorPrefix+"tags", res)
	} else {
		response.Fail(context, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, "")
		return
	}
	// STAGE -3: Real Insert - 1: ENC
	animals_id := data_transfer.GetFloat64Slice(context, "animals_id") // 由于是 Slice 就交给 EAlink 内部遍历时处理。
	// Real Insert - 2: EA LINK
	if encounter, ok := model.CreateEncounterFactory("").InsertDate(context); ok {
		// 2: EA Links; // TIP 感觉直接使用 go 会直接且清晰。
		go model.CreateEncounterAnimalLinkFactory("").Insert(encounter.Id, animals_id)

		// 3. ES speed // TODO 这里如何实现 不同 DB 之间的 “事务” 概念。
		if level := int(context.GetFloat64(consts.ValidatorPrefix + "level")); level > 1 {
			go model_es.CreateEncounterESFactory(&encounter).InsertDocument()
		}

		response.Success(context, consts.CurdStatusOkMsg, gin.H{
			"encounter_id": encounter.Id,
		})
	} else {
		response.Fail(context, consts.CurdCreatFailCode, consts.CurdCreatFailMsg+", 新增错误", "")
	}
}

func (e *Encounters) List(context *gin.Context) {
	num := context.GetFloat64(consts.ValidatorPrefix + "num")
	skip := context.GetFloat64(consts.ValidatorPrefix + "skip")
	user_id := context.GetFloat64(consts.ValidatorPrefix + "user_id")
	mode := context.GetString(consts.ValidatorPrefix + "mode")

	encounters := curd.CreateEncounterCurdFactory().List(int(num), int(skip), int(user_id), mode)
	if encounters != nil {
		response.Success(context, consts.CurdStatusOkMsg, encounters)
	} else {
		code := errcode.ErrEncounterNoData
		response.Fail(context, code, errcode.ErrMsg[code], errcode.ErrMsgForUser[code])
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
