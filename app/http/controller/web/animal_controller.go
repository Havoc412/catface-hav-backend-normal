package web

import (
	"catface/app/global/consts"
	"catface/app/global/errcode"
	"catface/app/service/animals/curd"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

type Animals struct { // INFO 起到一个标记的作用，这样 web.xxx 的时候不同模块就不会命名冲突了。
}

func (a *Animals) List(context *gin.Context) {
	// 1. Get Params
	attrs := context.GetString(consts.ValidatorPrefix + "attrs")
	gender := context.GetString(consts.ValidatorPrefix + "gender")
	breed := context.GetString(consts.ValidatorPrefix + "breed")
	sterilzation := context.GetString(consts.ValidatorPrefix + "sterilzation")
	status := context.GetString(consts.ValidatorPrefix + "status")
	num := context.GetFloat64(consts.ValidatorPrefix + "num")
	skip := context.GetFloat64(consts.ValidatorPrefix + "skip")

	animals := curd.CreateUserCurdFactory().List(attrs, gender, breed, sterilzation, status, int(num), int(skip))
	if animals != nil {
		response.Success(context, consts.CurdStatusOkMsg, animals)
	} else {
		response.Fail(context, errcode.AnimalNoFind, errcode.ErrMsg[errcode.AnimalNoFind], "")
	}
}

// v1
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

	animal := curd.CreateUserCurdFactory().Detail(anmId)
	if animal != nil {
		response.Success(context, consts.CurdStatusOkMsg, animal)
	} else {
		response.Fail(context, errcode.AnimalNoFind, errcode.ErrMsg[errcode.AnimalNoFind], "")
	}
}
