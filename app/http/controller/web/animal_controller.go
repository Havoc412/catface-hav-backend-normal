package web

import (
	"catface/app/global/consts"
	"catface/app/global/errcode"
	"catface/app/global/variable"
	"catface/app/http/validator/core/data_transfer"
	"catface/app/model"
	"catface/app/service/animals/curd"
	"catface/app/service/upload_file"
	"catface/app/utils/response"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

)

type Animals struct { // INFO 起到一个标记的作用，这样 web.xxx 的时候不同模块就不会命名冲突了。
}

func (a *Animals) List(context *gin.Context) {
	// 1. Get Params
	attrs := context.GetString(consts.ValidatorPrefix + "attrs")
	gender := context.GetString(consts.ValidatorPrefix + "gender")
	breed := context.GetString(consts.ValidatorPrefix + "breed")
	sterilization := context.GetString(consts.ValidatorPrefix + "sterilization")
	status := context.GetString(consts.ValidatorPrefix + "status")
	num := context.GetFloat64(consts.ValidatorPrefix + "num")
	skip := context.GetFloat64(consts.ValidatorPrefix + "skip")
	userId := context.GetFloat64(consts.ValidatorPrefix + "user_id")

	animals := curd.CreateAnimalsCurdFactory().List(attrs, gender, breed, sterilization, status, int(num), int(skip), int(userId))
	if animals != nil {
		response.Success(context, consts.CurdStatusOkMsg, animals)
	} else {
		response.Fail(context, errcode.AnimalNoFind, errcode.ErrMsg[errcode.AnimalNoFind], "")
	}
}

// v0.1
// func (a *Animals) Detail(context *gin.Context) {
// 	// 1. Get Params
// 	anmId, err := strconv.Atoi(context.Param("anm_id"))
// 	// 2. Select & Filter
// 	var animal model.Animal
// 	err = variable.GormDbMysql.Table("animals").Model(&animal).Where("id = ?", anmId).Scan(&animal).Error // TIP GORM.First 采取默认的
// 	if err != nil {
// 		response.Fail(context, errcode.ErrAnimalSqlFind, errcode.ErrMsg[errcode.ErrAnimalSqlFind], err) // UPDATE consts ?
// 	} else {
// 		response.Success(context, consts.CurdStatusOkMsg, animal)
// 	}
// }

func (a *Animals) Detail(context *gin.Context) {
	// 1. Get Params
	anmId := context.Param("anm_id")

	animal := curd.CreateAnimalsCurdFactory().Detail(anmId)
	if animal != nil {
		response.Success(context, consts.CurdStatusOkMsg, animal)
	} else {
		response.Fail(context, errcode.AnimalNoFind, errcode.ErrMsg[errcode.AnimalNoFind], "")
	}
}
func (a *Animals) Create(context *gin.Context) {
	userId := strconv.Itoa(int(context.GetFloat64(consts.ValidatorPrefix + "user_id")))
	// STAGE-1 Get Params
	photos := data_transfer.GetStringSlice(context, "photos")
	if len(photos) > 0 {
		avatar := photos[0]
		avatarWidth := variable.ConfigYml.GetFloat64("FileUploadSetting.AvatarWidth")

		srcPath := filepath.Join(variable.BasePath, variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath"), "catsPhotos", "hum_"+userId, avatar)
		dstPath := filepath.Join(variable.BasePath, variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath"), "catsAvatar", avatar)
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
	health := context.GetStringMap(consts.ValidatorPrefix + "health")
	if health != nil {
		sterilization := health["sterilization"].(float64)
		vaccination := health["vaccination"].(float64)
		deworming := health["deworming"].(float64)
		context.Set(consts.ValidatorPrefix+"sterilization", uint8(sterilization))
		context.Set(consts.ValidatorPrefix+"vaccination", uint8(vaccination))
		context.Set(consts.ValidatorPrefix+"deworming", uint8(deworming))
	}
	extra := context.GetStringMap(consts.ValidatorPrefix + "extra")
	var nickNames []string
	var tags []string
	if extra != nil {
		context.Set(consts.ValidatorPrefix+"nick_names", extra["nick_names"])
		context.Set(consts.ValidatorPrefix+"tags", extra["tags"])
		nickNames = data_transfer.GetStringSlice(context, "nick_names")
		tags = data_transfer.GetStringSlice(context, "tags")
		context.Set(consts.ValidatorPrefix+"nick_names", nickNames)
		context.Set(consts.ValidatorPrefix+"tags", tags) // UPDATE 有点冗余，但是不用复杂代码；
	}
	// STAGE-2
	if res, err := data_transfer.ConvertSliceToString(photos); err == nil {
		context.Set(consts.ValidatorPrefix+"photos", res)
	} else {
		response.Fail(context, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, "")
		return
	}
	if res, err := data_transfer.ConvertSliceToString(nickNames); err == nil {
		context.Set(consts.ValidatorPrefix+"nick_names", res)
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
	// STAGE-3
	if anm_id, ok := model.CreateAnimalFactory("").InsertDate(context); ok {
		// 转移 photos 到 anm；采用 rename dir 的方式
		oldName := filepath.Join(variable.BasePath, variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath"), "catsPhotos", "hum_"+userId)
		newName := filepath.Join(variable.BasePath, variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath"), "catsPhotos", "anm_"+strconv.FormatInt(anm_id, 10))
		err := os.Rename(oldName, newName)
		if err != nil {
			// TODO 特殊返回，成功了一半？或者需要清空原有的操作？不过感觉这一步几乎不会出错。
		}
		response.Success(context, consts.CurdStatusOkMsg, gin.H{
			"anm_id": anm_id,
		})
	} else {
		response.Fail(context, consts.CurdCreatFailCode, consts.CurdCreatFailMsg+",新增错误", "")
	}
}
