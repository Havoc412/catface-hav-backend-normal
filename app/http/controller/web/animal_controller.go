package web

import (
	"catface/app/global/consts"
	"catface/app/global/errcode"
	"catface/app/global/variable"
	"catface/app/http/validator/core/data_transfer"
	"catface/app/model"
	"catface/app/model_es"
	"catface/app/model_redis"
	"catface/app/service/animals/curd"
	"catface/app/service/upload_file"
	"catface/app/utils/query_handler"
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
	department := context.GetString(consts.ValidatorPrefix + "department")
	num := int(context.GetFloat64(consts.ValidatorPrefix + "num"))
	skip := int(context.GetFloat64(consts.ValidatorPrefix + "skip"))
	userId := context.GetFloat64(consts.ValidatorPrefix + "user_id")

	mode := context.GetString(consts.ValidatorPrefix + "mode")

	// TAG prefer MODE 查询模式。
	var redis_selctedCatsId model_redis.SelectedAnimal4Prefer
	var animalsWithLike []model.AnimalWithLikeList
	if mode == consts.AnimalPreferMode {
		key := int64(context.GetFloat64(consts.ValidatorPrefix + "key"))

		if key != 0 {
			if ok, err := redis_selctedCatsId.GetDataByKey(key); !ok {
				_ = err // TODO
			}
		} else {
			redis_selctedCatsId.GenerateKey()
		}

		if redis_selctedCatsId.Length() == skip {
			preferCats, _ := getPreferCats(int(userId), num, attrs, &redis_selctedCatsId)
			if len(preferCats) > 0 {
				animalsWithLike = preferCats
			}

			// TODO 刷新 Redis 有效期
			if ok, err := redis_selctedCatsId.SetDataByKey(); !ok {
				_ = err // TODO
			}
		}
	}

	// 计算还需要多少毛茸茸
	num -= len(animalsWithLike)
	skip = max(0, skip-redis_selctedCatsId.Length())
	if num > 0 {
		additionalAnimals := curd.CreateAnimalsCurdFactory().List(attrs, gender, breed, sterilization, status, department, redis_selctedCatsId.GetAllIds(), num, skip, int(userId))
		// 将 additionalAnimals 整合到 animalsWithLike 的后面
		animalsWithLike = append(animalsWithLike, additionalAnimals...)
	}

	if animalsWithLike != nil {
		response.Success(context, consts.CurdStatusOkMsg, gin.H{
			"animals": animalsWithLike,
			"key":     redis_selctedCatsId.GetKey(),
		})
	} else {
		response.Fail(context, errcode.AnimalNoFind, errcode.ErrMsg[errcode.AnimalNoFind], errcode.ErrMsgForUser[errcode.AnimalNoFind])
	}
}

/**
 * @description: 在常规条件过滤查询之前，通过一定规则获取偏好的目标。
 * @param {*} userId
 * @param {int} num
 * @param {string} attrs
 * @param {model_redis.SelectedAnimal4Prefer} redis 使用单独的结构体，更清晰的控制中间状态。
 * @return {*}
 */
func getPreferCats(userId, num int, attrs string, redis *model_redis.SelectedAnimal4Prefer) (list []model.AnimalWithLikeList, err error) {
	// STAGE #1 无视过滤条件，获取路遇“过”的 id 列表；先获取 ID，然后再去查询细节信息。
	ids, err := model.CreateEncounterFactory("").EncounteredCatsId(userId, num, redis.NumEnc(), redis.NewCatsId)

	// STAGE #2 获取近期新增的毛茸茸；只在第一次操作 && 数量不够时 启用。
	var idsNew []int64
	if redis.NotPassNew() { // UPDATE && len(ids) < num 调整为第一次访问时一定会次优先返回 NewCats 的推荐。
		// 获取近期新增的毛茸茸
		newCats, _ := model.CreateAnimalFactory("").NewCatsId(3, 0) // INFO 硬编码获取最新的 3 个。

		// 去重：默认只会在首次查询时启用，所以只需要去重 STAGE#1 的 ids。
		uniqueIds := make(map[int64]bool)
		for _, id := range ids {
			uniqueIds[id] = true
		}

		for _, id := range newCats {
			if _, ok := uniqueIds[id]; !ok {
				idsNew = append(idsNew, id)
			}
		}
	}

	// 3. 合并然后查询数据 && 处理 redis 之间的处理。
	redis.AppendEncIds(ids)
	if len(idsNew) > 0 {
		ids = append(ids, idsNew...)
		redis.NewCatsId = idsNew
	}

	if err == nil && len(ids) > 0 {
		attrsSlice := query_handler.StringToStringArray(attrs)
		attrsSlice = append(attrsSlice, "id")

		animalMap := make(map[int64]model.Animal, len(ids))
		animals := model.CreateAnimalFactory("").ShowByIDs(ids, attrsSlice...)

		for _, v := range animals {
			animalMap[v.Id] = v
		}

		// 根据 preferCatsId 的顺序重构最终结果列表
		for _, id := range ids {
			if animal, ok := animalMap[id]; ok {
				list = append(list, model.AnimalWithLikeList{Animal: animal})
			}
		}
	}

	return
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

		srcPath := filepath.Join(variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath"), "catsPhotos", "hum_"+userId, avatar)
		dstPath := filepath.Join(variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath"), "catsAvatar", avatar)
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
		context.Set(consts.ValidatorPrefix+"sterilization", sterilization)
		context.Set(consts.ValidatorPrefix+"vaccination", vaccination)
		context.Set(consts.ValidatorPrefix+"deworming", deworming)
	}
	extra := context.GetStringMap(consts.ValidatorPrefix + "extra")
	var nickNames []string
	var tags []string
	if extra != nil {
		context.Set(consts.ValidatorPrefix+"nick_names", extra["nick_names"])
		context.Set(consts.ValidatorPrefix+"tags", extra["tags"])
		nickNames = data_transfer.GetStringSlice(context, "nick_names")
		tags = data_transfer.GetStringSlice(context, "tags")
		context.Set(consts.ValidatorPrefix+"nick_names_list", nickNames)
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
	if animal, ok := model.CreateAnimalFactory("").InsertDate(context); ok {
		// 转移 photos 到 anm；采用 rename dir 的方式
		oldName := filepath.Join(variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath"), "catsPhotos", "hum_"+userId)
		newName := filepath.Join(variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath"), "catsPhotos", "anm_"+strconv.FormatInt(animal.Id, 10))
		err := os.Rename(oldName, newName)
		if err != nil {
			// TODO 特殊返回，成功了一半？或者需要清空原有的操作？不过感觉这一步几乎不会出错。
			// TODO 或许直接采用 go 会比较好呢？
		}

		// 2. 将部分数据插入 ES；
		go model_es.CreateAnimalESFactory(&animal).InsertDocument()

		response.Success(context, consts.CurdStatusOkMsg, gin.H{
			"anm_id": animal.Id,
		})
	} else {
		response.Fail(context, consts.CurdCreatFailCode, consts.CurdCreatFailMsg+",新增错误", "")
	}
}

func (a *Animals) Name(context *gin.Context) {
	attrs := context.GetString(consts.ValidatorPrefix + "attrs")
	name := context.GetString(consts.ValidatorPrefix + "name")

	animals := curd.CreateAnimalsCurdFactory().ShowByName(attrs, name)
	if animals != nil {
		response.Success(context, consts.CurdStatusOkMsg, animals)
	} else {
		response.Fail(context, errcode.AnimalNoFind, errcode.ErrMsg[errcode.AnimalNoFind], "")
	}
}
