package web

import (
	"catface/app/global/consts"
	"catface/app/model"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

type AnimalLike struct {
}

func (a *AnimalLike) Create(context *gin.Context) {
	userId := context.GetFloat64(consts.ValidatorPrefix + "user_id")
	animalId := context.GetFloat64(consts.ValidatorPrefix + "animal_id")
	if model.CreateAnimalLikeFactory("").Create(int(userId), int(animalId)) {
		response.Success(context, "关注成功", nil)
	} else {
		response.Fail(context, consts.CurdCreatFailCode, consts.CurdCreatFailMsg+",新增错误", "")
	}
}

func (a *AnimalLike) Delete(context *gin.Context) {
	userId := context.GetFloat64(consts.ValidatorPrefix + "user_id")
	animalId := context.GetFloat64(consts.ValidatorPrefix + "animal_id")
	if model.CreateAnimalLikeFactory("").SoftDelete(int(userId), int(animalId)) {
		response.Success(context, "取消关注成功", nil)
	} else {
		response.Fail(context, consts.CurdDeleteFailCode, consts.CurdDeleteFailMsg+",删除错误", "")
	}
}
