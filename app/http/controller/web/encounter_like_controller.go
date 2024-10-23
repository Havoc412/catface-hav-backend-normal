package web

import (
	"catface/app/global/consts"
	"catface/app/model"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

type EncounterLike struct {
}

func (e *EncounterLike) Create(context *gin.Context) {
	userId := context.GetFloat64(consts.ValidatorPrefix + "user_id")
	encounterId := context.GetFloat64(consts.ValidatorPrefix + "encounter_id")

	if model.CreateEncounterLikeFactory("").Create(int(userId), int(encounterId)) {
		response.Success(context, "点赞成功", nil)
	} else {
		response.Fail(context, consts.CurdCreatFailCode, consts.CurdCreatFailMsg+",新增错误", "")

	}
}

func (e *EncounterLike) Delete(context *gin.Context) {
	userId := context.GetFloat64(consts.ValidatorPrefix + "user_id")
	encounterId := context.GetFloat64(consts.ValidatorPrefix + "encounter_id")

	if model.CreateEncounterLikeFactory("").SoftDelete(int(userId), int(encounterId)) {
		response.Success(context, "取消点赞成功", nil)
	} else {
		response.Fail(context, consts.CurdDeleteFailCode, consts.CurdDeleteFailMsg+",删除错误", "")

	}
}
